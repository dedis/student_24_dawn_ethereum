// Copyright EPFL DEDIS

package f3b

import "testing"

var protocol Protocol

func SelectedProtocol() Protocol {
	if protocol == nil {
		p, err := ReadParams()
		if err != nil {
			panic(err)
		}

		switch p.Protocol {
		case "tibe":
			tibe, err := NewTIBE()
			if err != nil {
				panic(err)
			}
			protocol = tibe
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
			panic("unknown F3B protocol: " + p.Protocol)
		}
	}
	return protocol
}

// ForceSelectedProtocol forcefully sets the globally selected protocol.
// It is only meant to be used in tests.
func ForceSelectedProtocol(_ *testing.T, p Protocol) {
	protocol = p
}
