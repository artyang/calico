#! /bin/bash

echo "Uploading Artifacts from Global Epilogue"
# No set -e so that we continue to collect artifacts even if some fail.
set -x

# set nullglob so that '*' is not used literally if the artifacts dir is empty
shopt -s nullglob

cd "${REPO_DIR}"

my_dir=$(dirname "$0")
repo_dir=$my_dir/../..

cd "$repo_dir/artifacts" || exit 1

for file in *; do
  echo "Uploading file to Artifact storage - '${file}'"
  gsutil cp "$file" "${CI_ARTIFACT_STEP_STORAGE}/$file"
  if [[ ${file} =~ .*\.xml ]]; then
    echo "File is of type XML - Converting to HTML Test results and uploading"
    if ! command -v npm &> /dev/null; then
      curl -sL https://deb.nodesource.com/setup_16.x | sudo -E bash -
      apt-get install -y nodejs
    fi
    if ! command -v xunit-viewer &> /dev/null; then
      npm i -g xunit-viewer
    fi
    xunit-viewer -r ${file} -o "${file}.html"
    gsutil cp "${file}.html" "${CI_ARTIFACT_STEP_STORAGE}/${file}.html"
  fi
done