package pedersen

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"go.dedis.ch/dela/mino/minogrpc"
	"go.dedis.ch/dela/mino/router/tree"

	"go.dedis.ch/dela/dkg"

	"go.dedis.ch/dela/mino"
	"go.dedis.ch/dela/mino/minoch"

	"go.dedis.ch/kyber/v3"
)

// initDkgFailed error message indicating a DKG setup failure
const initDkgFailed = "failed to setup DKG."

// testMessage dummy message used in many tests
const testMessage = "Hello World"

// resharingUnsuccessful message showing that resharing didn't work
const resharingUnsuccessful = "Resharing was not successful"

func init() {
	rand.Seed(0)
}

// This test creates a dkg committee then creates another committee (that can
// share some nodes with the old committee) and then redistributes the secret to
// the new commitee. Using minoch as the underlying network
func TestResharing_minoch(t *testing.T) {

	// Setting up the first dkg
	nOld := 15
	thresholdOld := nOld

	minosOld := make([]mino.Mino, nOld)
	dkgsOld := make([]dkg.DKG, nOld)
	addrsOld := make([]mino.Address, nOld)
	pubkeysOld := make([]kyber.Point, len(minosOld))
	minoManager := minoch.NewManager()

	// Defining the addresses
	for i := 0; i < nOld; i++ {
		m := minoch.MustCreate(minoManager, fmt.Sprintf("addr %d", i))
		minosOld[i] = m
		addrsOld[i] = m.GetAddress()
	}

	// Initializing the pedersen
	for i, m := range minosOld {
		pdkg, pubkey := NewPedersen(m)
		dkgsOld[i] = pdkg
		pubkeysOld[i] = pubkey
	}

	fakeAuthority := NewAuthority(addrsOld, pubkeysOld)

	// Initializing the old committee actors
	actorsOld := make([]dkg.Actor, nOld)

	for i := 0; i < nOld; i++ {
		actor, err := dkgsOld[i].Listen()
		require.NoError(t, err)
		actorsOld[i] = actor
	}

	_, err := actorsOld[1].Setup(fakeAuthority, thresholdOld)
	require.NoError(t, err, initDkgFailed)

	t.Log("setup done")

	// Extract with the old committee public key.
	// It should verify
	message := []byte(testMessage)
	sig, err := actorsOld[0].Extract(message)
	require.NoError(t, err, "encrypting the message was not successful")
	_ = sig // TODO: verify

	// Setting up the second dkg nCommon is the number of nodes that are common
	// between the new and the old committee.
	nCommon := 5

	// The number of new added nodes. the new committee should have nCommon+nNew
	// nodes in total.
	nNew := 10
	thresholdNew := nCommon + nNew
	minosNew := make([]mino.Mino, nNew+nCommon)
	dkgsNew := make([]dkg.DKG, nNew+nCommon)
	addrsNew := make([]mino.Address, nNew+nCommon)

	// The first nCommon nodes of  committee are the same as the first nCommon
	// nodes of the old committee
	for i := 0; i < nCommon; i++ {
		minosNew[i] = minosOld[i]
		addrsNew[i] = minosOld[i].GetAddress()
	}

	pubkeysNew := make([]kyber.Point, len(minosNew))

	// Defining the address of the new nodes.
	for i := 0; i < nNew; i++ {
		m := minoch.MustCreate(minoManager, fmt.Sprintf("addr new %d", i))
		minosNew[i+nCommon] = m
		addrsNew[i+nCommon] = m.GetAddress()
	}

	// Initializing the pedersen of the new nodes. the common nodes already have
	// a pedersen
	for i := 0; i < nNew; i++ {
		pdkg, pubkey := NewPedersen(minosNew[i+nCommon])
		dkgsNew[i+nCommon] = pdkg
		pubkeysNew[i+nCommon] = pubkey
	}

	for i := 0; i < nCommon; i++ {
		dkgsNew[i] = dkgsOld[i]
		pubkeysNew[i] = pubkeysOld[i]
	}

	// Initializing the actor of the new nodes. the common nodes already have an
	// actor
	actorsNew := make([]dkg.Actor, nNew+nCommon)

	for i := 0; i < nCommon; i++ {
		actorsNew[i] = actorsOld[i]
	}

	for i := 0; i < nNew; i++ {
		actor, err := dkgsNew[i+nCommon].Listen()
		require.NoError(t, err)
		actorsNew[i+nCommon] = actor
	}

	// Resharing the committee secret among the new committee
	fakeAuthority = NewAuthority(addrsNew, pubkeysNew)
	err = actorsOld[0].Reshare(fakeAuthority, thresholdNew)
	require.NoError(t, err, resharingUnsuccessful)

	// Comparing the public key of the old and the new committee
	oldPubKey, err := actorsOld[0].GetPublicKey()
	require.NoError(t, err)

	for _, actorNew := range actorsNew {
		newPubKey, err := actorNew.GetPublicKey()

		// The public key should remain the same
		require.NoError(t, err, "the public key should remain the same")
		newPubKey.Equal(oldPubKey)
	}

}

