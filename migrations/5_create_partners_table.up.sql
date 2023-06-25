CREATE TABLE "partners" (
  "id" uuid DEFAULT gen_random_uuid() NOT NULL,
  "alt" VARCHAR(30) NOT NULL,
  "thumb" VARCHAR(255) NOT NULL,
  "created" Timestamp With Time Zone NOT NULL DEFAULT NOW(),
  "modified" Timestamp With Time Zone NOT NULL DEFAULT NOW(),
  PRIMARY KEY ("id"),
  CONSTRAINT "unique_partners_id" UNIQUE("id")
);

CREATE TRIGGER update_partners_modtime 
BEFORE UPDATE ON "partners" 
FOR EACH ROW EXECUTE PROCEDURE  update_modified_column();