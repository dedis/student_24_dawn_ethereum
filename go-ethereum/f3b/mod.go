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
			protocol, err = NewTIBE(NewSmcCli(p))
			if err != nil {
				panic(err)
			}
		case "tpke":
			protocol, err = NewTPKE(NewSmcCli(p))
			if err != nil {
				panic(err)
			}
		case "vdf":
			protocol = &VDF{defaultLog2t}
		case "null":
			protocol = NewNull()
		case "faketpke":
			protocol, err = NewFakeTPKE()
			if err != nil {
				panic(err)
			}
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
func ForceSelectedProtocol(_ testing.TB, p Protocol) {
	protocol = p
}
