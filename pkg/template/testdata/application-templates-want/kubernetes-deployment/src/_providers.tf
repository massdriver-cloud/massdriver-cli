




terraform {
  required_version = ">= 1.0"
  required_providers {
    massdriver = {
      source = "massdriver-cloud/massdriver"
    }
    jq = {
      source = "massdriver-cloud/jq"
    }
    google = {
      source = "hashicorp/google"
    }
    helm = {
      source = "hashicorp/helm"
    }
  }
}

locals {
  gcp_authentication = module.k8s_application.connections.gcp_authentication
  kubernetes_cluster = module.k8s_application.connections.kubernetes_cluster
  gcp_region         = split("/", local.kubernetes_cluster.data.infrastructure.grn)[3]
  gcp_project_id     = local.gcp_authentication.data.project_id

  k8s_host                  = local.kubernetes_cluster.data.authentication.cluster.server
  k8s_certificate_authority = base64decode(local.kubernetes_cluster.data.authentication.cluster.certificate-authority-data)
  k8s_token                 = local.kubernetes_cluster.data.authentication.user.token
}

provider "google" {
  project     = local.gcp_project_id
  credentials = jsonencode(local.gcp_authentication.data)
  region      = local.gcp_region
}

provider "helm" {
  kubernetes {
    host                   = local.k8s_host
    cluster_ca_certificate = local.k8s_certificate_authority
    token                  = local.k8s_token
  }
}

