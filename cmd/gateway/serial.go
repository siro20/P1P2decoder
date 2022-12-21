package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"

	goterm "github.com/google/goterm/term"
	serial "go.bug.st/serial.v1"
)

type SerialConfig struct {
	Enable   bool   `yaml:"enable"`
	Device   string `yaml:"device"`
	BaudRate int    `yaml:"baud_rate"`
	Parity   string `yaml:"parity"`
	StopBits int    `yaml:"stop_bits"`
}

func UpdateConfigFromArg(cfg *SerialConfig) {
	if *baudrate > 0 {
		cfg.BaudRate = *baudrate
	}
	if *stopbits == 1 || *stopbits == 2 {
		cfg.StopBits = *stopbits
	}
	if *parity == "even" {
		cfg.Parity = *parity
	} else if *parity == "odd" {
		cfg.Parity = *parity
	}
	if *ttyDev != "" {
		cfg.Device = *ttyDev
	}
}

func VirtualGetSerialFromCfg(file string) (rc io.ReadCloser, err error) {
	pty, err := goterm.OpenPTY()
	if err != nil {
		return
	}

	// Close the slave, we don't need it
	pty.Slave.Close()

	ptyName, err := pty.PTSName()
	if err != nil {
		return
	}

	go func() {
		var scanner *bufio.Scanner
		file, err := os.Open(file)
		if err != nil {
			return
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
		for scanner.Scan() {
			str := scanner.Text()
			pty.Master.Write([]byte(str + "\n"))
			time.Sleep(time.Millisecond * 100)
		}
	}()
	return GetSerialFromCfg(SerialConfig{
		Enable:   true,
		Device:   ptyName,
		BaudRate: 115200,
	})
}

func GetSerialFromCfg(cfg SerialConfig) (rc io.ReadCloser, err error) {
	var con serial.Port

	mode := serial.Mode{
		BaudRate: cfg.BaudRate,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}
	if cfg.StopBits == 2 {
		mode.StopBits = serial.TwoStopBits
	}
	if cfg.Parity == "even" {
		mode.Parity = serial.EvenParity
	} else if cfg.Parity == "odd" {
		mode.Parity = serial.OddParity
	}

	con, err = serial.Open(cfg.Device, &mode)
	if err != nil {
		if *verbose {
			fmt.Printf("Failed to open TTY: %v\n", err)
		}
		time.Sleep(time.Second)
	}
	rc = con
	return
}
