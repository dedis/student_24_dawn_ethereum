// Copyright EPFL DEDIS

package ibe

import (
	"encoding/hex"
	"go.dedis.ch/kyber/v3"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func getEnv(name string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		log.Fatalf("environment variable %s must be set", name)
	}
	return value
}

type SmcCli struct {
	configPath string
}

func NewSmcCli() *SmcCli {
	configPath := filepath.Clean(getEnv("F3B_DKG_PATH"))
	return &SmcCli{configPath: configPath}
}

func (d *SmcCli) GetPublicKey() (kyber.Point, error) {
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

func (d *SmcCli) Extract(label []byte) ([]byte, error) {
	return d.run("extract", "--label", hex.EncodeToString(label))
}

func (d *SmcCli) run(args ...string) ([]byte, error) {
	args = append([]string{"--config", d.configPath, "dkg"}, args...)
	output, err := exec.Command("smccli", args...).Output()

	if exitError, ok := err.(*exec.ExitError); ok {
		log.Printf("smccli stderr:\n%s", exitError.Stderr)
	}
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(strings.TrimSpace(string(output)))
}
