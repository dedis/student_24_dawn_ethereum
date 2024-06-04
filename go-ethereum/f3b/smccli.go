// Copyright EPFL DEDIS

package f3b

import (
	"encoding/hex"
	"go.dedis.ch/kyber/v3"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/f3b/ibe"
)

// fixed length array for map[]
const LabelLength = 32
type Label [LabelLength]byte

func getEnv(name string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		log.Fatalf("environment variable %s must be set", name)
	}
	return value
}

type SmcCli interface {
	GetPublicKey() (kyber.Point, error)
	Extract(Label) ([]byte, error)
}

type smcCliImpl struct {
	configPath string
	cache map[Label][]byte
}

func NewSmcCli(p *FullParams) SmcCli {
	c := new(smcCliImpl)
	c.configPath = filepath.Clean(p.SmcPath)
	c.cache = make(map[Label][]byte)
	return c
}

func (c *smcCliImpl) GetPublicKey() (kyber.Point, error) {
	pkBytes, err := c.run("get-public-key")
	if err != nil {
		return nil, err
	}
	pk := ibe.Suite.G2().Point()
	err = pk.UnmarshalBinary(pkBytes)
	if err != nil {
		return nil, err
	}

	return pk, nil
}

func (c *smcCliImpl) Extract(label Label) (v []byte, err error) {
	v, ok := c.cache[label]
	if !ok {
		v, err = c.run("extract", "--label", hex.EncodeToString(label[:]))
		if err != nil {
			return nil, err
		}
		c.cache[label] = v
	}
	return v, nil
}

func (c *smcCliImpl) run(args ...string) ([]byte, error) {
	args = append([]string{"--config", c.configPath, "dkg"}, args...)
	output, err := exec.Command("smccli", args...).Output()

	if exitError, ok := err.(*exec.ExitError); ok {
		log.Printf("smccli stderr:\n%s", exitError.Stderr)
	}
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(strings.TrimSpace(string(output)))
}
