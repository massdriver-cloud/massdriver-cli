
# resource "massdriver_artifact" "<name>" {
#   field                = "the field in the artifacts schema"
#   provider_resource_id = "AWS ARN or K8S SelfLink"
#   type                 = "file-name-from-artifacts"
#   name                 = "a contextual name for the artifact"
#   artifact = jsonencode(
#     {
#       # # Top-level is massdriver metadata
#       # #
#       # # The cloud provider's full resource ID. 
#       # # This is used as an idempotent key to ensure we are updating existing artifacts and no creating
#       # # new ones. This key should be unique to a project. It will be combined with
#       # # project_id in massdriver to ensure uniqueness.
#       # # Examples: AWS _full_ ARN, Kubernetes selfLink
#       # data = {
#       #   # This should match the aws-rds-arn.json schema file
#       #   arn = "aws::..."
#       # }
#       # specs = {
#       #   # Any existing spec in ./specs
#       #   # aws = {}
#       # }
#     }
#   )
# }
