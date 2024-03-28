DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'req_status') THEN
CREATE TYPE req_status AS ENUM ('in progress','approved','rejected');
    END IF;
END $$;
