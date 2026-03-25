-- +goose Up

-- +goose StatementBegin
DO $$
BEGIN
  CREATE TYPE collection_visibility AS ENUM ('private', 'shared_readonly', 'internal');
EXCEPTION
  WHEN duplicate_object THEN NULL;
END $$;
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS collections (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  owner_user_id uuid REFERENCES users(id) ON DELETE SET NULL,
  name text NOT NULL,
  slug text,
  description text,
  visibility collection_visibility NOT NULL DEFAULT 'private',
  color text,
  sort_order integer NOT NULL DEFAULT 0,
  shared_token text UNIQUE,
  metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz,
  CONSTRAINT collections_name_not_empty CHECK (length(trim(name)) > 0)
);

COMMENT ON TABLE collections IS 'Saved request groups shown in the app sidebar and collection pages.';

CREATE INDEX IF NOT EXISTS idx_collections_owner_user_id ON collections (owner_user_id);
CREATE INDEX IF NOT EXISTS idx_collections_visibility ON collections (visibility);

DROP TRIGGER IF EXISTS trg_collections_updated_at ON collections;
CREATE TRIGGER trg_collections_updated_at
BEFORE UPDATE ON collections
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- +goose Down

DROP TRIGGER IF EXISTS trg_collections_updated_at ON collections;
DROP TABLE IF EXISTS collections;
DROP TYPE IF EXISTS collection_visibility;
