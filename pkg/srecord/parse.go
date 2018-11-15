package srecord

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"
)

func Parse(b []byte) (Srecord, error) {
	rec := make(Srecord, 0, 10)

	src := bufio.NewScanner(bytes.NewReader(b))
	linenum := 0
	addrlen := 0

	for src.Scan() {
		linenum++

		line := src.Bytes()

		// Remove newline characters
		line = bytes.TrimSpace(line)

		// Sanity checks
		if len(line) < 10 || len(line) > 514 || line[0] != 'S' {
			return nil, fmt.Errorf("invalid record on line %d", linenum)
		}

		// Convert ASCII pairs to literal bytes
		// Index 0 and 1 are special, they aren't hex
		line = append(line[0:2], convertHex(line[2:])...)
		if len(line) == 2 {
			return nil, fmt.Errorf("invalid record length on line %d: data length is not even", linenum)
		}

		// Check line length
		bcount := line[2]
		if len(line[3:]) != int(bcount) {
			return nil, fmt.Errorf("invalid record length on line %d: data length doesn't match byte count", linenum)
		}

		checksum := line[len(line)-1]
		data := line[3 : len(line)-1] // Everything between count and checksum

		recLine := SrecLine{
			checksum: checksum,
			rtype:    srecTypes[line[1]],
		}

		// Get line address
		address := uint32(0)
		switch recLine.rtype {
		case SrecHeader:
			data = data[2:]
		// Data/termination records
		case SrecData16, SrecStart16:
			if addrlen == 0 {
				addrlen = 16
			} else if addrlen != 16 {
				return nil, fmt.Errorf("invalid record on line %d: address length doesn't match file", linenum)
			}

			address = (uint32(data[0]) << 8) + uint32(data[1])
			data = data[2:]
		case SrecData24, SrecStart24:
			if addrlen == 0 {
				addrlen = 24
			} else if addrlen != 24 {
				return nil, fmt.Errorf("invalid record on line %d: address length doesn't match file", linenum)
			}

			address = (uint32(data[0]) << 16) + (uint32(data[1]) << 8) + uint32(data[2])
			data = data[3:]
		case SrecData32, SrecStart32:
			if addrlen == 0 {
				addrlen = 32
			} else if addrlen != 32 {
				return nil, fmt.Errorf("invalid record on line %d: address length doesn't match file", linenum)
			}

			address = (uint32(data[0]) << 24) + (uint32(data[1]) << 16) + (uint32(data[2]) << 8) + uint32(data[3])
			data = data[4:]
		case SrecCount16:
			address = (uint32(data[0]) << 8) + uint32(data[1])
			data = nil
		case SrecCount24:
			address = (uint32(data[0]) << 16) + (uint32(data[1]) << 8) + uint32(data[2])
			data = nil
		}

		// Fill in remaining struct
		recLine.data = copyBytes(data)
		recLine.address = address

		// Count and unknown checks
		switch recLine.rtype {
		case SrecUnknown:
			return nil, fmt.Errorf("invalid record type on line %d", linenum)
		case SrecCount16, SrecCount24:
			count := int(recLine.address)
			if len(rec)-1 != count {
				return nil, fmt.Errorf("invalid count on line %d", linenum)
			}
		}

		rec = append(rec, recLine)
	}

	return rec, nil
}

func copyBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

func hexToByte(b []byte) byte {
	r, _ := strconv.ParseUint(string(b), 16, 8)
	return byte(r)
}

func formatHexBytes(b uint64, p int) string {
	if p == 2 {
		return fmt.Sprintf("%02X", b)
	}
	return fmt.Sprintf("%04X", b)
}

func convertHex(b []byte) []byte {
	if len(b)%2 != 0 { // Length must be even
		return nil
	}

	ret := make([]byte, len(b)/2)
	hex.Decode(ret, b)
	return ret
}
