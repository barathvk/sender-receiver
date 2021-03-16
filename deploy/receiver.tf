resource "kubernetes_deployment" "receiver" {
  metadata {
    name      = "${var.name}-receiver"
    namespace = var.namespace
  }
  spec {
    replicas = "1"
    selector {
      match_labels = {
        "app" = var.name
      }
    }
    template {
      metadata {
        labels = {
          "app" = var.name
        }
      }
      spec {
        container {
          image = var.image
          name  = "${var.name}-receiver"
          env {
            name  = "REDIS_ADDRESS"
            value = "redis-master.redis:6379"
          }
          env {
            name  = "APP_ID"
            value = "${var.namespace}-${var.name}"
          }
        }
      }
    }
  }
}
