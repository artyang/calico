#! /bin/bash

export REPO_DIR="$(pwd)"
mkdir artifacts
echo $DOCKERHUB_PASSWORD | docker login --username "$DOCKERHUB_USERNAME" --password-stdin

# Generate ~/.rnd file
cd ~
openssl rand -writerand .rnd
cd -

export CI_EXIT_CODE=0