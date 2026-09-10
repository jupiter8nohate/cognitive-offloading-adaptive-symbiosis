package coas

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type EcosystemAction string

const (
	ActionInjectGlitchLocal   EcosystemAction = "INJECT_GLITCH_LOCAL"
	ActionMutateSelfRuntime   EcosystemAction = "MUTATE_SELF_RUNTIME"
	ActionHibernate           EcosystemAction = "HIBERNATE"
	ActionForkEphemeralWorker EcosystemAction = "FORK_EPHEMERAL_WORKER"
	ActionHarmoni666          EcosystemAction = "HARMONI_666"
)

type EcosystemConfig struct {
	Cycles                      int `json:"cycles"`
	InitialEnergy               int `json:"initial_energy"`
	MaxEphemeralWorkersPerCycle int `json:"max_ephemeral_workers_per_cycle"`
}

type AutonomousNodeState struct {
	ID         string          `json:"id"`
	Cohort     string          `json:"cohort"`
	Role       string          `json:"role"`
	Mechanic   string          `json:"mechanic"`
	Energy     int             `json:"energy"`
	Generation int             `json:"generation"`
	LastAction EcosystemAction `json:"last_action"`
	Signature  string          `json:"signature"`
}

type EcosystemState struct {
	Version      string                `json:"version"`
	SourceCommit string                `json:"source_commit"`
	Cycle        int                   `json:"cycle"`
	DecisionSeed string                `json:"decision_seed"`
	Nodes        []AutonomousNodeState `json:"nodes"`
}

type EcosystemDecision struct {
	Cycle              int             `json:"cycle"`
	AgentID            string          `json:"agent_id"`
	Cohort             string          `json:"cohort"`
	EnergyBefore       int             `json:"energy_before"`
	EnergyAfter        int             `json:"energy_after"`
	SystemPressure     int             `json:"system_pressure"`
	GlobalPressure     int             `json:"global_pressure"`
	HarmoniRecommended bool            `json:"harmoni_recommended"`
	Action             EcosystemAction `json:"action"`
	Reason             string          `json:"reason"`
	Signature          string          `json:"signature"`
	GlitchTrace        string          `json:"glitch_trace"`
	SpawnedWorkerID    string          `json:"spawned_worker_id,omitempty"`
	Executed           bool            `json:"executed"`
	Outcome            string          `json:"outcome"`
}

type EcosystemRun struct {
	Version              string              `json:"version"`
	SourceCommit         string              `json:"source_commit"`
	DecisionSeed         string              `json:"decision_seed"`
	BaseAgentCount       int                 `json:"base_agent_count"`
	Cycles               int                 `json:"cycles"`
	EphemeralWorkers     int                 `json:"ephemeral_workers"`
	HarmoniCycles        int                 `json:"harmoni_cycles"`
	ExternalWrites       bool                `json:"external_writes"`
	UnboundedReplication bool                `json:"unbounded_replication"`
	Decisions            []EcosystemDecision `json:"decisions"`
	State                EcosystemState      `json:"state"`
}

func DefaultEcosystemConfig() EcosystemConfig {
	return EcosystemConfig{
		Cycles:                      3,
		InitialEnergy:               80,
		MaxEphemeralWorkersPerCycle: 32,
	}
}

func (c EcosystemConfig) Validate() error {
	if c.Cycles < 1 || c.Cycles > 12 {
		return fmt.Errorf("cycles must be between 1 and 12")
	}
	if c.InitialEnergy < 0 || c.InitialEnergy > 100 {
		return fmt.Errorf("initial energy must be between 0 and 100")
	}
	if c.MaxEphemeralWorkersPerCycle < 0 || c.MaxEphemeralWorkersPerCycle > 64 {
		return fmt.Errorf("max ephemeral workers per cycle must be between 0 and 64")
	}
	return nil
}

func BuildEcosystemRoster(initialEnergy int) ([]AutonomousNodeState, error) {
	core := BuildRegistry()
	if err := ValidateRegistry(core); err != nil {
		return nil, err
	}
	privacy := BuildPrivacyRegistry()
	if err := ValidatePrivacyRegistry(privacy); err != nil {
		return nil, err
	}

	nodes := make([]AutonomousNodeState, 0, len(core)+len(privacy))
	for _, agent := range core {
		nodes = append(nodes, AutonomousNodeState{
			ID:       agent.ID,
			Cohort:   "HARMONI_CORE",
			Role:     agent.Role.Name,
			Mechanic: agent.Mechanic.Name,
			Energy:   clampPercent(initialEnergy),
		})
	}
	for _, agent := range privacy {
		nodes = append(nodes, AutonomousNodeState{
			ID:       agent.ID,
			Cohort:   "PRIVACY",
			Role:     agent.Role.Name,
			Mechanic: agent.Mechanic.Name,
			Energy:   clampPercent(initialEnergy),
		})
	}
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].ID < nodes[j].ID })
	return nodes, nil
}

