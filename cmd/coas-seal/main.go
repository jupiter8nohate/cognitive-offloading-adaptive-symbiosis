package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"log"
	"os"

	coas "github.com/jupiter8nohate/cognitive-offloading-adaptive-symbiosis"
)

func main() {
	in := flag.String("in", "", "input plaintext file")
	out := flag.String("out", "", "output sealed JSON file")
	flag.Parse()
	if *in == "" || *out == "" {
		log.Fatal("-in and -out are required")
	}
	keyHex := os.Getenv("COAS_PRIVATE_DRAFT_KEY_HEX")
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		log.Fatal("COAS_PRIVATE_DRAFT_KEY_HEX must be valid hex")
	}
	data, err := os.ReadFile(*in)
	if err != nil {
		log.Fatal(err)
	}
	sealed, err := coas.SealDraft(data, key)
	if err != nil {
		log.Fatal(err)
	}
	encoded, err := json.MarshalIndent(sealed, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(*out, append(encoded, '\n'), 0o600); err != nil {
		log.Fatal(err)
	}
}
