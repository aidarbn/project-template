BEGIN;

CREATE TABLE IF NOT EXISTS "objects"
(
    "id"         uuid                 DEFAULT gen_random_uuid(),
    "data"       text,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY ("id")
);

CREATE OR REPLACE FUNCTION update_objects()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

CREATE TRIGGER update_objects
    BEFORE UPDATE
    ON "objects"
    FOR EACH ROW
EXECUTE PROCEDURE update_objects();

COMMIT;