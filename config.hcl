db {
  database = "games"
  user = "games"
  password = "games"
  host = "localhost"
  port = 5437
  sslmode = "disable"
  migrations = "./db/migrations"
}

app {
  listening = 8083
  prod = false
  disableStacktrace = true
}