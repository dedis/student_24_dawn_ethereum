// Copyright EPFL DEDIS

package f3b

import (
	"log"
	"encoding/hex"
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
func (d *DkgCli) run(extraArgs ...string) ([]byte, error) {
	var args = []string{"go", "run", "go.dedis.ch/dela/dkg/pedersen_bn256/dkgcli"}
	args = append(args, "--config", d.configPath, "dkg")
	args = append(args, extraArgs...)
	output, err := exec.Command(args[0], args[1:]...).Output()

	if exitError, ok := err.(*exec.ExitError); ok {
		log.Printf("dkgcli stderr:\n%s", exitError.Stderr)
	}
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(strings.TrimSpace(string(output)))
}
