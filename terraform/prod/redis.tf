resource "google_redis_instance" "playlist_manager_redis" {
  name           = "playlist-manager-prod-redis"
  tier           = "BASIC"
  memory_size_gb = "2"
  region         = local.region
  redis_version  = "REDIS_7_0"
}
