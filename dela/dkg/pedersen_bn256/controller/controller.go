package controller

import (
	"go.dedis.ch/dela"
	"go.dedis.ch/dela/cli"
	"go.dedis.ch/dela/cli/node"
	"go.dedis.ch/dela/dkg/pedersen_bn256"
	"go.dedis.ch/dela/mino"
	"golang.org/x/xerrors"
)

// NewMinimal returns a new minimal initializer
func NewMinimal() node.Initializer {
	return minimal{
		la: &listenAction{},
	}
}

// minimal is an initializer with the minimum set of commands. Indeed it only
// creates and injects a new DKG
//
// - implements node.Initializer
type minimal struct {
	la *listenAction
}

// Build implements node.Initializer. In this case we don't need any command.
func (m minimal) SetCommands(builder node.Builder) {
	cmd := builder.SetCommand("dkg")
	cmd.SetDescription("DKG service administration")

	sub := cmd.SetSubCommand("listen")
	sub.SetDescription("initialize DKG, create the actor and save the authority configuration")
	sub.SetAction(builder.MakeAction(m.la))

	sub = cmd.SetSubCommand("setup")
	sub.SetDescription("setup the DKG service")
	sub.SetFlags(
		cli.StringSliceFlag{
			Name:  "authority",
			Usage: "<ADDR>:<PK> string, where each token is encoded in base64",
		},
		cli.IntFlag{
			Name:  "threshold",
			Usage: "the threshold of the committee",
			Required: true,
		},
	)
	sub.SetAction(builder.MakeAction(setupAction{}))

	sub = cmd.SetSubCommand("get-public-key")
	sub.SetDescription("Query the collective public key. Outputs in hex")
	sub.SetAction(builder.MakeAction(getPublicKeyAction{}))

	sub = cmd.SetSubCommand("sign")
	sub.SetDescription("sign a message. Outputs signature in hex")
	sub.SetFlags(
		cli.StringFlag{
			Name:  "message",
			Usage: "the message to sign, encoded in hex",
		},
	)
	sub.SetAction(builder.MakeAction(signAction{}))

	sub = cmd.SetSubCommand("verify")
	sub.SetDescription("verify a threshold signature")
	sub.SetFlags(
		cli.StringFlag{
			Name:  "message",
			Usage: "the message that was sign, encoded in hex",
		},
		cli.StringFlag{
			Name:  "signature",
			Usage: "the signature, encoded in hex",
		},
	)
	sub.SetAction(builder.MakeAction(verifyAction{}))

	sub = cmd.SetSubCommand("encrypt")
	sub.SetDescription("encrypt a message. Outputs ciphertext in hex")
	sub.SetFlags(
		cli.StringFlag{
			Name:  "label",
			Usage: "the IBE label to encrypt to, encoded in hex",
		},
		cli.StringFlag{
			Name:  "message",
			Usage: "the message to encrypt, encoded in hex",
		},
	)
	sub.SetAction(builder.MakeAction(encryptAction{}))

	sub = cmd.SetSubCommand("decrypt")
	sub.SetDescription("decrypt a ciphertext. Outputs message in hex")
	sub.SetFlags(
		cli.StringFlag{
			Name:  "label",
			Usage: "the IBE label to encrypt to, encoded in hex",
		},
		cli.StringFlag{
			Name:  "ciphertext",
			Usage: "the ciphertext to decrypt, encoded in hex",
		},
	)
	sub.SetAction(builder.MakeAction(decryptAction{}))

	sub = cmd.SetSubCommand("reshare")
	sub.SetDescription("reshare the DKG secret")
	sub.SetFlags(
		cli.StringSliceFlag{
			Name:  "authority",
			Usage: "<ADDR>:<PK> string, where each token is encoded in base64",
		},
		cli.IntFlag{
			Name:     "thresholdNew",
			Usage:    "the threshold of the new committee",
			Required: true,
		},
	)
	sub.SetAction(builder.MakeAction(reshareAction{}))
}

// OnStart implements node.Initializer. It creates and registers a pedersen DKG.
func (m minimal) OnStart(ctx cli.Flags, inj node.Injector) error {
	var no mino.Mino
	err := inj.Resolve(&no)
	if err != nil {
		return xerrors.Errorf("failed to resolve mino: %v", err)
	}

	dkg, pubkey := pedersen.NewPedersen(no)

	inj.Inject(dkg)

	pubkeyBuf, err := pubkey.MarshalBinary()
	if err != nil {
		return xerrors.Errorf("failed to encode pubkey: %v", err)
	}

	dela.Logger.Info().
		Hex("public key", pubkeyBuf).
		Msg("perdersen public key")

	// the listen action is expecting the pubkey to be set
	m.la.pubkey = pubkey

	return nil
}

// OnStop implements node.Initializer.
func (minimal) OnStop(node.Injector) error {
	return nil
}
