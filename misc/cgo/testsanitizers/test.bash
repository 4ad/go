#!/usr/bin/env bash
# Copyright 2015 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# This directory is intended to test the use of Go with sanitizers
# like msan, asan, etc.  See https://github.com/google/sanitizers .

set -e

# The sanitizers were originally developed with clang, so prefer it.
CC=cc
if test -x "$(type -p clang)"; then
  CC=clang
fi
export CC

msan=yes

TMPDIR=${TMPDIR:-/tmp}
echo 'int main() { return 0; }' > ${TMPDIR}/testsanitizers$$.c
if $CC -fsanitize=memory -c ${TMPDIR}/testsanitizers$$.c -o ${TMPDIR}/testsanitizers$$.o 2>&1 | grep "unrecognized" >& /dev/null; then
  echo "skipping msan tests: -fsanitize=memory not supported"
  msan=no
fi
rm -f ${TMPDIR}/testsanitizers$$.*

tsan=yes

# The memory and thread sanitizers in versions of clang before 3.6
# don't work with Go.
if test "$msan" = "yes" && $CC --version | grep clang >& /dev/null; then
  ver=$($CC --version | sed -e 's/.* version \([0-9.-]*\).*/\1/')
  major=$(echo $ver | sed -e 's/\([0-9]*\).*/\1/')
  minor=$(echo $ver | sed -e 's/[0-9]*\.\([0-9]*\).*/\1/')
  if test "$major" -lt 3 || test "$major" -eq 3 -a "$minor" -lt 6; then
    echo "skipping msan/tsan tests: clang version $major.$minor (older than 3.6)"
    msan=no
    tsan=no
  fi

  # Clang before 3.8 does not work with Linux at or after 4.1.
  # golang.org/issue/12898.
  if test "$msan" = "yes" -a "$major" -lt 3 || test "$major" -eq 3 -a "$minor" -lt 8; then
    if test "$(uname)" = Linux; then
      linuxver=$(uname -r)
      linuxmajor=$(echo $linuxver | sed -e 's/\([0-9]*\).*/\1/')
      linuxminor=$(echo $linuxver | sed -e 's/[0-9]*\.\([0-9]*\).*/\1/')
      if test "$linuxmajor" -gt 4 || test "$linuxmajor" -eq 4 -a "$linuxminor" -ge 1; then
        echo "skipping msan/tsan tests: clang version $major.$minor (older than 3.8) incompatible with linux version $linuxmajor.$linuxminor (4.1 or newer)"
	msan=no
	tsan=no
      fi
    fi
  fi
fi

status=0

if test "$msan" = "yes"; then
    if ! go build -msan std; then
	echo "FAIL: build -msan std"
	status=1
    fi

    if ! go run -msan msan.go; then
	echo "FAIL: msan"
	status=1
    fi

    if ! CGO_LDFLAGS="-fsanitize=memory" CGO_CPPFLAGS="-fsanitize=memory" go run -msan -a msan2.go; then
	echo "FAIL: msan2 with -fsanitize=memory"
	status=1
    fi

    if ! go run -msan -a msan2.go; then
	echo "FAIL: msan2"
	status=1
    fi

    if ! go run -msan msan3.go; then
	echo "FAIL: msan3"
	status=1
    fi

    if ! go run -msan msan4.go; then
	echo "FAIL: msan4"
	status=1
    fi

    if go run -msan msan_fail.go 2>/dev/null; then
	echo "FAIL: msan_fail"
	status=1
    fi
fi

if test "$tsan" = "yes"; then
    echo 'int main() { return 0; }' > ${TMPDIR}/testsanitizers$$.c
    ok=yes
    if ! $CC -fsanitize=thread ${TMPDIR}/testsanitizers$$.c -o ${TMPDIR}/testsanitizers$$ &> ${TMPDIR}/testsanitizers$$.err; then
	ok=no
    fi
     if grep "unrecognized" ${TMPDIR}/testsanitizers$$.err >& /dev/null; then
	echo "skipping tsan tests: -fsanitize=thread not supported"
	tsan=no
     elif test "$ok" != "yes"; then
	 cat ${TMPDIR}/testsanitizers$$.err
	 echo "skipping tsan tests: -fsanitizer=thread build failed"
	 tsan=no
     fi
     rm -f ${TMPDIR}/testsanitizers$$*
fi

if test "$tsan" = "yes"; then
    err=${TMPDIR}/tsanerr$$.out

    if ! go run tsan.go 2>$err; then
	cat $err
	echo "FAIL: tsan"
	status=1
    elif grep -i warning $err >/dev/null 2>&1; then
	cat $err
	echo "FAIL: tsan"
	status=1
    fi

    if ! go run tsan2.go 2>$err; then
	cat $err
	echo "FAIL: tsan2"
	status=1
    elif grep -i warning $err >/dev/null 2>&1; then
	cat $err
	echo "FAIL: tsan2"
	status=1
    fi

    rm -f $err
fi

exit $status
