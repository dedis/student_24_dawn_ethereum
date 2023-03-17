// Copyright EPFL DEDIS

package f3b

import (
	"log"
       "os"
       "path/filepath"
)

func getEnv(name string) string {
value, ok := os.LookupEnv(name)
       if !ok {
               log.Fatalf("environment variable %s must be set", name)
       }
       return value
}

var (
       dkgPath = filepath.Clean(getEnv("F3B_DKG_PATH"))
       gBar = getEnv("F3B_GBAR")
)

func DkgPath() string {
	return dkgPath
}

func GBar() string {
	return gBar
}

// Choice of Chain for Drand Tlock encryption
// this is the unchained drand testnet
const DrandChain = "7672797f548f3f4748ac4bf3352fc6c6b6468c9ad40ad456a397545c6e2df5bf"

// Network for Drand Tlock encryption
const DrandNetwork = "http://pl-us.testnet.drand.sh/"

// Hardcoded Drand round number for testing purposes
const RoundNumber = 1337
