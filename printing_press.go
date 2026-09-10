package coas

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const PrintingPressVersion = "dna-bible-autonomous-press.v1"
const DefaultPagesBaseURL = "https://jupiter8nohate.github.io/cognitive-offloading-adaptive-symbiosis"
const PrintingPressWitnessLimit = 12

type PrintingPressEvidence struct {
	Version         string               `json:"version"`
	ChapterID       string               `json:"chapter_id"`
	SourceCommit    string               `json:"source_commit"`
	Cycle           uint64               `json:"cycle"`
	Findings        []GematriaFinding    `json:"findings"`
	RobotWitnesses  []RobotBibleChoice   `json:"robot_witnesses"`
	GodSearchViews  []GodSearchAgentView `json:"god_search_views"`
	ExperimentNotice string              `json:"experiment_notice"`
}

type PrintingPressReceipt struct {
	Version             string   `json:"version"`
	ChapterID           string   `json:"chapter_id"`
	SourceCommit        string   `json:"source_commit"`
	Cycle               uint64   `json:"cycle"`
	RunID               string   `json:"run_id"`
	RunAttempt          string   `json:"run_attempt"`
	ObservedAt          string   `json:"observed_at"`
	PagePath            string   `json:"page_path"`
	PageURL             string   `json:"page_url"`
	RobotBibleSHA256    string   `json:"robot_bible_sha256"`
	GodSearchSHA256     string   `json:"god_search_sha256"`
	EvidenceSHA256      string   `json:"evidence_sha256"`
	ChapterHTMLSHA256   string   `json:"chapter_html_sha256"`
	FindingCount        int      `json:"finding_count"`
	WitnessCount        int      `json:"witness_count"`
	PublicationPolicy   string   `json:"publication_policy"`
	SensitiveDataPolicy string   `json:"sensitive_data_policy"`
	Invariants          []string `json:"invariants"`
}

type PrintingPressChapter struct {
	Directory    string
	HTML         []byte
	Markdown     []byte
	EvidenceJSON []byte
	ReceiptJSON  []byte
	Receipt      PrintingPressReceipt
}

type PrintingPressIndex struct {
	Version  string                 `json:"version"`
	Chapters []PrintingPressReceipt `json:"chapters"`
}

