DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'role') THEN
CREATE TYPE role AS ENUM ('user', 'admin','superAdmin');
    END IF;
END $$;
