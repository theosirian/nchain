ALTER TABLE accounts ALTER COLUMN user_id TYPE text;
ALTER TABLE networks ALTER COLUMN user_id TYPE text;
ALTER TABLE nodes ALTER COLUMN user_id TYPE text;
ALTER TABLE transactions ALTER COLUMN user_id TYPE text;
ALTER TABLE wallets ALTER COLUMN user_id TYPE text;

ALTER TABLE connectors ALTER COLUMN organization_id TYPE text;
ALTER TABLE load_balancers ALTER COLUMN organization_id TYPE text;
ALTER TABLE nodes ALTER COLUMN organization_id TYPE text;

ALTER TABLE accounts ALTER COLUMN organization_id TYPE text;
ALTER TABLE wallets ALTER COLUMN organization_id TYPE text;

ALTER TABLE contracts ALTER COLUMN organization_id TYPE text;

ALTER TABLE transactions ALTER COLUMN organization_id TYPE text;
ALTER TABLE tokens ALTER COLUMN organization_id TYPE text;
