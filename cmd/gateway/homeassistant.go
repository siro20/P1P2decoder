package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/siro20/p1p2decoder/pkg/p1p2"
)

type HomeAssistant struct {
	cfg     HomeAssistantConfig
	httpCon map[string]*http.Request
	Alive   bool
}

type HomeAssistantConfig struct {
	Enable      bool   `yaml:"enable"`
	Hostname    string `yaml:"hostname"`
	Port        int    `yaml:"port"`
	BearerToken string `yaml:"bearer"`
}

type HomeAssistantAttributes struct {
	UnitOfMeasurement string `json:"unit_of_measurement,omitempty"`
	FriendlyName      string `json:"friendly_name,omitempty"`
	Icon              string `json:"icon,omitempty"`
	DeviceClass       string `json:"device_class,omitempty"`
}

type HomeAssistantState struct {
	State      string                  `json:"state"`
	Attributes HomeAssistantAttributes `json:"attributes,omitempty"`
}

func NewHomeAssistant(c HomeAssistantConfig) (h *HomeAssistant, err error) {
	h = &HomeAssistant{cfg: c,
		httpCon: map[string]*http.Request{}}

	return h, err
}

func (h *HomeAssistant) Serve() {
	for {
		h.checkAlive()
		time.Sleep(time.Minute * 10)
	}
}

func (h *HomeAssistant) checkAlive() {
	var err error
	h.Alive, err = h.CheckCfgValid(h.cfg)

	if !h.Alive || err != nil {
		for i := range h.httpCon {
			h.httpCon[i].Close = true
		}
		h.httpCon = map[string]*http.Request{}
	}
}

func (h *HomeAssistant) CheckCfgValid(cfg HomeAssistantConfig) (alive bool, err error) {
	// Create a new request using http
	req, err := http.NewRequest("GET", "http://"+cfg.Hostname+":"+strconv.Itoa(cfg.Port)+"/api/", nil)
	if err != nil {
		log.Printf("Error on NewRequest %v", err)
		return false, err
	}
	// add authorization header to the req
	req.Header.Add("Authorization", "Bearer "+cfg.BearerToken)
	req.Header.Set("Content-Type", "application/json")
	req.Close = true
	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
		return false, err
	}
	resp.Body.Close()
	if err == nil {
		return true, nil
	}

	return false, nil
}

func (h *HomeAssistant) DoConnect(id string) (con *http.Request, err error) {
	h.httpCon[id], err = http.NewRequest("POST", "http://"+h.cfg.Hostname+":"+strconv.Itoa(h.cfg.Port)+"/api/states/"+id, nil)
	if err != nil {
		h.httpCon = nil
		return nil, err

	}
	// add authorization header to the req
	h.httpCon[id].Header.Add("Authorization", "Bearer "+h.cfg.BearerToken)
	h.httpCon[id].Header.Set("Content-Type", "application/json")

	con = h.httpCon[id]

	return
}

func (h *HomeAssistant) DoPost(id string, payload []byte) (err error) {
	con, ok := h.httpCon[id]
	if !ok || con == nil {
		con, err = h.DoConnect(id)
		if err != nil {
			return
		}
	}
	if con != nil {
		con.Body = io.NopCloser(bytes.NewBuffer(payload))
		// Send req using http Client
		client := &http.Client{
			Timeout: 10 * time.Second,
		}
		resp, err := client.Do(con)
		if err != nil {
			con.Close = true
			h.httpCon[id] = nil
			go h.checkAlive()
			return err
		}
		resp.Body.Close()
	}
	return
}

func (h *HomeAssistant) SendSensor(name string, binary bool, s HomeAssistantState) (err error) {
	if !h.Alive {
		err = fmt.Errorf("Remote seems offline. Skipping POST...")
		return
	}
	entity_id := "sensor."
	if binary {
		entity_id = "binary_sensor."
		if strings.ToLower(s.State) == "true" || strings.ToLower(s.State) == "on" || s.State == "1" {
			s.State = "on"
		} else {
			s.State = "off"
		}
	}
	entity_id += name

	jsonStr, err := json.Marshal(s)
	if err != nil {
		return err
	}

	err = h.DoPost(entity_id, jsonStr)
	fmt.Printf("HA: Sending sensor %s. State: %s Error: %v\n", name, s.State, err)

	return
}

func DeviceClassFromSensor(s p1p2.Sensor) string {
	if s.Type() == "gauge" {
		return "temperature"
	} else if s.Type() == "valve" {
		return "opening"
	} else if s.Type() == "state" {
		if s.ID() == p1p2.StatePower.ID() || s.ID() == p1p2.StateDHW.ID() || s.ID() == p1p2.StateDHWEnable.ID() {
			return "power"
		}
		if s.ID() == p1p2.StateCompressor.ID() || s.ID() == p1p2.StateGas.ID() || s.ID() == p1p2.PumpDHWCirculation.ID() {
			return "running"
		}
		if s.ID() == p1p2.PumpMain.ID() {
			return "running"
		}
	} else if s.Type() == "time" {
		return "timestamp"
	}
	return ""
}

func HomeAssistantAddSensors(ha *HomeAssistant) {
	f := func(s p1p2.Sensor, value interface{}) {
		deviceClass := DeviceClassFromSensor(s)
		state := HomeAssistantState{
			Attributes: HomeAssistantAttributes{
				UnitOfMeasurement: s.Unit(),
				FriendlyName:      s.Name(),
				Icon:              s.Icon(),
				DeviceClass:       deviceClass,
			},
		}

		if s.Type() == "gauge" {
			newVal, ok := value.(float32)
			if !ok {
				return
			}
			state.State = fmt.Sprintf("%.1f", newVal)
			ha.SendSensor(strings.ToLower(s.Name()), false, state)
		} else if s.Type() == "valve" || s.Type() == "state" {
			newVal, ok := value.(bool)
			if !ok {
				return
			}
			state.State = strconv.FormatBool(newVal)
			ha.SendSensor(strings.ToLower(s.Name()), true, state)
		}
	}

	for i := range p1p2.Sensors {
		if p1p2.Sensors[i].Type() == "gauge" {
			if p1p2.Sensors[i].Name() == "DomesticHotWater" {
				p1p2.Sensors[i].RegisterStateChangedWithHysteresisCallback(0.1, f)
			} else {
				//
				// Temperature sensors precision is less than 0.5°C
				//
				p1p2.Sensors[i].RegisterStateChangedWithHysteresisCallback(1.0, f)
			}
		} else {
			p1p2.Sensors[i].RegisterStateChangedCallback(f)
		}
	}
}
