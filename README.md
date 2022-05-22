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

Build the bundle locally:

```shell
mass bundle build
```

### Configuring a bundle

TBD
