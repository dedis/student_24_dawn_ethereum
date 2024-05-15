package controller

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"go.dedis.ch/dela/cli/node"
	"go.dedis.ch/dela/crypto"
	"go.dedis.ch/dela/testing/fake"
	"go.dedis.ch/dela/mino"
	"go.dedis.ch/f3b/smc/dkg"
	"go.dedis.ch/kyber/v3"
)

func TestSetupAction_NoActor(t *testing.T) {
	a := setupAction{}

	inj := node.NewInjector()

	ctx := node.Context{
		Injector: inj,
	}

	err := a.Execute(ctx)
	require.EqualError(t, err, "failed to resolve actor, did you call listen?: "+
		"couldn't find dependency for 'dkg.Actor'")
}

func TestSetupAction_NoCollectiveAuth(t *testing.T) {
	a := setupAction{}

	inj := node.NewInjector()
	inj.Inject(&fakeActor{})

	flags := node.FlagSet{
		"authority": []interface{}{"fake"},
	}

	ctx := node.Context{
		Injector: inj,
		Flags:    flags,
	}

	err := a.Execute(ctx)
	require.EqualError(t, err, "failed to get collective authority: "+
		"failed to decode authority: invalid identity base64 string")
}

func TestSetupAction_badSetup(t *testing.T) {
	a := setupAction{}

	inj := node.NewInjector()
	inj.Inject(&fakeActor{
		setupErr: fake.GetError(),
	})

	flags := node.FlagSet{}

	ctx := node.Context{
		Injector: inj,
		Flags:    flags,
	}

	err := a.Execute(ctx)
	require.EqualError(t, err, fake.Err("failed to setup"))
}

func TestSetupAction_OK(t *testing.T) {
	a := setupAction{}

	inj := node.NewInjector()
	inj.Inject(&fakeActor{})

	flags := node.FlagSet{}

	out := &bytes.Buffer{}

	ctx := node.Context{
		Injector: inj,
		Flags:    flags,
		Out:      out,
	}

	err := a.Execute(ctx)
	require.NoError(t, err)

	require.Regexp(t, "^✅ Setup done", out.String())
}

func TestListenAction_NoDKG(t *testing.T) {
	a := listenAction{}

	inj := node.NewInjector()

	ctx := node.Context{
		Injector: inj,
	}

	err := a.Execute(ctx)
	require.EqualError(t, err, "failed to resolve dkg: couldn't find dependency for 'dkg.DKG'")
}

func TestListenAction_listenFail(t *testing.T) {
	a := listenAction{}

	inj := node.NewInjector()
	inj.Inject(fakeDKG{
		err: fake.GetError(),
	})

	ctx := node.Context{
		Injector: inj,
	}

	err := a.Execute(ctx)
	require.EqualError(t, err, fake.Err("failed to listen"))
}

func TestListenAction_encodeFail(t *testing.T) {
	a := listenAction{}

	inj := node.NewInjector()
	inj.Inject(fakeDKG{
		actor: fakeActor{},
	})

	ctx := node.Context{
		Injector: inj,
		Out:      io.Discard,
	}

	err := a.Execute(ctx)
	require.EqualError(t, err, "failed to encode authority: failed to resolve"+
		" mino: couldn't find dependency for 'mino.Mino'")
}

func TestListenAction_writeFail(t *testing.T) {
	a := listenAction{
		pubkey: suite.Point(),
	}

	inj := node.NewInjector()
	inj.Inject(fakeDKG{
		actor: fakeActor{},
	})
	inj.Inject(fake.Mino{})

	flags := node.FlagSet{
		"config": "/fake/fake", // wrong configuration path
	}

	ctx := node.Context{
		Injector: inj,
		Out:      io.Discard,
		Flags:    flags,
	}

	err := a.Execute(ctx)
	require.Regexp(t, "^failed to write authority configuration", err.Error())
}

func TestListenAction_OK(t *testing.T) {
	tmpDir := t.TempDir()

	a := listenAction{
		pubkey: suite.Point(),
	}

	inj := node.NewInjector()
	inj.Inject(fakeDKG{
		actor: fakeActor{},
	})
	inj.Inject(fake.Mino{})

	flags := node.FlagSet{
		"config": tmpDir,
	}

	out := &bytes.Buffer{}

	ctx := node.Context{
		Injector: inj,
		Out:      out,
		Flags:    flags,
	}

	err := a.Execute(ctx)
	require.NoError(t, err)

	require.Regexp(t, "^✅  Listen done, actor is created.📜 Config file written in", out.String())
}

func TestEncodeAuthority_marshalFail(t *testing.T) {
	inj := node.NewInjector()
	inj.Inject(fakeDKG{
		actor: fakeActor{},
	})
	inj.Inject(fakeMino{
		addr: fake.NewBadAddress(),
	})

	ctx := node.Context{
		Injector: inj,
	}

	_, err := encodeAuthority(ctx, badPoint{})
	require.EqualError(t, err, fake.Err("failed to marshal address"))
}

func TestEncodeAuthority_pkFail(t *testing.T) {
	inj := node.NewInjector()
	inj.Inject(fakeDKG{
		actor: fakeActor{},
	})
	inj.Inject(fakeMino{
		addr: fake.NewAddress(0),
	})

	ctx := node.Context{
		Injector: inj,
	}

	pk := badPoint{
		err: fake.GetError(),
	}

	_, err := encodeAuthority(ctx, pk)
	require.EqualError(t, err, fake.Err("failed to marshall pubkey"))
}

