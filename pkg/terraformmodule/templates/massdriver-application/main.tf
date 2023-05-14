locals {
  # use locals for commonly used values
  # region = var.virtual_network.specs.aws.region

  # or to compute / transform variables
  # id_set = toset(flatten([for val in var.things : val.id]))
}

module "application" {
  source  = "github.com/massdriver-cloud/terraform-modules//massdriver-application?ref=fc5f7b1"
  name    = var.md_metadata.name_prefix
  # this can be one of many values, the complete list is here:
  # https://github.com/massdriver-cloud/terraform-modules/blob/main/massdriver-application/main.tf#L32-L35
  service = "function"
}
