package parser

type LabelReplace struct {
	Label  string
	Offset int16
}

type LabelMap map[string]uint16
type LabelLinkMap map[uint16]LabelReplace

type Program struct {
	Code    []uint8
	Labels  LabelMap
	LinkMap LabelLinkMap
	pc      uint16
}

func (p *Program) incPC() { p.pc++ }

func (p *Program) appendCode(b ...byte) {
	p.Code = append(p.Code, b...)
	p.pc += uint16(len(b))
}

func (p *Program) addLabel(name string) { p.Labels[name] = p.pc }

func (p *Program) addLink(pcoffset uint16, name string, offset int16) {
	p.LinkMap[p.pc+pcoffset] = LabelReplace{
		Label:  name,
		Offset: offset,
	}
}
