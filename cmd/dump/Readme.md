# P1P2 raw packet hexdump decoder

Decodes a file containing P1P2 raw captures packets, one on each line.
The raw packets are expected to be in ASCII hexdump format:

[prefix:][0x] 00 [,][ ][0x] 01 [,][ ][0x] 02 ...

Examples:
 - `000001: 0x01, 0x02, 0x03, 0x04`
 - `0xff 0xff 0x00 0x00`
 - `de ed be ef`

## Building

    go build

## Running

   ./dump p1p2.hex

## Example output

```
Decoding packet '40001128102f8c04002a7c24b0296c0e80040100000000b5':
Temperature      ExternalSensor        : 4.003906
Temperature      ActualRoom            : 14.500000
Temperature      Refrigerant           : 41.421875
Temperature      GasBoiler             : 36.687500
```