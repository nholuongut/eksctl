#!/bin/sh -ex

if [ -z "${CIRCLE_PULL_REQUEST}" ] && [ -n "${CIRCLE_TAG}" ] && [ "${CIRCLE_PROJECT_USERNAME}" = "nholuongut" ] ; then
  export RELEASE_DESCRIPTION="${CIRCLE_TAG} (permalink)"
  RELEASE_NOTES_FILE="docs/release_notes/${CIRCLE_TAG}.md"

  if [[ ! -f "${RELEASE_NOTES_FILE}" ]]; then
    echo "Release notes ${RELEASE_NOTES_FILE} not found. Exiting..."
    return
  fi

  cat ./.goreleaser.yml ./.goreleaser.brew.yml > .goreleaser.brew.combined.yml
  goreleaser release --skip-validate --config=./.goreleaser.brew.combined.yml --release-notes="${RELEASE_NOTES_FILE}"

  sleep 90 # GitHub API resolves the time to the nearest minute, so in order to control the sorting oder we need this

  git tag --delete "${CIRCLE_TAG}"
  git tag --force latest_release

  if github-release info --user nholuongut --repo eksctl --tag latest_release > /dev/null 2>&1 ; then
    github-release delete --user nholuongut --repo eksctl --tag latest_release
  fi

  export RELEASE_DESCRIPTION="${CIRCLE_TAG}"
  goreleaser release --skip-validate --rm-dist --config=./.goreleaser.yml --release-notes="${RELEASE_NOTES_FILE}"

else
  echo "Not a tag release, skip publish"
fi