func TestDecodeAuthority_noMino(t *testing.T) {
	pubKey := "RjEyNy4wLjAuMToyMDAx:RcTqbSXCIkZmmaGjLVAZs8TdTvq7b3SFr14F89h6ID8="

	inj := node.NewInjector()

	ctx := node.Context{
		Injector: inj,
	}

	_, _, err := decodeAuthority(ctx, pubKey)
	require.EqualError(t, err, "injector: couldn't find dependency for 'mino.Mino'")
}

func TestDecodeAuthority_noBase64Addr(t *testing.T) {
	pubKey := "aa:RcTqbSXCIkZmmaGjLVAZs8TdTvq7b3SFr14F89h6ID8="

	inj := node.NewInjector()
	inj.Inject(fake.Mino{})

	ctx := node.Context{
		Injector: inj,
	}

	_, _, err := decodeAuthority(ctx, pubKey)
	require.EqualError(t, err, "base64 address: illegal base64 data at input byte 0")
}

func TestDecodeAuthority_noBase64Pubkey(t *testing.T) {
	pubKey := "RjEyNy4wLjAuMToyMDAx:aa"

	inj := node.NewInjector()
	inj.Inject(fake.Mino{})

	ctx := node.Context{
		Injector: inj,
	}

	_, _, err := decodeAuthority(ctx, pubKey)
	require.EqualError(t, err, "base64 public key: illegal base64 data at input byte 0")
}

func TestDecodeAuthority_badUnmarshalPubkey(t *testing.T) {
	pubKey := "RjEyNy4wLjAuMToyMDAx:aaa="

	inj := node.NewInjector()
	inj.Inject(fake.Mino{})

	ctx := node.Context{
		Injector: inj,
	}

	_, _, err := decodeAuthority(ctx, pubKey)
	require.EqualError(t, err, "failed to decode pubkey: invalid Ed25519 curve point")
}

func TestReshareAction_noActor(t *testing.T) {
	a := reshareAction{}

	inj := node.NewInjector()

	ctx := node.Context{
		Injector: inj,
	}

	err := a.Execute(ctx)
	require.EqualError(t, err, "failed to resolve actor, did you call listen?: couldn't find dependency for 'dkg.Actor'")
}

func TestReshareAction_NoCollectiveAuth(t *testing.T) {
	a := reshareAction{}

	inj := node.NewInjector()
	inj.Inject(&fakeActor{})

	flags := node.FlagSet{
		"authority": []interface{}{"fake"},
	}

	ctx := node.Context{
		Injector: inj,
		Flags:    flags,
	}

	err := a.Execute(ctx)
	require.EqualError(t, err, "failed to get collective authority: "+
		"failed to decode authority: invalid identity base64 string")
}

func TestReshareAction_reshareFail(t *testing.T) {
	a := reshareAction{}

	inj := node.NewInjector()
	inj.Inject(&fakeActor{
		reshareErr: fake.GetError(),
	})

	flags := node.FlagSet{}

	ctx := node.Context{
		Injector: inj,
		Flags:    flags,
	}

	err := a.Execute(ctx)
	require.EqualError(t, err, fake.Err("failed to reshare"))
}

func TestReshareAction_OK(t *testing.T) {
	a := reshareAction{}

	inj := node.NewInjector()
	inj.Inject(&fakeActor{})

	flags := node.FlagSet{}

	out := &bytes.Buffer{}

	ctx := node.Context{
		Injector: inj,
		Flags:    flags,
		Out:      out,
	}

	err := a.Execute(ctx)
	require.NoError(t, err)

	require.Equal(t, "✅ Reshare done.\n", out.String())
}

// -----------------------------------------------------------------------------
// Utility functions

type fakeActor struct {
	dkg.Actor

	setupErr    error
	encryptErr  error
	decryptErr  error
	vencryptErr error
	vdecryptErr error
	reshareErr  error

	k kyber.Point
	c kyber.Point

	decryptData  []byte
	vdecryptData [][]byte
}

func (f fakeActor) Setup(co crypto.CollectiveAuthority, threshold int) (pubKey kyber.Point, err error) {
	return suite.Point(), f.setupErr
}

func (f fakeActor) Encrypt(message []byte) (K, C kyber.Point, remainder []byte, err error) {
	return f.k, f.c, nil, f.encryptErr
}

func (f fakeActor) Decrypt(K, C kyber.Point) ([]byte, error) {
	return f.decryptData, f.decryptErr
}

func (f fakeActor) Reshare(co crypto.CollectiveAuthority, newThreshold int) error {
	return f.reshareErr
}

type fakeDKG struct {
	dkg.DKG

	err   error
	actor dkg.Actor
}

func (f fakeDKG) Listen() (dkg.Actor, error) {
	return f.actor, f.err
}

type badPoint struct {
	kyber.Point

	err  error
	data string
}

func (b badPoint) MarshalBinary() (data []byte, err error) {
	return []byte(b.data), b.err
}

type badScalar struct {
	kyber.Scalar

	err  error
	data string
}

func (b badScalar) MarshalBinary() (data []byte, err error) {
	return []byte(b.data), b.err
}

type fakeMino struct {
	mino.Mino

	addr mino.Address
}

func (f fakeMino) GetAddress() mino.Address {
	return f.addr
}