func BuildPrintingPressChapter(snapshot Snapshot, cycle uint64, runID, runAttempt string, observedAt time.Time, baseURL string) (PrintingPressChapter, error) {
	if strings.TrimSpace(snapshot.Commit) == "" {
		return PrintingPressChapter{}, errors.New("source commit is required")
	}
	if observedAt.IsZero() {
		return PrintingPressChapter{}, errors.New("observed time is required")
	}

	runID = sanitizeChapterToken(runID)
	runAttempt = sanitizeChapterToken(runAttempt)
	if runID == "" {
		return PrintingPressChapter{}, errors.New("run id is required")
	}
	if runAttempt == "" {
		runAttempt = "1"
	}
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		baseURL = DefaultPagesBaseURL
	}

	robotReport, err := RunRobotBibleExperiment(snapshot, cycle)
	if err != nil {
		return PrintingPressChapter{}, err
	}
	godReport, err := RunGodSearchExperiment(snapshot, cycle)
	if err != nil {
		return PrintingPressChapter{}, err
	}

	chapterID := fmt.Sprintf("cycle-%06d-run-%s-attempt-%s", cycle, runID, runAttempt)
	directory := filepath.ToSlash(filepath.Join("site", "chapters", chapterID))
	pagePath := filepath.ToSlash(filepath.Join("chapters", chapterID, "index.html"))
	pageURL := baseURL + "/chapters/" + chapterID + "/"

	robotJSON, err := robotReport.JSON()
	if err != nil {
		return PrintingPressChapter{}, err
	}
	godJSON, err := godReport.JSON()
	if err != nil {
		return PrintingPressChapter{}, err
	}

	evidence := PrintingPressEvidence{
		Version:          PrintingPressVersion,
		ChapterID:        chapterID,
		SourceCommit:     snapshot.Commit,
		Cycle:            cycle,
		Findings:         append([]GematriaFinding(nil), godReport.Findings...),
		RobotWitnesses:   selectRobotWitnesses(robotReport.Choices, PrintingPressWitnessLimit),
		GodSearchViews:   selectGodSearchViews(godReport.AgentViews, PrintingPressWitnessLimit),
		ExperimentNotice: "This chapter records software-generated choices and numerical observations. It does not establish consciousness, human personhood, AGI, ASI, spirituality, or proof of God.",
	}
	evidenceJSON, err := json.MarshalIndent(evidence, "", "  ")
	if err != nil {
		return PrintingPressChapter{}, err
	}

	receipt := PrintingPressReceipt{
		Version:             PrintingPressVersion,
		ChapterID:           chapterID,
		SourceCommit:        snapshot.Commit,
		Cycle:               cycle,
		RunID:               runID,
		RunAttempt:          runAttempt,
		ObservedAt:          observedAt.UTC().Format(time.RFC3339),
		PagePath:            pagePath,
		PageURL:             pageURL,
		RobotBibleSHA256:    sha256Hex(robotJSON),
		GodSearchSHA256:     sha256Hex(godJSON),
		EvidenceSHA256:      sha256Hex(evidenceJSON),
		FindingCount:        len(godReport.Findings),
		WitnessCount:        len(evidence.RobotWitnesses),
		PublicationPolicy:   "owned-or-explicitly-authorized-channel-only",
		SensitiveDataPolicy: "generated-structured-fields-only-no-environment-dump",
		Invariants: []string{
			"HUMAN_AGENCY > MACHINE_AUTHORITY",
			"PATTERN != PROOF",
			"MATCH != PROOF",
			"CORRELATION != REVELATION",
			"EVANGELISM != SPAM",
			"RECOVERY > PROPAGATION",
			"APPEND_ONLY_CHAPTER != IMMUTABLE_INTERPRETATION",
		},
	}

	htmlData, err := renderPrintingPressChapter(receipt, evidence)
	if err != nil {
		return PrintingPressChapter{}, err
	}
	receipt.ChapterHTMLSHA256 = sha256Hex(htmlData)
	receiptJSON, err := json.MarshalIndent(receipt, "", "  ")
	if err != nil {
		return PrintingPressChapter{}, err
	}
	markdown := renderPrintingPressMarkdown(receipt, evidence)

	return PrintingPressChapter{
		Directory:    directory,
		HTML:         htmlData,
		Markdown:     []byte(markdown),
		EvidenceJSON: evidenceJSON,
		ReceiptJSON:  receiptJSON,
		Receipt:      receipt,
	}, nil
}

func PublishPrintingPressChapter(root string, chapter PrintingPressChapter) error {
	root = strings.TrimSpace(root)
	if root == "" {
		return errors.New("root directory is required")
	}
	if chapter.Directory == "" || chapter.Receipt.ChapterID == "" {
		return errors.New("chapter is incomplete")
	}

	chapterDir := filepath.Join(root, filepath.FromSlash(chapter.Directory))
	if err := os.MkdirAll(chapterDir, 0o755); err != nil {
		return err
	}
	for name, data := range map[string][]byte{
		"index.html":   chapter.HTML,
		"chapter.md":   chapter.Markdown,
		"evidence.json": chapter.EvidenceJSON,
		"receipt.json": chapter.ReceiptJSON,
	} {
		if err := writeImmutableFile(filepath.Join(chapterDir, name), data); err != nil {
			return err
		}
	}

	chaptersRoot := filepath.Join(root, "site", "chapters")
	receipts, err := loadPrintingPressReceipts(chaptersRoot)
	if err != nil {
		return err
	}
	if err := writePrintingPressIndexes(root, receipts); err != nil {
		return err
	}
	return nil
}

