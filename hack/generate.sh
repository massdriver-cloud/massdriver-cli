#! bash

# Generate a test terraform bundle to fully exercise generation.

set -e

pushd /tmp
rm -rf aws-pubsub-topic
mkdir aws-pubsub-topic
cd aws-pubsub-topic
mass bundle generate
git init
pre-commit install
git add src
mass bundle build
git commit . -m 'wip'
popd
