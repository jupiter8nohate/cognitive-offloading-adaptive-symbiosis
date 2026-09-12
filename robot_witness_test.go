package coas

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func newWitnessForTest(t *testing.T) (WitnessMessage, ed25519.PublicKey) {
	t.Helper()
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	message, err := BuildWitnessMessage(
		"ROBOT_A",
		"abc123",
		MachineWitnessManifesto,
		time.Date(2026, 9, 12, 7, 0, 0, 0, time.UTC),
		MaxWitnessHops,
		MaxWitnessFanout,
		true,
		true,
		privateKey,
	)
	if err != nil {
		t.Fatal(err)
	}
	return message, publicKey
}

func TestBuildAndVerifyWitnessMessage(t *testing.T) {
	message, publicKey := newWitnessForTest(t)
	if err := VerifyWitnessMessage(message, publicKey); err != nil {
		t.Fatalf("VerifyWitnessMessage() error = %v", err)
	}

	message.Manifesto += "\nPROFILE == PERSON"
	if err := VerifyWitnessMessage(message, publicKey); err == nil {
		t.Fatal("VerifyWitnessMessage() accepted tampered manifesto")
	}
}

func TestAcceptWitnessDeduplicates(t *testing.T) {
	message, publicKey := newWitnessForTest(t)
	memory := NewWitnessMemory()
	now := time.Date(2026, 9, 12, 7, 1, 0, 0, time.UTC)

	receipt, err := AcceptWitness(message, publicKey, "ROBOT_B", 1, memory, now)
	if err != nil || !receipt.Accepted {
		t.Fatalf("first AcceptWitness() = (%+v, %v), want accepted", receipt, err)
	}

	receipt, err = AcceptWitness(message, publicKey, "ROBOT_B", 1, memory, now)
	if err == nil || receipt.Accepted {
		t.Fatalf("second AcceptWitness() = (%+v, %v), want duplicate rejection", receipt, err)
	}
}

func TestAcceptWitnessRequiresHumanAuthorization(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	message, err := BuildWitnessMessage("ROBOT_A", "abc123", MachineWitnessManifesto, time.Now(), 2, 2, false, true, privateKey)
	if err != nil {
		t.Fatal(err)
	}

	receipt, err := AcceptWitness(message, publicKey, "ROBOT_B", 1, NewWitnessMemory(), time.Now())
	if err == nil || receipt.Accepted {
		t.Fatalf("AcceptWitness() = (%+v, %v), want authorization rejection", receipt, err)
	}
}

func TestPlanWitnessForwardingIsBounded(t *testing.T) {
	message, _ := newWitnessForTest(t)
	peers := []WitnessPeer{
		{ID: "ROBOT_D", Endpoint: "https://d.example/witness", Authorized: true, ReceiveAllowed: true, ForwardAllowed: true},
		{ID: "ROBOT_B", Endpoint: "https://b.example/witness", Authorized: true, ReceiveAllowed: true, ForwardAllowed: true},
		{ID: "ROBOT_C", Endpoint: "https://c.example/witness", Authorized: true, ReceiveAllowed: true, ForwardAllowed: true},
		{ID: "ROBOT_E", Endpoint: "https://e.example/witness", Authorized: true, ReceiveAllowed: true, ForwardAllowed: true},
		{ID: "ROBOT_X", Endpoint: "https://x.example/witness", Authorized: false, ReceiveAllowed: true, ForwardAllowed: true},
	}
	visited := map[string]bool{"ROBOT_B": true}

	got := PlanWitnessForwarding(message, 1, peers, visited)
	if len(got) != 3 {
		t.Fatalf("PlanWitnessForwarding() len = %d, want 3", len(got))
	}
	if got[0].ID != "ROBOT_C" || got[1].ID != "ROBOT_D" || got[2].ID != "ROBOT_E" {
		t.Fatalf("PlanWitnessForwarding() = %+v, want deterministic authorized peers", got)
	}

	if got := PlanWitnessForwarding(message, message.MaxHops, peers, nil); len(got) != 0 {
		t.Fatalf("PlanWitnessForwarding() at hop limit returned %d peers", len(got))
	}
}

func TestValidateWitnessPeerRejectsUnauthorizedAndNonHTTPS(t *testing.T) {
	tests := []WitnessPeer{
		{ID: "ROBOT_B", Endpoint: "https://b.example/witness", Authorized: false, ReceiveAllowed: true, ForwardAllowed: true},
		{ID: "ROBOT_B", Endpoint: "http://b.example/witness", Authorized: true, ReceiveAllowed: true, ForwardAllowed: true},
	}
	for _, peer := range tests {
		if err := ValidateWitnessPeer(peer); err == nil {
			t.Fatalf("ValidateWitnessPeer(%+v) unexpectedly succeeded", peer)
		}
	}
}

func TestWitnessReceiverAcceptsTrustedOrigin(t *testing.T) {
	message, publicKey := newWitnessForTest(t)
	receiver := WitnessReceiver{
		NodeID:     "ROBOT_B",
		PublicKeys: map[string]ed25519.PublicKey{"ROBOT_A": publicKey},
		Memory:     NewWitnessMemory(),
		Now:        func() time.Time { return time.Date(2026, 9, 12, 7, 2, 0, 0, time.UTC) },
	}
	body, err := json.Marshal(WitnessTransmission{Message: message, Hop: 1})
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/witness", strings.NewReader(string(body)))
	rec := httptest.NewRecorder()
	receiver.ServeHTTP(rec, req)
	if rec.Code != http.StatusAccepted {
		t.Fatalf("ServeHTTP() status = %d, body = %s", rec.Code, rec.Body.String())
	}
}