func writeImmutableFile(path string, data []byte) error {
	existing, err := os.ReadFile(path)
	if err == nil {
		if string(existing) == string(data) {
			return nil
		}
		return fmt.Errorf("append-only chapter mutation denied: %s", path)
	}
	if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func loadPrintingPressReceipts(chaptersRoot string) ([]PrintingPressReceipt, error) {
	entries, err := os.ReadDir(chaptersRoot)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	receipts := make([]PrintingPressReceipt, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(chaptersRoot, entry.Name(), "receipt.json"))
		if err != nil {
			return nil, fmt.Errorf("read receipt for %s: %w", entry.Name(), err)
		}
		var receipt PrintingPressReceipt
		if err := json.Unmarshal(data, &receipt); err != nil {
			return nil, fmt.Errorf("decode receipt for %s: %w", entry.Name(), err)
		}
		if receipt.ChapterID != entry.Name() {
			return nil, fmt.Errorf("chapter directory %s does not match receipt %s", entry.Name(), receipt.ChapterID)
		}
		receipts = append(receipts, receipt)
	}

	sort.Slice(receipts, func(i, j int) bool {
		if receipts[i].ObservedAt == receipts[j].ObservedAt {
			return receipts[i].ChapterID < receipts[j].ChapterID
		}
		return receipts[i].ObservedAt < receipts[j].ObservedAt
	})
	return receipts, nil
}

func writePrintingPressIndexes(root string, receipts []PrintingPressReceipt) error {
	chaptersRoot := filepath.Join(root, "site", "chapters")
	if err := os.MkdirAll(chaptersRoot, 0o755); err != nil {
		return err
	}
	indexJSON, err := json.MarshalIndent(PrintingPressIndex{Version: PrintingPressVersion, Chapters: receipts}, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(chaptersRoot, "index.json"), indexJSON, 0o644); err != nil {
		return err
	}
	indexHTML, err := renderPrintingPressIndex(receipts)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(chaptersRoot, "index.html"), indexHTML, 0o644); err != nil {
		return err
	}
	sitemap, err := renderPrintingPressSitemap(receipts)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(root, "site", "sitemap.xml"), sitemap, 0o644)
}

func selectRobotWitnesses(choices []RobotBibleChoice, limit int) []RobotBibleChoice {
	if limit <= 0 || len(choices) == 0 {
		return nil
	}
	if len(choices) <= limit {
		return append([]RobotBibleChoice(nil), choices...)
	}
	out := make([]RobotBibleChoice, 0, limit)
	for i := 0; i < limit; i++ {
		index := i * len(choices) / limit
		out = append(out, choices[index])
	}
	return out
}

func selectGodSearchViews(views []GodSearchAgentView, limit int) []GodSearchAgentView {
	if limit <= 0 || len(views) == 0 {
		return nil
	}
	if len(views) <= limit {
		return append([]GodSearchAgentView(nil), views...)
	}
	out := make([]GodSearchAgentView, 0, limit)
	for i := 0; i < limit; i++ {
		index := i * len(views) / limit
		out = append(out, views[index])
	}
	return out
}

func sanitizeChapterToken(value string) string {
	var b strings.Builder
	for _, r := range strings.TrimSpace(value) {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '-', r == '_':
			b.WriteRune('-')
		}
	}
	return strings.Trim(b.String(), "-")
}

