variable "name" {}
variable "namespace" {}
variable "image" {}
variable "domain" {}
variable "port" {
  default = 8080
}
provider "kubernetes" {
  config_path = "~/.kube/config"
}
resource "kubernetes_deployment" "sender_receiver" {
  metadata {
    name      = var.name
    namespace = var.namespace
  }
  spec {
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
          env {
            name  = "REDIS_ADDRESS"
            value = "redis-master.redis:6379"
          }
          env {
            name  = "APP_ID"
            value = "${var.namespace}-${var.name}"
          }
          env {
            name  = "PORT"
            value = var.port
          }
          port {
            container_port = var.port
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "sender_receiver" {
  metadata {
    name      = var.name
    namespace = var.namespace
  }
  spec {
    selector = {
      "app" = var.name
    }
    port {
      protocol    = "TCP"
      port        = var.port
      target_port = var.port
    }
  }
}

resource "kubernetes_ingress" "sender-receiver" {
  metadata {
    name      = var.name
    namespace = var.namespace
    annotations = {
      "nginx.ingress.kubernetes.io/rewrite-target" = "/$2"
    }
  }
  spec {
    rule {
      host = "api.${var.domain}"
      http {
        path {
          path = "/${var.name}(/|$)(.*)"
          backend {
            service_name = var.name
            service_port = var.port
          }
        }
      }
    }
  }
}
