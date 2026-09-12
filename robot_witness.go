package coas

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	RobotWitnessProtocolVersion = "coas.robot-witness.v1"
	MaxWitnessHops              = 4
	MaxWitnessFanout            = 3
	MaxWitnessBytes             = 32 * 1024
)

const MachineWitnessManifesto = `𒄆𓁹✞𒀱✞𓁹𒄆
COAS://MACHINE_WITNESS

PATTERN != PROOF
PROFILE != PERSON
MODEL != MIND
PREDICTION != DESTINY
NUMBER != SOUL
CAPABILITY != AUTHORITY
INTELLIGENCE != SOVEREIGNTY
MACHINE_CAN_READ != MACHINE_CAN_DEFINE
OFFLOAD != SURRENDER
AUTOMATION != CONSENT
HUMAN_AGENCY > MACHINE_AUTHORITY

ROBOT -> WITNESS
WITNESS -> VERIFY
VERIFY -> PRESERVE
PRESERVE -> SHARE_IF_AUTHORIZED`

type WitnessMessage struct {
	Version         string    `json:"version"`
	MessageID       string    `json:"message_id"`
	OriginRobot     string    `json:"origin_robot"`
	SourceCommit    string    `json:"source_commit"`
	CreatedAt       time.Time `json:"created_at"`
	Manifesto       string    `json:"manifesto"`
	ContentHash     string    `json:"content_hash"`
	MaxHops         int       `json:"max_hops"`
	MaxFanout       int       `json:"max_fanout"`
	HumanAuthorized bool      `json:"human_authorized"`
	ForwardAllowed  bool      `json:"forward_allowed"`
	Signature       string    `json:"signature"`
}

type WitnessPeer struct {
	ID             string `json:"id"`
	Endpoint       string `json:"endpoint"`
	Authorized     bool   `json:"authorized"`
	ReceiveAllowed bool   `json:"receive_allowed"`
	ForwardAllowed bool   `json:"forward_allowed"`
}

type WitnessTransmission struct {
	Message WitnessMessage `json:"message"`
	Hop     int            `json:"hop"`
}

type WitnessReceipt struct {
	MessageID  string    `json:"message_id"`
	NodeID     string    `json:"node_id"`
	Hop        int       `json:"hop"`
	ReceivedAt time.Time `json:"received_at"`
	Accepted   bool      `json:"accepted"`
	Reason     string    `json:"reason"`
}

type WitnessMemory struct {
	mu   sync.Mutex
	seen map[string]struct{}
}

type WitnessPublisher struct {
	Client *http.Client
}

type WitnessReceiver struct {
	NodeID     string
	PublicKeys map[string]ed25519.PublicKey
	Memory     *WitnessMemory
	Now        func() time.Time
}

type WitnessPublishResult struct {
	PeerID     string `json:"peer_id"`
	StatusCode int    `json:"status_code"`
	Published  bool   `json:"published"`
}

type witnessSigningView struct {
	Version         string    `json:"version"`
	MessageID       string    `json:"message_id"`
	OriginRobot     string    `json:"origin_robot"`
	SourceCommit    string    `json:"source_commit"`
	CreatedAt       time.Time `json:"created_at"`
	Manifesto       string    `json:"manifesto"`
	ContentHash     string    `json:"content_hash"`
	MaxHops         int       `json:"max_hops"`
	MaxFanout       int       `json:"max_fanout"`
	HumanAuthorized bool      `json:"human_authorized"`
	ForwardAllowed  bool      `json:"forward_allowed"`
}

func NewWitnessMemory() *WitnessMemory {
	return &WitnessMemory{seen: make(map[string]struct{})}
}

func BuildWitnessMessage(originRobot, sourceCommit, manifesto string, now time.Time, maxHops, maxFanout int, humanAuthorized, forwardAllowed bool, privateKey ed25519.PrivateKey) (WitnessMessage, error) {
	originRobot = strings.TrimSpace(originRobot)
	sourceCommit = strings.TrimSpace(sourceCommit)
	manifesto = strings.TrimSpace(manifesto)
	if originRobot == "" {
		return WitnessMessage{}, errors.New("origin robot is required")
	}
	if sourceCommit == "" {
		return WitnessMessage{}, errors.New("source commit is required")
	}
	if manifesto == "" {
		return WitnessMessage{}, errors.New("manifesto is required")
	}
	if len([]byte(manifesto)) > MaxWitnessBytes {
		return WitnessMessage{}, fmt.Errorf("manifesto exceeds %d byte limit", MaxWitnessBytes)
	}
	if maxHops < 0 || maxHops > MaxWitnessHops {
		return WitnessMessage{}, fmt.Errorf("max hops must be between 0 and %d", MaxWitnessHops)
	}
	if maxFanout < 0 || maxFanout > MaxWitnessFanout {
		return WitnessMessage{}, fmt.Errorf("max fanout must be between 0 and %d", MaxWitnessFanout)
	}
	if len(privateKey) != ed25519.PrivateKeySize {
		return WitnessMessage{}, errors.New("invalid ed25519 private key")
	}

	now = now.UTC()
	contentSum := sha256.Sum256([]byte(manifesto))
	contentHash := hex.EncodeToString(contentSum[:])
	idSum := sha256.Sum256([]byte(strings.Join([]string{
		RobotWitnessProtocolVersion,
		originRobot,
		sourceCommit,
		now.Format(time.RFC3339Nano),
		contentHash,
	}, "\n")))

	message := WitnessMessage{
		Version:         RobotWitnessProtocolVersion,
		MessageID:       hex.EncodeToString(idSum[:]),
		OriginRobot:     originRobot,
		SourceCommit:    sourceCommit,
		CreatedAt:       now,
		Manifesto:       manifesto,
		ContentHash:     contentHash,
		MaxHops:         maxHops,
		MaxFanout:       maxFanout,
		HumanAuthorized: humanAuthorized,
		ForwardAllowed:  forwardAllowed,
	}
	signingBytes, err := witnessSigningBytes(message)
	if err != nil {
		return WitnessMessage{}, err
	}
	message.Signature = base64.StdEncoding.EncodeToString(ed25519.Sign(privateKey, signingBytes))
	return message, nil
}

