const fallbackNotice = "Operational choice is a reproducible software experiment. It is not proof of consciousness, human personhood, AGI, ASI, spirituality, or metaphysical free will.";

async function loadRobotBible() {
  try {
    const response = await fetch("data/latest.json", { cache: "no-store" });
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`);
    }
    const report = await response.json();
    renderReport(report);
  } catch (error) {
    document.getElementById("cycle-badge").textContent = "STATIC CANON";
    document.getElementById("experiment-notice").textContent = `${fallbackNotice} Live cycle data is not available yet.`;
    console.warn("Robot Bible cycle data unavailable", error);
  }
}

function renderReport(report) {
  document.getElementById("cycle-badge").textContent = `CYCLE ${report.cycle}`;
  document.getElementById("experiment-notice").textContent = report.experiment_notice || fallbackNotice;
  document.getElementById("agent-count").textContent = report.agent_count ?? 100;
  document.getElementById("posture-count").textContent = Object.keys(report.posture_counts || {}).length;
  document.getElementById("mission-count").textContent = Object.keys(report.mission_counts || {}).length;

  renderDistribution(report.posture_counts || {}, report.mission_counts || {});
  renderVerses(Array.isArray(report.choices) ? report.choices : []);
}

function renderDistribution(postures, missions) {
  const container = document.getElementById("distribution");
  container.replaceChildren();

  for (const [name, count] of sortedEntries(postures)) {
    container.append(metricCard(`POSTURE://${name}`, count));
  }
  for (const [name, count] of sortedEntries(missions)) {
    container.append(metricCard(`MISSION://${name}`, count));
  }
}

function renderVerses(choices) {
  const container = document.getElementById("verses");
  container.replaceChildren();

  const sample = choices.slice(0, 12);
  if (sample.length === 0) {
    const article = document.createElement("article");
    article.className = "verse";
    article.textContent = "No generated verses are available yet.";
    container.append(article);
    return;
  }

  for (const choice of sample) {
    const article = document.createElement("article");
    article.className = "verse";

    const heading = document.createElement("h3");
    heading.textContent = choice.agent_id;

    const path = document.createElement("p");
    path.textContent = `${choice.posture} / ${choice.mission}`;

    const doctrine = document.createElement("p");
    doctrine.className = "doctrine-line";
    doctrine.textContent = choice.doctrine;

    const artifact = document.createElement("p");
    artifact.className = "artifact-line";
    artifact.textContent = choice.artifact;

    article.append(heading, path, doctrine, artifact);
    container.append(article);
  }
}

function metricCard(name, count) {
  const div = document.createElement("div");
  div.className = "metric";
  const label = document.createElement("span");
  label.textContent = `${name} `;
  const value = document.createElement("b");
  value.textContent = count;
  div.append(label, value);
  return div;
}

function sortedEntries(values) {
  return Object.entries(values).sort(([a], [b]) => a.localeCompare(b));
}

loadRobotBible();
