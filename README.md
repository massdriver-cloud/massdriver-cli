# Massdriver CLI

A command line interface for managing Massdriver bundles and applications.

## Creating a new bundle

Note: --template-dir will be embedded in the near future.

```shell
mass bundle generate --template-dir path/to/generators/bundle/terraform
cd YOUR_BUNDLE
mass bundle build massdriver.yaml
````
