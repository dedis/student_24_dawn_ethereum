package pedersen

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"go.dedis.ch/dela"
	"go.dedis.ch/dela/crypto"
	"go.dedis.ch/dela/crypto/bls"
	"go.dedis.ch/dela/dkg"
	"go.dedis.ch/dela/dkg/pedersen_bn256/types"
	"go.dedis.ch/dela/internal/testing/fake"
	"go.dedis.ch/dela/mino"
	"go.dedis.ch/dela/mino/minogrpc"
	"go.dedis.ch/dela/mino/router/tree"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/pairing/bn256"
	"go.dedis.ch/kyber/v3/share"
	"go.dedis.ch/kyber/v3/sign/tbls"
)

func TestPedersen_Listen(t *testing.T) {
	pedersen, _ := NewPedersen(fake.Mino{})

	actor, err := pedersen.Listen()
	require.NoError(t, err)

	require.NotNil(t, actor)
}

func TestPedersen_Setup(t *testing.T) {
	actor := Actor{
		rpc:      fake.NewBadRPC(),
		startRes: &state{},
	}

	fakeAuthority := fake.NewAuthority(1, fake.NewSigner)

	_, err := actor.Setup(fakeAuthority, 0)
	require.EqualError(t, err, fake.Err("failed to stream"))

	rpc := fake.NewStreamRPC(fake.NewReceiver(), fake.NewBadSender())
	actor.rpc = rpc

	_, err = actor.Setup(fakeAuthority, 0)
	require.EqualError(t, err, "expected bls.PublicKey, got 'fake.PublicKey'")

	fakeAuthority = fake.NewAuthority(2, bls.Generate)

	_, err = actor.Setup(fakeAuthority, 1)
	require.EqualError(t, err, fake.Err("failed to send start"))

	rpc = fake.NewStreamRPC(fake.NewBadReceiver(), fake.Sender{})
	actor.rpc = rpc

	_, err = actor.Setup(fakeAuthority, 1)
	require.EqualError(t, err, fake.Err("got an error from '%!s(<nil>)' while receiving"))

	recv := fake.NewReceiver(fake.NewRecvMsg(fake.NewAddress(0), nil))

	actor.rpc = fake.NewStreamRPC(recv, fake.Sender{})

	_, err = actor.Setup(fakeAuthority, 1)
	require.EqualError(t, err, "expected to receive a Done message, but go the following: <nil>")

	rpc = fake.NewStreamRPC(fake.NewReceiver(
		fake.NewRecvMsg(fake.NewAddress(0), types.NewStartDone(suite.Point())),
		fake.NewRecvMsg(fake.NewAddress(0), types.NewStartDone(suite.Point().Pick(suite.RandomStream()))),
	), fake.Sender{})
	actor.rpc = rpc

	_, err = actor.Setup(fakeAuthority, 1)
	require.Error(t, err)
	require.Regexp(t, "^the public keys does not match:", err)
}

func TestPedersen_GetPublicKey(t *testing.T) {
	actor := Actor{
		startRes: &state{},
	}

	_, err := actor.GetPublicKey()
	require.EqualError(t, err, "DKG has not been initialized")

	actor.startRes = &state{dkgState: certified}
	_, err = actor.GetPublicKey()
	require.NoError(t, err)
}

func TestPedersen_Sign(t *testing.T) {
	priShares := []*share.PriShare{
		&share.PriShare{0, suite.Scalar().Pick(suite.RandomStream())},
		&share.PriShare{1, suite.Scalar().Pick(suite.RandomStream())},
	}
	priPoly, err := share.RecoverPriPoly(bn256.NewSuite().G2(), priShares, 2, 2)
	require.NoError(t, err)
	pubPoly := priPoly.Commit(nil)
	_, commits := pubPoly.Info()

	actor := Actor{
		rpc: fake.NewBadRPC(),
		startRes: &state{dkgState: certified,
			participants: []mino.Address{fake.NewAddress(0), fake.NewAddress(1)}, Commits: commits, threshold: 2},
	}

	msg := []byte("merry christmas")
	var tsigs [][]byte
	for i := range priShares {
		tsig, err := tbls.Sign(pairingSuite, priShares[i], msg)
		require.NoError(t, err)
		tsigs = append(tsigs, tsig)
	}

	_, err = tbls.Recover(pairingSuite, pubPoly, msg, tsigs, 2, 2)
	require.NoError(t, err)

	recv := fake.NewReceiver(
		fake.NewRecvMsg(fake.NewAddress(0), types.NewSignReply(tsigs[0])),
		fake.NewRecvMsg(fake.NewAddress(1), types.NewSignReply(tsigs[1])),
	)

	rpc := fake.NewStreamRPC(recv, fake.Sender{})
	actor.rpc = rpc

	sig, err := actor.Sign(msg)
	require.NoError(t, err)

	// Expect a valid signature
	err = bls.NewPublicKeyFromPoint(pubPoly.Commit()).Verify(msg, bls.NewSignature(sig))
	require.NoError(t, err)
}

