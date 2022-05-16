# Massdriver CLI

A command line interface for managing Massdriver bundles and applications.

## Usage

### Creating a new bundle

Note: --template-dir will be embedded in the near future.

```shell
mass bundle generate --template-dir path/to/generators/bundle/terraform
cd YOUR_BUNDLE
mass bundle build
```

## Development

### Precommit

This repo uses precommit. Don't skip this step.

1. [Install precommmit](https://pre-commit.com/index.html#installation) on your system.
2. Run `pre-commit install` to run hooks on `git commit`
