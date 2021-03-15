data "google_client_config" "default" {}
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

resource "kubernetes_namespace" "redis" {
  metadata {
    name = "redis"
  }
}

resource "helm_release" "redis" {
  repository = "https://charts.bitnami.com/bitnami"
  chart      = "redis"
  name       = "redis"
  namespace  = kubernetes_namespace.redis.metadata.0.name
}
