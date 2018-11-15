package srecord

import "fmt"

type SrecLine struct {
	rtype    SrecType
	address  uint32
	data     []byte
	checksum byte
}

func (s SrecLine) byteCount() uint8 {
	switch s.rtype {
	case SrecHeader, SrecData16:
		return uint8(len(s.data) + 3)
	case SrecData24:
		return uint8(len(s.data) + 4)
	case SrecData32:
		return uint8(len(s.data) + 5)
	case SrecStart16, SrecCount16:
		return 3
	case SrecStart24, SrecCount24:
		return 4
	case SrecStart32:
		return 5
	}
	return 0
}

func (s SrecLine) Checksum() byte {
	if s.checksum != 0 {
		return s.checksum
	}

	sum := s.byteCount()

	sum += byte(s.address)
	sum += byte(s.address >> 8)

	switch s.rtype {
	case SrecData24, SrecStart24:
		sum += byte(s.address >> 16)
	case SrecData32, SrecStart32:
		sum += byte(s.address >> 16)
		sum += byte(s.address >> 24)
	}

	for _, d := range s.data {
		sum += d
	}

	s.checksum = ^sum
	return s.checksum
}

func (s SrecLine) print() {
	switch s.rtype {
	case SrecHeader:
		fmt.Printf("Header (%d): %s\n", len(s.data), s.data)
	case SrecData16, SrecData24, SrecData32:
		fmt.Printf("Data @ 0x%04X: %v\n", s.address, s.data)
	case SrecCount16, SrecCount24:
		fmt.Printf("Count: %v\n", s.data)
	case SrecStart16, SrecStart24, SrecStart32:
		fmt.Printf("Start Address: 0x%X\n", s.address)
	}
}

func (s SrecLine) String() string {
	return fmt.Sprintf(
		"S%s%s%s%X%X",
		string(srecTypesToByte[s.rtype]),
		formatHexBytes(uint64(s.byteCount()), 2),
		s.formatAddress(),
		s.data,
		s.Checksum(),
	)
}

func (s SrecLine) formatAddress() string {
	switch s.rtype {
	case SrecHeader:
		return "0000"
	case SrecData16, SrecData24, SrecData32, SrecStart16, SrecStart24, SrecStart32, SrecCount16, SrecCount24:
		return formatHexBytes(uint64(s.address), 4)
	}
	return ""
}
