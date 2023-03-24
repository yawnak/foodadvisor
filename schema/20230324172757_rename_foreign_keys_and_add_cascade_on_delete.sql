-- Modify "permissions_to_roles" table
ALTER TABLE "public"."permissions_to_roles" DROP CONSTRAINT "permissions_to_role_permission_fkey",
    DROP CONSTRAINT "permissions_to_role_role_fkey",
    ADD CONSTRAINT "permissions_to_roles_permission_fkey" FOREIGN KEY ("permission") REFERENCES "public"."permissions" ("name") ON UPDATE NO ACTION ON DELETE CASCADE,
    ADD CONSTRAINT "permissions_to_roles_role_fkey" FOREIGN KEY ("role") REFERENCES "public"."roles" ("name") ON UPDATE NO ACTION ON DELETE CASCADE;