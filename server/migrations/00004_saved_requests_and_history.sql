-- +goose Up

-- +goose StatementBegin
DO $$
BEGIN
  CREATE TYPE request_method AS ENUM ('GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'HEAD', 'OPTIONS');
EXCEPTION
  WHEN duplicate_object THEN NULL;
END $$;
-- +goose StatementEnd

-- +goose StatementBegin
DO $$
BEGIN
  CREATE TYPE request_body_mode AS ENUM ('none', 'raw', 'json', 'form_urlencoded', 'form_data');
EXCEPTION
  WHEN duplicate_object THEN NULL;
END $$;
-- +goose StatementEnd

-- +goose StatementBegin
DO $$
BEGIN
  CREATE TYPE auth_scheme AS ENUM ('none', 'basic', 'bearer', 'api_key', 'oauth2', 'custom');
EXCEPTION
  WHEN duplicate_object THEN NULL;
END $$;
-- +goose StatementEnd

-- +goose StatementBegin
DO $$
BEGIN
  CREATE TYPE run_source AS ENUM ('guest', 'authenticated', 'template', 'manual_replay', 'import');
EXCEPTION
  WHEN duplicate_object THEN NULL;
END $$;
-- +goose StatementEnd

-- +goose StatementBegin
DO $$
BEGIN
  CREATE TYPE run_status AS ENUM ('queued', 'running', 'succeeded', 'failed', 'blocked', 'timed_out', 'canceled');
EXCEPTION
  WHEN duplicate_object THEN NULL;
END $$;
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS saved_requests (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  collection_id uuid REFERENCES collections(id) ON DELETE CASCADE,
  owner_user_id uuid REFERENCES users(id) ON DELETE SET NULL,
  name text NOT NULL,
  description text,
  method request_method NOT NULL DEFAULT 'GET',
  url text NOT NULL,
  query_params jsonb NOT NULL DEFAULT '{}'::jsonb,
  headers jsonb NOT NULL DEFAULT '{}'::jsonb,
  auth_scheme auth_scheme NOT NULL DEFAULT 'none',
  auth_config jsonb NOT NULL DEFAULT '{}'::jsonb,
  body_mode request_body_mode NOT NULL DEFAULT 'none',
  body_config jsonb NOT NULL DEFAULT '{}'::jsonb,
  example_response jsonb NOT NULL DEFAULT '{}'::jsonb,
  metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz,
  CONSTRAINT saved_requests_name_not_empty CHECK (length(trim(name)) > 0),
  CONSTRAINT saved_requests_url_http CHECK (url ~* '^https?://')
);

CREATE INDEX IF NOT EXISTS idx_saved_requests_collection_id ON saved_requests (collection_id);
CREATE INDEX IF NOT EXISTS idx_saved_requests_owner_user_id ON saved_requests (owner_user_id);

CREATE TABLE IF NOT EXISTS request_runs (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid REFERENCES users(id) ON DELETE SET NULL,
  collection_id uuid REFERENCES collections(id) ON DELETE SET NULL,
  saved_request_id uuid REFERENCES saved_requests(id) ON DELETE SET NULL,
  source run_source NOT NULL DEFAULT 'authenticated',
  status run_status NOT NULL DEFAULT 'queued',
  method request_method NOT NULL,
  url text NOT NULL,
  final_url text,
  target_host text,
  request_headers jsonb NOT NULL DEFAULT '{}'::jsonb,
  request_query_params jsonb NOT NULL DEFAULT '{}'::jsonb,
  request_auth jsonb NOT NULL DEFAULT '{}'::jsonb,
  request_body jsonb NOT NULL DEFAULT '{}'::jsonb,
  response_status integer,
  response_headers jsonb NOT NULL DEFAULT '{}'::jsonb,
  response_body_preview text,
  response_size_bytes integer,
  response_time_ms integer,
  response_content_type text,
  redirect_count integer NOT NULL DEFAULT 0,
  blocked_reason text,
  error_code text,
  error_message text,
  started_at timestamptz,
  completed_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now(),
  metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
  CONSTRAINT request_runs_url_http CHECK (url ~* '^https?://'),
  CONSTRAINT request_runs_final_url_http CHECK (final_url IS NULL OR final_url ~* '^https?://')
);

CREATE INDEX IF NOT EXISTS idx_request_runs_user_id_created_at ON request_runs (user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_request_runs_saved_request_id_created_at ON request_runs (saved_request_id, created_at DESC);

-- +goose Down

DROP TABLE IF EXISTS request_runs;
DROP TABLE IF EXISTS saved_requests;
DROP TYPE IF EXISTS run_status;
DROP TYPE IF EXISTS run_source;
DROP TYPE IF EXISTS auth_scheme;
DROP TYPE IF EXISTS request_body_mode;
DROP TYPE IF EXISTS request_method;
