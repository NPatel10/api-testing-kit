-- +goose Up

DO $$
BEGIN
  CREATE TYPE entitlement_source AS ENUM ('plan', 'trial', 'manual', 'comped', 'subscription');
EXCEPTION
  WHEN duplicate_object THEN NULL;
END $$;

DO $$
BEGIN
  CREATE TYPE subscription_status AS ENUM ('trialing', 'active', 'past_due', 'paused', 'canceled', 'incomplete', 'unpaid');
EXCEPTION
  WHEN duplicate_object THEN NULL;
END $$;

DO $$
BEGIN
  CREATE TYPE billing_provider AS ENUM ('stripe', 'paddle', 'lemonsqueezy', 'manual');
EXCEPTION
  WHEN duplicate_object THEN NULL;
END $$;

DO $$
BEGIN
  CREATE TYPE billing_interval AS ENUM ('month', 'year', 'one_time');
EXCEPTION
  WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE IF NOT EXISTS user_identities (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  provider text NOT NULL,
  provider_user_id text NOT NULL,
  provider_email text,
  access_token_encrypted text,
  refresh_token_encrypted text,
  token_expires_at timestamptz,
  profile jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (provider, provider_user_id)
);

COMMENT ON TABLE user_identities IS 'Optional OAuth or external identity links. Token fields should be encrypted in the application layer.';

