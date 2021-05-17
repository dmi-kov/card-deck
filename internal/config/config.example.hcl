db {
  database = "test"
  user = "test"
  password = "test"
  host = "localhost"
  port = 5437
  sslmode = "disable"
  migrations = "./db/migrations"
}

app {
  listening = 8083
  prod = true
  disableStacktrace = true
}