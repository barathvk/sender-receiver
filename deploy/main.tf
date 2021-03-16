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
