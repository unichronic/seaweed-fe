CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id             TEXT PRIMARY KEY,
    name           TEXT NOT NULL,
	email          TEXT NOT NULL UNIQUE,
    usn            TEXT NOT NULL UNIQUE,
    mobile_number  TEXT,
    joining_year   INT NOT NULL,
    department     TEXT NOT NULL
);

CREATE TABLE admin (
    user_id    TEXT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE contests (
    id                      TEXT PRIMARY KEY,
    name                    TEXT NOT NULL,
    description             TEXT,
    eligible_to             TEXT,
    registration_status     TEXT NOT NULL DEFAULT 'open',
    registration_start_time BIGINT NOT NULL,
    registration_end_time   BIGINT NOT NULL,
    start_time              BIGINT NOT NULL,
    end_time                BIGINT NOT NULL,
    finalized               BOOLEAN NOT NULL DEFAULT FALSE,
    CHECK (registration_status IN ('open', 'closed', 'invite-only'))
);

CREATE TABLE contest_registrations (
    contest_id    TEXT NOT NULL REFERENCES contests(id) ON DELETE CASCADE,
    user_id       TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    registered_at BIGINT NOT NULL,
    PRIMARY KEY (contest_id, user_id)
);

CREATE TABLE problems (
    id          TEXT PRIMARY KEY,
    contest_id  TEXT NOT NULL REFERENCES contests(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    description TEXT,
    score       INT NOT NULL,
    test_cases  JSONB NOT NULL DEFAULT '[]'::jsonb
);

CREATE TYPE submission_status AS ENUM (
    'pending', 'failed_to_process', 'accepted',
    'tle', 'mle', 'rte', 'wrong_answer'
);

CREATE TABLE submissions (
    id          TEXT PRIMARY KEY,
    user_id     TEXT NOT NULL REFERENCES users(id),
    contest_id  TEXT NOT NULL REFERENCES contests(id) ON DELETE CASCADE,
    problem_id  TEXT NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    language    TEXT NOT NULL,
    s3_key      TEXT NOT NULL,
    status      submission_status NOT NULL DEFAULT 'pending',
    created_at  BIGINT NOT NULL,
    updated_at  BIGINT
);

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

CREATE TABLE rankings (
    contest_id    TEXT NOT NULL REFERENCES contests(id) ON DELETE CASCADE,
    user_id       TEXT NOT NULL REFERENCES users(id),
    score         INT NOT NULL DEFAULT 0,
    hidden        BOOLEAN NOT NULL DEFAULT FALSE,
    disqualified  BOOLEAN NOT NULL DEFAULT FALSE,
    shortlisted   BOOLEAN NOT NULL DEFAULT FALSE,
    correct_attempts   INT NOT NULL DEFAULT 0,
    incorrect_attempts INT NOT NULL DEFAULT 0,
    PRIMARY KEY (contest_id, user_id)
);

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
    r.correct_attempts,
    r.incorrect_attempts,
    (RANK() OVER (PARTITION BY r.contest_id ORDER BY r.score DESC))::INT AS rank
FROM rankings r
JOIN users u ON u.id = r.user_id
WHERE r.disqualified = FALSE
ORDER BY r.contest_id, rank;

CREATE UNIQUE INDEX ON ranking_mv (contest_id, user_id);
