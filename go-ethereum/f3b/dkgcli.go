// Copyright EPFL DEDIS

package f3b

import (
	"log"
	"encoding/hex"
       "os"
       "os/exec"
       "path/filepath"
       "strings"
       "go.dedis.ch/kyber/v3"
)

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

func (d *DkgCli) GetPublicKey() (kyber.Point, error) {
	pkBytes, err := d.run("get-public-key")
	if err != nil {
		return nil, err
	}
	pk := Suite.G2().Point()
	err = pk.UnmarshalBinary(pkBytes)
	if err != nil {
		return nil, err
	}

	return pk, nil
}

func (d *DkgCli) Extract(label []byte) ([]byte, error) {
	return d.run("extract", "--label", hex.EncodeToString(label))
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

