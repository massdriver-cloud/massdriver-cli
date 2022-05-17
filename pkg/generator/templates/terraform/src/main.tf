resource "random_pet" "name" {
  keepers = {
    # An example resource w/ JSON Schema input
    ami_id = "${var.name}"
  }
}
