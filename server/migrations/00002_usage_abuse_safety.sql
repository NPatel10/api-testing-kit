-- +goose Up

-- +goose StatementBegin
DO $$
BEGIN
  CREATE TYPE usage_bucket AS ENUM ('minute', 'hour', 'day', 'month');
EXCEPTION
  WHEN duplicate_object THEN NULL;
END $$;
-- +goose StatementEnd

-- +goose StatementBegin
DO $$
BEGIN
  CREATE TYPE abuse_severity AS ENUM ('low', 'medium', 'high', 'critical');
EXCEPTION
  WHEN duplicate_object THEN NULL;
END $$;
-- +goose StatementEnd

-- +goose StatementBegin
DO $$
BEGIN
  CREATE TYPE blocked_target_type AS ENUM ('domain', 'host', 'ip', 'cidr', 'url_pattern');
EXCEPTION
  WHEN duplicate_object THEN NULL;
END $$;
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS usage_events (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid REFERENCES users(id) ON DELETE SET NULL,
  session_id uuid REFERENCES sessions(id) ON DELETE SET NULL,
  request_run_id uuid,
  bucket usage_bucket NOT NULL DEFAULT 'month',
  event_key text NOT NULL,
  quantity integer NOT NULL DEFAULT 1,
  dimensions jsonb NOT NULL DEFAULT '{}'::jsonb,
  occurred_at timestamptz NOT NULL DEFAULT now()
);

COMMENT ON TABLE usage_events IS 'Time-series usage records for quotas, analytics, and plan enforcement.';

CREATE INDEX IF NOT EXISTS idx_usage_events_user_id_occurred_at ON usage_events (user_id, occurred_at DESC);
CREATE INDEX IF NOT EXISTS idx_usage_events_event_key_occurred_at ON usage_events (event_key, occurred_at DESC);
CREATE INDEX IF NOT EXISTS idx_usage_events_occurred_at ON usage_events (occurred_at DESC);

CREATE TABLE IF NOT EXISTS abuse_events (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid REFERENCES users(id) ON DELETE SET NULL,
  session_id uuid REFERENCES sessions(id) ON DELETE SET NULL,
  request_run_id uuid,
  severity abuse_severity NOT NULL DEFAULT 'low',
  category text NOT NULL,
  source_ip inet,
  target text,
  rule_key text,
  action_taken text NOT NULL,
  details jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at timestamptz NOT NULL DEFAULT now()
);

COMMENT ON TABLE abuse_events IS 'Audit trail for blocked requests, suspicious behavior, and automated enforcement.';

CREATE INDEX IF NOT EXISTS idx_abuse_events_user_id_created_at ON abuse_events (user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_abuse_events_category_created_at ON abuse_events (category, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_abuse_events_source_ip_created_at ON abuse_events (source_ip, created_at DESC);

CREATE TABLE IF NOT EXISTS blocked_targets (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  target_type blocked_target_type NOT NULL,
  target_value text NOT NULL,
  reason text NOT NULL,
  source text NOT NULL DEFAULT 'manual',
  is_active boolean NOT NULL DEFAULT true,
  expires_at timestamptz,
  created_by_user_id uuid REFERENCES users(id) ON DELETE SET NULL,
  metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (target_type, target_value)
);

COMMENT ON TABLE blocked_targets IS 'Domain/IP/path denylist used to prevent SSRF and abuse.';

CREATE INDEX IF NOT EXISTS idx_blocked_targets_is_active ON blocked_targets (is_active);
CREATE INDEX IF NOT EXISTS idx_blocked_targets_expires_at ON blocked_targets (expires_at);

DROP TRIGGER IF EXISTS trg_blocked_targets_updated_at ON blocked_targets;
CREATE TRIGGER trg_blocked_targets_updated_at
BEFORE UPDATE ON blocked_targets
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- +goose Down

DROP TRIGGER IF EXISTS trg_blocked_targets_updated_at ON blocked_targets;
DROP TABLE IF EXISTS blocked_targets;
DROP TABLE IF EXISTS abuse_events;
DROP TABLE IF EXISTS usage_events;

DROP TYPE IF EXISTS blocked_target_type;
DROP TYPE IF EXISTS abuse_severity;
DROP TYPE IF EXISTS usage_bucket;
