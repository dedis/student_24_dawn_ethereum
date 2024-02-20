package types

import (
	"bytes"
	"go.dedis.ch/dela/mino"
	"go.dedis.ch/dela/serde"
	"go.dedis.ch/dela/serde/registry"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/sign/tbls"
	"golang.org/x/xerrors"
)

var msgFormats = registry.NewSimpleRegistry()

// RegisterMessageFormat register the engine for the provided format.
func RegisterMessageFormat(c serde.Format, f serde.FormatEngine) {
	msgFormats.Register(c, f)
}

// Start is the message the initiator of the DKG protocol should send to all the
// nodes.
//
// - implements serde.Message
type Start struct {
	// threshold
	thres int
	// the full list of addresses that will participate in the DKG
	addresses []mino.Address
	// the corresponding kyber.Point pub keys of the addresses
	pubkeys []kyber.Point
}

// NewStart creates a new start message.
func NewStart(thres int, addrs []mino.Address, pubkeys []kyber.Point) Start {
	return Start{
		thres:     thres,
		addresses: addrs,
		pubkeys:   pubkeys,
	}
}

// GetThreshold returns the threshold.
func (s Start) GetThreshold() int {
	return s.thres
}

// GetAddresses returns the list of addresses.
func (s Start) GetAddresses() []mino.Address {
	return emptyIfNil(s.addresses)
}

// GetPublicKeys returns the list of public keys.
func (s Start) GetPublicKeys() []kyber.Point {
	return emptyIfNil(s.pubkeys)
}

// Serialize implements serde.Message. It looks up the format and returns the
// serialized data for the start message.
func (s Start) Serialize(ctx serde.Context) ([]byte, error) {
	format := msgFormats.Get(ctx.GetFormat())

	data, err := format.Encode(ctx, s)
	if err != nil {
		return nil, xerrors.Errorf("couldn't encode message: %v", err)
	}

	return data, nil
}

// StartResharing is the message the initiator of the resharing protocol
// should send to all the old nodes.
//
// - implements serde.Message
type StartResharing struct {
	// New threshold
	tNew int
	// Old threshold
	tOld int
	// The full list of addresses that will participate in the new DKG
	addrsNew []mino.Address
	// The full list of addresses of old dkg members
	addrsOld []mino.Address
	// The corresponding kyber.Point pub keys of the new addresses
	pubkeysNew []kyber.Point
	// The corresponding kyber.Point pub keys of the old addresses
	pubkeysOld []kyber.Point
}

// NewStartResharing creates a new start resharing message.
func NewStartResharing(tNew int, tOld int, addrsNew []mino.Address, addrsOld []mino.Address,
	pubkeysNew []kyber.Point, pubkeysOld []kyber.Point) StartResharing {
	return StartResharing{
		tNew:       tNew,
		tOld:       tOld,
		addrsNew:   addrsNew,
		addrsOld:   addrsOld,
		pubkeysNew: pubkeysNew,
		pubkeysOld: pubkeysOld,
	}
}

// GetTNew returns the new threshold.
func (r StartResharing) GetTNew() int {
	return r.tNew
}

// GetTOld returns the old threshold.
func (r StartResharing) GetTOld() int {
	return r.tOld
}

// GetAddrsNew returns the list of new addresses.
func (r StartResharing) GetAddrsNew() []mino.Address {
	return emptyIfNil(r.addrsNew)
}

// GetAddrsOld returns the list of old addresses.
func (r StartResharing) GetAddrsOld() []mino.Address {
	return emptyIfNil(r.addrsOld)
}

// GetPubkeysNew returns the list of new public keys.
func (r StartResharing) GetPubkeysNew() []kyber.Point {
	return emptyIfNil(r.pubkeysNew)
}

// GetPubkeysOld returns the list of old public keys.
func (r StartResharing) GetPubkeysOld() []kyber.Point {
	return emptyIfNil(r.pubkeysOld)
}

// Serialize implements serde.Message. It looks up the format and returns the
// serialized data for the resharingRequest message.
func (r StartResharing) Serialize(ctx serde.Context) ([]byte, error) {
	format := msgFormats.Get(ctx.GetFormat())

	data, err := format.Encode(ctx, r)
	if err != nil {
		return nil, xerrors.Errorf("couldn't encode message: %v", err)
	}

	return data, nil
}

