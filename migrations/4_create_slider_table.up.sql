CREATE TABLE "slider" (
  "id" uuid DEFAULT gen_random_uuid() NOT NULL,
  "title" VARCHAR(30) NOT NULL,
  "thumb" VARCHAR(255) NOT NULL,
  "created" Timestamp With Time Zone NOT NULL DEFAULT NOW(),
  "modified" Timestamp With Time Zone NOT NULL DEFAULT NOW(),
  PRIMARY KEY ("id"),
  CONSTRAINT "unique_slider_id" UNIQUE("id")
);

CREATE TRIGGER update_slider_modtime 
BEFORE UPDATE ON "slider" 
FOR EACH ROW EXECUTE PROCEDURE  update_modified_column();