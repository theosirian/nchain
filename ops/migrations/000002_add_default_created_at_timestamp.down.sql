ALTER TABLE ONLY accounts ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE ONLY bridges ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE ONLY connectors ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE ONLY contracts ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE ONLY filters ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE ONLY load_balancers ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE ONLY networks ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE ONLY nodes ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE ONLY oracles ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE ONLY tokens ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE ONLY transactions ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE ONLY wallets ALTER COLUMN created_at DROP DEFAULT;
