package json

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.dedis.ch/f3b/smc/dkg/pedersen/types"
	"go.dedis.ch/dela/testing/fake"
	"go.dedis.ch/dela/mino"
	"go.dedis.ch/dela/serde"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/suites"
)

// suite is the Kyber suite for Pedersen.
var suite = suites.MustFind("Ed25519")

func TestMessageFormat_Start_Encode(t *testing.T) {
	start := types.NewStart(1, []mino.Address{fake.NewAddress(0)}, []kyber.Point{suite.Point()})

	format := newMsgFormat()
	ctx := serde.NewContext(fake.ContextEngine{})

	data, err := format.Encode(ctx, start)
	require.NoError(t, err)
	regexp := `{"Start":{"Threshold":1,"Addresses":\["AAAAAA=="\],"PublicKeys":\["[^"]+"\]}}`
	require.Regexp(t, regexp, string(data))

	start = types.NewStart(0, []mino.Address{fake.NewBadAddress()}, nil)
	_, err = format.Encode(ctx, start)
	require.EqualError(t, err, fake.Err("failed to encode message: couldn't marshal address"))

	start = types.NewStart(0, nil, []kyber.Point{badPoint{}})
	_, err = format.Encode(ctx, start)
	require.EqualError(t, err, fake.Err("failed to encode message: couldn't marshal public key"))

	_, err = format.Encode(fake.NewBadContext(), types.Start{})
	require.EqualError(t, err, fake.Err("couldn't marshal"))

	_, err = format.Encode(ctx, fake.Message{})
	require.EqualError(t, err, "unsupported message of type 'fake.Message'")
}

func TestMessageFormat_StartResharing_Encode(t *testing.T) {
	start := types.NewStartResharing(1, 1, []mino.Address{fake.NewAddress(0)},
		[]mino.Address{fake.NewAddress(1)}, []kyber.Point{suite.Point()},
		[]kyber.Point{suite.Point()})

	format := newMsgFormat()
	ctx := serde.NewContext(fake.ContextEngine{})

	data, err := format.Encode(ctx, start)
	require.NoError(t, err)
	regexp := `{"StartResharing":{"TNew":1,"TOld":1,"AddrsNew":\["AAAAAA=="\],"AddrsOld":\["AQAAAA=="\],"PubkeysNew":\["[^"]+"\],"PubkeysOld":\["[^"]+"\]}}`
	require.Regexp(t, regexp, string(data))

	start = types.NewStartResharing(1, 1, []mino.Address{fake.NewBadAddress()}, nil, nil, nil)
	_, err = format.Encode(ctx, start)
	require.EqualError(t, err, fake.Err("failed to encode message: couldn't marshal new address"))

	start = types.NewStartResharing(1, 1, nil, []mino.Address{fake.NewBadAddress()}, nil, nil)
	_, err = format.Encode(ctx, start)
	require.EqualError(t, err, fake.Err("failed to encode message: couldn't marshal old address"))

	start = types.NewStartResharing(1, 1, nil, nil, []kyber.Point{badPoint{}}, nil)
	_, err = format.Encode(ctx, start)
	require.EqualError(t, err, fake.Err("failed to encode message: couldn't marshal new public key"))

	start = types.NewStartResharing(1, 1, nil, nil, nil, []kyber.Point{badPoint{}})
	_, err = format.Encode(ctx, start)
	require.EqualError(t, err, fake.Err("failed to encode message: couldn't marshal old public key"))
}

func TestMessageFormat_Deal_Encode(t *testing.T) {
	deal := types.NewDeal(1, []byte{1}, types.EncryptedDeal{})

	format := newMsgFormat()
	ctx := serde.NewContext(fake.ContextEngine{})

	data, err := format.Encode(ctx, deal)
	require.NoError(t, err)
	expected := `{"Deal":{"Index":1,"Signature":"AQ==","EncryptedDeal":{"DHKey":"","Signature":"","Nonce":"","Cipher":""}}}`
	require.Equal(t, expected, string(data))
}

