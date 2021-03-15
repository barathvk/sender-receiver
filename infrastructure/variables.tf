variable "project" {}
variable "cloudflare_email" {}
variable "cloudflare_api_key" {}
variable "name" {
  default = "sender-receiver"
}

variable "location" {
  default = "europe-west3-a"
}

variable "node_count" {
  default = 1
}

variable "domain" {
  default = "barathvk.dev"
}
