-- Modify "users" table
ALTER TABLE "public"."users"
ADD COLUMN "role" character varying(30) NOT NULL DEFAULT 'user',
    ADD CONSTRAINT "users_role_fkey" FOREIGN KEY ("role") REFERENCES "public"."roles" ("name") ON UPDATE NO ACTION ON DELETE NO ACTION;