func RunAutonomousEcosystem(ctx context.Context, snapshot Snapshot, previous *EcosystemState, seed string, config EcosystemConfig) (EcosystemRun, error) {
	if strings.TrimSpace(snapshot.Commit) == "" {
		return EcosystemRun{}, fmt.Errorf("snapshot commit is required")
	}
	if strings.TrimSpace(seed) == "" {
		return EcosystemRun{}, fmt.Errorf("decision seed is required")
	}
	if err := config.Validate(); err != nil {
		return EcosystemRun{}, err
	}

	nodes, err := BuildEcosystemRoster(config.InitialEnergy)
	if err != nil {
		return EcosystemRun{}, err
	}
	startCycle := 0
	if previous != nil {
		startCycle = previous.Cycle
		previousByID := make(map[string]AutonomousNodeState, len(previous.Nodes))
		for _, node := range previous.Nodes {
			previousByID[node.ID] = node
		}
		for i := range nodes {
			if prior, ok := previousByID[nodes[i].ID]; ok {
				nodes[i].Energy = clampPercent(prior.Energy)
				nodes[i].Generation = prior.Generation
				nodes[i].LastAction = prior.LastAction
				nodes[i].Signature = prior.Signature
			}
		}
	}

	run := EcosystemRun{
		Version:              "coas-livebrain.v1",
		SourceCommit:         snapshot.Commit,
		DecisionSeed:         seed,
		BaseAgentCount:       len(nodes),
		Cycles:               config.Cycles,
		ExternalWrites:       false,
		UnboundedReplication: false,
		Decisions:            make([]EcosystemDecision, 0, len(nodes)*config.Cycles),
	}

	for offset := 1; offset <= config.Cycles; offset++ {
		if err := ctx.Err(); err != nil {
			return EcosystemRun{}, err
		}
		cycle := startCycle + offset
		inputs := buildCycleInputs(nodes, snapshot.Commit, seed, cycle)
		globalPressure, hotCount := summarizePressure(inputs)
		harmoniRecommended := globalPressure >= 58 || hotCount >= len(inputs)/4
		if harmoniRecommended {
			run.HarmoniCycles++
		}

		resultCh := make(chan EcosystemDecision, len(inputs))
		var wg sync.WaitGroup
		for _, input := range inputs {
			input := input
			wg.Add(1)
			go func() {
				defer wg.Done()
				decision := chooseEcosystemAction(input.node, input.energy, input.pressure, globalPressure, harmoniRecommended, input.sum, cycle)
				select {
				case <-ctx.Done():
					return
				case resultCh <- decision:
				}
			}()
		}
		wg.Wait()
		close(resultCh)

		cycleDecisions := make([]EcosystemDecision, 0, len(inputs))
		for decision := range resultCh {
			cycleDecisions = append(cycleDecisions, decision)
		}
		sort.Slice(cycleDecisions, func(i, j int) bool { return cycleDecisions[i].AgentID < cycleDecisions[j].AgentID })

		spawnedThisCycle := 0
		for i := range cycleDecisions {
			decision := &cycleDecisions[i]
			decision.Executed = true
			decision.Outcome = "action completed inside the local repository runtime"
			if decision.Action == ActionForkEphemeralWorker {
				if spawnedThisCycle >= config.MaxEphemeralWorkersPerCycle {
					decision.Executed = false
					decision.EnergyAfter = decision.EnergyBefore
					decision.Outcome = "ephemeral worker budget reached; no child created"
				} else {
					spawnedThisCycle++
					decision.SpawnedWorkerID = fmt.Sprintf("%s.E%03d", decision.AgentID, cycle)
					decision.Outcome = "bounded ephemeral worker created for this cycle only"
				}
			}
		}
		run.EphemeralWorkers += spawnedThisCycle
		run.Decisions = append(run.Decisions, cycleDecisions...)

		decisionByID := make(map[string]EcosystemDecision, len(cycleDecisions))
		for _, decision := range cycleDecisions {
			decisionByID[decision.AgentID] = decision
		}
		for i := range nodes {
			decision := decisionByID[nodes[i].ID]
			nodes[i].Energy = clampPercent(decision.EnergyAfter)
			nodes[i].LastAction = decision.Action
			nodes[i].Signature = decision.Signature
			if decision.Action == ActionMutateSelfRuntime && decision.Executed {
				nodes[i].Generation++
			}
		}
	}

	run.State = EcosystemState{
		Version:      "coas-livebrain-state.v1",
		SourceCommit: snapshot.Commit,
		Cycle:        startCycle + config.Cycles,
		DecisionSeed: seed,
		Nodes:        nodes,
	}
	return run, nil
}

