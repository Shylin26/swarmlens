# SwarmLens

**Real-time observability and control plane for multi-agent LLM systems.**

Multi-agent systems built with AutoGen, CrewAI, LangGraph, or custom frameworks are notoriously hard to debug. Agent-to-agent messages fly past with no single log to tail, and failures — infinite negotiation loops, one agent silently dominating the workload, runaway API costs — are invisible until they've already cost you money or trust.

SwarmLens is a Kafka-native observability layer that sits alongside your agent swarm. Every message flows through a live event bus, gets checked in real time against detectors for loops, cost spikes, and role collapse, and is fully replayable after the fact — so you can see exactly what happened, when, and why.

![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go)
![Kafka](https://img.shields.io/badge/Kafka-KRaft-231F20?logo=apachekafka)
![Redis](https://img.shields.io/badge/Redis-7-DC382D?logo=redis)
![Grafana](https://img.shields.io/badge/Grafana-live_dashboards-F46800?logo=grafana)

---

## Why this exists

Most multi-agent tooling stops at orchestration — it helps you *build* a swarm, but gives you nothing once it's running. SwarmLens exists for the moment after that: your swarm is live, something's wrong, and you need to know what, without manually grepping through scattered logs.

It's built around a simple idea from current multi-agent research: healthy swarms show **behavioral differentiation** — agents settle into distinct roles. Unhealthy ones don't — they loop, one agent takes over everyone else's work, or costs spiral with no clear cause. SwarmLens turns those failure patterns into concrete, running detectors instead of vague warning signs.

## What it does

- **Live event bus** — every agent action (message, tool call, decision) flows through Kafka, keyed by swarm so ordering and replay stay correct.
- **Three real-time detectors**, running against live traffic:
  - **Loop detection** — catches agents stuck in repeating back-and-forth patterns.
  - **Cost-anomaly detection** — EWMA-based tracking flags per-event cost spikes against a swarm's own rolling baseline.
  - **Role-collapse detection** — flags when one agent ends up completing a disproportionate share of tasks, the opposite of healthy role differentiation.
- **Live Grafana dashboards** — event throughput, active agents, and alerts, fed by Prometheus metrics.
- **Replay** — reconstruct any swarm's exact message history from a point in time, human-readable, no Kafka knowledge required.
- **Diff** — compare two swarm runs side by side: alert counts and role distribution, at a glance.
- **Multi-language SDKs** — emit events from Go or Python; any framework that can produce JSON over Kafka can plug in.

## Quickstart

```bash
git clone https://github.com/Shylin26/swarmlens.git
cd swarmlens
docker compose -f deploy/docker-compose.yml up -d
```

This brings up Kafka, Redis, Prometheus, and Grafana locally.

Run the control plane:

```bash
go run ./cmd/swarmlens
```

Run the example swarm (3 agents — planner, worker, reviewer — collaborating on a real task via a local LLM):

```bash
go run ./examples/swarm-demo
```

Watch it live in Grafana at `http://localhost:3000`, or replay it from the CLI:

```bash
go run ./cmd/swarmlens-cli replay --swarm-id swarm-demo-1 --from 10m
```

## Architecture

```
 Agents (Go/Python SDK)
        │
        ▼
   Kafka (agent.messages)
        │
        ▼
  Go control plane ──► Redis (live state, sliding windows, EWMA)
        │
        ├──► Prometheus ──► Grafana (live dashboards)
        │
        └──► Detectors (loop / cost / role-collapse) ──► alerts
```

Every topic is keyed by `swarm_id`, so all of one swarm's events land on the same Kafka partition — this is what makes both live ordering and after-the-fact replay correct by construction, not by extra bookkeeping.

## Detectors, in detail

**Loop detection** — maintains a rolling window of the last N `sender→recipient` message shapes per swarm, and flags when a short pattern repeats back-to-back beyond a threshold. Catches the most common real failure: two agents stuck volleying the same clarification request.

**Cost anomaly** — tracks an exponentially-weighted moving average (EWMA) of per-event cost per swarm, and flags any event costing well above that swarm's own recent baseline (default: 5x). Reacts to a swarm's actual spending pattern, not a fixed global threshold.

**Role collapse** — tracks which agent completes each `parent_task_id`-linked subtask, and flags when one agent's share crosses a threshold (default: 70%). Operationalizes the idea that healthy multi-agent systems distribute work, not concentrate it.

All three detectors are unit-tested against synthetic scenarios independent of any live LLM, and proven against real Kafka traffic — see `internal/detect/`.

## CLI

```bash
swarmlens-cli replay --swarm-id <id> --from <duration>   # e.g. --from 30m, --from 2h
swarmlens-cli diff --swarm-a <id> --swarm-b <id> --from <duration>
```

`diff` runs both swarms through the same detectors offline and reports alert counts and role distribution side by side — the fastest way to answer "what went wrong compared to a healthy run."

## SDKs

**Go** — see `examples/swarm-demo` for a full working reference (LLM-backed agents publishing real events).

**Python:**
```python
from client import SwarmLensClient

client = SwarmLensClient(brokers="localhost:19092")
client.emit(
    swarm_id="my-swarm",
    agent_id="my-agent",
    content="hello",
    recipient_agent_id="other-agent",
)
```

Both SDKs emit the same JSON event schema (see `docs/event-schema.md`), so Go and Python agents can coexist in the same swarm and be observed identically.

## Project status

Actively developed. Core pipeline (Kafka, Redis, Prometheus, Grafana), all three detectors, replay/diff, and Go + Python SDKs are complete and tested. Next: AutoGen/CrewAI adapters for plugging in existing frameworks without code changes.

## License

MIT
