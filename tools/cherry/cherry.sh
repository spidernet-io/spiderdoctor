#!/bin/bash
set -e

SCRIPT_DIR="$( cd "$( dirname "$0" )" && pwd )"
source ${SCRIPT_DIR}/common.sh

require_linux

ORG=${ORG:-"spidernet-io"}
REPO=${REPO:-"spiderdoctor"}
CHERRY_FROM_BRANCH=${CHERRY_FROM_BRANCH:-""}

[ -z "${ORG}" ] && echo "error, miss ORG" && exit 1
[ -z "${REPO}" ] && echo "error, miss REPO" && exit 1
[ -z "${CHERRY_FROM_BRANCH}" ] && echo "error, miss CHERRY_FROM_BRANCH" && exit 1




cleanup () {
  if [ -n "$TMPF" ]; then
    rm $TMPF
  fi
}

trap cleanup EXIT

cherry_pick () {
  CID=$1
  if ! commit_in_upstream "$CID" "$CHERRY_FROM_BRANCH" "${ORG}" "${REPO}"; then
    echo "Commit $CID not in $REM/$CHERRY_FROM_BRANCH!"
    exit 1
  fi
  TMPF=`mktemp cp.XXXXXX`
  FROM=`git show --pretty=email $CID | head -n 2 | grep "From: "`
  FULL_ID=`git show $CID | head -n 1 | cut -f 2 -d ' '`
  git format-patch -1 $FULL_ID --stdout | sed '/^$/Q' > $TMPF
  echo "" >> $TMPF
  echo "[ upstream commit $FULL_ID ]" >> $TMPF
  git format-patch -1 $FULL_ID --stdout | sed -n '/^$/,$p' >> $TMPF
  echo "Applying: $(git log -1 --oneline $FULL_ID)"
  git am --quiet -3 --signoff $TMPF
}

main () {
  REM="$(get_remote "${ORG}" "${REPO}")"
  for CID in "$@"; do
    cherry_pick "$CID"
  done
}

if [ $# -lt 1 ]; then
  echo "Usage: $0 <commit-id> [commit-id ...]"
  exit 1
fi

main "$@"
