// Copyright EPFL DEDIS

package types_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

// generated with log2t = 15 and label = ""
var n, _ = new(big.Int).SetString("16873719031177440519420385353106311245093946912391915790543906994565350025529232199472547339055776710866492481819771315144412970832229183934291815798702632106821315162951185682079512996913404945531313130281985163398042645504012093098965581863271589044155303951042531773225687430069494478147601856923588941378744120531835850079447506367508785797647767049817223368432657437296272691901615313485629426591413948425770414521812753240878028054404936665247636079215607325558168818824733925227663071545200559010908182889634972380886256831976586027608801802938790511380555598031436316324028712854286898938060528327070233534626", 10)
var l, _ = new(big.Int).SetString("642819411475522714700591476625638417019", 10)
var π, _ = new(big.Int).SetString("715815299172175155292255397374330261696895950777068115030450371188782653805122929976354223279518914817072864186313365527312942935726557502163741773600183700053063173561979463857677867578983741651598598545772327722021238587501308833053356002193949325899869919000817648955165908622698700985941228922833019902867948095541224984895767447893297132253276940533453925163065551486752223623110552177748807269492752599715382775666723885557809787121137944894683264956315827222346740526817158864444785715396222889158223634582808682883754971734771539154669596594539485278523404239489539599269927154086816478456426095930304271949", 10)

func TestHashEncryptedTx(t *testing.T) {
	reveal, err := rlp.EncodeToBytes([]*big.Int{l, π})
	if err != nil {
		t.Fatal(err)
	}
	dec_tx := types.NewTx(&types.DecryptedTx{
		ChainID:    big.NewInt(1),
		Nonce:      0,
		GasFeeCap:  big.NewInt(1),
		Gas:        100000,
		Value:      big.NewInt(0),
		To:         &common.Address{},
		Data:       []byte("hello world"),
		EncKey:     n.Bytes(),
		Reveal:     reveal,
	})
	enc_tx, err := dec_tx.Reencrypt()
	if err != nil {
		t.Fatal(err)
	}
	dec_hash := dec_tx.Hash().String()
	enc_hash := enc_tx.Hash().String()
	if dec_hash != enc_hash {
		t.Fatalf("Expected %s to equal %s", dec_hash, enc_hash)
	}
}

func TestSignEncryptedTx(t *testing.T) {
	reveal, err := rlp.EncodeToBytes([]*big.Int{l, π})
	if err != nil {
		t.Fatal(err)
	}
	dec_tx := types.NewTx(&types.DecryptedTx{
		ChainID:    big.NewInt(1),
		Nonce:      0,
		GasFeeCap:  big.NewInt(1),
		Gas:        100000,
		Value:      big.NewInt(0),
		To:         &common.Address{},
		Data:       []byte("hello world"),
		EncKey:     n.Bytes(),
		Reveal:     reveal,
	})
	signer := types.NewLausanneSigner(big.NewInt(1))
	acct, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	dec_tx, err = types.SignTx(dec_tx, signer, acct)
	if err != nil {
		t.Fatal(err)
	}
	enc_tx, err := dec_tx.Reencrypt()
	if err != nil {
		t.Fatal(err)
	}
	enc_sender, err := types.Sender(signer, enc_tx)
	if err != nil {
		t.Fatal(err)
	}
	dec_sender, err := types.Sender(signer, dec_tx)
	if err != nil {
		t.Fatal(err)
	}
	if enc_sender != dec_sender {
		t.Fatal("Expected sender to be the same")
	}
}
