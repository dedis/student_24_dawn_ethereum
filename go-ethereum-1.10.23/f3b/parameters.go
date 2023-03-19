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