// EncryptedDeal contains the different parameters and data of an encrypted
// deal.
type EncryptedDeal struct {
	dhkey     []byte
	signature []byte
	nonce     []byte
	cipher    []byte
}

// NewEncryptedDeal creates a new encrypted deal message.
func NewEncryptedDeal(dhkey, sig, nonce, cipher []byte) EncryptedDeal {
	return EncryptedDeal{
		dhkey:     dhkey,
		signature: sig,
		nonce:     nonce,
		cipher:    cipher,
	}
}

// GetDHKey returns the Diffie-Helmann key in bytes.
func (d EncryptedDeal) GetDHKey() []byte {
	return emptyIfNil(d.dhkey)
}

// GetSignature returns the signatures in bytes.
func (d EncryptedDeal) GetSignature() []byte {
	return emptyIfNil(d.signature)
}

// GetNonce returns the nonce in bytes.
func (d EncryptedDeal) GetNonce() []byte {
	return emptyIfNil(d.nonce)
}

// GetCipher returns the cipher in bytes.
func (d EncryptedDeal) GetCipher() []byte {
	return emptyIfNil(d.cipher)
}

// Deal matches the attributes defined in kyber dkg.Deal.
//
// - implements serde.Message
type Deal struct {
	index     uint32
	signature []byte

	encryptedDeal EncryptedDeal
}

// NewDeal creates a new deal.
func NewDeal(index uint32, sig []byte, e EncryptedDeal) Deal {
	return Deal{
		index:         index,
		signature:     sig,
		encryptedDeal: e,
	}
}

// GetIndex returns the index.
func (d Deal) GetIndex() uint32 {
	return d.index
}

// GetSignature returns the signature in bytes.
func (d Deal) GetSignature() []byte {
	return emptyIfNil(d.signature)
}

// GetEncryptedDeal returns the encrypted deal.
func (d Deal) GetEncryptedDeal() EncryptedDeal {
	return d.encryptedDeal
}

// Serialize implements serde.Message.
func (d Deal) Serialize(ctx serde.Context) ([]byte, error) {
	format := msgFormats.Get(ctx.GetFormat())

	data, err := format.Encode(ctx, d)
	if err != nil {
		return nil, xerrors.Errorf("couldn't encode deal: %v", err)
	}

	return data, nil
}

// Reshare messages for resharing process
// - implements serde.Message
type Reshare struct {
	deal        Deal
	publicCoeff []kyber.Point
}

// NewReshare creates a new deal.
func NewReshare(deal Deal, publicCoeff []kyber.Point) Reshare {
	return Reshare{
		deal:        deal,
		publicCoeff: publicCoeff,
	}
}

// GetDeal returns the deal.
func (d Reshare) GetDeal() Deal {
	return d.deal
}

// GetPublicCoeffs returns the public coeff.
func (d Reshare) GetPublicCoeffs() []kyber.Point {
	return d.publicCoeff
}

// Serialize implements serde.Message.
func (d Reshare) Serialize(ctx serde.Context) ([]byte, error) {
	format := msgFormats.Get(ctx.GetFormat())

	data, err := format.Encode(ctx, d)
	if err != nil {
		return nil, xerrors.Errorf("couldn't encode deal: %v", err)
	}

	return data, nil
}

// DealerResponse is a response of a single dealer.
type DealerResponse struct {
	sessionID []byte
	// Index of the verifier issuing this Response from the new set of
	// nodes.
	index     uint32
	status    bool
	signature []byte
}

// NewDealerResponse creates a new dealer response.
func NewDealerResponse(index uint32, status bool, sessionID, sig []byte) DealerResponse {
	return DealerResponse{
		sessionID: sessionID,
		index:     index,
		status:    status,
		signature: sig,
	}
}

// GetSessionID returns the session ID in bytes.
func (dresp DealerResponse) GetSessionID() []byte {
	return emptyIfNil(dresp.sessionID)
}

// GetIndex returns the index.
func (dresp DealerResponse) GetIndex() uint32 {
	return dresp.index
}

// GetStatus returns the status.
func (dresp DealerResponse) GetStatus() bool {
	return dresp.status
}

