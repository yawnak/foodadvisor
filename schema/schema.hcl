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

table "hintstousers" {
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
  foreign_key "hintstousers_mealid_fkey" {
    columns     = [column.mealid]
    ref_columns = [table.meals.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  foreign_key "hintstousers_hintid_fkey" {
    columns     = [column.hintid]
    ref_columns = [table.hints.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
}

table "mealstousers" {
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
  foreign_key "foodtousers_mealid_fkey" {
    columns     = [column.mealid]
    ref_columns = [table.meals.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  foreign_key "foodtousers_userid_fkey" {
    columns     = [column.userid]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
}
table "meals" {
  schema = schema.public
  column "id" {
    null = false
    type = serial
  }
  column "name" {
    null = false
    type = character_varying(30)
  }
  column "cooktime" {
    null = false
    type = minute
  }
  primary_key {
    columns = [column.id]
  }
  index "meals_name_key" {
    unique  = true
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
