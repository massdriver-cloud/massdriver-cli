resource "random_pet" "server" {
  keepers = {
    app_name = "renders-dependencies"
  }
}
