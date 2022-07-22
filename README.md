# Massdriver CLI

A command line interface for managing Massdriver bundles and applications.

## Using the CLI

### Writing Your Own Bundle

The fastest way to get started is to use the _massdriver-cli_ to generate a new bundle. To do this, run `mass bundle generate`. You can read more about bundle development in our [docs](https://docs-tawny-nu.vercel.app/bundles/walk-through).

### Generate An Application

To use Massdriver to manage your application, we'll use the _massdriver-cli_ again. Simply run `mass app generate`. You'll be able to choose an [application template](https://github.com/massdriver-cloud/application-templates) and configure it to your needs. You can find application docs [here](https://docs-tawny-nu.vercel.app/applications).

## Contributing

### Precommit

This repo uses precommit. Don't skip this step.

1. [Install precommmit](https://pre-commit.com/index.html#installation) on your system.
2. Run `pre-commit install` to run hooks on `git commit`