CREATE TABLE IF NOT EXISTS email_verification_tokens (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token_hash text NOT NULL UNIQUE,
  email text NOT NULL,
  expires_at timestamptz NOT NULL,
  consumed_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_email_verification_tokens_user_id ON email_verification_tokens (user_id);
CREATE INDEX IF NOT EXISTS idx_email_verification_tokens_expires_at ON email_verification_tokens (expires_at);

CREATE TABLE IF NOT EXISTS password_reset_tokens (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token_hash text NOT NULL UNIQUE,
  expires_at timestamptz NOT NULL,
  consumed_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now()
);

COMMENT ON TABLE password_reset_tokens IS 'Optional local-auth support if password login is enabled.';

CREATE TABLE IF NOT EXISTS environments (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  owner_user_id uuid REFERENCES users(id) ON DELETE CASCADE,
  collection_id uuid REFERENCES collections(id) ON DELETE CASCADE,
  name text NOT NULL,
  description text,
  is_default boolean NOT NULL DEFAULT false,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz,
  CONSTRAINT environments_one_owner_or_collection CHECK (
    (owner_user_id IS NOT NULL AND collection_id IS NULL)
    OR
    (owner_user_id IS NULL AND collection_id IS NOT NULL)
  )
);

CREATE INDEX IF NOT EXISTS idx_environments_owner_user_id ON environments (owner_user_id);
CREATE INDEX IF NOT EXISTS idx_environments_collection_id ON environments (collection_id);

CREATE TABLE IF NOT EXISTS environment_variables (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  environment_id uuid NOT NULL REFERENCES environments(id) ON DELETE CASCADE,
  name text NOT NULL,
  value_encrypted text,
  value_preview text,
  is_secret boolean NOT NULL DEFAULT false,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (environment_id, name)
);

COMMENT ON TABLE environment_variables IS 'Store secrets encrypted in the application layer. Keep previews non-sensitive.';
COMMENT ON COLUMN environment_variables.value_encrypted IS 'Application-layer encrypted secret value.';

CREATE INDEX IF NOT EXISTS idx_environment_variables_environment_id ON environment_variables (environment_id);

CREATE TABLE IF NOT EXISTS plans (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  code text NOT NULL UNIQUE,
  name text NOT NULL,
  description text,
  is_active boolean NOT NULL DEFAULT true,
  sort_order integer NOT NULL DEFAULT 0,
  metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

COMMENT ON TABLE plans IS 'Canonical product plan definitions independent of the billing provider.';

CREATE TABLE IF NOT EXISTS plan_prices (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  plan_id uuid NOT NULL REFERENCES plans(id) ON DELETE CASCADE,
  provider billing_provider NOT NULL,
  provider_price_id text NOT NULL,
  billing_interval billing_interval NOT NULL DEFAULT 'month',
  amount_cents integer NOT NULL DEFAULT 0 CHECK (amount_cents >= 0),
  currency text NOT NULL DEFAULT 'usd',
  is_active boolean NOT NULL DEFAULT true,
  metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (provider, provider_price_id)
);

COMMENT ON TABLE plan_prices IS 'Provider-specific sellable price points for a plan. Supports one provider now and leaves room for another later.';

CREATE INDEX IF NOT EXISTS idx_plan_prices_plan_id ON plan_prices (plan_id);
CREATE INDEX IF NOT EXISTS idx_plan_prices_provider_is_active ON plan_prices (provider, is_active);

CREATE TABLE IF NOT EXISTS plan_entitlements (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  plan_id uuid NOT NULL REFERENCES plans(id) ON DELETE CASCADE,
  key text NOT NULL,
  value jsonb NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (plan_id, key)
);

COMMENT ON TABLE plan_entitlements IS 'Plan-level feature limits and capabilities such as quotas, request execution, and retention.';

CREATE INDEX IF NOT EXISTS idx_plan_entitlements_plan_id ON plan_entitlements (plan_id);

CREATE TABLE IF NOT EXISTS billing_customers (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  provider billing_provider NOT NULL,
  provider_customer_id text NOT NULL,
  provider_email text,
  tax_country text,
  tax_id_last4 text,
  metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (provider, provider_customer_id),
  UNIQUE (user_id, provider)
);

COMMENT ON TABLE billing_customers IS 'Maps internal users to billing provider customers. Store only non-sensitive provider metadata here.';

CREATE TABLE IF NOT EXISTS subscriptions (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  billing_customer_id uuid REFERENCES billing_customers(id) ON DELETE SET NULL,
  plan_id uuid NOT NULL REFERENCES plans(id) ON DELETE RESTRICT,
  plan_price_id uuid REFERENCES plan_prices(id) ON DELETE SET NULL,
  provider billing_provider NOT NULL,
  provider_subscription_id text NOT NULL,
  status subscription_status NOT NULL DEFAULT 'incomplete',
  quantity integer NOT NULL DEFAULT 1,
  trial_ends_at timestamptz,
  current_period_start timestamptz,
  current_period_end timestamptz,
  cancel_at_period_end boolean NOT NULL DEFAULT false,
  canceled_at timestamptz,
  paused_at timestamptz,
  past_due_since timestamptz,
  payment_method_hint text,
  metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (provider, provider_subscription_id)
);

COMMENT ON TABLE subscriptions IS 'Provider-backed subscription state. Webhooks should be the source of truth.';

CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id_status ON subscriptions (user_id, status);
CREATE INDEX IF NOT EXISTS idx_subscriptions_billing_customer_id ON subscriptions (billing_customer_id);

CREATE TABLE IF NOT EXISTS subscription_events (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  subscription_id uuid NOT NULL REFERENCES subscriptions(id) ON DELETE CASCADE,
  provider billing_provider NOT NULL,
  provider_event_id text NOT NULL,
  event_type text NOT NULL,
  payload jsonb NOT NULL,
  processed_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (provider, provider_event_id)
);

COMMENT ON TABLE subscription_events IS 'Webhook/event inbox for Stripe, Paddle, or Lemon integrations.';

CREATE INDEX IF NOT EXISTS idx_subscription_events_subscription_id_created_at ON subscription_events (subscription_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_subscription_events_event_type_created_at ON subscription_events (event_type, created_at DESC);

CREATE TABLE IF NOT EXISTS invoices (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  subscription_id uuid REFERENCES subscriptions(id) ON DELETE SET NULL,
  provider billing_provider NOT NULL,
  provider_invoice_id text NOT NULL,
  invoice_number text,
  status text NOT NULL,
  amount_due_cents integer NOT NULL DEFAULT 0 CHECK (amount_due_cents >= 0),
  amount_paid_cents integer NOT NULL DEFAULT 0 CHECK (amount_paid_cents >= 0),
  currency text NOT NULL DEFAULT 'usd',
  hosted_invoice_url text,
  pdf_url text,
  issued_at timestamptz,
  due_at timestamptz,
  paid_at timestamptz,
  metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (provider, provider_invoice_id)
);

COMMENT ON TABLE invoices IS 'Optional invoice history for billing portal and support workflows.';

CREATE INDEX IF NOT EXISTS idx_invoices_subscription_id ON invoices (subscription_id);
CREATE INDEX IF NOT EXISTS idx_invoices_status_issued_at ON invoices (status, issued_at DESC);

CREATE TABLE IF NOT EXISTS user_entitlements (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  plan_id uuid REFERENCES plans(id) ON DELETE SET NULL,
  source entitlement_source NOT NULL DEFAULT 'manual',
  key text NOT NULL,
  value jsonb NOT NULL,
  starts_at timestamptz,
  ends_at timestamptz,
  is_active boolean NOT NULL DEFAULT true,
  metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (user_id, key, source)
);

COMMENT ON TABLE user_entitlements IS 'Flattened effective entitlements for gating features and quotas at runtime.';

CREATE INDEX IF NOT EXISTS idx_user_entitlements_user_id_is_active ON user_entitlements (user_id, is_active);
CREATE INDEX IF NOT EXISTS idx_user_entitlements_key_is_active ON user_entitlements (key, is_active);

CREATE TABLE IF NOT EXISTS audit_log_entries (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  actor_user_id uuid REFERENCES users(id) ON DELETE SET NULL,
  entity_type text NOT NULL,
  entity_id uuid,
  action text NOT NULL,
  before_state jsonb NOT NULL DEFAULT '{}'::jsonb,
  after_state jsonb NOT NULL DEFAULT '{}'::jsonb,
  metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at timestamptz NOT NULL DEFAULT now()
);

COMMENT ON TABLE audit_log_entries IS 'Generic admin audit trail for account, billing, and abuse operations.';

CREATE INDEX IF NOT EXISTS idx_audit_log_entries_actor_user_id_created_at ON audit_log_entries (actor_user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_log_entries_entity_type_entity_id ON audit_log_entries (entity_type, entity_id);

DROP TRIGGER IF EXISTS trg_user_identities_updated_at ON user_identities;
CREATE TRIGGER trg_user_identities_updated_at
BEFORE UPDATE ON user_identities
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_environments_updated_at ON environments;
CREATE TRIGGER trg_environments_updated_at
BEFORE UPDATE ON environments
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_environment_variables_updated_at ON environment_variables;
CREATE TRIGGER trg_environment_variables_updated_at
BEFORE UPDATE ON environment_variables
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_plans_updated_at ON plans;
CREATE TRIGGER trg_plans_updated_at
BEFORE UPDATE ON plans
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_plan_prices_updated_at ON plan_prices;
CREATE TRIGGER trg_plan_prices_updated_at
BEFORE UPDATE ON plan_prices
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_billing_customers_updated_at ON billing_customers;
CREATE TRIGGER trg_billing_customers_updated_at
BEFORE UPDATE ON billing_customers
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_subscriptions_updated_at ON subscriptions;
CREATE TRIGGER trg_subscriptions_updated_at
BEFORE UPDATE ON subscriptions
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_invoices_updated_at ON invoices;
CREATE TRIGGER trg_invoices_updated_at
BEFORE UPDATE ON invoices
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_user_entitlements_updated_at ON user_entitlements;
CREATE TRIGGER trg_user_entitlements_updated_at
BEFORE UPDATE ON user_entitlements
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- +goose Down

DROP TRIGGER IF EXISTS trg_user_entitlements_updated_at ON user_entitlements;
DROP TRIGGER IF EXISTS trg_invoices_updated_at ON invoices;
DROP TRIGGER IF EXISTS trg_subscriptions_updated_at ON subscriptions;
DROP TRIGGER IF EXISTS trg_billing_customers_updated_at ON billing_customers;
DROP TRIGGER IF EXISTS trg_plan_prices_updated_at ON plan_prices;
DROP TRIGGER IF EXISTS trg_plans_updated_at ON plans;
DROP TRIGGER IF EXISTS trg_environment_variables_updated_at ON environment_variables;
DROP TRIGGER IF EXISTS trg_environments_updated_at ON environments;
DROP TRIGGER IF EXISTS trg_user_identities_updated_at ON user_identities;

DROP TABLE IF EXISTS audit_log_entries;
DROP TABLE IF EXISTS user_entitlements;
DROP TABLE IF EXISTS invoices;
DROP TABLE IF EXISTS subscription_events;
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS billing_customers;
DROP TABLE IF EXISTS plan_entitlements;
DROP TABLE IF EXISTS plan_prices;
DROP TABLE IF EXISTS plans;
DROP TABLE IF EXISTS environment_variables;
DROP TABLE IF EXISTS environments;
DROP TABLE IF EXISTS password_reset_tokens;
DROP TABLE IF EXISTS email_verification_tokens;
DROP TABLE IF EXISTS user_identities;

DROP TYPE IF EXISTS billing_interval;
DROP TYPE IF EXISTS billing_provider;
DROP TYPE IF EXISTS subscription_status;
DROP TYPE IF EXISTS entitlement_source;
