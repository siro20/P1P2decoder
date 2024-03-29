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
	"sync"
	"time"

	"github.com/siro20/p1p2decoder/pkg/p1p2"
)

var domain string = "daikin_p1p2_"

type HomeAssistant struct {
	cfg     HomeAssistantConfig
	httpCon map[string]*http.Request
	Alive   bool
	Lock    sync.Mutex
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
		h.Lock.Lock()
		for i := range h.httpCon {
			if h.httpCon[i] != nil {
				h.httpCon[i].Close = true
				h.httpCon[i] = nil
			}
		}
		h.httpCon = map[string]*http.Request{}
		h.Lock.Unlock()
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
	con, err = http.NewRequest("POST", "http://"+h.cfg.Hostname+":"+strconv.Itoa(h.cfg.Port)+"/api/states/"+id, nil)
	if err != nil {
		return nil, err
	}

	// add authorization header to the req
	con.Header.Add("Authorization", "Bearer "+h.cfg.BearerToken)
	con.Header.Set("Content-Type", "application/json")

	return
}

func (h *HomeAssistant) DoPost(id string, payload []byte) (err error) {
	h.Lock.Lock()
	con, ok := h.httpCon[id]
	if !ok || con == nil {
		con, err = h.DoConnect(id)
		if err != nil {
			h.Lock.Unlock()
			return
		}
		h.httpCon[id] = con
	}
	h.Lock.Unlock()
	if con != nil {
		con.Body = io.NopCloser(bytes.NewBuffer(payload))
		// Send req using http Client
		client := &http.Client{
			Timeout: 10 * time.Second,
		}
		resp, err := client.Do(con)
		if err != nil {
			con.Close = true
			h.Lock.Lock()
			h.httpCon[id] = nil
			h.Lock.Unlock()
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
	entity_id := "sensor." + domain
	if binary {
		entity_id = "binary_sensor." + domain
		if strings.ToLower(s.State) == "true" || strings.ToLower(s.State) == "on" || s.State == "1" {
			s.State = "on"
		} else if strings.ToLower(s.State) == "false" || strings.ToLower(s.State) == "off" || s.State == "0" {
			s.State = "off"
		}
	}
	entity_id += name

	jsonStr, err := json.Marshal(s)
	if err != nil {
		return err
	}

	err = h.DoPost(entity_id, jsonStr)
	fmt.Printf("HA: Sending sensor %s. State: %s Error: %v\n", entity_id, s.State, err)

	return
}

func DeviceClassFromSensor(s p1p2.Sensor) string {
	if s.Type() == "temperature" {
		return "temperature"
	} else if s.Type() == "gauge" {
		return "gauge"
	} else if s.Type() == "valve" {
		return "opening"
	} else if s.Type() == "energy" {
		return "energy"
	} else if s.Type() == "working_hours" {
		return "duration"
	} else if s.Type() == "count" {
		return ""
	} else if s.Type() == "state" {
		if s.ID() == p1p2.StateHeatingEnabled.ID() || s.ID() == p1p2.StateDHW.ID() || s.ID() == p1p2.StateDHWEnable.ID() || s.ID() == p1p2.StateGasEnabled.ID() {
			return "power"
		}
		if s.ID() == p1p2.StateCompressor.ID() || s.ID() == p1p2.PumpDHWCirculation.ID() || s.ID() == p1p2.StateBoilerRunning.ID() || s.ID() == p1p2.PumpMain.ID() {
			return "running"
		}
	} else if s.Type() == "time" {
		return "timestamp"
	}
	return ""
}

func PrettyNameFromSensor(s p1p2.Sensor) string {
	prefix := "Sensor"
	switch s.Type() {
	case "temperature":
		prefix = "Temperature"
	case "gauge":
		prefix = "Gauge"
	case "valve":
		prefix = "Valve"
	case "pump":
		prefix = "Pump"
	case "software":
		prefix = "Software Version"
	case "time":
		prefix = "Time"
	case "state":
		prefix = ""
	case "count":
		prefix = "Number of"
	case "working_hours":
		prefix = "Working hours"
	case "energy":
		prefix = "Energy"
	}
	return fmt.Sprintf("%s %s", prefix, s.Name())
}

func EntityNameFromSensor(s p1p2.Sensor) string {
	str := fmt.Sprintf("%s_%s", s.Type(), s.Name())
	str = strings.ReplaceAll(str, " ", "_")
	str = strings.ToLower(str)

	return str
}

func UnitFromSensor(s p1p2.Sensor) string {
	unit := s.Unit()
	if strings.ToLower(unit) == "boolean" || strings.ToLower(unit) == "bool" {
		unit = ""
	}
	return unit
}

func SensorToHomeAssistant(ha *HomeAssistant, s p1p2.Sensor, value interface{}) {
	state := HomeAssistantState{
		State: "",
		Attributes: HomeAssistantAttributes{
			UnitOfMeasurement: UnitFromSensor(s),
			FriendlyName:      PrettyNameFromSensor(s),
			Icon:              s.Icon(),
			DeviceClass:       DeviceClassFromSensor(s),
		},
	}
	if s.Type() == "temperature" || s.Type() == "gauge" {
		newVal, ok := value.(float32)
		if ok {
			state.State = fmt.Sprintf("%.1f", newVal)
		}
		ha.SendSensor(EntityNameFromSensor(s), false, state)
	} else if s.Type() == "valve" || s.Type() == "state" || s.Type() == "pump" {
		newVal, ok := value.(bool)
		if ok {
			state.State = strconv.FormatBool(newVal)
		}
		ha.SendSensor(EntityNameFromSensor(s), true, state)
	} else if s.Type() == "software" {
		newVal, ok := value.(string)
		if ok {
			state.State = newVal
		}
		state.State = newVal
		ha.SendSensor(EntityNameFromSensor(s), false, state)
	} else if s.Type() == "working_hours" || s.Type() == "count" || s.Type() == "energy" {
		newVal, ok := value.(int)
		if ok {
			state.State = fmt.Sprintf("%d", newVal)
		}
		ha.SendSensor(EntityNameFromSensor(s), false, state)
	}
}

func HomeAssistantAddSensors(ha *HomeAssistant) {
	f := func(s p1p2.Sensor, value interface{}) {
		SensorToHomeAssistant(ha, s, value)
	}

	ha.checkAlive()
	if ha.Alive {
		for i := range p1p2.Sensors {
			// Sent all sensors to HA
			// Use a nil value to create the sensor but don't provide data
			SensorToHomeAssistant(ha, p1p2.Sensors[i], nil)
		}
	}
	// Register change event

	for i := range p1p2.Sensors {
		if p1p2.Sensors[i].Type() == "temperature" {
			if p1p2.Sensors[i].ID() == p1p2.ValveDomesticHotWater.ID() {
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
	go func(ha *HomeAssistant) {
		for {
			if ha.Alive {
				// Sent all sensors to HA
				for i := range p1p2.Sensors {
					if p1p2.Sensors[i].Value() != nil {
						SensorToHomeAssistant(ha, p1p2.Sensors[i], p1p2.Sensors[i].Value())
					}
				}
				time.Sleep(time.Hour * 12)
			} else {
				time.Sleep(time.Minute * 5)
			}
		}
	}(ha)
}
