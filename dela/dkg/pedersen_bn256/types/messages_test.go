package types

import (
	"bytes"
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/require"
	"go.dedis.ch/dela/internal/testing/fake"
	"go.dedis.ch/dela/mino"
	"go.dedis.ch/dela/serde"
	"go.dedis.ch/kyber/v3"
)

var testCalls = &fake.Call{}

func init() {
	RegisterMessageFormat(fake.GoodFormat, fake.Format{Msg: Start{}, Call: testCalls})
	RegisterMessageFormat(fake.BadFormat, fake.NewBadFormat())
}

func TestStart_GetThreshold(t *testing.T) {
	start := NewStart(5, nil, nil)

	require.Equal(t, 5, start.GetThreshold())
}

func TestStart_GetAddresses(t *testing.T) {
	start := NewStart(0, []mino.Address{fake.NewAddress(0)}, nil)

	require.Len(t, start.GetAddresses(), 1)
}

func TestStart_GetPublicKeys(t *testing.T) {
	start := NewStart(0, nil, []kyber.Point{nil, nil})

	require.Len(t, start.GetPublicKeys(), 2)
}

func TestStart_Serialize(t *testing.T) {
	start := Start{}

	data, err := start.Serialize(fake.NewContext())
	require.NoError(t, err)
	require.Equal(t, fake.GetFakeFormatValue(), data)

	_, err = start.Serialize(fake.NewBadContext())
	require.EqualError(t, err, fake.Err("couldn't encode message"))
}

func TestResharingStart_GetTNew(t *testing.T) {
	start := NewStartResharing(5, 0, nil, nil, nil, nil)

	require.Equal(t, 5, start.GetTNew())
}

func TestResharingStart_GetTOld(t *testing.T) {
	start := NewStartResharing(0, 5, nil, nil, nil, nil)

	require.Equal(t, 5, start.GetTOld())
}

func TestResharingStart_GetAddrsNew(t *testing.T) {
	addrs := []mino.Address{fake.NewAddress(0)}
	start := NewStartResharing(0, 0, addrs, nil, nil, nil)

	require.Equal(t, addrs, start.GetAddrsNew())
}

func TestResharingStart_GetAddrsOld(t *testing.T) {
	addrs := []mino.Address{fake.NewAddress(0)}
	start := NewStartResharing(0, 0, nil, addrs, nil, nil)

	require.Equal(t, addrs, start.GetAddrsOld())
}

func TestResharingStart_GetPubKeysNew(t *testing.T) {
	start := NewStartResharing(0, 0, nil, nil, []kyber.Point{nil, nil}, nil)

	require.Len(t, start.GetPubkeysNew(), 2)
}

func TestResharingStart_GetPubKeysOld(t *testing.T) {
	start := NewStartResharing(0, 0, nil, nil, nil, []kyber.Point{nil, nil})

	require.Len(t, start.GetPubkeysOld(), 2)
}

func TestResharingStart_Serialize(t *testing.T) {
	start := StartResharing{}

	data, err := start.Serialize(fake.NewContext())
	require.NoError(t, err)
	require.Equal(t, fake.GetFakeFormatValue(), data)

	_, err = start.Serialize(fake.NewBadContext())
	require.EqualError(t, err, fake.Err("couldn't encode message"))
}

func TestEncryptedDeal_Getters(t *testing.T) {
	f := func(key, sig, nonce, cipher []byte) bool {
		e := NewEncryptedDeal(key, sig, nonce, cipher)

		require.Equal(t, key, e.GetDHKey())
		require.Equal(t, sig, e.GetSignature())
		require.Equal(t, nonce, e.GetNonce())
		require.Equal(t, cipher, e.GetCipher())

		return true
	}

	err := quick.Check(f, nil)
	require.NoError(t, err)
}

func TestDeal_GetIndex(t *testing.T) {
	f := func(index uint32) bool {
		deal := NewDeal(index, nil, EncryptedDeal{})

		return deal.GetIndex() == index
	}

	err := quick.Check(f, nil)
	require.NoError(t, err)
}

func TestDeal_GetSignature(t *testing.T) {
	f := func(sig []byte) bool {
		deal := NewDeal(0, sig, EncryptedDeal{})

		return bytes.Equal(deal.GetSignature(), sig)
	}

	err := quick.Check(f, nil)
	require.NoError(t, err)
}