// GetSignature returns the signature in bytes.
func (dresp DealerResponse) GetSignature() []byte {
	return emptyIfNil(dresp.signature)
}

// Response matches the attributes defined in kyber pedersen.Response.
//
// - implements serde.Message
type Response struct {
	// Index of the Dealer this response is for.
	index    uint32
	response DealerResponse
}

// NewResponse creates a new response.
func NewResponse(index uint32, r DealerResponse) Response {
	return Response{
		index:    index,
		response: r,
	}
}

// GetIndex returns the index.
func (r Response) GetIndex() uint32 {
	return r.index
}

// GetResponse returns the dealer response.
func (r Response) GetResponse() DealerResponse {
	return r.response
}

// Serialize implements serde.Message.
func (r Response) Serialize(ctx serde.Context) ([]byte, error) {
	format := msgFormats.Get(ctx.GetFormat())

	data, err := format.Encode(ctx, r)
	if err != nil {
		return nil, xerrors.Errorf("couldn't encode response: %v", err)
	}

	return data, nil
}

// StartDone should be sent by all the nodes to the initiator of the DKG when
// the DKG setup is done.
//
// - implements serde.Message
type StartDone struct {
	pubkey kyber.Point
}

// NewStartDone creates a new start done message.
func NewStartDone(pubkey kyber.Point) StartDone {
	return StartDone{
		pubkey: pubkey,
	}
}

// GetPublicKey returns the public key of the LTS.
func (s StartDone) GetPublicKey() kyber.Point {
	return s.pubkey
}

// Serialize implements serde.Message.
func (s StartDone) Serialize(ctx serde.Context) ([]byte, error) {
	format := msgFormats.Get(ctx.GetFormat())

	data, err := format.Encode(ctx, s)
	if err != nil {
		return nil, xerrors.Errorf("couldn't encode ack: %v", err)
	}

	return data, nil
}

// SignRequest is a message sent to request a threshold signature share.
//
// - implements serde.Message
type SignRequest struct {
	msg []byte
}

// NewSignRequest creates a new signature request.
func NewSignRequest(msg []byte) SignRequest {
	return SignRequest{
		msg: bytes.Clone(msg),
	}
}

// GetMsg returns the message being signed.
func (req SignRequest) GetMsg() []byte {
	return bytes.Clone(req.msg)
}

// Serialize implements serde.Message.
func (req SignRequest) Serialize(ctx serde.Context) ([]byte, error) {
	format := msgFormats.Get(ctx.GetFormat())

	data, err := format.Encode(ctx, req)
	if err != nil {
		return nil, xerrors.Errorf("couldn't encode sign request: %v", err)
	}

	return data, nil
}

// SignReply is the response of a decryption request.
//
// - implements serde.Message
type SignReply struct {
	Share tbls.SigShare
}

// NewSignReply returns a new decryption reply.
func NewSignReply(share tbls.SigShare) SignReply {
	return SignReply{
		Share: bytes.Clone(share),
	}
}

// Serialize implements serde.Message.
func (resp SignReply) Serialize(ctx serde.Context) ([]byte, error) {
	format := msgFormats.Get(ctx.GetFormat())

	data, err := format.Encode(ctx, resp)
	if err != nil {
		return nil, xerrors.Errorf("couldn't encode sign reply: %v", err)
	}

	return data, nil
}

// AddrKey is the key for the address factory.
type AddrKey struct{}

// MessageFactory is a message factory for the different DKG messages.
//
// - implements serde.Factory
type MessageFactory struct {
	addrFactory mino.AddressFactory
}

// NewMessageFactory returns a message factory for the DKG.
func NewMessageFactory(f mino.AddressFactory) MessageFactory {
	return MessageFactory{
		addrFactory: f,
	}
}

// Deserialize implements serde.Factory.
func (f MessageFactory) Deserialize(ctx serde.Context, data []byte) (serde.Message, error) {
	format := msgFormats.Get(ctx.GetFormat())

	ctx = serde.WithFactory(ctx, AddrKey{}, f.addrFactory)

	msg, err := format.Decode(ctx, data)
	if err != nil {
		return nil, xerrors.Errorf("couldn't decode message: %v", err)
	}

	return msg, nil
}

func emptyIfNil[T any](slice []T) []T {
	if slice == nil {
		return make([]T, 0)
	}

	return slice
}
