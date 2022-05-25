# Massdriver CLI

A command line interface for managing Massdriver bundles and applications.

## Precommit

This repo uses precommit. Don't skip this step.

1. [Install precommmit](https://pre-commit.com/index.html#installation) on your system.
2. Run `pre-commit install` to run hooks on `git commit`

## Developing a bundle

### Creating a new bundle

Generate a new bundle:

```shell
mass bundle generate -o ./my-bundle
```

### Params & Connections

TBD

### Building a bundle

Next you'll need to build your bundle. This will convert your params and connections into terraform variables for local development. This is also run during CI/CD to publish your bundle to Massdriver.

```shell
mass bundle build
```
