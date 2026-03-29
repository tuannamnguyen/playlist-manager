resource "aws_secretsmanager_secret" "db_host" {
  name        = "db-host"
  description = "Host value for Postgres DB"
}

resource "aws_secretsmanager_secret_version" "db_host" {
  secret_id     = aws_secretsmanager_secret.db_host.id
  secret_string = aws_db_instance.db.address
}
