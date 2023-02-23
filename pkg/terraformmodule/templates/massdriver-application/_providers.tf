terraform {
  required_version = ">= 1.0"
  required_providers {
    # Used in the massdriver-application module
    # https://registry.terraform.io/providers/massdriver-cloud/mdxc/latest/docs
    mdxc = {
      source  = "massdriver-cloud/mdxc"
      version = "~> 0.10"
    }
    # Used in the massdriver-application module
    # https://registry.terraform.io/providers/massdriver-cloud/jq/latest/docs
    jq = {
      source  = "massdriver-cloud/jq"
      version = "~> 0.2"
    }
    # Useful if you do anything with CIDR ranges
    # https://registry.terraform.io/providers/massdriver-cloud/utility/latest/docs/resources/available_cidr
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
