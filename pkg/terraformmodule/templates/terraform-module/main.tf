locals {
  # use locals for commonly used values
  region = var.virtual_network.specs.aws.region

  # or to compute / transform varaibles
}

# resource "cloud_resource" "main" {
#   name = var.name
#   tags = var.tags
# }
