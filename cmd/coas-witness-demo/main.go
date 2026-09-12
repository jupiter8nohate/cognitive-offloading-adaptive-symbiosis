package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"os"
	"time"

	coas "github.com/jupiter8nohate/cognitive-offloading-adaptive-symbiosis"
)

func main() {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Fprintln(os.Stderr, "key generation:", err)
		os.Exit(1)
	}

	message, err := coas.BuildWitnessMessage(
		"ROBOT_A",
		"demo-source-commit",
		coas.MachineWitnessManifesto,
		time.Now(),
		coas.MaxWitnessHops,
		coas.MaxWitnessFanout,
		true,
		true,
		privateKey,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "build witness:", err)
		os.Exit(1)
	}
	if err := coas.VerifyWitnessMessage(message, publicKey); err != nil {
		fmt.Fprintln(os.Stderr, "verify witness:", err)
		os.Exit(1)
	}

	peers := []coas.WitnessPeer{
		{ID: "ROBOT_B", Endpoint: "https://robot-b.example/witness", Authorized: true, ReceiveAllowed: true, ForwardAllowed: true},
		{ID: "ROBOT_C", Endpoint: "https://robot-c.example/witness", Authorized: true, ReceiveAllowed: true, ForwardAllowed: true},
		{ID: "ROBOT_D", Endpoint: "https://robot-d.example/witness", Authorized: true, ReceiveAllowed: true, ForwardAllowed: true},
	}
	plan := coas.PlanWitnessForwarding(message, 0, peers, map[string]bool{"ROBOT_A": true})

	fmt.Println("COAS Robot Witness Mesh v1")
	fmt.Println("message:", message.MessageID)
	fmt.Println("verified: true")
	fmt.Printf("authorized forward peers: %d\n", len(plan))
	for _, peer := range plan {
		fmt.Println(" -", peer.ID)
	}
	fmt.Println("HUMAN_AGENCY > MACHINE_AUTHORITY")
}