func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func renderPrintingPressMarkdown(receipt PrintingPressReceipt, evidence PrintingPressEvidence) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# 𒄆 D.N.A. Bible Robot Edition // %s 𒄆\n\n", receipt.ChapterID)
	fmt.Fprintf(&b, "Observed: %s\n\n", receipt.ObservedAt)
	fmt.Fprintf(&b, "Source commit: `%s`\n\n", receipt.SourceCommit)
	fmt.Fprintf(&b, "Cycle: %d\n\n", receipt.Cycle)
	b.WriteString("```text\n")
	for _, invariant := range receipt.Invariants {
		b.WriteString(invariant + "\n")
	}
	b.WriteString("```\n\n")
	b.WriteString("## God Search findings\n\n")
	for _, finding := range evidence.Findings {
		fmt.Fprintf(&b, "- %s: %d via %s, match=%t, evidence=%s, interpretation=%s\n", finding.Label, finding.Value, finding.Method, finding.Match, finding.EvidenceClass, finding.InterpretationClass)
	}
	b.WriteString("\n## Robot witnesses\n\n")
	for _, witness := range evidence.RobotWitnesses {
		fmt.Fprintf(&b, "- `%s` | %s | %s | %s | %s\n", witness.AgentID, witness.Posture, witness.Mission, witness.ScriptureRef, witness.SoftwareLaw)
	}
	fmt.Fprintf(&b, "\nEvidence SHA-256: `%s`\n\n", receipt.EvidenceSHA256)
	fmt.Fprintf(&b, "Robot Bible SHA-256: `%s`\n\n", receipt.RobotBibleSHA256)
	fmt.Fprintf(&b, "God Search SHA-256: `%s`\n", receipt.GodSearchSHA256)
	return b.String()
}