func TestMessageFormat_Reshare_Encode(t *testing.T) {
	reshare := types.NewReshare(types.Deal{}, []kyber.Point{suite.Point()})
	format := newMsgFormat()
	ctx := serde.NewContext(fake.ContextEngine{})

	data, err := format.Encode(ctx, reshare)
	require.NoError(t, err)
	regexp := `{"Reshare":{"Deal":{"Index":0,"Signature":"","EncryptedDeal":{"DHKey":"","Signature":"","Nonce":"","Cipher":""}},"PublicCoeff":\["[^"]+"\]}}`
	require.Regexp(t, regexp, string(data))

	reshare = types.NewReshare(types.Deal{}, []kyber.Point{badPoint{}})
	_, err = format.Encode(ctx, reshare)
	require.EqualError(t, err, fake.Err("failed to encode message: couldn't marshal public coefficient"))
}

func TestMessageFormat_Response_Encode(t *testing.T) {
	resp := types.NewResponse(1, types.DealerResponse{})

	format := newMsgFormat()
	ctx := serde.NewContext(fake.ContextEngine{})

	data, err := format.Encode(ctx, resp)
	require.NoError(t, err)
	expected := `{"Response":{"Index":1,"Response":{"SessionID":"","Index":0,"Status":false,"Signature":""}}}`
	require.Equal(t, expected, string(data))
}

func TestMessageFormat_StartDone_Encode(t *testing.T) {
	done := types.NewStartDone(suite.Point())

	format := newMsgFormat()
	ctx := serde.NewContext(fake.ContextEngine{})

	data, err := format.Encode(ctx, done)
	require.NoError(t, err)
	require.Regexp(t, `{(("StartDone":{"PublicKey":"[^"]+"}|"\w+":null),?)+}`, string(data))

	done = types.NewStartDone(badPoint{})
	_, err = format.Encode(ctx, done)
	require.EqualError(t, err, fake.Err("failed to encode message: couldn't marshal public key"))
}

func TestMessageFormat_ExtractRequest_Encode(t *testing.T) {
	req := types.NewExtractRequest([]byte{1,2,3,4})

	format := newMsgFormat()
	ctx := serde.NewContext(fake.ContextEngine{})

	data, err := format.Encode(ctx, req)
	require.NoError(t, err)
	require.Regexp(t, `{(("ExtractRequest":{"Label":"[^"]+"}|"\w+":null),?)+}`, string(data))
}

func TestMessageFormat_ExtractReply_Encode(t *testing.T) {
	resp := types.NewExtractReply([]byte{1,2,3,4})

	format := newMsgFormat()
	ctx := serde.NewContext(fake.ContextEngine{})

	data, err := format.Encode(ctx, resp)
	require.NoError(t, err)
	require.Regexp(t, `{(("ExtractReply":{"Share":"[^"]+"}|"\w+":null),?)+}`, string(data))
}

func TestMessageFormat_Decode(t *testing.T) {
	format := newMsgFormat()
	ctx := serde.NewContext(fake.ContextEngine{})
	ctx = serde.WithFactory(ctx, types.AddrKey{}, fake.AddressFactory{})

	// Decode start messages.
	expected := types.NewStart(
		5,
		[]mino.Address{fake.NewAddress(0)},
		[]kyber.Point{suite.Point()},
	)

	data, err := format.Encode(ctx, expected)
	require.NoError(t, err)

	start, err := format.Decode(ctx, data)
	require.NoError(t, err)
	require.Equal(t, expected.GetThreshold(), start.(types.Start).GetThreshold())
	require.Len(t, start.(types.Start).GetAddresses(), len(expected.GetAddresses()))
	require.Len(t, start.(types.Start).GetPublicKeys(), len(expected.GetPublicKeys()))

	_, err = format.Decode(ctx, []byte(`{"Start":{"PublicKeys":[[]]}}`))
	require.EqualError(t, err,
		"couldn't unmarshal public key: invalid Ed25519 curve point")

	badCtx := serde.WithFactory(ctx, types.AddrKey{}, nil)
	_, err = format.Decode(badCtx, []byte(`{"Start":{}}`))
	require.EqualError(t, err, "invalid factory of type '<nil>'")

	// Decode deal messages.
	deal, err := format.Decode(ctx, []byte(`{"Deal":{}}`))
	require.NoError(t, err)
	require.Equal(t, types.NewDeal(0, nil, types.EncryptedDeal{}), deal)

	// Decode response messages.
	resp, err := format.Decode(ctx, []byte(`{"Response":{}}`))
	require.NoError(t, err)
	require.Equal(t, types.NewResponse(0, types.DealerResponse{}), resp)

	// Decode start done messages.
	data = []byte(fmt.Sprintf(`{"StartDone":{"PublicKey":"%s"}}`, testPoint))
	done, err := format.Decode(ctx, data)
	require.NoError(t, err)
	require.IsType(t, types.StartDone{}, done)

	data = []byte(`{"StartDone":{"PublicKey":[]}}`)
	_, err = format.Decode(ctx, data)
	require.EqualError(t, err,
		"couldn't unmarshal public key: invalid Ed25519 curve point")

	// Decode sign request messages.
	data = []byte(`{"ExtractRequest":{}}`)
	req, err := format.Decode(ctx, data)
	require.NoError(t, err)
	require.IsType(t, types.ExtractRequest{}, req)

	// Decode sign reply messages.
	data = []byte(`{"ExtractReply":{}}`)
	resp, err = format.Decode(ctx, data)
	require.NoError(t, err)
	require.IsType(t, types.ExtractReply{}, resp)

	_, err = format.Decode(fake.NewBadContext(), []byte(`{}`))
	require.EqualError(t, err, fake.Err("couldn't deserialize message"))

	_, err = format.Decode(ctx, []byte(`{}`))
	require.EqualError(t, err, "message is empty")
}