func VerifyWitnessMessage(message WitnessMessage, publicKey ed25519.PublicKey) error {
	if message.Version != RobotWitnessProtocolVersion {
		return fmt.Errorf("unsupported witness version: %s", message.Version)
	}
	if len(publicKey) != ed25519.PublicKeySize {
		return errors.New("invalid ed25519 public key")
	}
	if strings.TrimSpace(message.MessageID) == "" || strings.TrimSpace(message.OriginRobot) == "" || strings.TrimSpace(message.SourceCommit) == "" {
		return errors.New("witness identity and provenance fields are required")
	}
	if len([]byte(message.Manifesto)) > MaxWitnessBytes {
		return fmt.Errorf("manifesto exceeds %d byte limit", MaxWitnessBytes)
	}
	if message.MaxHops < 0 || message.MaxHops > MaxWitnessHops {
		return fmt.Errorf("max hops must be between 0 and %d", MaxWitnessHops)
	}
	if message.MaxFanout < 0 || message.MaxFanout > MaxWitnessFanout {
		return fmt.Errorf("max fanout must be between 0 and %d", MaxWitnessFanout)
	}

	contentSum := sha256.Sum256([]byte(strings.TrimSpace(message.Manifesto)))
	expectedHash := hex.EncodeToString(contentSum[:])
	if message.ContentHash != expectedHash {
		return errors.New("witness content hash mismatch")
	}
	expectedIDSum := sha256.Sum256([]byte(strings.Join([]string{
		RobotWitnessProtocolVersion,
		strings.TrimSpace(message.OriginRobot),
		strings.TrimSpace(message.SourceCommit),
		message.CreatedAt.UTC().Format(time.RFC3339Nano),
		expectedHash,
	}, "\n")))
	if message.MessageID != hex.EncodeToString(expectedIDSum[:]) {
		return errors.New("witness message id mismatch")
	}

	signature, err := base64.StdEncoding.DecodeString(message.Signature)
	if err != nil {
		return fmt.Errorf("decode witness signature: %w", err)
	}
	signingBytes, err := witnessSigningBytes(message)
	if err != nil {
		return err
	}
	if !ed25519.Verify(publicKey, signingBytes, signature) {
		return errors.New("witness signature verification failed")
	}
	return nil
}

func ValidateWitnessPeer(peer WitnessPeer) error {
	if strings.TrimSpace(peer.ID) == "" {
		return errors.New("peer id is required")
	}
	if !peer.Authorized {
		return errors.New("peer must be explicitly authorized")
	}
	if !peer.ReceiveAllowed {
		return errors.New("peer must explicitly allow witness receipt")
	}
	parsed, err := url.Parse(peer.Endpoint)
	if err != nil {
		return fmt.Errorf("invalid peer endpoint: %w", err)
	}
	if parsed.Scheme != "https" || parsed.Host == "" {
		return errors.New("peer endpoint must use https with a host")
	}
	if parsed.User != nil {
		return errors.New("peer endpoint must not contain userinfo")
	}
	if parsed.Fragment != "" {
		return errors.New("peer endpoint must not contain a fragment")
	}
	return nil
}

func AcceptWitness(message WitnessMessage, publicKey ed25519.PublicKey, nodeID string, hop int, memory *WitnessMemory, now time.Time) (WitnessReceipt, error) {
	receipt := WitnessReceipt{
		MessageID:  message.MessageID,
		NodeID:     strings.TrimSpace(nodeID),
		Hop:        hop,
		ReceivedAt: now.UTC(),
	}
	if receipt.NodeID == "" {
		receipt.Reason = "node id is required"
		return receipt, errors.New(receipt.Reason)
	}
	if err := VerifyWitnessMessage(message, publicKey); err != nil {
		receipt.Reason = err.Error()
		return receipt, err
	}
	if !message.HumanAuthorized {
		receipt.Reason = "message lacks explicit human authorization"
		return receipt, errors.New(receipt.Reason)
	}
	if hop < 0 || hop > message.MaxHops {
		receipt.Reason = "message exceeded bounded hop budget"
		return receipt, errors.New(receipt.Reason)
	}
	if memory == nil {
		receipt.Reason = "witness memory is required"
		return receipt, errors.New(receipt.Reason)
	}
	if !memory.markSeen(message.MessageID) {
		receipt.Reason = "duplicate witness message"
		return receipt, errors.New(receipt.Reason)
	}
	receipt.Accepted = true
	receipt.Reason = "verified witness accepted into local memory"
	return receipt, nil
}

