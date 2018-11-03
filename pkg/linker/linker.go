package linker

import (
	"fmt"

	"github.com/lfkeitel/asml-sim/pkg/parser"
)

func Link(program *parser.Program) error {
	for loc, label := range program.LinkMap {
		memloc, exists := program.Labels[label.Label]
		if !exists {
			return fmt.Errorf("label %s not defined", label.Label)
		}

		newloc := memloc + uint16(label.Offset)

		program.Code[loc] = uint8(newloc >> 8)
		program.Code[loc+1] = uint8(newloc)
	}

	return nil
}
