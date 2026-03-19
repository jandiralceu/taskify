-- Restore legacy role enum type and column
CREATE TYPE user_role AS ENUM ('admin', 'employee');
ALTER TABLE users ADD COLUMN role user_role NOT NULL DEFAULT 'employee';

-- Optional: try to restore data if needed, but for now simple restore is enough for rollback structure
