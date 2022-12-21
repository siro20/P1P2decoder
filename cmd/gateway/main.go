package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/siro20/p1p2decoder/pkg/p1p2"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	verbose          = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()
	ttyDev           = kingpin.Flag("chardev", "Path to tty device to read data from").Short('c').String()
	dumpFile         = kingpin.Flag("dump", "Path to regular file to read data from").Short('d').String()
	baudrate         = kingpin.Flag("baud", "Serial baud rate").Short('b').Int()
	parity           = kingpin.Flag("parity", "Serial parity: none/even/odd").Default("none").Short('p').String()
	stopbits         = kingpin.Flag("stop", "Serial stop bits: 1/2").Short('s').Int()
	database         = kingpin.Flag("db", "Path to database for non-volatile storage").Short('e').String()
	htmlServer       = kingpin.Flag("html", "Run HTTP server").Bool()
	htmlServerAssets = kingpin.Flag("html-assets-path", "The path to HTML assets.").String()
	prometheusServer = kingpin.Flag("prometheus", "Run prometheus server").Bool()
	prometheusAddr   = kingpin.Flag("prometheus-listen-address", "The address to listen on for prometheus requests.").String()
)

var db *p1p2.DB

func main() {
	var cfg *Config
	var err error
	var path string
	kingpin.Parse()

	path, err = os.Getwd()
	if err == nil {
		cfg, err = ReadConfig(path + "/p1p2.yaml")
	}
	if err != nil {
		cfg, err = ReadConfig("/etc/p1p2_gateway/p1p2.yaml")
	}
	if err != nil {
		dirname, err := os.UserHomeDir()
		if err == nil {
			cfg, err = ReadConfig(dirname + "/.config/p1p2_gateway.yaml")
		}
	}
	if err != nil {
		fmt.Printf("WARN: Could not find a config file\n")
		cfg = &Config{
			Prometheus: PrometheusConfig{
				Enable: true,
			},
		}
	}
	UpdateConfigFromArg(&cfg.Serial)

	if *database != "" {
		db, err = p1p2.OpenDB(*database)
		if err != nil {
			fmt.Printf("Error opening database: %v\n", err)
			os.Exit(1)
			return
		}
	}

	if cfg.Serial.Device == "" && (dumpFile == nil || *dumpFile == "") {
		log.Print("No input specified")
		os.Exit(1)
	}

	if cfg.Html.Enable {
		go runHtml(cfg.Html)
	}

	if cfg.Prometheus.Enable {
		go runPrometheusServer(cfg.Prometheus)
	}

	if cfg.HomeAssistant.Enable {
		ha, _ := NewHomeAssistant(cfg.HomeAssistant)
		HomeAssistantAddSensors(ha)
		go ha.Serve()
	}

	for {
		var scanner *bufio.Scanner
		var rc io.ReadCloser
		// Poll on Serial to open (Testing)
		if cfg.Serial.Device != "" {
			rc, err = GetSerialFromCfg(cfg.Serial)
		} else if dumpFile != nil && *dumpFile != "" {
			rc, err = VirtualGetSerialFromCfg(*dumpFile)
		}
		scanner = bufio.NewScanner(rc)
		for scanner.Err() == nil {
			for scanner.Scan() {
				var s string
				if strings.Contains(scanner.Text(), ":") {
					s = strings.Split(scanner.Text(), ":")[1]
				}

				// Remove whitespace
				s = strings.ReplaceAll(s, ", 0x", "")
				s = strings.ReplaceAll(s, " 0x", "")
				s = strings.ReplaceAll(s, " ", "")

				buf, err := hex.DecodeString(s)
				if err != nil {
					if *verbose {
						fmt.Printf("Skipping invalid line in file: %s", scanner.Text())
					}
					continue
				}

				_, err = p1p2.Decode(buf)
				if err != nil {
					if *verbose {
						fmt.Printf("Error decoding packet '%s': %v\n", s, err)
					}
					continue
				}
			}
			if scanner.Err() != nil {
				if scanner.Err() != syscall.EINTR {
					fmt.Printf("Reading from serial failed with: %v\n", scanner.Err())
				} else {
					// Clear scanner error by creating new scanner...
					scanner = bufio.NewScanner(rc)
				}
			}
		}

		rc.Close()
		if dumpFile != nil && *dumpFile != "" {
			fmt.Printf("Now waiting...")
			time.Sleep(time.Second * 99999) // Keep running
		}
	}

}
