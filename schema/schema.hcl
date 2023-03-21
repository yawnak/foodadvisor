table "hints" {
  schema = schema.public
  column "id" {
    null = false
    type = integer
  }
  column "name" {
    null = false
    type = text
  }
  primary_key {
    columns = [column.id]
  }
}
table "hints_to_meals" {
  schema = schema.public
  column "hintid" {
    null = false
    type = integer
  }
  column "mealid" {
    null = false
    type = integer
  }
  primary_key {
    columns = [column.mealid, column.hintid]
  }
  foreign_key "hints_to_meals_hintid_fkey" {
    columns     = [column.hintid]
    ref_columns = [table.hints.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  foreign_key "hints_to_meals_mealid_fkey" {
    columns     = [column.mealid]
    ref_columns = [table.meals.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
}
table "meals_to_users" {
  schema = schema.public
  column "userid" {
    null = false
    type = integer
  }
  column "mealid" {
    null = false
    type = integer
  }
  column "lasteaten" {
    null = false
    type = date
  }
  primary_key {
    columns = [column.userid, column.mealid]
  }
  foreign_key "meals_to_users_mealid_fkey" {
    columns     = [column.mealid]
    ref_columns = [table.meals.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  foreign_key "meals_to_users_userid_fkey" {
    columns     = [column.userid]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
}

table permissions {
  schema = schema.public
  column "name" {
    null = false
    type = character_varying(30)
  }
  primary_key {
    columns = [column.name]
  }
}

table permissions_to_users {
  schema = schema.public
  column "permission" {
    null = false
    type = character_varying(30)
  }
  column "role" {
    null = false
    type = character_varying(30)
  }
  primary_key {
    columns = [column.role, column.permission]
  }
    foreign_key "permissions_to_role_permission_fkey" {
    columns     = [column.permission]
    ref_columns = [table.permissions.column.name]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  foreign_key "permissions_to_role_role_fkey" {
    columns     = [column.role]
    ref_columns = [table.roles.column.name]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
}

table "roles" {
  schema = schema.public
  column "name" {
    null = false
    type = character_varying(30)
  }
  primary_key {
    columns = [column.name]
  }
}

table "users" {
  schema = schema.public
  column "id" {
    null = false
    type = serial
  }
  column "username" {
    null = false
    type = character_varying(30)
  }
  column "passhash" {
    null = false
    type = text
  }
  column "expiration" {
    null    = true
    type    = day
    default = sql("'00:00:00'::interval")
  }
  primary_key {
    columns = [column.id]
  }
  index "users_username_key" {
    unique  = true
    columns = [column.username]
  }
}
schema "public" {
}
