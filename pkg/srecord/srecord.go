package srecord

import (
	"strings"
)

type SrecType int

const (
	SrecUnknown SrecType = iota
	SrecHeader
	SrecData16
	SrecData24
	SrecData32
	SrecReserved
	SrecCount16
	SrecCount24
	SrecStart32
	SrecStart24
	SrecStart16
)

var srecTypes = map[byte]SrecType{
	'0': SrecHeader,
	'1': SrecData16,
	'2': SrecData24,
	'3': SrecData32,
	'4': SrecReserved,
	'5': SrecCount16,
	'6': SrecCount24,
	'7': SrecStart32,
	'8': SrecStart24,
	'9': SrecStart16,
}

var srecTypesToByte = map[SrecType]byte{
	SrecHeader:   '0',
	SrecData16:   '1',
	SrecData24:   '2',
	SrecData32:   '3',
	SrecReserved: '4',
	SrecCount16:  '5',
	SrecCount24:  '6',
	SrecStart32:  '7',
	SrecStart24:  '8',
	SrecStart16:  '9',
}

type Srecord []SrecLine

func New() Srecord {
	return make(Srecord, 0, 10)
}

func (s *Srecord) AddHeader(data string) {
	*s = append(*s, SrecLine{
		rtype: SrecHeader,
		data:  []byte(data),
	})
}

func (s *Srecord) AddRecord16(t SrecType, address uint16, data []byte) {
	*s = append(*s, SrecLine{
		rtype:   t,
		address: uint32(address),
		data:    data,
	})
}

func (s *Srecord) AddRecord32(t SrecType, address uint32, data []byte) {
	*s = append(*s, SrecLine{
		rtype:   t,
		address: address,
		data:    data,
	})
}

func (s Srecord) String() string {
	var b strings.Builder
	for _, l := range s {
		b.WriteString(l.String())
		b.WriteByte('\n')
	}
	return b.String()
}
