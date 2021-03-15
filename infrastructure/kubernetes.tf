data "google_client_config" "default" {}
resource "kubernetes_namespace" "redis" {
  metadata {
    name = "redis"
  }
}
resource "kubernetes_namespace" "ingress_nginx" {
  metadata {
    name = "ingress-nginx"
  }
}

resource "helm_release" "redis" {
  repository = "https://charts.bitnami.com/bitnami"
  chart      = "redis"
  name       = "redis"
  namespace  = kubernetes_namespace.redis.metadata.0.name
}

resource "helm_release" "ingress_nginx" {
  repository = "https://kubernetes.github.io/ingress-nginx"
  chart      = "ingress-nginx"
  name       = "ingress-nginx"
  namespace  = kubernetes_namespace.ingress_nginx.metadata.0.name
}
resource "helm_release" "external_dns" {
  repository = "https://charts.bitnami.com/bitnami"
  chart      = "external-dns"
  name       = "external-dns"
  namespace  = kubernetes_namespace.ingress_nginx.metadata.0.name
  values = [
    <<EOF
    sources:
      - ingress
    provider: cloudflare
    cloudflare:
      email: ${var.cloudflare_email}
      apiKey: ${var.cloudflare_api_key}
      proxied: true
    domainFilters:
      - ${var.domain}
    txtOwnerId: ${var.name}
    logLevel: debug
    policy: sync
    EOF
  ]
}
