terraform {
  required_version = ">= 0.13"
  backend "gcs" {
    bucket = "tf-state-sender-receiver"
    prefix = "terraform/infrastructure"
  }
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "2.11.0"
    }
  }
}

provider "kubernetes" {
  token                  = data.google_client_config.default.access_token
  client_certificate     = base64decode(google_container_cluster.main.master_auth.0.client_certificate)
  cluster_ca_certificate = base64decode(google_container_cluster.main.master_auth.0.cluster_ca_certificate)
  client_key             = base64decode(google_container_cluster.main.master_auth.0.client_key)
  host                   = google_container_cluster.main.endpoint
}

provider "helm" {
  kubernetes {
    token                  = data.google_client_config.default.access_token
    client_certificate     = base64decode(google_container_cluster.main.master_auth.0.client_certificate)
    cluster_ca_certificate = base64decode(google_container_cluster.main.master_auth.0.cluster_ca_certificate)
    client_key             = base64decode(google_container_cluster.main.master_auth.0.client_key)
    host                   = google_container_cluster.main.endpoint
  }
}

provider "cloudflare" {
  email   = var.cloudflare_email
  api_key = var.cloudflare_api_key
}
