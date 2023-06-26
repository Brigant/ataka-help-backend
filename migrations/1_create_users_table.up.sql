CREATE TABLE "users" (
   "id" uuid DEFAULT gen_random_uuid() NOT NULL,
   "email" VARCHAR(252) NOT NULL,
   "password" VARCHAR(255) NOT NULL,
   "firs_name" VARCHAR(255) NOT NULL,
   "last_name" VARCHAR(255) NOT NULL,
   "created" Timestamp With Time Zone NOT NULL DEFAULT NOW(),
   "modified" Timestamp With Time Zone NOT NULL DEFAULT NOW(),
   PRIMARY KEY ("id"),
   CONSTRAINT "unique_users_id" UNIQUE("id"),
   CONSTRAINT "unique_users_email" UNIQUE("email")
);

CREATE FUNCTION update_modified_column()   
RETURNS TRIGGER AS $$
BEGIN
   IF row(NEW.*) IS DISTINCT FROM row(OLD.*) THEN
      NEW.modified = now(); 
      RETURN NEW;
   ELSE
      RETURN OLD;
   END IF;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_modtime 
BEFORE UPDATE ON "users" 
FOR EACH ROW EXECUTE PROCEDURE  update_modified_column();
