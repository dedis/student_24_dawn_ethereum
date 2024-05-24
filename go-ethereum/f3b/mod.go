// Copyright EPFL DEDIS

package f3b

var protocol Protocol

func SelectedProtocol() Protocol {
	if protocol == nil {
		p, err := ReadParams()
		if err != nil {
			panic(err)
		}

		switch p.Protocol {
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
