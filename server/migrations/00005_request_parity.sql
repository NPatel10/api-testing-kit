-- +goose Up

ALTER TABLE saved_requests
  ADD COLUMN IF NOT EXISTS path_template text,
  ADD COLUMN IF NOT EXISTS pre_request_script text,
  ADD COLUMN IF NOT EXISTS post_request_script text,
  ADD COLUMN IF NOT EXISTS assertions jsonb NOT NULL DEFAULT '[]'::jsonb,
  ADD COLUMN IF NOT EXISTS tags text[] NOT NULL DEFAULT '{}'::text[],
  ADD COLUMN IF NOT EXISTS is_template boolean NOT NULL DEFAULT false,
  ADD COLUMN IF NOT EXISTS sort_order integer NOT NULL DEFAULT 0;

ALTER TABLE saved_requests
  ALTER COLUMN query_params SET DEFAULT '[]'::jsonb,
  ALTER COLUMN headers SET DEFAULT '[]'::jsonb;

CREATE INDEX IF NOT EXISTS idx_saved_requests_is_template ON saved_requests (is_template);

COMMENT ON COLUMN saved_requests.path_template IS 'Optional path template for reusable request variables.';
COMMENT ON COLUMN saved_requests.pre_request_script IS 'Optional script executed before a request is sent.';
COMMENT ON COLUMN saved_requests.post_request_script IS 'Optional script executed after a response is received.';
COMMENT ON COLUMN saved_requests.assertions IS 'Saved response assertions for later replay and debugging.';
COMMENT ON COLUMN saved_requests.tags IS 'Lightweight labels for grouping saved requests and templates.';
COMMENT ON COLUMN saved_requests.is_template IS 'Marks a saved request as template material for reuse.';

ALTER TABLE request_runs
  ADD COLUMN IF NOT EXISTS target_ip inet,
  ADD COLUMN IF NOT EXISTS request_body_preview text,
  ADD COLUMN IF NOT EXISTS request_size_bytes integer,
  ADD COLUMN IF NOT EXISTS response_body text,
  ADD COLUMN IF NOT EXISTS response_body_json jsonb;

ALTER TABLE request_runs
  ALTER COLUMN request_headers SET DEFAULT '[]'::jsonb,
  ALTER COLUMN request_query_params SET DEFAULT '[]'::jsonb,
  ALTER COLUMN request_auth SET DEFAULT '{}'::jsonb,
  ALTER COLUMN response_headers SET DEFAULT '[]'::jsonb;

CREATE INDEX IF NOT EXISTS idx_request_runs_collection_id_created_at ON request_runs (collection_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_request_runs_status_created_at ON request_runs (status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_request_runs_target_host ON request_runs (target_host);

DO $$
BEGIN
  ALTER TABLE usage_events
    ADD CONSTRAINT fk_usage_events_request_run_id
    FOREIGN KEY (request_run_id) REFERENCES request_runs(id) ON DELETE SET NULL;
EXCEPTION
  WHEN duplicate_object THEN NULL;
END $$;

DO $$
BEGIN
  ALTER TABLE abuse_events
    ADD CONSTRAINT fk_abuse_events_request_run_id
    FOREIGN KEY (request_run_id) REFERENCES request_runs(id) ON DELETE SET NULL;
EXCEPTION
  WHEN duplicate_object THEN NULL;
END $$;

-- +goose Down

DO $$
BEGIN
  ALTER TABLE abuse_events DROP CONSTRAINT IF EXISTS fk_abuse_events_request_run_id;
EXCEPTION
  WHEN undefined_object THEN NULL;
END $$;

DO $$
BEGIN
  ALTER TABLE usage_events DROP CONSTRAINT IF EXISTS fk_usage_events_request_run_id;
EXCEPTION
  WHEN undefined_object THEN NULL;
END $$;

DROP INDEX IF EXISTS idx_request_runs_target_host;
DROP INDEX IF EXISTS idx_request_runs_status_created_at;
DROP INDEX IF EXISTS idx_request_runs_collection_id_created_at;
DROP INDEX IF EXISTS idx_saved_requests_is_template;

ALTER TABLE request_runs
  ALTER COLUMN request_headers SET DEFAULT '{}'::jsonb,
  ALTER COLUMN request_query_params SET DEFAULT '{}'::jsonb,
  ALTER COLUMN request_auth SET DEFAULT '{}'::jsonb,
  ALTER COLUMN response_headers SET DEFAULT '{}'::jsonb,
  DROP COLUMN IF EXISTS response_body_json,
  DROP COLUMN IF EXISTS response_body,
  DROP COLUMN IF EXISTS request_size_bytes,
  DROP COLUMN IF EXISTS request_body_preview,
  DROP COLUMN IF EXISTS target_ip;

ALTER TABLE saved_requests
  ALTER COLUMN query_params SET DEFAULT '{}'::jsonb,
  ALTER COLUMN headers SET DEFAULT '{}'::jsonb,
  DROP COLUMN IF EXISTS sort_order,
  DROP COLUMN IF EXISTS is_template,
  DROP COLUMN IF EXISTS tags,
  DROP COLUMN IF EXISTS assertions,
  DROP COLUMN IF EXISTS post_request_script,
  DROP COLUMN IF EXISTS pre_request_script,
  DROP COLUMN IF EXISTS path_template;
