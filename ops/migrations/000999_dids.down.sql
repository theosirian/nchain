ALTER TABLE accounts ALTER COLUMN user_id TYPE uuid;
ALTER TABLE networks ALTER COLUMN user_id TYPE uuid;
ALTER TABLE nodes ALTER COLUMN user_id TYPE uuid;
ALTER TABLE transactions ALTER COLUMN user_id TYPE uuid;
ALTER TABLE wallets ALTER COLUMN user_id TYPE uuid;

ALTER TABLE connectors ALTER COLUMN organization_id TYPE uuid;
ALTER TABLE load_balancers ALTER COLUMN organization_id TYPE uuid;
ALTER TABLE nodes ALTER COLUMN organization_id TYPE uuid;

ALTER TABLE accounts ALTER COLUMN organization_id TYPE uuid;
ALTER TABLE wallets ALTER COLUMN organization_id TYPE uuid;

ALTER TABLE contracts ALTER COLUMN organization_id TYPE uuid;

ALTER TABLE transactions ALTER COLUMN organization_id TYPE uuid;
ALTER TABLE tokens ALTER COLUMN organization_id TYPE uuid;
