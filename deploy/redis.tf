resource "kubernetes_deployment" "redis" {
  metadata {
    name      = "${var.name}-redis"
    namespace = var.namespace
  }
  spec {
    replicas = "1"
    selector {
      match_labels = {
        "app" = "${var.name}-redis"
      }
    }
    template {
      metadata {
        labels = {
          "app" = "${var.name}-redis"
        }
      }
      spec {
        container {
          image = "redis:alpine"
          name  = "${var.name}-redis"
          port {
            container_port = 6379
          }
        }
      }
    }
  }
}
resource "kubernetes_service" "redis" {
  metadata {
    name      = "${var.name}-redis"
    namespace = var.namespace
  }
  spec {
    selector = {
      "app" = "${var.name}-redis"
    }
    port {
      protocol    = "TCP"
      port        = 6379
      target_port = 6379
    }
  }
}
