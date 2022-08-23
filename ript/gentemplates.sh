#!/bin/bash -e

script_dir="$( cd "$(dirname "$( readlink -f "${BASH_SOURCE[0]}" )" )" && pwd )"

errcho() { echo "$@" 1>&2; }
die() { errcho "$@"; exit 1; }

#
# TODO: use git ls-files to list the files that go into the tar file.
#

[[ $script_dir == `pwd` ]] || die "Must run from dir that holds this script ($script_dir)"

cd templates

template_sub_dirs="$(find ./ -maxdepth 1 -mindepth 1 -type d)"
for d in $template_sub_dirs; do
  name="$(basename "$d")"

  echo ""
  echo "----------"
  echo tar -cvf "${name}2.tar" "${name}"
  tar -cvf "${name}.tar" "${name}"
done