func TestMessageFormat_Decode_StartResharing(t *testing.T) {
	format := newMsgFormat()
	ctx := serde.NewContext(fake.ContextEngine{})
	ctx = serde.WithFactory(ctx, types.AddrKey{}, fake.AddressFactory{})

	expected := types.NewStartResharing(
		5,
		6,
		[]mino.Address{fake.NewAddress(0)},
		[]mino.Address{fake.NewAddress(1)},
		[]kyber.Point{suite.Point()},
		[]kyber.Point{suite.Point()},
	)

	data, err := format.Encode(ctx, expected)
	require.NoError(t, err)

	start, err := format.Decode(ctx, data)
	require.NoError(t, err)
	require.Equal(t, expected.GetTNew(), start.(types.StartResharing).GetTNew())
	require.Equal(t, expected.GetTOld(), start.(types.StartResharing).GetTOld())
	require.Len(t, start.(types.StartResharing).GetAddrsNew(), len(expected.GetAddrsNew()))
	require.Len(t, start.(types.StartResharing).GetAddrsOld(), len(expected.GetAddrsOld()))

	badCtx := serde.WithFactory(ctx, types.AddrKey{}, nil)
	_, err = format.Decode(badCtx, []byte(`{"StartResharing":{}}`))
	require.EqualError(t, err, "invalid factory of type '<nil>'")

	_, err = format.Decode(ctx, []byte(`{"StartResharing":{"PubkeysNew":[[]]}}`))
	require.EqualError(t, err,
		"couldn't unmarshal new public key: invalid Ed25519 curve point")

	_, err = format.Decode(ctx, []byte(`{"StartResharing":{"PubkeysOld":[[]]}}`))
	require.EqualError(t, err,
		"couldn't unmarshal old public key: invalid Ed25519 curve point")
}

func TestMessageFormat_Decode_Reshare(t *testing.T) {
	format := newMsgFormat()
	ctx := serde.NewContext(fake.ContextEngine{})
	ctx = serde.WithFactory(ctx, types.AddrKey{}, fake.AddressFactory{})

	expected := types.NewReshare(
		types.NewDeal(3, []byte{}, types.NewEncryptedDeal([]byte{}, []byte{}, []byte{}, []byte{})),
		[]kyber.Point{suite.Point()},
	)

	data, err := format.Encode(ctx, expected)
	require.NoError(t, err)

	reshare, err := format.Decode(ctx, data)
	require.NoError(t, err)
	require.True(t, expected.GetPublicCoeffs()[0].Equal(reshare.(types.Reshare).GetPublicCoeffs()[0]))
	require.Equal(t, expected.GetDeal(), reshare.(types.Reshare).GetDeal())

	_, err = format.Decode(ctx, []byte(`{"Reshare":{"PublicCoeff":[[]]}}`))
	require.EqualError(t, err, "couldn't unmarshal public coeff key: invalid Ed25519 curve point")
}

// -----------------------------------------------------------------------------
// Utility functions

const testPoint = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="

type badPoint struct {
	kyber.Point
}

func (p badPoint) MarshalBinary() ([]byte, error) {
	return nil, fake.GetError()
}

type badScallar struct {
	kyber.Scalar
}

func (s badScallar) MarshalBinary() ([]byte, error) {
	return nil, fake.GetError()
}