func TestDeal_GetEncryptedDeal(t *testing.T) {
	deal := NewDeal(0, nil, EncryptedDeal{nonce: []byte{1}})

	require.Equal(t, []byte{1}, deal.GetEncryptedDeal().GetNonce())
}

func TestDeal_Serialize(t *testing.T) {
	deal := Deal{}

	data, err := deal.Serialize(fake.NewContext())
	require.NoError(t, err)
	require.Equal(t, fake.GetFakeFormatValue(), data)

	_, err = deal.Serialize(fake.NewBadContext())
	require.EqualError(t, err, fake.Err("couldn't encode deal"))
}

func TestDealerResponse_Getters(t *testing.T) {
	f := func(index uint32, status bool, sessionID, sig []byte) bool {
		resp := NewDealerResponse(index, status, sessionID, sig)

		require.Equal(t, index, resp.GetIndex())
		require.Equal(t, status, resp.GetStatus())
		require.Equal(t, sessionID, resp.GetSessionID())
		require.Equal(t, sig, resp.GetSignature())

		return true
	}

	err := quick.Check(f, nil)
	require.NoError(t, err)
}

func TestReshare_GetDeal(t *testing.T) {
	deal := Deal{index: 5}
	reshare := NewReshare(deal, nil)

	require.Equal(t, deal, reshare.GetDeal())
}

func TestReshare_GetPublicCoeffs(t *testing.T) {
	reshare := NewReshare(Deal{}, []kyber.Point{nil, nil})

	require.Len(t, reshare.GetPublicCoeffs(), 2)
}

func TestReshare_Serialize(t *testing.T) {
	reshare := Reshare{}

	data, err := reshare.Serialize(fake.NewContext())
	require.NoError(t, err)
	require.Equal(t, fake.GetFakeFormatValue(), data)

	_, err = reshare.Serialize(fake.NewBadContext())
	require.EqualError(t, err, fake.Err("couldn't encode deal"))
}

func TestResponse_GetIndex(t *testing.T) {
	f := func(index uint32) bool {
		resp := NewResponse(index, DealerResponse{})

		return index == resp.GetIndex()
	}

	err := quick.Check(f, nil)
	require.NoError(t, err)
}

func TestResponse_GetResponse(t *testing.T) {
	resp := NewResponse(0, DealerResponse{index: 1})

	require.Equal(t, uint32(1), resp.GetResponse().GetIndex())
}

func TestResponse_Serialize(t *testing.T) {
	resp := Response{}

	data, err := resp.Serialize(fake.NewContext())
	require.NoError(t, err)
	require.Equal(t, fake.GetFakeFormatValue(), data)

	_, err = resp.Serialize(fake.NewBadContext())
	require.EqualError(t, err, fake.Err("couldn't encode response"))
}

func TestStartDone_GetPublicKey(t *testing.T) {
	ack := NewStartDone(fakePoint{})

	require.Equal(t, fakePoint{}, ack.GetPublicKey())
}

func TestStartDone_Serialize(t *testing.T) {
	ack := StartDone{}

	data, err := ack.Serialize(fake.NewContext())
	require.NoError(t, err)
	require.Equal(t, fake.GetFakeFormatValue(), data)

	_, err = ack.Serialize(fake.NewBadContext())
	require.EqualError(t, err, fake.Err("couldn't encode ack"))
}

func TestSignRequest_Serialize(t *testing.T) {
	req := SignRequest{}

	data, err := req.Serialize(fake.NewContext())
	require.NoError(t, err)
	require.Equal(t, fake.GetFakeFormatValue(), data)

	_, err = req.Serialize(fake.NewBadContext())
	require.EqualError(t, err, fake.Err("couldn't encode sign request"))
}

func TestMessageFactory(t *testing.T) {
	factory := NewMessageFactory(fake.AddressFactory{})

	testCalls.Clear()

	msg, err := factory.Deserialize(fake.NewContext(), nil)
	require.NoError(t, err)
	require.Equal(t, Start{}, msg)

	require.Equal(t, 1, testCalls.Len())
	ctx := testCalls.Get(0, 0).(serde.Context)
	require.Equal(t, fake.AddressFactory{}, ctx.GetFactory(AddrKey{}))

	_, err = factory.Deserialize(fake.NewBadContext(), nil)
	require.EqualError(t, err, fake.Err("couldn't decode message"))
}

// -----------------------------------------------------------------------------
// Utility functions

type fakePoint struct {
	kyber.Point
}
