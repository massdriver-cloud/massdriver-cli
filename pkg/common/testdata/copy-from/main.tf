resource "mdxc_application_identity" "main" {
  name = var.name

  gcp_configuration   = local.is_gcp ? local.gcp_identity : null
  azure_configuration = local.is_azure ? local.azure_identity : null
  aws_configuration   = local.is_aws ? local.aws_identity : null
}
