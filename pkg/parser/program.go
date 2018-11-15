package parser

import (
	"fmt"
	"sort"
)

type LabelReplace struct {
	Label  string
	Offset int16
}

type LabelMap map[string]uint16
type LabelLinkMap map[uint16]LabelReplace

type CodePart struct {
	Bytes   []uint8
	StartPC uint16
	PC      uint16
	LinkMap LabelLinkMap
}

func newCodePart(pc uint16) CodePart {
	return CodePart{
		Bytes:   make([]uint8, 0, 100),
		LinkMap: make(LabelLinkMap),
		StartPC: pc,
		PC:      pc,
	}
}

type Program struct {
	Parts     []CodePart
	partIndex int
	Labels    LabelMap
}

func (p *Program) incPC() { p.Parts[p.partIndex].PC++ }

func (p *Program) pc() uint16 { return p.Parts[p.partIndex].PC }

func (p *Program) appendCode(b ...byte) {
	p.Parts[p.partIndex].Bytes = append(p.Parts[p.partIndex].Bytes, b...)
	p.Parts[p.partIndex].PC += uint16(len(b))
}

func (p *Program) addLabel(name string) { p.Labels[name] = p.Parts[p.partIndex].PC }

func (p *Program) addLink(pcoffset uint16, name string, offset int16) {
	pc := p.Parts[p.partIndex].PC - p.Parts[p.partIndex].StartPC

	p.Parts[p.partIndex].LinkMap[pc+pcoffset] = LabelReplace{
		Label:  name,
		Offset: offset,
	}
}

func (p *Program) addCodePart(pc uint16) {
	p.Parts = append(p.Parts, newCodePart(pc))
	p.partIndex++
}

func (p *Program) validate() error {
	sort.Slice(p.Parts, func(i, j int) bool {
		return p.Parts[i].StartPC < p.Parts[j].StartPC
	})

	for i, code := range p.Parts {
		if i == len(p.Parts)-1 {
			break
		}

		if code.StartPC+uint16(len(code.Bytes)) > p.Parts[i+1].StartPC {
			return fmt.Errorf(`overlapping address regions:
Origin 0x%04X goes to 0x%04X
Origin 0x%04X begins inside region`, code.StartPC, code.StartPC+uint16(len(code.Bytes)), p.Parts[i+1].StartPC)
		}
	}
	return nil
}
