resource "random_pet" "server" {
  keepers = {
    app_name = "copies-successfully"
  }
}
