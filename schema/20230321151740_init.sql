-- Create "hints" table
CREATE TABLE "public"."hints" ("id" integer NOT NULL, "name" text NOT NULL, PRIMARY KEY ("id"));
-- Create "meals" table
CREATE TABLE "public"."meals" ("id" serial NOT NULL, "name" character varying(30) NOT NULL, "cooktime" interval minute NOT NULL, PRIMARY KEY ("id"));
-- Create index "meals_name_key" to table: "meals"
CREATE UNIQUE INDEX "meals_name_key" ON "public"."meals" ("name");
-- Create "hints_to_meals" table
CREATE TABLE "public"."hints_to_meals" ("hintid" integer NOT NULL, "mealid" integer NOT NULL, PRIMARY KEY ("mealid", "hintid"), CONSTRAINT "hints_to_meals_hintid_fkey" FOREIGN KEY ("hintid") REFERENCES "public"."hints" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT "hints_to_meals_mealid_fkey" FOREIGN KEY ("mealid") REFERENCES "public"."meals" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- Create "users" table
CREATE TABLE "public"."users" ("id" serial NOT NULL, "username" character varying(30) NOT NULL, "passhash" text NOT NULL, "expiration" interval day NULL DEFAULT '00:00:00'::interval, PRIMARY KEY ("id"));
-- Create index "users_username_key" to table: "users"
CREATE UNIQUE INDEX "users_username_key" ON "public"."users" ("username");
-- Create "meals_to_users" table
CREATE TABLE "public"."meals_to_users" ("userid" integer NOT NULL, "mealid" integer NOT NULL, "lasteaten" date NOT NULL, PRIMARY KEY ("userid", "mealid"), CONSTRAINT "meals_to_users_mealid_fkey" FOREIGN KEY ("mealid") REFERENCES "public"."meals" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT "meals_to_users_userid_fkey" FOREIGN KEY ("userid") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- Create "permissions" table
CREATE TABLE "public"."permissions" ("name" character varying(30) NOT NULL, PRIMARY KEY ("name"));
-- Create "roles" table
CREATE TABLE "public"."roles" ("name" character varying(30) NOT NULL, PRIMARY KEY ("name"));
-- Create "permissions_to_users" table
CREATE TABLE "public"."permissions_to_users" ("permission" character varying(30) NOT NULL, "role" character varying(30) NOT NULL, PRIMARY KEY ("role", "permission"), CONSTRAINT "permissions_to_role_permission_fkey" FOREIGN KEY ("permission") REFERENCES "public"."permissions" ("name") ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT "permissions_to_role_role_fkey" FOREIGN KEY ("role") REFERENCES "public"."roles" ("name") ON UPDATE NO ACTION ON DELETE NO ACTION);
