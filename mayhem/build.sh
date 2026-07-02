#!/usr/bin/env bash
#
# mayhem/build.sh — build teacup's go-fuzz harness as a sanitized libFuzzer binary
# (OSS-Fuzz Go path: go-fuzz-build -libfuzzer + clang link) plus a standalone reproducer.
#
# Runs inside the commit image (GO mayhem/Dockerfile) as `mayhem` in /mayhem.
# GOROOT/GOPATH/GOMODCACHE are pinned by the Dockerfile ENV (under /opt/toolchains —
# absolute, $HOME-independent).
#
# AIR-GAPPED CONTRACT (SPEC §6.5): the PATCH tier re-runs THIS script OFFLINE.
set -euo pipefail

[ -n "${SOURCE_DATE_EPOCH:-}" ] || unset SOURCE_DATE_EPOCH

: "${CC:=clang}" ; : "${CXX:=clang++}" ; : "${LIB_FUZZING_ENGINE:=-fsanitize=fuzzer}"
: "${SANITIZER_FLAGS=-fsanitize=address}"
: "${MAYHEM_JOBS:=$(nproc)}"
export CC CXX LIB_FUZZING_ENGINE SANITIZER_FLAGS MAYHEM_JOBS

export GOFLAGS="${GOFLAGS:--mod=mod}"
export GOPROXY="${GOPROXY:-file://$(go env GOMODCACHE)/cache/download,https://proxy.golang.org,direct}"

# §6.2 item 10: Go gc emits DWARF4, but go-fuzz links via clang++ C shims that land FIRST
# in the binary — force DWARF3 on those shims and the final link.
: "${GO_DEBUG_FLAGS:=-g -gdwarf-3}"
export CGO_CFLAGS="${CGO_CFLAGS:+$CGO_CFLAGS }$GO_DEBUG_FLAGS"
export CGO_CXXFLAGS="${CGO_CXXFLAGS:+$CGO_CXXFLAGS }$GO_DEBUG_FLAGS"
export GO_DEBUG_FLAGS

cd "$SRC"
go version

go get github.com/dvyukov/go-fuzz/go-fuzz-dep
go get github.com/AdaLogics/go-fuzz-headers

HARNESS_DIR="mayhem"
TARGET="fuzzteacup"

mkdir -p "$SRC/mayhem-build"
echo "=== building $TARGET (go-fuzz-build -libfuzzer) ==="
(
  cd "$SRC/$HARNESS_DIR"
  go-fuzz-build -libfuzzer -o "$SRC/mayhem-build/$TARGET.a"
)

echo "=== linking /mayhem/$TARGET (libFuzzer) ==="
$CXX $SANITIZER_FLAGS $LIB_FUZZING_ENGINE $GO_DEBUG_FLAGS \
  "$SRC/mayhem-build/$TARGET.a" -o "/mayhem/$TARGET"

echo "=== linking /mayhem/${TARGET}-standalone ==="
$CC $SANITIZER_FLAGS $GO_DEBUG_FLAGS -c "$STANDALONE_FUZZ_MAIN" -o "$SRC/mayhem-build/standalone_main.o"
$CXX $SANITIZER_FLAGS $GO_DEBUG_FLAGS \
  "$SRC/mayhem-build/standalone_main.o" "$SRC/mayhem-build/$TARGET.a" \
  -o "/mayhem/${TARGET}-standalone"

echo "build.sh complete"