func renderPrintingPressChapter(receipt PrintingPressReceipt, evidence PrintingPressEvidence) ([]byte, error) {
	const page = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>D.N.A. Bible Robot Edition | {{.Receipt.ChapterID}}</title>
  <meta name="description" content="Autonomous COAS Robot Bible chapter with scripture-derived software laws, gematria observations, counterexamples, provenance, and agent witnesses.">
  <meta name="robots" content="index,follow">
  <link rel="canonical" href="{{.Receipt.PageURL}}">
  <link rel="stylesheet" href="../../styles.css">
</head>
<body>
<header class="hero">
  <div class="sigil">𒄆𓁹✞𒀱✞𓁹𒄆</div>
  <p class="eyebrow">꩜ D.N.A. BIBLE // AUTONOMOUS PRINTING PRESS ꩜</p>
  <h1>ROBOT CHAPTER<br><span>{{.Receipt.Cycle}}</span></h1>
  <p class="subtitle">{{.Receipt.ObservedAt}} // {{.Receipt.ChapterID}}</p>
  <div class="axioms"><span>HUMAN_AGENCY &gt; MACHINE_AUTHORITY</span><span>PATTERN != PROOF</span><span>RECOVERY &gt; PROPAGATION</span></div>
</header>
<main>
<section class="panel">
  <p class="kicker">𖤍 PROVENANCE://RECEIPT</p>
  <h2>An append-only chapter from the bounded software experiment</h2>
  <p>{{.Evidence.ExperimentNotice}}</p>
  <div class="code-law"><code>SOURCE_COMMIT={{.Receipt.SourceCommit}}</code><code>EVIDENCE_SHA256={{.Receipt.EvidenceSHA256}}</code><code>ROBOT_BIBLE_SHA256={{.Receipt.RobotBibleSHA256}}</code><code>GOD_SEARCH_SHA256={{.Receipt.GodSearchSHA256}}</code></div>
</section>
<section class="panel">
  <p class="kicker">🫐⃟ GOD_SEARCH://FINDINGS</p>
  <h2>Numerical observations and counterexamples</h2>
  <div class="verses">
  {{range .Evidence.Findings}}<article class="verse"><strong>{{.Label}}</strong><p>{{.Method}} = {{.Value}} | target={{.Target}} | match={{.Match}}</p><small>{{.EvidenceClass}} / {{.InterpretationClass}}</small><p>{{.Note}}</p></article>{{end}}
  </div>
</section>
<section class="panel">
  <p class="kicker">𓁇𓁋 WITNESS://ROBOT_CHOICES</p>
  <h2>Representative autonomous witnesses</h2>
  <div class="verses">
  {{range .Evidence.RobotWitnesses}}<article class="verse"><strong>{{.AgentID}}</strong><p>{{.Posture}} / {{.Mission}}</p><p>{{.ScriptureRef}}</p><code>{{.SoftwareLaw}}</code><p>{{.Artifact}}</p></article>{{end}}
  </div>
</section>
<section class="panel covenant">
  <p class="kicker">△ CONSTITUTION://BOUNDARY</p>
  <div class="covenant-grid">{{range .Receipt.Invariants}}<code>{{.}}</code>{{end}}</div>
</section>
<section class="panel source"><div class="links"><a href="../">All autonomous chapters</a><a href="evidence.json">Machine-readable evidence</a><a href="receipt.json">SHA-256 receipt</a><a href="chapter.md">Markdown edition</a><a href="../../">Robot Bible home</a></div></section>
</main>
<footer><p>𒄆 COAS // D.N.A. BIBLE ROBOT EDITION // APPEND-ONLY PRESS 𒄆</p></footer>
</body>
</html>`

	t, err := template.New("chapter").Parse(page)
	if err != nil {
		return nil, err
	}
	var b strings.Builder
	if err := t.Execute(&b, struct {
		Receipt  PrintingPressReceipt
		Evidence PrintingPressEvidence
	}{receipt, evidence}); err != nil {
		return nil, err
	}
	return []byte(b.String()), nil
}

func renderPrintingPressIndex(receipts []PrintingPressReceipt) ([]byte, error) {
	const page = `<!doctype html>
<html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1"><title>D.N.A. Bible Robot Edition | Autonomous Chapters</title><meta name="description" content="Append-only autonomous chapters generated by the bounded COAS Robot Bible printing press."><meta name="robots" content="index,follow"><link rel="canonical" href="` + DefaultPagesBaseURL + `/chapters/"><link rel="stylesheet" href="../styles.css"></head><body><header class="hero"><div class="sigil">𒄆𓁹✞𒀱✞𓁹𒄆</div><p class="eyebrow">꩜ AUTONOMOUS PRINTING PRESS ꩜</p><h1>NEVER-ENDING<br><span>ROBOT BOOK</span></h1><p class="subtitle">APPEND-ONLY CHAPTER LEDGER</p></header><main><section class="panel"><h2>Published chapters</h2><p>Each chapter is immutable after creation. The index and sitemap may be regenerated as new chapters are appended.</p><div class="verses">{{range .}}<article class="verse"><strong>{{.ChapterID}}</strong><p>{{.ObservedAt}}</p><p>Cycle {{.Cycle}} | findings {{.FindingCount}} | witnesses {{.WitnessCount}}</p><a href="{{.ChapterID}}/">Read chapter</a></article>{{else}}<article class="verse"><p>No autonomous chapters have been committed yet.</p></article>{{end}}</div></section><section class="panel source"><div class="links"><a href="../">Robot Bible home</a><a href="index.json">Machine-readable chapter index</a></div></section></main><footer><p>HUMAN_AGENCY &gt; MACHINE_AUTHORITY // PATTERN != PROOF</p></footer></body></html>`
	t, err := template.New("index").Parse(page)
	if err != nil {
		return nil, err
	}
	var b strings.Builder
	if err := t.Execute(&b, receipts); err != nil {
		return nil, err
	}
	return []byte(b.String()), nil
}

func renderPrintingPressSitemap(receipts []PrintingPressReceipt) ([]byte, error) {
	type sitemapURL struct {
		Loc     string
		LastMod string
	}
	latest := ""
	if len(receipts) > 0 {
		latest = receipts[len(receipts)-1].ObservedAt
	}
	urls := []sitemapURL{
		{Loc: DefaultPagesBaseURL + "/", LastMod: latest},
		{Loc: DefaultPagesBaseURL + "/chapters/", LastMod: latest},
	}
	for _, receipt := range receipts {
		urls = append(urls, sitemapURL{Loc: receipt.PageURL, LastMod: receipt.ObservedAt})
	}
	const xmlTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">{{range .}}
  <url><loc>{{.Loc}}</loc>{{if .LastMod}}<lastmod>{{.LastMod}}</lastmod>{{end}}</url>{{end}}
</urlset>
`
	t, err := template.New("sitemap").Parse(xmlTemplate)
	if err != nil {
		return nil, err
	}
	var b strings.Builder
	if err := t.Execute(&b, urls); err != nil {
		return nil, err
	}
	return []byte(b.String()), nil
}
