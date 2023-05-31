CREATE TABLE "cards" (
  "id" uuid DEFAULT gen_random_uuid() NOT NULL,
  "title" VARCHAR(252) NOT NULL,
  "thumb" VARCHAR(255) NOT NULL,
  "alt" VARCHAR(255),
  "description" JSONB NOT NULL,
  "created" Timestamp With Time Zone NOT NULL DEFAULT NOW(),
  "modified" Timestamp With Time Zone NOT NULL DEFAULT NOW(),
  PRIMARY KEY ("id"),
  CONSTRAINT "unique_cards_id" UNIQUE("id"),
  CONSTRAINT "unique_cards_title" UNIQUE("title")
);

CREATE TRIGGER update_cards_modtime 
BEFORE UPDATE ON "cards" 
FOR EACH ROW EXECUTE PROCEDURE  update_modified_column();