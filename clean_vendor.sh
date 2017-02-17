#!/bin/bash
# Copyright (c) 2017 Arista Networks, Inc.  All rights reserved.
# Use of this source code is governed by the Apache License 2.0
# that can be found in the LICENSE file.

# This script removes unneeded files from vendored dependencies.

case $# in
  (0)
    targetdir=vendor
    ;;
  (1)
    targetdir=$1
    ;;
  (*)
    echo >&2 'usage: $0 [directory]'
    exit 1
    ;;
esac

case `uname -s` in
  (Darwin)
    DARWIN_FIND_FLAGS='-E'
    ;;
  (Linux)
    LINUX_FIND_FLAGS='-regextype posix-extended'
    ;;
  (*)
    echo >&2 'unsupported platform'
    exit 1
    ;;
esac

find "$targetdir" \( \
  ! -path 'vendor/github.com/elastic/beats/*' \
  -a \( \
    -name .travis.yml \
    -o -name 'README*' \
    -o -name '*.pdf' \
    -o -name '*.md' \
    -o -name '*.p[ly]' \
    -o -name Makefile \
  \) \
  -a ! \( \
    -name 'AUTHORS*' \
    -o -name 'CONTRIBUTORS*' \
    -o -name 'COPY*' \
    -o -name 'LICEN*' \
    -o -name 'NOTICE*' \
  \) \) -print0 | xargs -0 rm -v
