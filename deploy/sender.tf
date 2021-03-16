resource "kubernetes_deployment" "sender" {
  metadata {
    name      = "${var.name}-sender"
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
          name  = "${var.name}-sender"
          args  = ["--sender"]
          env {
            name  = "REDIS_ADDRESS"
            value = "${var.name}-redis:6379"
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
