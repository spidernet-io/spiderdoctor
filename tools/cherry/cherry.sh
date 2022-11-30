#!/bin/bash
set -e

#require_linux

ORG=${ORG:-"spidernet-io"}
REPO=${REPO:-"spiderdoctor"}
CHERRY_FROM_BRANCH=${CHERRY_FROM_BRANCH:-""}

[ -z "${ORG}" ] && echo "error, miss ORG" && exit 1
[ -z "${REPO}" ] && echo "error, miss REPO" && exit 1
[ -z "${CHERRY_FROM_BRANCH}" ] && echo "error, miss CHERRY_FROM_BRANCH" && exit 1


#==============================================

get_remote () {
  local remote
  local org=${1:-cilium}
  local repo=${2:-cilium}
  remote=$(git remote -v | \
    grep "github.com[/:]${org}/${repo}" | \
    head -n1 | cut -f1)
  if [ -z "$remote" ]; then
      echo "No remote git@github.com:${org}/${repo}.git or https://github.com/${org}/${repo} found" 1>&2
      return 1
  fi
  echo "$remote"
}

commit_in_upstream() {
    local commit="$1"
    local branch="$2"
    local org="${3:-""}"
    local repo="${4:-""}"
    local remote="$(get_remote ${org} ${repo})"
    local branches="$(git branch -q -r --contains $commit $remote/$branch 2> /dev/null)"
    echo "$branches" | grep -q ".*$remote/$branch"
}
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
