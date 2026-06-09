# Autonomous Business Monitor (ABM)

A production Go backend that wakes up every morning at 6:45am, pulls data from your business tools, detects anomalies, and delivers an AI-written briefing to your inbox before you start your day.

No dashboards to check. No manual reporting. One intelligent email, every morning.

---

## What It Does

Every morning at 6:45am, ABM:

1. Fetches data from Stripe, Gmail, and Calendly simultaneously
2. Saves raw responses to a PostgreSQL database
3. Computes daily metrics (revenue, failed payments)
4. Compares today's numbers against a 7-day rolling average to detect anomalies
5. Generates a structured briefing using Google Gemini
6. Delivers the briefing via n8n (email + Slack for critical alerts)

---

## Sample Briefing

```
1. NEEDS ATTENTION
* High Email Volume: Gmail shows 201 pending messages. 
  Action: Triage inbox before your afternoon meeting.

2. LOOKING GOOD
* Stripe revenue steady at $7,599.70 — 0% delta from 7-day average.
* Zero failed payments recorded.

3. TODAY AT A GLANCE
Revenue is stable. Primary focus: clear the email backlog before 
your 14:00 meeting with Ayomide Alabi.
```

---

## Tech Stack

| Area | Tool |
|------|------|
| Language | Go 1.21+ |
| Database | PostgreSQL 16 |
| Local DB | Docker + docker-compose |
| Migrations | Goose |
| DB Queries | sqlc |
| LLM | Google Gemini |
| Delivery | n8n |
| Deployment | Railway |

---

## Prerequisites

- Go 1.21+
- Docker Desktop
- A Stripe account (test or live)
- A Google Cloud project with Gmail API enabled
- A Calendly account with a personal API token
- A Google AI Studio account (for Gemini API key)
- An n8n instance (cloud or self-hosted)
- Goose CLI: `go install github.com/pressly/goose/v3/cmd/goose@latest`
- sqlc CLI: `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`

---

## Setup

### 1. Clone the repository

```bash
git clone https://github.com/yourusername/abm.git
cd abm
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Create your `.env` file

```bash
cp .env.example .env
```

Fill in all values — see Environment Variables below.

### 4. Start PostgreSQL

```bash
docker compose up -d
```

### 5. Run database migrations

```bash
goose -dir migrations postgres "$DATABASE_URL" up
```

### 6. Generate sqlc code

```bash
sqlc generate
```

### 7. Run the server

```bash
go run ./cmd/server
```

To trigger an immediate run without waiting for 6:45am:

```bash
RUN_NOW=true go run ./cmd/server
```

---

## Environment Variables

| Variable | Description |
|----------|-------------|
| `DATABASE_URL` | PostgreSQL connection string |
| `STRIPE_SECRET_KEY` | Stripe secret key (sk_test_... or sk_live_...) |
| `GMAIL_CLIENT_ID` | Google OAuth2 client ID |
| `GMAIL_CLIENT_SECRET` | Google OAuth2 client secret |
| `GMAIL_REFRESH_TOKEN` | Gmail OAuth2 refresh token |
| `CALENDLY_API_KEY` | Calendly personal access token |
| `CALENDLY_USER_URI` | Your Calendly user URI (https://api.calendly.com/users/...) |
| `LLM_API_KEY` | Google AI Studio API key |
| `LLM_MODEL` | Gemini model (e.g. gemini-2.5-flash) |
| `N8N_WEBHOOK_URL` | n8n webhook URL for digest delivery |
| `CRON_SCHEDULE` | Cron expression (default: `45 6 * * *`) |
| `PORT` | API server port (default: 8080) |
| `RUN_NOW` | Set to `true` to trigger immediately on startup |

---

## Project Structure

```
abm/
├── cmd/server/main.go          — entry point, wires everything together
├── config/config.go            — loads all env vars into one typed struct
├── internal/
│   ├── api/handler.go          — REST API handlers
│   ├── anomaly/engine.go       — rolling average + anomaly detection
│   ├── fetcher/
│   │   ├── stripe.go           — Stripe API client
│   │   ├── gmail.go            — Gmail API client (OAuth2)
│   │   └── calendly.go         — Calendly API client
│   ├── llm/
│   │   ├── client.go           — Gemini LLM client
│   │   └── prompt.go           — template-based prompt builder
│   ├── notify/webhook.go       — n8n webhook delivery
│   ├── scheduler/scheduler.go  — cron job + collection pipeline
│   └── storage/
│       ├── storage.go          — database connection pool
│       ├── queries/            — SQL query files (you edit these)
│       └── db/                 — sqlc generated Go code (never edit)
├── migrations/                 — Goose migration files
├── prompts/morning_digest.txt  — LLM prompt template (edit to tune output)
├── sqlc.yaml                   — sqlc configuration
├── docker-compose.yml          — local PostgreSQL container
└── .env                        — local secrets (never commit this)
```

---

## REST API

### `GET /digest/history`

Returns the last 30 AI-generated briefings.

```bash
curl http://localhost:8080/digest/history
```

### `GET /metrics/trend?source=stripe`

Returns the last 30 days of daily metrics for a given source.

```bash
curl "http://localhost:8080/metrics/trend?source=stripe"
```

---

## Tuning the Briefing

The LLM prompt lives in `prompts/morning_digest.txt`. Edit it to change the tone, format, or focus of the briefing. No code changes needed — just edit the file and restart the server.

---

## Deployment (Railway)

1. Push your repository to GitHub
2. Create a new project on railway.app
3. Connect your GitHub repository
4. Add all environment variables from the table above
5. Railway will deploy automatically on every push

---

## Gmail OAuth2 Setup

1. Create a Google Cloud project at console.cloud.google.com
2. Enable the Gmail API
3. Create OAuth2 credentials (Web application type)
4. Add `https://developers.google.com/oauthplayground` as an authorized redirect URI
5. Add your email as a test user under OAuth consent screen
6. Use the OAuth2 Playground to exchange credentials for a refresh token
7. Add the refresh token to your `.env`

---

*Built with Go, PostgreSQL, and Google Gemini.*