-- Drop legacy role column and enum type
ALTER TABLE users DROP COLUMN role;
DROP TYPE IF EXISTS user_role;
