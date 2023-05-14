
variable "md_metadata" {
  # TODO: does TF care if this has other keys as well?
  type = object({
    name_prefix = string
    default_tags = any
  })
}
