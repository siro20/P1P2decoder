# Daikin P1/P2 bus decoder

This repository contains tools to
 - capture packets from a serial device
 - verify integrity of the received packets
 - decode the P1/P2 packets
 - display the data using a HTTP server

It's planned to extend this library to MQTT and other home automation
protocols.
It's also planned to add gateway support.

## Applications

Applications are located in `cmd` folder.

## Library

The decoding library is located in `pkg` folder.
