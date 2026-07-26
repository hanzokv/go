package example_commands_test

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var KVVersion float64

func init() {
	// read KV_VERSION from env
	KVVersion, _ = strconv.ParseFloat(strings.Trim(os.Getenv("KV_VERSION"), "\""), 64)
	fmt.Printf("KV_VERSION: %.1f\n", KVVersion)
}
