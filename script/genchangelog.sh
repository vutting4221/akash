#!/usr/bin/env bash

# generated changelog depends on the tag type
# only SEMVER tags are accounted
# there are two type of tag distinguished by minor part of the semver:
#   - even number: mainnet
#   - odd number: testnet
# net detection is done in section s1

# there are two type release notes generated
#   - prerelease: changelog between current and nearest lower prerelease (or previous release)
#     for example current tag v0.1.1-rc.10 and previous was v0.1.1-rc.9, so changelog is generated between
#   - release: changelog between current and previous release tags
# mainnet status is taken care as well. if current tag is edgenet (e.g. v0.1.1-rc.10) it
# will be generated to edgenet changes only

PATH=$PATH:$(pwd)/.cache/bin
export PATH=$PATH

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"

if [[ $# -ne 1 ]]; then
	echo "illegal number of parameters"
	exit 1
fi

curr_tag=$1

# s1
is_mainnet=$("${SCRIPT_DIR}"/mainnet-from-tag.sh "$curr_tag")
if [[ $is_mainnet == false ]]; then
	version_rel="^[v|V]?(0|[1-9][0-9]*)\\.(\\d*[13579])\\.(0|[1-9][0-9]*)$"
	version_prerel="^[v|V]?(0|[1-9][0-9]*)\\.(\\d*[13579])\\.(0|[1-9][0-9]*)(\\-[0-9A-Za-z-]+(\\.[0-9A-Za-z-]+)*)?(\\+[0-9A-Za-z-]+(\\.[0-9A-Za-z-]+)*)?$"
else
	version_rel="^[v|V]?(0|[1-9][0-9]*)\.(\d*[02468])\.(0|[1-9][0-9]*)$"
	version_prerel="^[v|V]?(0|[1-9][0-9]*)\.(\d*[02468])\.(0|[1-9][0-9]*)(\-[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?(\+[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?$"
fi

# s2
if [[ -z $("${SCRIPT_DIR}"/semver.sh get prerel "$curr_tag") ]]; then
	tag_regexp=$version_rel
else
	tag_regexp=$version_prerel
fi

prev_tag=

# it sucks, slow, but works
# shellcheck disable=SC2046
for tag in $(git tag --sort=-version:refname); do
	if [[ "$tag" =~ $tag_regexp ]] && [[ $("${SCRIPT_DIR}"/semver.sh compare "$tag" "$curr_tag") -eq -1 ]]; then
		prev_tag=$tag
		break
	fi
done

echo "$prev_tag"
if [[ -z $prev_tag ]]; then
	echo "couldn't detect previous revision"
	exit 1
fi

git-chglog --tag-filter-pattern="$version_prerel" "$prev_tag..$curr_tag"
