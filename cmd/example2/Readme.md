# P1P2 decoder - Example2

The `Decode` function returns the raw decoded message on success.
By casting it to one of the supported message struct you can access the raw
decoded register values.

```
	msg, err := p1p2.Decode( ... )
	if err != nil {
		fmt.Printf("Error decoding packet: %v\n", err)
	}
	if p, ok := msg.(p1p2.Packet10Req); ok {
		fmt.Printf("Received message was Packet 10 Request: %+v\n", p)

	} else if p, ok := msg.(p1p2.Packet10Resp); ok {
		fmt.Printf("Received message was Packet 10 Response: %+v\n", p)
	}
```

## Output:

```
Received message was Packet 10 Request: {Heating:1 OperationMode:129 DHWTank:1 Reserved:[0 0 0 0] TargetRoomTemperature:21 Reserved1:0 Flags:64 QuietMode:0 Reserved2:[0 8 0 0 24 0] DWHTankMode:64 DHWTankTargetTemperature:12288}
```