variable "name" {}
variable "namespace" {}
variable "image" {}
variable "domain" {}
terraform {
  backend "gcs" {
    bucket = "tf-state-sender-receiver"
  }
}
provider "kubernetes" {
  config_path = "~/.kube/config"
}
