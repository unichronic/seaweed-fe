# Technical Design Document
## Competitive Programming Recruitment Platform

**Informed by:** PRD v1.0 + pb-recruitment-server source analysis  
**Constraint:** AWS $300 student credits — zero waste  
**Last Updated:** April 2026

---

## Table of Contents

1. [Architecture Decision Summary](#1-architecture-decision-summary)
2. [System Architecture](#2-system-architecture)
3. [Tech Stack](#3-tech-stack)
4. [AWS Infrastructure (Budget-Optimized)](#4-aws-infrastructure-budget-optimized)
5. [Backend Design (Go)](#5-backend-design-go)
6. [Database Schema](#6-database-schema)
7. [API Reference](#7-api-reference)
8. [Frontend Design (Next.js)](#8-frontend-design-nextjs)
9. [Submission & Judging Pipeline](#9-submission--judging-pipeline)
10. [Authentication](#10-authentication)
11. [Leaderboard Design](#11-leaderboard-design)
12. [CI/CD Pipeline](#12-cicd-pipeline)
13. [Observability](#13-observability)
14. [Local Development](#14-local-development)
15. [PRD → TDD Decision Log](#15-prd--tdd-decision-log)

---

## 1. Architecture Decision Summary

Key decisions made reconciling the PRD, the pb-recruitment-server patterns, and the $300 AWS student credit constraint:

| PRD Spec | Decision | Reason |
|---|---|---|
| AWS RDS | **Neon (serverless Postgres)** | RDS ~$15/mo minimum. Neon free tier is generous, standard Postgres, zero ops. pb-recruitment-server already uses it in dev |
| ECS EC2 Nodegroup | **ECS Fargate** | No instance cost when idle. Only pay for task runtime. Free tier: 2M req + compute |
| DynamoDB for submissions | **Postgres (Neon)** | pb-recruitment-server proves Postgres handles this fine at our scale. Eliminates a second DB entirely |
| Redis for leaderboard | **Postgres Materialized View** | pb-recruitment-server uses this exact pattern. Free, no extra service |
| Separate Submission Server | **Single Go service** | Two ECS services doubles cost. Submission handling fits cleanly in one service with async goroutines |
| Prometheus + Loki + Grafana | **AWS CloudWatch** | Included in free tier. No self-hosted monitoring infra to manage |
| Secrets Manager | **ECS task env vars + GitHub Secrets** | Secrets Manager costs $0.40/secret/month. Not worth it at this scale |
| Consolidator as always-on service | **EventBridge scheduled rule → Lambda** | Lambda free tier is massive (1M req/month). Fires once at contest end |
| S3 for problem statements | **Inline `description` TEXT in Postgres** | pb-recruitment-server uses this. S3 adds latency + complexity for text content |
| S3 for submission code | **Keep S3** | Correct separation. Code is binary-ish content, not DB material |
| Rankings Recalculator service | **Postgres trigger + MV refresh** | Recalculating on each submission write via DB trigger. No extra service |
| Plain React | **Next.js 14 (App Router)** | Better SSE/streaming support, RSC for leaderboard, free on Cloudflare Pages |

---

## 2. System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Cloudflare Pages                          │
│              Next.js 14 (App Router + RSC)                  │
└────────────────────────┬────────────────────────────────────┘
                         │ HTTPS
                         ▼
┌─────────────────────────────────────────────────────────────┐
│            AWS Application Load Balancer (ALB)              │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│           ECS Fargate — Go API Server (Echo)                │
│                                                             │
│  ┌──────────┐  ┌──────────┐  ┌───────────┐  ┌──────────┐  │
│  │  /users  │  │ /contests│  │/submission│  │  /admin  │  │
│  └──────────┘  └──────────┘  └───────────┘  └──────────┘  │
│                        │                                    │
│              Firebase Auth (middleware)                     │
└────────────┬───────────┴──────────────────────────────────-┘
             │                         │
     ┌───────▼──────┐         ┌────────▼────────┐
     │  Neon Postgres│        │    AWS S3       │
     │  (serverless) │        │ (submissions    │
     │               │        │  bucket)        │
     │  • users      │        └─────────────────┘
     │  • contests   │
     │  • problems   │              │
     │  • submissions│        ┌─────▼──────────┐
     │  • rankings   │        │   Judge0       │
     │  • admin      │        │ (self-hosted   │
     │  • ranking_mv │        │  on EC2 t3.micro│
     └───────────────┘        │  or public API)│
                              └────────────────┘
                                      │
                         ┌────────────▼───────────┐
                         │  EventBridge + Lambda   │
                         │  (Consolidator — fires  │
                         │   at contest end_time)  │
                         └────────────────────────┘
```

**Firebase Auth** handles identity. The Go server verifies Firebase ID tokens on every protected request — no session storage needed.

**SSE** for live submission status flows directly from the Go API server to the browser, via the ALB with HTTP/1.1 keep-alive. No WebSocket needed.

---

## 3. Tech Stack

### Backend
| Layer | Choice | Why |
|---|---|---|
| Language | **Go 1.22+** | Fast, small binaries, great concurrency for SSE |
| HTTP Framework | **Echo v4** | Proven in pb-recruitment-server. Middleware composition is clean |
| Dependency Injection | **Uber FX** | Exact pattern from pb-recruitment-server. Testable, no globals |
| DB Driver | **pgx/v5** (pgxpool) | Best Go Postgres driver. Native pooling |
| Migrations | **golang-migrate** | SQL-first, no ORM magic. Prod-safe with up/down files |
| Validation | **go-playground/validator v10** | Struct tag validation on all DTOs |
| Auth | **Firebase Admin SDK** | Server-side token verification |
| AWS | **aws-sdk-go-v2** | S3 uploads for submission code |
| ID Generation | **google/uuid** (contests, problems) + **oklog/ulid** (submissions) | ULIDs are sortable by time — good for submissions |
| Logging | **go.uber.org/zap** | Structured JSON logs → CloudWatch |
| Hot reload (dev) | **air** | File-watch + auto-rebuild |

### Frontend
| Layer | Choice | Why |
|---|---|---|
| Framework | **Next.js 14** (App Router) | RSC for leaderboard, native SSE support, Cloudflare Pages compatible |
| Language | **TypeScript** | Type safety across DTOs |
| Styling | **Tailwind CSS + shadcn/ui** | Fast, accessible components. No design system to build from scratch |
| Code Editor | **Monaco Editor** (via `@monaco-editor/react`) | Industry standard. Syntax highlighting for 40+ languages |
| Markdown + LaTeX | **react-markdown + remark-math + rehype-katex** | KaTeX renders LaTeX in problem statements client-side |
| Auth | **Firebase JS SDK v10** | Handles login UI, token refresh, session persistence |
| State | **Zustand** | Lightweight. No Redux overhead for this app size |
| Package Manager | **pnpm** | Faster installs, disk-efficient |

### Infrastructure
| Component | Service | Free tier / cost |
|---|---|---|
| Frontend hosting | **Cloudflare Pages** | Free — unlimited bandwidth |
| API compute | **ECS Fargate** | ~750 vCPU-hours free/month |
| Container registry | **AWS ECR** | 500MB free storage |
| Database | **Neon (serverless Postgres)** | Free tier: 0.5GB storage, 1 compute unit |
| Object storage | **AWS S3** | 5GB free, 20K GET / 2K PUT free |
| Load balancer | **AWS ALB** | 750 hours free (Classic LB) — use ALB ~$0.008/LCU |
| Code judge | **Judge0** (self-hosted on EC2 t3.micro) | t3.micro free for 12 months |
| Post-contest job | **AWS Lambda + EventBridge** | 1M invocations free/month |
| Monitoring | **AWS CloudWatch** | 5GB logs free, 10 metrics free |
| DNS + CDN | **Cloudflare** | Free |
| CI/CD | **GitHub Actions** | 2000 min/month free |
| Secrets | **GitHub Secrets + ECS env vars** | Free |

---

## 4. AWS Infrastructure (Budget-Optimized)

### ECS Fargate Task — Go API
```
CPU:    0.25 vCPU
Memory: 512 MB
Count:  1 (scale to 2 during contest)
Cost:   ~$0.01/hour → ~$7/month running 24/7
```

Fargate scales to 0 when not needed (stop the service between contests). Only run continuously during active contest windows.

### Neon Postgres
- **Free tier** covers the full dataset at our scale (500 users, ~50K submissions per contest)
- Connection pooling via **PgBouncer** is built into Neon — no need to manage it
- **Serverless scaling** — auto-scales compute on query load, scales to zero when idle
- No RDS instance to pay for even when idle

### Judge0 on EC2 t3.micro
- **Free tier:** 750 hours/month for 12 months
- Self-hosted Judge0 via Docker on the t3.micro instance
- Accessible only within the VPC — not public-facing
- After 12 months or if heavier usage: switch to the **Judge0 hosted API** (~$0 for low volume)

### Lambda — Consolidator
```python
# Triggered by EventBridge rule at contest end_time
# Refreshes the ranking materialized view + sets final ranks
# Marks shortlisted candidates based on score threshold
```
EventBridge cron rule is created/deleted per contest by the admin API when a contest is created. Lambda free tier: 1M requests, 400,000 GB-seconds — more than enough.

### S3 — Submissions Bucket
```
Bucket: recruitment-submissions-{env}
Structure: submissions/{contestId}/{userId}/{submissionId}.{ext}
Access: Private. Pre-signed URLs for retrieval (valid 15 min)
```

### ALB
- Single ALB in front of ECS
- HTTP → HTTPS redirect
- Target group: ECS Fargate task on port 8080
- Health check: `GET /health` → 200

### VPC Layout
```
VPC: 10.0.0.0/16
  Public subnet:  ALB, NAT Gateway
  Private subnet: ECS tasks, EC2 (Judge0)
  
Security Groups:
  alb-sg:    inbound 443 from 0.0.0.0/0
  ecs-sg:    inbound 8080 from alb-sg only
  judge0-sg: inbound 2358 from ecs-sg only
```

---

## 5. Backend Design (Go)

### Project Structure (following pb-recruitment-server pattern)
```
cmd/app/main.go              # FX wiring — identical pattern to pb-recruitment-server
internal/
  server.go                  # Echo init, global middleware, /health
  boot/
    env.go                   # godotenv loader
    firebase.go              # Firebase Admin SDK init
  common/
    constants.go             # shared enums, keys
    errors.go                # typed HTTP errors
  controllers/               # HTTP layer — bind, validate, call service, respond
    contest-controller.go
    submission-controller.go
    user-controller.go
  db/
    db.go                    # pgxpool init with env-configurable pool settings
  middleware/
    firebase-auth-middleware.go    # RequireFirebaseAuth + OptionalFirebaseAuth
    admin-auth-middleware.go       # RequireAdminRole (queries admin table)
    validator-middleware.go        # ValidateRequest(dto)
  migrations/                # Sequential .up.sql / .down.sql files
  models/
    *.go                     # DB row structs
    dto/                     # Request/response structs with validator tags
  routes/
    admin-routes.go
    contest-routes.go
    submission-routes.go
    user-routes.go
  s3/
    s3.go                    # PutObject, GetPresignedURL
  services/                  # Business logic
    admin-service.go
    contest-service.go
    submission-service.go
    user-service.go
  stores/                    # DB queries via pgx
    storage.go               # Aggregator struct
    admin-store.go
    contest-store.go
    problem-store.go
    ranking-store.go
    submission-store.go
    user-store.go
```

### FX Wiring (main.go)
```go
fx.New(
    fx.Provide(
        boot.NewFirebaseAuth,
        db.NewDBConn,
        s3.NewS3Client,
        stores.NewStorage,
        services.NewContestService,
        services.NewUserService,
        services.NewSubmissionService,
        services.NewAdminService,
        controllers.NewContestController,
        controllers.NewUserController,
        controllers.NewSubmissionController,
        internal.NewEchoServer,
    ),
    fx.Invoke(routes.AddUserRoutes),
    fx.Invoke(routes.AddContestRoutes),
    fx.Invoke(routes.AddSubmissionRoutes),
    fx.Invoke(routes.AddAdminRoutes),
    fx.Invoke(internal.StartEchoServer),
).Run()
```

### Environment Variables
```env
STAGE=dev|prod
FIREBASE_SERVICE_ACCOUNT_PATH=./firebase-service-account.json
DB_ADDR=postgres://user:pass@ep-xxx.neon.tech/neondb?sslmode=require
AWS_REGION=ap-south-1
AWS_ACCESS_KEY_ID=...           # local dev only; ECS uses task role
AWS_SECRET_ACCESS_KEY=...       # local dev only
S3_SUBMISSIONS_BUCKET=recruitment-submissions-prod
JUDGE0_BASE_URL=http://10.0.1.x:2358   # internal VPC endpoint
PORT=8080

# DB Pool (optional)
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=5m
```

---

## 6. Database Schema

All data lives in **Neon Postgres**. No DynamoDB, no Redis. Timestamps stored as `BIGINT` (Unix milliseconds) — consistent with pb-recruitment-server and JS-friendly.

### `users`
```sql
CREATE TABLE users (
    id             TEXT PRIMARY KEY,        -- Firebase UID
    name           TEXT NOT NULL,
    email          TEXT NOT NULL UNIQUE,
    usn            TEXT NOT NULL UNIQUE,    -- University serial number
    mobile_number  TEXT,
    joining_year   INT NOT NULL,
    department     TEXT NOT NULL
);
```

### `admin`
```sql
CREATE TABLE admin (
    user_id    TEXT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```
Admin is a role, not a separate user type. Any user can be promoted by inserting into this table.

### `contests`
```sql
CREATE TABLE contests (
    id                      TEXT PRIMARY KEY,   -- UUID
    name                    TEXT NOT NULL,
    description             TEXT,
    eligible_to             TEXT,               -- e.g. '2,3' for 2nd and 3rd year only
    registration_start_time BIGINT NOT NULL,
    registration_end_time   BIGINT NOT NULL,
    start_time              BIGINT NOT NULL,
    end_time                BIGINT NOT NULL
);
```

### `contest_registrations`
```sql
CREATE TABLE contest_registrations (
    contest_id    TEXT NOT NULL REFERENCES contests(id) ON DELETE CASCADE,
    user_id       TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    registered_at BIGINT NOT NULL,
    PRIMARY KEY (contest_id, user_id)
);
```

### `problems`
```sql
CREATE TABLE problems (
    id          TEXT PRIMARY KEY,   -- UUID
    contest_id  TEXT NOT NULL REFERENCES contests(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    description TEXT,               -- Markdown + LaTeX source stored inline
    score       INT NOT NULL
);
```

Problem statements are stored as Markdown text directly in Postgres. The frontend renders them with KaTeX. No S3 for text content.

### `submissions`
```sql
CREATE TYPE submission_status AS ENUM (
    'pending', 'failed_to_process', 'accepted',
    'tle', 'mle', 'rte', 'wrong_answer'
);

CREATE TABLE submissions (
    id          TEXT PRIMARY KEY,    -- ULID (sortable, time-ordered)
    user_id     TEXT NOT NULL REFERENCES users(id),
    contest_id  TEXT NOT NULL REFERENCES contests(id) ON DELETE CASCADE,
    problem_id  TEXT NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    language    TEXT NOT NULL,
    s3_key      TEXT NOT NULL,       -- s3://bucket/submissions/{contestId}/{userId}/{submissionId}.{ext}
    status      submission_status NOT NULL DEFAULT 'pending',
    created_at  BIGINT NOT NULL
);
```

Raw code is stored in S3. The `s3_key` references it. The `code` column never lives in Postgres (follows pb-recruitment-server migration 017 lesson).

### `test_case_results`
```sql
CREATE TYPE test_case_status AS ENUM ('pass', 'wrong_answer', 'tle', 'mle', 'rte');

CREATE TABLE test_case_results (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    submission_id TEXT NOT NULL REFERENCES submissions(id) ON DELETE CASCADE,
    test_case_id  TEXT NOT NULL,
    status        test_case_status NOT NULL,
    runtime_ms    BIGINT NOT NULL,
    memory_kb     BIGINT NOT NULL,
    created_at    BIGINT NOT NULL
);
```

One row per test case per submission. Powers the detailed submission view.

### `rankings`
```sql
CREATE TABLE rankings (
    contest_id    TEXT NOT NULL REFERENCES contests(id) ON DELETE CASCADE,
    user_id       TEXT NOT NULL REFERENCES users(id),
    score         INT NOT NULL DEFAULT 0,
    hidden        BOOLEAN NOT NULL DEFAULT FALSE,       -- admin can hide from public board
    disqualified  BOOLEAN NOT NULL DEFAULT FALSE,       -- DQ flag
    shortlisted   BOOLEAN NOT NULL DEFAULT FALSE,       -- recruitment shortlist
    PRIMARY KEY (contest_id, user_id)
);
```

The `shortlisted`, `hidden`, and `disqualified` flags directly support the recruitment workflow — admins can act on candidates from the leaderboard without a separate tool.

### `ranking_mv` (Materialized View)
```sql
CREATE MATERIALIZED VIEW ranking_mv AS
SELECT
    r.contest_id,
    r.user_id,
    u.name,
    u.usn,
    u.department,
    r.score,
    r.hidden,
    r.disqualified,
    r.shortlisted,
    RANK() OVER (PARTITION BY r.contest_id ORDER BY r.score DESC) AS rank
FROM rankings r
JOIN users u ON u.id = r.user_id
WHERE r.disqualified = FALSE
ORDER BY r.contest_id, rank;

CREATE UNIQUE INDEX ON ranking_mv (contest_id, user_id);
```

The MV is refreshed concurrently on every submission verdict update (`REFRESH MATERIALIZED VIEW CONCURRENTLY ranking_mv`). This keeps leaderboard reads instant without Redis or a separate recalculator service.

---

## 7. API Reference

### Auth conventions
- `∅` — No auth
- `🔒` — Firebase ID token required (`Authorization: Bearer <token>`)
- `👑` — Firebase token + admin table membership required

### User Routes
| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/auth/signup` | ∅ | Firebase signup — creates Firebase account |
| POST | `/users/create` | 🔒 | Create DB user record post-signup |
| GET | `/users/profile` | 🔒 | Get current user's profile |
| POST | `/users/profile` | 🔒 | Update profile (name, mobile, etc.) |

### Contest Routes
| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/contests/list` | Optional | List all contests. Authed users see registration status |
| GET | `/contests/:id` | Optional | Contest details. Authed users see if registered |
| POST | `/contests/:id/registration` | 🔒 | Register or unregister (`action: register\|unregister` in body) |
| GET | `/contests/:id/problems` | 🔒 | List problems (no descriptions — just names + scores) |
| GET | `/contests/:id/problems/:problemId` | 🔒 | Full problem with Markdown description |

### Submission Routes
| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/submission/submit` | 🔒 | Submit code. Returns `submissionId` immediately |
| GET | `/submission/:id/status` | 🔒 | Poll or SSE-compatible status endpoint |
| GET | `/submission/:id/details` | 🔒 | Per-test-case breakdown |
| GET | `/submission/list` | 🔒 | List user's submissions. Required query param: `problem_id` |

### Leaderboard Route
| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/contests/:id/leaderboard` | Optional | Read from `ranking_mv`. Hides `hidden=true` rows for non-admins |

### Admin Routes (`/admin/*` — requires 👑)
| Method | Path | Description |
|---|---|---|
| GET | `/admin/` | Role check — returns 200 if admin |
| GET | `/admin/contests/list` | List all contests (including unpublished) |
| POST | `/admin/contest` | Create contest |
| PUT | `/admin/contest/:id` | Update contest |
| DELETE | `/admin/contest/:id` | Delete contest |
| POST | `/admin/:contestId/problem` | Add problem to contest |
| PUT | `/admin/:contestId/:problemId` | Update problem |
| DELETE | `/admin/:contestId/:problemId` | Delete problem |
| PUT | `/admin/:contestId/leaderboard/:userId` | Set `hidden`, `disqualified`, `shortlisted` on a ranking entry |
| GET | `/admin/contests/:contestId/registrations` | List all registered users for a contest |

### System
| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/health` | ∅ | ALB health check — returns 200 |

---

## 8. Frontend Design (Next.js)

### Hosting
**Cloudflare Pages** — free, global CDN, zero config for Next.js. Deploy via `@cloudflare/next-on-pages`.

### App Router Structure
```
app/
  page.tsx                        # Landing page
  auth/
    login/page.tsx
    signup/page.tsx
  contests/
    page.tsx                      # Contest list (RSC — server fetch)
    [contestId]/
      page.tsx                    # Contest landing page
      problems/
        page.tsx                  # Problem list (requires auth)
        [problemId]/page.tsx      # Problem + Monaco editor
      leaderboard/page.tsx        # Live/final leaderboard
  admin/
    layout.tsx                    # Admin auth guard
    contests/page.tsx
    contests/[contestId]/page.tsx
components/
  editor/MonacoEditor.tsx         # Client component
  leaderboard/LeaderboardTable.tsx
  problems/MarkdownRenderer.tsx   # react-markdown + KaTeX
  ui/                             # shadcn/ui components
lib/
  firebase.ts                     # Firebase client init
  api.ts                          # Typed fetch wrappers
  auth.ts                         # Auth hooks + context
```

### Key Frontend Patterns

**Leaderboard — SSE streaming**
The leaderboard page opens an SSE connection to `GET /contests/:id/leaderboard` with `Accept: text/event-stream`. The Go server streams ranking updates as server-sent events, refreshed after each MV update. On contest end, the stream closes and the page displays the final frozen rankings.

**Problem statement rendering**
```tsx
// components/problems/MarkdownRenderer.tsx
import ReactMarkdown from 'react-markdown'
import remarkMath from 'remark-math'
import rehypeKatex from 'rehype-katex'
import 'katex/dist/katex.min.css'

export function MarkdownRenderer({ content }: { content: string }) {
  return (
    <ReactMarkdown
      remarkPlugins={[remarkMath]}
      rehypePlugins={[rehypeKatex]}
    >
      {content}
    </ReactMarkdown>
  )
}
```

**Monaco Editor — language selection + submission**
```tsx
// Client component — not SSR
'use client'
import Editor from '@monaco-editor/react'
```
Language selector dropdown → changes Monaco language mode. On submit: POST to `/submission/submit` with `{ contestId, problemId, language, code }`. Then poll `GET /submission/:id/status` every 2 seconds until verdict is non-pending.

**Two-step auth signup**
```
1. Firebase createUserWithEmailAndPassword()
2. POST /users/create  { name, usn, department, joining_year }
   with Authorization: Bearer <firebase-id-token>
```

---

## 9. Submission & Judging Pipeline

```
Browser
  │  POST /submission/submit { problemId, language, code }
  ▼
Go API (submission-service)
  │  1. Generate ULID for submissionId
  │  2. Upload code to S3: submissions/{contestId}/{userId}/{submissionId}.{lang}
  │  3. INSERT into submissions (status='pending', s3_key=...)
  │  4. Return { submissionId } immediately to browser
  │  5. Fire goroutine → call Judge0
  ▼
Judge0 (EC2 t3.micro, VPC-internal)
  │  POST /submissions { source_code, language_id, stdin: test_cases }
  │  Poll until verdict
  ▼
Go API (goroutine callback)
  │  1. INSERT test_case_results for each test case
  │  2. UPDATE submissions SET status = <verdict>
  │  3. UPDATE rankings SET score = score + problem.score (if accepted + first AC)
  │  4. REFRESH MATERIALIZED VIEW CONCURRENTLY ranking_mv
  ▼
Browser (polling GET /submission/:id/status every 2s)
  │  Returns verdict when status != 'pending'
```

**Why polling instead of SSE for submission status:**
SSE for submission status requires the Go server to hold an open connection per in-flight submission — at 500 concurrent users, that's 500 goroutines blocked. Polling every 2 seconds is simpler, cheaper, and perfectly acceptable UX for a 2–10 second judge latency.

**SSE is used only for the leaderboard** — one SSE stream per leaderboard page view, not per submission.

**Judge0 integration notes:**
- Language IDs follow Judge0's numbering (e.g., 71 = Python 3, 62 = Java, 54 = C++17)
- Each submission is sent with all test cases as `stdin` batches, or individual calls per test case depending on Judge0 config
- Self-hosted Judge0 on EC2 t3.micro is sufficient for ~50 concurrent judge evaluations
- Judge0 is only accessible within the VPC — not exposed publicly

---

## 10. Authentication

### Flow
```
Browser                Firebase             Go API
   │                      │                    │
   │── signInWithEmail ──►│                    │
   │◄── idToken ──────────│                    │
   │                                           │
   │── POST /users/create ──────────────────►  │
   │   Authorization: Bearer <idToken>          │
   │                    VerifyIDToken(idToken) ─►Firebase
   │                    ◄─ {uid, email}         │
   │◄── 200 OK ─────────────────────────────── │
```

### Middleware stack on protected routes
```go
// RequireFirebaseAuth — hard block
e.GET("/users/profile", userController.GetUserProfile,
    middleware.RequireFirebaseAuth(authClient),
)

// OptionalFirebaseAuth — enriches context, doesn't block
e.GET("/contests/list", contestController.ListContests,
    middleware.OptionalFirebaseAuth(authClient),
)

// RequireAdminRole — stacked on top of RequireFirebaseAuth
adminGroup := e.Group("/admin")
adminGroup.Use(middleware.RequireFirebaseAuth(authClient))
adminGroup.Use(middleware.RequireAdminRole(userService, adminService))
```

### Admin promotion
Admins are regular users in the `admin` table. Promotion is done directly in the DB — no Firebase custom claims needed. `RequireAdminRole` middleware queries the `admin` table on every admin request (cached in-process with a 60-second TTL to avoid hammering Neon).

---

## 11. Leaderboard Design

### Live leaderboard (during contest)
- `GET /contests/:id/leaderboard` returns the current `ranking_mv` snapshot
- Frontend polls every 10 seconds (simple interval) — no SSE needed at this scale
- MV is refreshed concurrently after every accepted submission (sub-100ms overhead on Neon)
- Non-admins never see rows where `hidden = true`

### Final leaderboard (post-contest)
- EventBridge rule fires Lambda at `end_time` of each contest
- Lambda calls `REFRESH MATERIALIZED VIEW CONCURRENTLY ranking_mv` one final time
- Lambda sets a `finalized = true` flag on the contest row
- Frontend detects `finalized` and stops polling — displays static final rankings
- Admins can mark candidates as `shortlisted` from the leaderboard UI, updating the `rankings` table

### Ranking score logic
- Score = sum of `problems.score` for first accepted submission per problem per user
- Penalty time is not tracked in v1.0 (can be added as a `penalty_minutes` column later)
- `RANK()` window function in the MV handles ties correctly

---

## 12. CI/CD Pipeline

### GitHub Actions — Two workflows

**`build.yaml` — CI on every PR**
```yaml
- go build ./...
- go vet ./...
- go test ./...
- pnpm install && pnpm build   # frontend build check
```

**`deploy.yml` — CD on merge to main**
```yaml
- Configure AWS credentials (OIDC → github-actions-role IAM role)
- Login to ECR
- docker build + docker push to ECR
- aws ecs update-service --force-new-deployment
```

### Dockerfile (multi-stage, non-root)
```dockerfile
FROM golang:1.22-alpine AS builder
# CGO_ENABLED=0 static binary
RUN go build -ldflags="-s -w" -o /app ./cmd/app/main.go

FROM alpine:latest
RUN adduser -D -u 1000 appuser
COPY --from=builder /app /app/app
COPY internal/migrations /app/internal/migrations
USER appuser
EXPOSE 8080
HEALTHCHECK CMD wget --spider http://localhost:8080/health
CMD ["/app/app"]
```

### ECS Task Definition (key fields)
```json
{
  "cpu": "256",
  "memory": "512",
  "image": "<accountId>.dkr.ecr.ap-south-1.amazonaws.com/recruitment-server:latest",
  "environment": [
    { "name": "STAGE", "value": "prod" },
    { "name": "DB_ADDR", "value": "<neon-connection-string>" },
    { "name": "S3_SUBMISSIONS_BUCKET", "value": "recruitment-submissions-prod" },
    { "name": "JUDGE0_BASE_URL", "value": "http://10.0.1.x:2358" },
    { "name": "AWS_REGION", "value": "ap-south-1" }
  ]
}
```

Sensitive values (DB password, Firebase service account JSON) are stored as **ECS task secrets** pointing to **AWS Systems Manager Parameter Store** (free tier, unlike Secrets Manager).

---

## 13. Observability

**CloudWatch Logs** — all `zap` structured JSON logs are shipped to CloudWatch via the `awslogs` ECS log driver. Zero config.

**CloudWatch Metrics** — ALB emits request count, latency, 5xx rate automatically. Free.

**CloudWatch Alarms** — set alarms on:
- ALB 5xx rate > 1% → SNS notification
- ECS task count drops to 0 → SNS notification
- Judge0 EC2 CPU > 80% sustained → SNS notification

**No self-hosted Prometheus/Loki/Grafana** — would cost ~$30/month for the EC2 instances to run them. CloudWatch covers 90% of the observability needs for free.

---

## 14. Local Development

```bash
# Prerequisites: Go 1.22+, Node 20+, pnpm, Docker, air

# Backend
cp .env.example .env          # fill in Neon DB URL, Firebase JSON path
air                           # hot reload on localhost:8080

# Run migrations
make migrate-up

# Frontend
cd frontend
pnpm install
pnpm dev                      # localhost:3000

# Judge0 locally (optional — use public Judge0 API for dev)
docker run -d -p 2358:2358 judge0/judge0
```

**Neon for local dev** — use the same Neon database in dev (with a separate `dev` branch on Neon — Neon supports database branching like Git). No need to run local Postgres.

---

## 15. PRD → TDD Decision Log

Explicit record of where the TDD diverges from the PRD and why:

| PRD Spec | TDD Decision | Rationale |
|---|---|---|
| DynamoDB for submissions | **Postgres** | pb-recruitment-server proves Postgres handles this. Eliminates entire AWS service. Saves ~$5/month |
| Redis for leaderboard cache | **Postgres Materialized View** | Same pattern from pb-recruitment-server. `REFRESH CONCURRENTLY` is non-blocking. Saves ~$15/month (ElastiCache minimum) |
| Separate Submission Server | **Goroutine within API server** | Two ECS tasks = double Fargate cost. Goroutine-based async is idiomatic Go and sufficient at 500 users |
| Consolidator as always-on service | **Lambda + EventBridge** | Lambda is effectively free. No need to run a service 24/7 just to fire once at contest end |
| RDS for SQL | **Neon (serverless Postgres)** | Neon is standard Postgres. Free tier is sufficient. No idle instance cost |
| Prometheus + Loki + Grafana | **CloudWatch** | Saves ~$30/month in EC2. Adequate observability for this scale |
| Secrets Manager | **SSM Parameter Store** | Free tier vs $0.40/secret/month. Functionally identical for our use case |
| S3 for problem statements | **Postgres TEXT column** | Text content doesn't benefit from S3. Reduces latency. pb-recruitment-server proves this works |
| Plain React | **Next.js 14** | RSC improves initial load. Better SSE support. Still free on Cloudflare Pages |
| Architecture diagram SSE for submission status | **Polling (2s interval)** | Simpler. SSE for submissions would hold 500 goroutines open. Polling is fine for 2–10s judge latency |
| `rank` field in rankings schema | **Computed via `RANK()` in MV** | No need to store rank — it's computed correctly in the view and updates automatically |
| No admin system specified | **`admin` table + `RequireAdminRole`** | Adopted from pb-recruitment-server. Necessary for contest management and recruitment shortlisting |
| Rankings schema: contestId + userId + score + rank | **+ hidden + disqualified + shortlisted** | Adopted from pb-recruitment-server. These flags make the leaderboard directly useful for recruitment without a separate tool |

---

*This TDD is the implementation source of truth. The PRD defines what to build; this defines how.*
