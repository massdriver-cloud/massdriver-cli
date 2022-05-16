
resource "massdriver_artifact" "<name>" {
  artifact = jsonencode(
    {
      # # Top-level is massdriver metadata
      # #
      # # The cloud provider's full resource ID. 
      # # This is used as an idempotent key to ensure we are updating existing artifacts and no creating
      # # new ones. This key should be unique to a project. It will be combined with
      # # project_id in massdriver to ensure uniqueness.
      # # Examples: AWS _full_ ARN, Kubernetes selfLink
      # metadata = {
      #   field = "the field in the artifacts schema"
      #   provider_resource_id = "AWS ARN or K8S SelfLink"
      #   #
      #   # The artifact type this creates. 
      #   # You will need to define the schema the first time in ./definitions/artifacts
      #   type = "file-name-from-artifacts",
      #   name = "a contextual name for the user"
      # }
      # data = {
      #   # This should match the aws-rds-arn.json schema file
      #   arn = "aws::..."
      # }
      # specs = {
      #   # Any existing spec in ./specs
      #   # aws = {}
      # }
    }
  )
}
