// Copyright EPFL DEDIS

package f3b

import "os"

var protocol Protocol

func SelectedProtocol() Protocol {
	if protocol == nil {
		switch choice := os.Getenv("F3B_PROTOCOL"); choice {
		case "tpke":
			tpke, err := NewTPKE()
			if err != nil {
				panic(err)
			}
			protocol = tpke
		case "vdf":
			protocol = &VDF{}
		case "":
			protocol = nil
		default:
			panic("unknown F3B_PROTOCOL: " + choice)
		}
	}
	return protocol
}
