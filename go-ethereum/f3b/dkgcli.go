// Copyright EPFL DEDIS

package f3b

import (
	"log"
	"encoding/hex"
       "os"
       "os/exec"
       "path/filepath"
       "strings"
       "go.dedis.ch/kyber/v3/pairing"
       "go.dedis.ch/kyber/v3/pairing/bn256"
       "go.dedis.ch/kyber/v3"
)

var Suite pairing.Suite = bn256.NewSuite()

func getEnv(name string) string {
value, ok := os.LookupEnv(name)
       if !ok {
               log.Fatalf("environment variable %s must be set", name)
       }
       return value
}

type DkgCli struct {
	configPath string
}

func NewDkgCli() *DkgCli {
	configPath := filepath.Clean(getEnv("F3B_DKG_PATH"))
	return &DkgCli{configPath: configPath}
}

func (d *DkgCli) Encrypt(label []byte, plaintext []byte) ([]byte, error) {
	return d.run("encrypt", "--label", hex.EncodeToString(label), "--message", hex.EncodeToString(plaintext))
}
func (d *DkgCli) Decrypt(label []byte, ciphertext []byte) ([]byte, error) {
	return d.run("decrypt", "--label", hex.EncodeToString(label), "--ciphertext", hex.EncodeToString(ciphertext))
}
func (d *DkgCli) Extract(label []byte) (kyber.Point, error) {
	identityBytes, err := d.run("decrypt", "--label", hex.EncodeToString(label))
	if err != nil {
		return nil, err
	}

	identity := Suite.G1().Point()
	err = identity.UnmarshalBinary(identityBytes)
	if err != nil {
		return nil, err
	}

	return identity, nil
}

func (d *DkgCli) run(args ...string) ([]byte, error) {
	args = append([]string{"--config", d.configPath, "dkg"}, args...)
	output, err := exec.Command("dkgcli", args...).Output()

	if exitError, ok := err.(*exec.ExitError); ok {
		log.Printf("dkgcli stderr:\n%s", exitError.Stderr)
	}
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(strings.TrimSpace(string(output)))
}

