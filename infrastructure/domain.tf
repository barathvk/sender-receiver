data "cloudflare_zones" "main" {
  filter {
    name = var.domain
  }
}

resource "cloudflare_zone_settings_override" "domain_settings" {
  zone_id = data.cloudflare_zones.main.zones.0.id
  settings {
    brotli                   = "on"
    security_level           = "medium"
    always_use_https         = "on"
    opportunistic_encryption = "on"
    automatic_https_rewrites = "on"
    minify {
      css  = "on"
      js   = "on"
      html = "on"
    }
    security_header {
      enabled = true
    }
  }
}
