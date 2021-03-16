resource "kubernetes_deployment" "receiver" {
  metadata {
    name      = var.name
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
          name  = var.name
          args  = ["--sender"]
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
