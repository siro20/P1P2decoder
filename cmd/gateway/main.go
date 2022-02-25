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
	serial "go.bug.st/serial.v1"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	verbose  = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()
	ttyDev   = kingpin.Flag("chardev", "Path to tty device to read data from").Short('c').String()
	dumpFile = kingpin.Flag("dump", "Path to regular file to read data from").Short('d').String()
	baudrate = kingpin.Flag("baud", "Serial baud rate").Short('b').Int()
	parity   = kingpin.Flag("parity", "Serial parity: none/even/odd").Default("none").Short('p').String()
	stopbits = kingpin.Flag("stop", "Serial stop bits: 1/2").Short('s').Int()
)

func main() {
	kingpin.Parse()

	mode := serial.Mode{
		BaudRate: 115200,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}
	if *baudrate > 0 {
		mode.BaudRate = *baudrate
	}
	if *stopbits == 2 {
		mode.StopBits = serial.TwoStopBits
	}
	if *parity == "even" {
		mode.Parity = serial.EvenParity
	} else if *parity == "odd" {
		mode.Parity = serial.OddParity
	}

	if (ttyDev == nil || *ttyDev == "") && (dumpFile == nil || *dumpFile == "") {
		log.Print("No input specified")
		os.Exit(1)
	}

	go runHtml(p1p2.Sys)

	for {
		var scanner *bufio.Scanner
		var closer io.Closer
		// Poll on Serial to open (Testing)
		if ttyDev != nil && *ttyDev != "" {
			con, err := serial.Open(*ttyDev, &mode)
			if err != nil {
				if *verbose {
					fmt.Printf("Failed to open TTY: %v\n", err)
				}
				time.Sleep(time.Second)
				continue
			}
			scanner = bufio.NewScanner(con)
			closer = con
		} else if dumpFile != nil && *dumpFile != "" {

			file, err := os.Open(*dumpFile)
			if err != nil {
				log.Fatal(err)
			}
			scanner = bufio.NewScanner(file)
			closer = file
		}

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
		if err := scanner.Err(); err != nil {
			if err != syscall.EINTR {
				fmt.Print(err)
			}
		}
		closer.Close()
		if dumpFile != nil && *dumpFile != "" {
			time.Sleep(time.Second * 99999) // Keep running
		}
	}

}
