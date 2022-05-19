# {{ .Name }}

{{ .Description }}

## Development
### Enabling Pre-commit

This repo includes Terraform pre-commit hooks.

[Install precommmit](https://pre-commit.com/index.html#installation) on your system.

```shell
git init
pre-commit install
```

Terraform hooks will now be run on each commit.
### Configuring a bundle

`massdriver.yaml` ...

Build the bundle locally:

```shell
mass bundle build
```
