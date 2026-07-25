#!/usr/bin/env bash

set -euo pipefail

SPEC="$(dirname "$0")/nxtools.spec"
DATE=$(date +"%a %b %d %Y")
VERSION=$(rpmspec -q --qf '%{version}\n' "$SPEC" | head -1)
REL=$(rpmspec -q --qf '%{release}\n' "$SPEC" | head -1)
MAINTAINER="Binary package builder <builder@famillegratton.net>"

LAST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || git rev-list --max-parents=0 HEAD)

# If HEAD itself is already tagged (e.g. the release was tagged before being
# built), fall back to the previous tag so LAST_TAG..HEAD is not empty.
if [[ "$(git rev-list -n1 "$LAST_TAG")" == "$(git rev-parse HEAD)" ]]; then
  LAST_TAG=$(git describe --tags --abbrev=0 "${LAST_TAG}^" 2>/dev/null || git rev-list --max-parents=0 HEAD)
fi

ENTRIES=$(git log --oneline "${LAST_TAG}..HEAD" | sed 's/^[a-f0-9]\+ /- /')

if [[ -z "$ENTRIES" ]]; then
  echo "No new commits since last tag, nothing to append."
  exit 0
fi

# Guard against duplicate entries
if grep -qF "${VERSION}-${REL}" "$SPEC"; then
  echo "Changelog entry for ${VERSION}-${REL} already present, skipping."
  exit 0
fi

NEW_ENTRY="* ${DATE} ${MAINTAINER} ${VERSION}-${REL}"$'\n'"${ENTRIES}"

awk -v entry="$NEW_ENTRY" '
  /^%changelog$/ { print; print entry; print ""; next }
  { print }
' "$SPEC" > "$SPEC.tmp" && mv "$SPEC.tmp" "$SPEC"

echo "Changelog updated for ${VERSION}-${REL}"
