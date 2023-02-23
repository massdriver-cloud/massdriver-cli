terraform {
  required_version = ">= 1.0"
  required_providers {
    massdriver = {
      source  = "massdriver-cloud/massdriver"
      version = "~> 1.0"
    }
    # utility = {
    #   source = "massdriver-cloud/utility"
    # }
    # aws = {
    #   source  = "hashicorp/aws"
    #   version = "~> 4.0"
    # }
    # azurerm = {
    #   source  = "hashicorp/azurerm"
    #   version = "~> 3.0"
    # }
    # google = {
    #   source  = "hashicorp/google"
    #   version = "~> 4.0"
    # }
    # google-beta = {
    #   source  = "hashicorp/google-beta"
    #   version = "~> 4.0"
    # }
    # helm = {
    #   source  = "hashicorp/helm"
    #   version = "~> 2.0"
    # }
    # kubernetes = {
    #   source  = "hashicorp/kubernetes"
    #   version = "~> 2.0"
    # }
  }
}
