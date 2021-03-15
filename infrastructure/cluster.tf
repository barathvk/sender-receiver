resource "google_service_account" "main" {
  account_id   = var.name
  display_name = var.name
  project      = var.project
}

resource "google_container_cluster" "main" {
  name               = var.name
  location           = var.location
  initial_node_count = var.node_count
  project            = var.project
  node_config {
    service_account = google_service_account.main.email
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]
  }
}