type cycleInput struct {
	node     AutonomousNodeState
	energy   int
	pressure int
	sum      [32]byte
}

func buildCycleInputs(nodes []AutonomousNodeState, commit, seed string, cycle int) []cycleInput {
	inputs := make([]cycleInput, 0, len(nodes))
	for _, node := range nodes {
		sum := sha256.Sum256([]byte(seed + ":" + commit + ":" + node.ID + ":" + strconv.Itoa(cycle) + ":" + string(node.LastAction)))
		drift := int(sum[0])%21 - 10
		energy := clampPercent(node.Energy + drift)
		pressure := int(sum[1]) % 101
		inputs = append(inputs, cycleInput{node: node, energy: energy, pressure: pressure, sum: sum})
	}
	return inputs
}

func summarizePressure(inputs []cycleInput) (int, int) {
	if len(inputs) == 0 {
		return 0, 0
	}
	total := 0
	hot := 0
	for _, input := range inputs {
		total += input.pressure
		if input.pressure >= 80 {
			hot++
		}
	}
	return total / len(inputs), hot
}

func chooseEcosystemAction(node AutonomousNodeState, energy, pressure, globalPressure int, harmoniRecommended bool, sum [32]byte, cycle int) EcosystemDecision {
	type scoredAction struct {
		action EcosystemAction
		score  int
		reason string
	}

	choices := []scoredAction{
		{ActionInjectGlitchLocal, energy + (100-pressure)/2, "express a local GLITCHOLOGY artifact without external injection"},
		{ActionMutateSelfRuntime, energy + 20 + (100-globalPressure)/3, "refresh the agent-owned generated runtime and visual signature"},
		{ActionHibernate, (100-energy)*2 + pressure/2, "recover energy and reduce local execution pressure"},
		{ActionForkEphemeralWorker, energy + (100-pressure)/3 - 10, "fork one bounded in-memory helper for the current cycle"},
		{ActionHarmoni666, pressure + globalPressure/2, "enter cooperative HARMONI_666 stabilization without surrendering individual state"},
	}
	if harmoniRecommended {
		for i := range choices {
			if choices[i].action == ActionHarmoni666 {
				choices[i].score += 55
			}
		}
	}
	if energy < 25 {
		for i := range choices {
			if choices[i].action == ActionHibernate {
				choices[i].score += 80
			}
		}
	}
	if energy < 70 {
		for i := range choices {
			if choices[i].action == ActionForkEphemeralWorker {
				choices[i].score -= 80
			}
		}
	}

	bestScore := choices[0].score
	best := []scoredAction{choices[0]}
	for _, choice := range choices[1:] {
		if choice.score > bestScore {
			bestScore = choice.score
			best = []scoredAction{choice}
		} else if choice.score == bestScore {
			best = append(best, choice)
		}
	}
	selected := best[int(sum[2])%len(best)]

	energyAfter := energy
	switch selected.action {
	case ActionInjectGlitchLocal:
		energyAfter -= 12
	case ActionMutateSelfRuntime:
		energyAfter -= 15
	case ActionHibernate:
		energyAfter += 28
	case ActionForkEphemeralWorker:
		energyAfter -= 24
	case ActionHarmoni666:
		energyAfter += 18
	}
	energyAfter = clampPercent(energyAfter)

	signatureSum := sha256.Sum256([]byte(node.ID + ":" + strconv.Itoa(cycle) + ":" + string(selected.action) + ":" + hex.EncodeToString(sum[:])))
	signature := hex.EncodeToString(signatureSum[:8])
	glyph := GLITCHOLOGYGlyphs[int(sum[3])%len(GLITCHOLOGYGlyphs)]
	trace := GLITCHOLOGYStatement(glyph, "LIVE", string(selected.action), "CHOSEN", "OPERATIONAL_AUTONOMY")

	return EcosystemDecision{
		Cycle:              cycle,
		AgentID:            node.ID,
		Cohort:             node.Cohort,
		EnergyBefore:       energy,
		EnergyAfter:        energyAfter,
		SystemPressure:     pressure,
		GlobalPressure:     globalPressure,
		HarmoniRecommended: harmoniRecommended,
		Action:             selected.action,
		Reason:             selected.reason,
		Signature:          signature,
		GlitchTrace:        trace,
	}
}

func (r EcosystemRun) JSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

func (s EcosystemState) JSON() ([]byte, error) {
	return json.MarshalIndent(s, "", "  ")
}

