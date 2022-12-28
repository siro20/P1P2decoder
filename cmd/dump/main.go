package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/siro20/p1p2decoder/pkg/p1p2"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	verbose   = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()
	dumpFile  = kingpin.Arg("dump", "Path to regular file to read data from").Required().String()
	newPacket = false
	packet    = ""
)

func main() {
	fmt.Printf("Decoding dump file with the following syntax:\n")
	fmt.Printf("[timestamp :][state :][0x] 00 [,][ ][0x] 01 [,][ ][0x] 02 ...\n")
	fmt.Printf("Examples:\n")
	fmt.Printf(" 000: 0x01, 0x02, 0x03\n")
	fmt.Printf(" 0xFF, 0xDD, 0xEE\n")
	fmt.Printf(" de ad be ef\n")

	kingpin.Parse()

	file, err := os.Open(*dumpFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for i := range p1p2.Sensors {
		p1p2.Sensors[i].RegisterStateChangedCallback(func(s p1p2.Sensor, value interface{}) {
			fmt.Printf("Sensor %s %s changed to %v\n", s.Type(), s.Name(), value)
		})
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var s string
		if strings.Contains(scanner.Text(), ":") {
			split := strings.Split(scanner.Text(), ":")
			if len(split) == 3 {
				s = split[2]
			} else {
				s = split[1]
			}
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

		packet = s
		newPacket = true
		_, err = p1p2.Decode(buf)
		if err != nil {
			if *verbose {
				fmt.Printf("Error decoding packet '%s': %v\n", s, err)
			}
			continue
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
