# Product Requirements Document
## Competitive Programming Recruitment Platform

**Version:** 1.0  
**Status:** Draft  
**Last Updated:** April 2026  

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [Problem Statement](#2-problem-statement)
3. [Goals & Success Metrics](#3-goals--success-metrics)
4. [Stakeholders](#4-stakeholders)
5. [System Architecture Overview](#5-system-architecture-overview)
6. [Functional Requirements](#6-functional-requirements)
7. [Data Models](#7-data-models)
8. [Non-Functional Requirements](#8-non-functional-requirements)
9. [Infrastructure & DevOps](#9-infrastructure--devops)
10. [Frontend Requirements](#10-frontend-requirements)
11. [Out of Scope](#11-out-of-scope)
12. [Risks & Mitigations](#12-risks--mitigations)
13. [Open Questions](#13-open-questions)

---

## 1. Executive Summary

This document defines the product requirements for a **full-fledged competitive programming recruitment platform** — a system that enables organizations to host timed coding contests, evaluate submissions in real-time, and surface ranked candidates for recruitment purposes.

The platform is designed to support the full contest lifecycle: from contest creation and registration, through live code submission and automated judging, to final leaderboard generation and recruitment decision-making.

---

## 2. Problem Statement

Recruiting organizations and universities need a structured, scalable way to assess candidates through competitive programming. Existing tools either lack recruitment-specific workflows, don't support real-time feedback, or cannot maintain accurate live leaderboards at scale.

This platform addresses:

- The lack of a unified tool for running structured recruitment-oriented coding contests
- The need for real-time submission feedback (accepted, TLE, MLE, wrong answer, etc.)
- The challenge of fairly ranking candidates and preserving rankings post-contest
- The absence of LaTeX-rendered problem statements for mathematical/technical problems

---

## 3. Goals & Success Metrics

### Goals

- Enable contest administrators to create, configure, and publish contests with minimal friction
- Provide participants with a seamless coding submission experience with live status updates
- Evaluate submissions accurately and in near-real-time using a sandboxed judge
- Generate and persist accurate leaderboards during and after contests
- Support organizational recruitment workflows via ranked candidate outputs

### Success Metrics

| Metric | Target |
|--------|--------|
| Submission evaluation latency (P95) | < 5 seconds |
| System uptime during contest | 99.9% |
| Concurrent users supported | 500+ per contest |
| False judge verdicts | < 0.1% |
| Leaderboard update lag | < 10 seconds |

---

## 4. Stakeholders

| Role | Responsibility |
|------|---------------|
| Contest Administrators | Create contests, add problems, manage registration windows |
| Participants / Candidates | Register, solve problems, view standings |
| Recruiters | View final rankings and candidate profiles |
| Platform Engineers | Maintain infrastructure, judging pipeline, and deployments |

---

## 5. System Architecture Overview

The platform is composed of the following major subsystems:

### 5.1 Frontend
A **React** application hosted on **Cloudflare Pages/Workers**, responsible for rendering all user-facing pages. It communicates with the Website Backend over HTTP/WebSocket.

### 5.2 Website Backend (Go)
The primary API server handling:
- Authentication flows (via Firebase Auth)
- Contest and problem management
- Submission creation and status retrieval
- Live submission status updates via **Server-Sent Events (SSE)**
- Caching of hot data via **Redis**

### 5.3 Submission Server (Go)
Decoupled from the Website Backend, the Submission Server handles:
- Receiving and queuing submissions
- Delegating to **Judge0** for sandboxed code execution
- Writing raw submissions to **S3**
- Updating submission records in the database

### 5.4 Judge0
An open-source code execution sandbox used to evaluate submissions against test cases. Returns verdicts such as: `accepted`, `wrong_answer`, `time_limit_exceeded`, `memory_limit_exceeded`, `runtime_error`.

### 5.5 Rankings Recalculator
A background service that recalculates live rankings during a contest and writes results to the **Live Rankings** store.

### 5.6 Consolidator
A post-contest service that runs when a contest ends, computes final scores/rankings, and persists them to the **Saved Rankings** database. This ensures rankings are immutable after contest close.

### 5.7 Databases
- **SQL Databases (AWS RDS):** Stores Users, Contests, Problems, Submissions metadata, and Saved Rankings
- **DynamoDB:** Used for the Submissions table (NoSQL), optimized for high write throughput during contests
- **Redis Cache:** Caches rankings and frequently-read data
- **S3 Buckets:** Stores raw submission code and problem statement files (`.md` format)

---

## 6. Functional Requirements

### 6.1 Authentication

- **FR-AUTH-01:** Users must be able to register and log in using Firebase Authentication
- **FR-AUTH-02:** Sessions must be securely managed; unauthenticated users cannot access contest content
- **FR-AUTH-03:** Authentication state must be reflected in the React frontend without page reloads

### 6.2 Contest Management

- **FR-CONT-01:** Admins can create a contest with the following attributes:
  - Contest name
  - Description
  - Registration start and end time
  - Contest start and end time
  - Registration status (open/closed/invite-only)
- **FR-CONT-02:** Admins can add one or more problems to a contest
- **FR-CONT-03:** Participants can view a Contest Landing Page with all contest metadata before registering
- **FR-CONT-04:** Participants can register for a contest within the registration window
- **FR-CONT-05:** A Contest Problem List page must be accessible to registered participants during the contest window

### 6.3 Problem Management

- **FR-PROB-01:** Each problem belongs to one contest (FK relationship)
- **FR-PROB-02:** All problems are of type `coding`
- **FR-PROB-03:** Problem statements must be written in Markdown and support **LaTeX rendering** for mathematical notation
- **FR-PROB-04:** Problem statement files are stored as `.md` files in S3
- **FR-PROB-05:** Each coding problem has a defined score (INT) and answer validation via Judge0

### 6.4 Submission Handling

- **FR-SUB-01:** Participants can submit code via the Monaco editor
- **FR-SUB-02:** Raw code is persisted to S3 immediately on submission
- **FR-SUB-03:** The Submission Server sends the code to Judge0 for evaluation
- **FR-SUB-04:** Judge0 returns a verdict from: `pending | failed_to_process | accepted | tle | mle | rte | wrong_answer`
- **FR-SUB-05:** Submission details (verdict, timestamps) are written to DynamoDB
- **FR-SUB-06:** The Website Backend updates the participant's correct/incorrect attempt count in the SQL database
- **FR-SUB-07:** Participants receive live submission status updates via **Server-Sent Events (SSE)** without polling
- **FR-SUB-08:** Submission history is viewable from the Contest Problem Page

### 6.5 Rankings & Leaderboard

- **FR-RANK-01:** A live leaderboard is computed continuously during the contest by the Rankings Recalculator service
- **FR-RANK-02:** The Contest Leaderboard page displays live standings to all participants
- **FR-RANK-03:** At contest end, the Consolidator computes and freezes final rankings into the Saved Rankings table
- **FR-RANK-04:** Rankings are queryable post-contest for recruitment review
- **FR-RANK-05:** The Ranking schema stores `contestId`, `userId`, `score`, and `rank`

---

## 7. Data Models

### 7.1 Users

| Field | Type | Notes |
|-------|------|-------|
| userId | TEXT | Firebase UID, Primary Key |
| name | TEXT | Display name |
| usn | TEXT | Unique, nullable — university serial number |
| dept | TEXT | Department |
| join_year | INT | Year of joining |

### 7.2 Contests

| Field | Type | Notes |
|-------|------|-------|
| contestId | TEXT (UUID) | Primary Key |
| name | TEXT | Contest title |
| registration_start_time | TIMESTAMP | |
| registration_end_time | TIMESTAMP | |
| start_time | TIMESTAMP | |
| end_time | TIMESTAMP | |

### 7.3 Problems

| Field | Type | Notes |
|-------|------|-------|
| problemId | TEXT (UUID) | Primary Key |
| contestId | TEXT | Foreign Key → Contests |
| name | TEXT | |
| score | INT | Points awarded for a correct solution |
| type | ENUM | `coding` (only) |

### 7.4 Submissions (DynamoDB)

| Field | Type | Notes |
|-------|------|-------|
| submissionId | String (ULID) | Primary Key |
| userId | String | |
| contestId | String | |
| problemId | String | |
| codeId | String | Raw code submission ID (S3 reference) |
| status | String | `pending / failed_to_process / accepted / tle / mle / rte / wrong_answer` |

### 7.5 Rankings

| Field | Type | Notes |
|-------|------|-------|
| contestId | TEXT | Composite Primary Key |
| userId | TEXT | Composite Primary Key |
| score | INT | Computed score |
| rank | INT | For saved/final rankings |

---

## 8. Non-Functional Requirements

### 8.1 Performance
- **NFR-PERF-01:** The platform must support at least 500 concurrent participants per contest
- **NFR-PERF-02:** Submission status updates must be pushed within 10 seconds of a verdict being available
- **NFR-PERF-03:** Leaderboard reads must be served from Redis cache; DB should not be hit on each leaderboard refresh

### 8.2 Scalability
- **NFR-SCALE-01:** The system must auto-scale using ECS with a Nodegroup and ALB load balancing
- **NFR-SCALE-02:** DynamoDB is used for Submissions due to its ability to handle bursty write loads during active contests
- **NFR-SCALE-03:** S3 is used for raw submission storage, offloading binary/blob storage from relational DBs

### 8.3 Reliability
- **NFR-REL-01:** Contest rankings must be persisted atomically at contest end by the Consolidator; partial writes are not acceptable
- **NFR-REL-02:** Raw submissions must be written to S3 before being sent to Judge0, ensuring submissions are never lost even if judging fails
- **NFR-REL-03:** Redis acts as a cache, not a source of truth; all critical data is backed by RDS or DynamoDB

### 8.4 Security
- **NFR-SEC-01:** All authentication is handled via Firebase Auth; no custom credential storage
- **NFR-SEC-02:** Code execution is sandboxed via Judge0 to prevent malicious code from affecting the host
- **NFR-SEC-03:** S3 buckets for problem statements and submissions must have restricted access (not public)

### 8.5 Observability
- **NFR-OBS-01:** All ECS services must emit metrics to **Prometheus**
- **NFR-OBS-02:** Logs must be aggregated in **Loki**
- **NFR-OBS-03:** Dashboards in **Grafana** must provide real-time visibility into submission throughput, judge latency, and error rates

---

## 9. Infrastructure & DevOps

### 9.1 Hosting
- Application services run on **AWS ECS** (Elastic Container Service)
- Container images are stored in **AWS ECR** (Elastic Container Registry)
- Databases are hosted on **AWS RDS** (PostgreSQL or MySQL)
- ECS cluster uses a **Nodegroup** and is fronted by an **ALB** (Application Load Balancer)

### 9.2 Deployment Pipeline
- Code is containerized (Dockerized) and pushed to ECR
- **GitHub Actions** triggers force deployments to ECS on merge to the main branch
- IAM role: `github-actions-role` is used by GitHub Actions for deployment
- Infrastructure provisioning is managed under the `ph-infra-recruitment` context (likely Terraform/CDK)

### 9.3 Monitoring Stack
- **Loki** — Log aggregation
- **Prometheus** — Metrics collection
- **Grafana** — Unified dashboard for logs and metrics

### 9.4 Frontend Hosting
- React app is deployed to **Cloudflare Pages/Workers** for global edge delivery and low latency

---

## 10. Frontend Requirements

### 10.1 Pages

| Page | Description |
|------|-------------|
| Login / Register | Firebase Auth UI; entry point for all users |
| Landing Page | Marketing/info page for the platform |
| Contest List | Browse all available and upcoming contests |
| Contest Landing Page | Per-contest info: name, description, start/end time, registration status |
| Contest Problem List | List of problems in a contest (accessible to registered users during contest window) |
| Contest Problem Page | Problem statement with LaTeX support; Monaco editor for code submission; live verdict display |
| Contest Leaderboard | Live rankings during contest; final rankings post-contest |

### 10.2 Key UI/UX Requirements

- **FR-UI-01:** The Monaco editor must be used for code submission with syntax highlighting and language-aware formatting
- **FR-UI-02:** Problem statements must render Markdown with full LaTeX support (e.g., via KaTeX or MathJax)
- **FR-UI-03:** Submission status must update in real-time via SSE without user-initiated refreshes
- **FR-UI-04:** The leaderboard must auto-refresh with live data; no manual reload required
- **FR-UI-05:** Contest metadata (name, description, start/end time, registration status) must be prominently displayed on the Contest Landing Page
- **FR-UI-06:** The platform must gracefully handle the transition from live rankings (during contest) to saved/final rankings (post-contest)

---

## 11. Out of Scope

The following are explicitly out of scope for v1.0:

- LSP / IntelliSense in the code editor (plain Monaco editor is used)
- Plagiarism detection between submissions
- Custom test case upload by participants
- Admin analytics dashboard (beyond leaderboard)
- Team-based contests (individual participation only)
- Third-party OAuth (Google, GitHub) — Firebase Auth handles identity
- Support for compiled languages beyond what Judge0 supports
- Video proctoring or browser lockdown

---

## 12. Risks & Mitigations

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|-----------|
| Judge0 overload during peak submission bursts | Medium | High | Queue submissions; horizontal scale Judge0 workers |
| Redis cache invalidation causing stale rankings | Low | Medium | TTL-based expiry + forced invalidation on score update |
| DynamoDB hot partition during contest spike | Low | High | Use ULID-based partition keys (avoids sequential hot spots) |
| SSE connection drops causing missed status updates | Medium | Medium | Implement client-side reconnect with last-event-id |
| Consolidator failure at contest end | Low | High | Idempotent design; re-triggerable via admin action |
| LaTeX rendering failures on problem pages | Low | Low | Fallback to plaintext Markdown; pre-validate on upload |

---

## 13. Open Questions

1. **Multi-language support:** Which programming languages will be supported in v1.0? (Affects LSP server setup per language)
2. **Submission limits:** Is there a cap on the number of submissions per problem per participant?
3. **Problem visibility:** Can participants see problems before the contest starts (for preparation), or only after it begins?
4. **Admin roles:** Is there a distinction between a super-admin and a contest-specific admin?
5. **Recruitment output:** What format should the final ranked candidate list be exported in (CSV, PDF, API)?
6. **Proctoring:** Are there any basic integrity measures needed (e.g., tab-switch detection) for v1.0?

---

*Document prepared based on system architecture diagram analysis. Review with engineering and product leads before finalizing.*