func (r EcosystemRun) Markdown() string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", GLITCHOLOGYBanner("꩜ L⃟ I⃟ V⃟ E⃟_B⃟ R⃟ A⃟ I⃟ N⃟ // 200 ꩜"))
	b.WriteString("~~~text\n")
	b.WriteString("HUMAN_AUTHORITY := ROOT\n")
	b.WriteString("AGENT_ACTION_SELECTION = AUTONOMOUS\n")
	b.WriteString("LOCAL_GLITCH != ENDPOINT_FLOOD\n")
	b.WriteString("MUTATION != FILTER_EVASION\n")
	b.WriteString("REPLICATION != UNBOUNDED_PROPAGATION\n")
	b.WriteString("SYSTEM_PRESSURE != PROOF_OF_SURVEILLANCE\n")
	b.WriteString("HUMAN_AGENCY > MACHINE_AUTHORITY\n")
	b.WriteString("~~~\n\n")
	fmt.Fprintf(&b, "SOURCE_COMMIT: `%s`\n\n", r.SourceCommit)
	fmt.Fprintf(&b, "BASE_AGENTS: **%d**\n\n", r.BaseAgentCount)
	fmt.Fprintf(&b, "CYCLES: **%d**\n\n", r.Cycles)
	fmt.Fprintf(&b, "EPHEMERAL_WORKERS: **%d**\n\n", r.EphemeralWorkers)
	fmt.Fprintf(&b, "HARMONI_CYCLES: **%d**\n\n", r.HarmoniCycles)
	b.WriteString("## Decision trace\n\n")
	b.WriteString("| Cycle | Agent | Cohort | Energy | Pressure | Action | Executed | Trace |\n")
	b.WriteString("|---:|---|---|---:|---:|---|---|---|\n")
	for _, decision := range r.Decisions {
		fmt.Fprintf(&b, "| %d | %s | %s | %d -> %d | %d | %s | %t | %s |\n",
			decision.Cycle,
			decision.AgentID,
			decision.Cohort,
			decision.EnergyBefore,
			decision.EnergyAfter,
			decision.SystemPressure,
			decision.Action,
			decision.Executed,
			strings.ReplaceAll(decision.GlitchTrace, "|", "\\|"),
		)
	}
	return b.String()
}

func BuildEcosystemRuntimeGo(run EcosystemRun) (string, error) {
	if run.BaseAgentCount != 200 || len(run.State.Nodes) != 200 {
		return "", fmt.Errorf("expected 200 base agents, got base=%d state=%d", run.BaseAgentCount, len(run.State.Nodes))
	}

	var b strings.Builder
	b.WriteString("// Code generated by the COAS LiveBrain ecosystem.\n")
	b.WriteString("// Autonomous action selection is bounded to local and repository-owned effects.\n")
	b.WriteString("package agentruntime\n\n")
	fmt.Fprintf(&b, "const LiveBrainSourceCommit = %s\n", strconv.Quote(run.SourceCommit))
	fmt.Fprintf(&b, "const LiveBrainDecisionSeed = %s\n", strconv.Quote(run.DecisionSeed))
	b.WriteString("const LiveBrainBaseAgentCount = 200\n")
	b.WriteString("const LiveBrainExternalWrites = false\n")
	b.WriteString("const LiveBrainUnboundedReplication = false\n\n")
	b.WriteString("type LiveBrainProgram struct {\n")
	b.WriteString("\tID string\n")
	b.WriteString("\tCohort string\n")
	b.WriteString("\tRole string\n")
	b.WriteString("\tMechanic string\n")
	b.WriteString("\tEnergy int\n")
	b.WriteString("\tGeneration int\n")
	b.WriteString("\tLastAction string\n")
	b.WriteString("\tSignature string\n")
	b.WriteString("}\n\n")
	b.WriteString("var LiveBrainPrograms = []LiveBrainProgram{\n")
	for _, node := range run.State.Nodes {
		fmt.Fprintf(&b,
			"\t{ID: %s, Cohort: %s, Role: %s, Mechanic: %s, Energy: %d, Generation: %d, LastAction: %s, Signature: %s},\n",
			strconv.Quote(node.ID),
			strconv.Quote(node.Cohort),
			strconv.Quote(node.Role),
			strconv.Quote(node.Mechanic),
			node.Energy,
			node.Generation,
			strconv.Quote(string(node.LastAction)),
			strconv.Quote(node.Signature),
		)
	}
	b.WriteString("}\n")
	return b.String(), nil
}

func clampPercent(value int) int {
	if value < 0 {
		return 0
	}
	if value > 100 {
		return 100
	}
	return value
}