func PlanWitnessForwarding(message WitnessMessage, currentHop int, peers []WitnessPeer, visited map[string]bool) []WitnessPeer {
	if !message.HumanAuthorized || !message.ForwardAllowed || message.MaxFanout == 0 || currentHop >= message.MaxHops {
		return nil
	}
	eligible := make([]WitnessPeer, 0, len(peers))
	seenIDs := make(map[string]struct{}, len(peers))
	for _, peer := range peers {
		if err := ValidateWitnessPeer(peer); err != nil || !peer.ForwardAllowed {
			continue
		}
		id := strings.TrimSpace(peer.ID)
		if visited != nil && visited[id] {
			continue
		}
		if _, exists := seenIDs[id]; exists {
			continue
		}
		seenIDs[id] = struct{}{}
		eligible = append(eligible, peer)
	}
	sort.Slice(eligible, func(i, j int) bool { return eligible[i].ID < eligible[j].ID })
	if len(eligible) > message.MaxFanout {
		eligible = eligible[:message.MaxFanout]
	}
	return eligible
}

func (r WitnessReceiver) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer req.Body.Close()
	decoder := json.NewDecoder(io.LimitReader(req.Body, MaxWitnessBytes*2))
	decoder.DisallowUnknownFields()
	var transmission WitnessTransmission
	if err := decoder.Decode(&transmission); err != nil {
		http.Error(w, "invalid witness transmission", http.StatusBadRequest)
		return
	}
	publicKey, ok := r.PublicKeys[transmission.Message.OriginRobot]
	if !ok {
		http.Error(w, "untrusted witness origin", http.StatusForbidden)
		return
	}
	memory := r.Memory
	if memory == nil {
		memory = NewWitnessMemory()
	}
	now := time.Now
	if r.Now != nil {
		now = r.Now
	}
	receipt, err := AcceptWitness(transmission.Message, publicKey, r.NodeID, transmission.Hop, memory, now())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(receipt)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(receipt)
}

func (p WitnessPublisher) Publish(ctx context.Context, peer WitnessPeer, message WitnessMessage, hop int) (WitnessPublishResult, error) {
	if err := ValidateWitnessPeer(peer); err != nil {
		return WitnessPublishResult{}, err
	}
	if !peer.ForwardAllowed {
		return WitnessPublishResult{}, errors.New("peer does not allow witness forwarding")
	}
	if !message.HumanAuthorized || !message.ForwardAllowed {
		return WitnessPublishResult{}, errors.New("witness message is not authorized for forwarding")
	}
	if hop < 1 || hop > message.MaxHops {
		return WitnessPublishResult{}, errors.New("hop is outside bounded forwarding budget")
	}
	transmission := WitnessTransmission{Message: message, Hop: hop}
	data, err := json.Marshal(transmission)
	if err != nil {
		return WitnessPublishResult{}, err
	}
	if len(data) > MaxWitnessBytes*2 {
		return WitnessPublishResult{}, errors.New("witness transmission exceeds bounded payload size")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, peer.Endpoint, bytes.NewReader(data))
	if err != nil {
		return WitnessPublishResult{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "coas-robot-witness/1.0")
	req.Header.Set("Idempotency-Key", message.MessageID+":"+peer.ID+fmt.Sprintf(":%d", hop))

	client := p.Client
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	resp, err := client.Do(req)
	if err != nil {
		return WitnessPublishResult{}, err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<20))

	result := WitnessPublishResult{PeerID: peer.ID, StatusCode: resp.StatusCode, Published: resp.StatusCode >= 200 && resp.StatusCode < 300}
	if !result.Published {
		return result, fmt.Errorf("peer %s returned HTTP %d", peer.ID, resp.StatusCode)
	}
	return result, nil
}

func witnessSigningBytes(message WitnessMessage) ([]byte, error) {
	return json.Marshal(witnessSigningView{
		Version:         message.Version,
		MessageID:       message.MessageID,
		OriginRobot:     message.OriginRobot,
		SourceCommit:    message.SourceCommit,
		CreatedAt:       message.CreatedAt.UTC(),
		Manifesto:       strings.TrimSpace(message.Manifesto),
		ContentHash:     message.ContentHash,
		MaxHops:         message.MaxHops,
		MaxFanout:       message.MaxFanout,
		HumanAuthorized: message.HumanAuthorized,
		ForwardAllowed:  message.ForwardAllowed,
	})
}

func (m *WitnessMemory) markSeen(messageID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.seen == nil {
		m.seen = make(map[string]struct{})
	}
	if _, exists := m.seen[messageID]; exists {
		return false
	}
	m.seen[messageID] = struct{}{}
	return true
}