func TestPedersen_Scenario(t *testing.T) {
	// Use with MINO_TRAFFIC=log
	// traffic.LogItems = false
	// traffic.LogEvent = false
	// defer func() {
	// 	traffic.SaveItems("graph.dot", true, false)
	// 	traffic.SaveEvents("events.dot")
	// }()

	oldLog := dela.Logger
	defer func() {
		dela.Logger = oldLog
	}()

	dela.Logger = dela.Logger.Level(zerolog.WarnLevel)

	n := 32

	minos := make([]mino.Mino, n)
	dkgs := make([]dkg.DKG, n)
	addrs := make([]mino.Address, n)

	for i := 0; i < n; i++ {
		addr := minogrpc.ParseAddress("127.0.0.1", 0)

		m, err := minogrpc.NewMinogrpc(addr, nil, tree.NewRouter(minogrpc.NewAddressFactory()))
		require.NoError(t, err)

		defer m.GracefulStop()

		minos[i] = m
		addrs[i] = m.GetAddress()
	}

	pubkeys := make([]kyber.Point, len(minos))

	for i, mi := range minos {
		for _, m := range minos {
			mi.(*minogrpc.Minogrpc).GetCertificateStore().Store(m.GetAddress(), m.(*minogrpc.Minogrpc).GetCertificateChain())
		}

		d, pubkey := NewPedersen(mi.(*minogrpc.Minogrpc))

		dkgs[i] = d
		pubkeys[i] = pubkey
	}

	fakeAuthority := NewAuthority(addrs, pubkeys)

	message := []byte("Hello world")
	actors := make([]dkg.Actor, n)
	for i := 0; i < n; i++ {
		actor, err := dkgs[i].Listen()
		require.NoError(t, err)

		actors[i] = actor
	}

	// trying to call sign before a setup
	_, err := actors[0].Sign(message)
	require.EqualError(t, err, "you must first initialize DKG. Did you call setup() first?")

	_, err = actors[0].Setup(fakeAuthority, n)
	require.NoError(t, err)

	_, err = actors[0].Setup(fakeAuthority, n)
	require.EqualError(t, err, "startRes is already done, only one setup call is allowed")

	// every node should be able to sign
	for i := 0; i < n; i++ {
		_, err := actors[i].Sign(message)
		require.NoError(t, err)
	}
}

func Test_Reshare_NotDone(t *testing.T) {
	a := Actor{
		startRes: &state{dkgState: initial},
	}

	err := a.Reshare(nil, 0)
	require.EqualError(t, err, "you must first initialize DKG. Did you call setup() first?")
}

func Test_Reshare_WrongPK(t *testing.T) {
	a := Actor{
		startRes: &state{dkgState: certified},
	}

	co := fake.NewAuthority(1, fake.NewSigner)

	err := a.Reshare(co, 0)
	require.EqualError(t, err, "expected bls.PublicKey, got 'fake.PublicKey'")
}

func Test_Reshare_BadRPC(t *testing.T) {
	a := Actor{
		startRes: &state{dkgState: certified},
		rpc:      fake.NewBadRPC(),
	}

	co := NewAuthority(nil, nil)

	err := a.Reshare(co, 0)
	require.EqualError(t, err, fake.Err("failed to create stream"))
}

func Test_Reshare_BadSender(t *testing.T) {
	a := Actor{
		startRes: &state{dkgState: certified},
		rpc:      fake.NewStreamRPC(nil, fake.NewBadSender()),
	}

	co := NewAuthority(nil, nil)

	err := a.Reshare(co, 0)
	require.EqualError(t, err, fake.Err("failed to send resharing request"))
}

func Test_Reshare_BadReceiver(t *testing.T) {
	a := Actor{
		startRes: &state{dkgState: certified},
		rpc:      fake.NewStreamRPC(fake.NewBadReceiver(), fake.Sender{}),
	}

	co := NewAuthority([]mino.Address{fake.NewAddress(0)}, []kyber.Point{suite.Point()})

	err := a.Reshare(co, 0)
	require.EqualError(t, err, fake.Err("stream stopped unexpectedly"))
}

// -----------------------------------------------------------------------------
// Utility functions

//
// Collective authority
//

// CollectiveAuthority is a fake implementation of the cosi.CollectiveAuthority
// interface.
type CollectiveAuthority struct {
	crypto.CollectiveAuthority
	addrs   []mino.Address
	pubkeys []kyber.Point
	signers []crypto.Signer
}

// NewAuthority returns a new collective authority of n members with new signers
// generated by g.
func NewAuthority(addrs []mino.Address, pubkeys []kyber.Point) CollectiveAuthority {
	signers := make([]crypto.Signer, len(pubkeys))
	for i, pubkey := range pubkeys {
		signers[i] = newFakeSigner(pubkey)
	}

	return CollectiveAuthority{
		pubkeys: pubkeys,
		addrs:   addrs,
		signers: signers,
	}
}

// GetPublicKey implements cosi.CollectiveAuthority.
func (ca CollectiveAuthority) GetPublicKey(addr mino.Address) (crypto.PublicKey, int) {

	for i, address := range ca.addrs {
		if address.Equal(addr) {
			return bls.NewPublicKeyFromPoint(ca.pubkeys[i]), i
		}
	}
	return nil, -1
}

// Len implements mino.Players.
func (ca CollectiveAuthority) Len() int {
	return len(ca.pubkeys)
}

// AddressIterator implements mino.Players.
func (ca CollectiveAuthority) AddressIterator() mino.AddressIterator {
	return fake.NewAddressIterator(ca.addrs)
}

func (ca CollectiveAuthority) PublicKeyIterator() crypto.PublicKeyIterator {
	return fake.NewPublicKeyIterator(ca.signers)
}

func newFakeSigner(pubkey kyber.Point) fakeSigner {
	return fakeSigner{
		pubkey: pubkey,
	}
}

// fakeSigner is a fake signer
//
// - implements crypto.Signer
type fakeSigner struct {
	crypto.Signer
	pubkey kyber.Point
}

// GetPublicKey implements crypto.Signer
func (s fakeSigner) GetPublicKey() crypto.PublicKey {
	return bls.NewPublicKeyFromPoint(s.pubkey)
}