// This test creats a dkg committee then creats another committee (that can
// share some nodes with the old committee) and then redistributes the secret to
// the new commitee. Using minogrpc as the underlying network
func TestResharing_minogrpc(t *testing.T) {

	// Setting up the first dkg
	nOld := 10
	thresholdOld := 10

	minosOld := make([]mino.Mino, nOld)
	dkgsOld := make([]dkg.DKG, nOld)
	addrsOld := make([]mino.Address, nOld)
	pubkeysOld := make([]kyber.Point, len(minosOld))

	// Defining the addresses
	for i := 0; i < nOld; i++ {
		addr := minogrpc.ParseAddress("127.0.0.1", 0)
		m, err := minogrpc.NewMinogrpc(addr, nil, tree.NewRouter(minogrpc.NewAddressFactory()))
		require.NoError(t, err)
		defer func() {
			_ = m.GracefulStop()
		}()

		minosOld[i] = m
		addrsOld[i] = m.GetAddress()
	}

	// Initializing the pedersen
	for i, mi := range minosOld {
		for _, mj := range minosOld {
			err := mi.(*minogrpc.Minogrpc).GetCertificateStore().Store(mj.GetAddress(),
				mj.(*minogrpc.Minogrpc).GetCertificateChain())
			require.NoError(t, err)
		}

		pdkg, pubkey := NewPedersen(mi.(*minogrpc.Minogrpc))

		dkgsOld[i] = pdkg
		pubkeysOld[i] = pubkey
	}

	fakeAuthority := NewAuthority(addrsOld, pubkeysOld)

	// Initializing the old committee actors
	actorsOld := make([]dkg.Actor, nOld)

	for i := 0; i < nOld; i++ {
		actor, err := dkgsOld[i].Listen()
		require.NoError(t, err)
		actorsOld[i] = actor
	}

	_, err := actorsOld[1].Setup(fakeAuthority, thresholdOld)
	require.NoError(t, err, initDkgFailed)

	// Extract with the old committee public key. the new committee
	// should be able to decrypt it successfully
	message := []byte(testMessage)
	sig, err := actorsOld[0].Extract(message)
	require.NoError(t, err, "signing the message was not successful")
	_ = sig // TODO: verify

	// Setting up the second dkg. nCommon is the number of nodes that are common
	// between the new and the old committee
	nCommon := 5

	// The number of new added nodes. the new committee should have nCommon+nNew
	// nodes in totatl
	nNew := 20
	thresholdNew := nCommon + nNew

	minosNew := make([]mino.Mino, nNew+nCommon)
	dkgsNew := make([]dkg.DKG, nNew+nCommon)
	addrsNew := make([]mino.Address, nNew+nCommon)

	// The first nCommon nodes of  committee are the same as the first nCommon
	// nodes of the old committee
	for i := 0; i < nCommon; i++ {
		minosNew[i] = minosOld[i]
		addrsNew[i] = minosOld[i].GetAddress()
	}

	pubkeysNew := make([]kyber.Point, len(minosNew))

	// Defining the address of the new nodes.
	for i := 0; i < nNew; i++ {
		addr := minogrpc.ParseAddress("127.0.0.1", 0)
		m, err := minogrpc.NewMinogrpc(addr, nil, tree.NewRouter(minogrpc.NewAddressFactory()))
		require.NoError(t, err)
		defer func() {
			err := m.GracefulStop()
			require.NoError(t, err)
		}()

		minosNew[i+nCommon] = m
		addrsNew[i+nCommon] = m.GetAddress()
	}

	// The first nCommon nodes of  committee are the same as the first nCommon
	// nodes of the old committee
	for i := 0; i < nCommon; i++ {
		dkgsNew[i] = dkgsOld[i]
		pubkeysNew[i] = pubkeysOld[i]
	}

	// Initializing the pedersen of the new nodes. the common nodes already have
	// a pedersen
	for i, mi := range minosNew[nCommon:] {
		for _, mj := range minosNew {
			err := mi.(*minogrpc.Minogrpc).GetCertificateStore().Store(mj.GetAddress(),
				mj.(*minogrpc.Minogrpc).GetCertificateChain())
			require.NoError(t, err)
			err = mj.(*minogrpc.Minogrpc).GetCertificateStore().Store(mi.GetAddress(),
				mi.(*minogrpc.Minogrpc).GetCertificateChain())
			require.NoError(t, err)
		}
		for _, mk := range minosOld[nCommon:] {
			err := mi.(*minogrpc.Minogrpc).GetCertificateStore().Store(mk.GetAddress(),
				mk.(*minogrpc.Minogrpc).GetCertificateChain())
			require.NoError(t, err)
			err = mk.(*minogrpc.Minogrpc).GetCertificateStore().Store(mi.GetAddress(),
				mi.(*minogrpc.Minogrpc).GetCertificateChain())
			require.NoError(t, err)
		}
		pdkg, pubkey := NewPedersen(mi.(*minogrpc.Minogrpc))
		dkgsNew[i+nCommon] = pdkg
		pubkeysNew[i+nCommon] = pubkey
	}

	// Initializing the actor of the new nodes. the common nodes already have an
	// actor
	actorsNew := make([]dkg.Actor, nNew+nCommon)

	for i := 0; i < nCommon; i++ {
		actorsNew[i] = actorsOld[i]
	}
	for i := 0; i < nNew; i++ {
		actor, err := dkgsNew[i+nCommon].Listen()
		require.NoError(t, err)
		actorsNew[i+nCommon] = actor
	}

	// Resharing the committee secret among the new committee
	fakeAuthority = NewAuthority(addrsNew, pubkeysNew)
	err = actorsOld[0].Reshare(fakeAuthority, thresholdNew)
	require.NoError(t, err, resharingUnsuccessful)

	// Comparing the public key of the old and the new committee
	oldPubKey, err := actorsOld[0].GetPublicKey()
	require.NoError(t, err)

	for _, actorNew := range actorsNew {
		newPubKey, err := actorNew.GetPublicKey()

		// The public key should remain the same
		require.NoError(t, err, "the public key should remain the same")
		newPubKey.Equal(oldPubKey)
	}
}
