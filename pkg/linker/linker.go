package linker

import (
	"fmt"

	"github.com/lfkeitel/asml-sim/pkg/parser"
)

func Link(program *parser.Program) error {
	for _, part := range program.Parts {
		for loc, label := range part.LinkMap {
			memloc, exists := program.Labels[label.Label]
			if !exists {
				return fmt.Errorf("label %s not defined", label.Label)
			}

			newloc := memloc + uint16(label.Offset)

			part.Bytes[loc] = uint8(newloc >> 8)
			part.Bytes[loc+1] = uint8(newloc)
		}
	}

	return nil
}
