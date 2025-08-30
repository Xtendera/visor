#!/usr/bin/env bash
set -euo pipefail

APP_NAME="visor"
OUTPUT_DIR="build"

mkdir -p "$OUTPUT_DIR"

TARGETS=(
  "linux/amd64"
  "linux/arm64"
  "linux/arm"
  "windows/amd64"
  "windows/arm64"
  "darwin/amd64"
  "darwin/arm64"
  "freebsd/amd64"
)

echo "Compiling $APP_NAME for ${#TARGETS[@]} targets..."
for target in "${TARGETS[@]}"; do
  GOOS="${target%/*}"
  GOARCH="${target#*/}"
  EXT=""
  if [ "$GOOS" = "windows" ]; then
    EXT=".exe"
  fi

  echo " → $GOOS/$GOARCH"
  CGO_ENABLED=0 GOOS="$GOOS" GOARCH="$GOARCH" \
    go build -o "$OUTPUT_DIR/${APP_NAME}-${GOOS}-${GOARCH}${EXT}" .
done

echo "All builds complete. Binaries are in '$OUTPUT_DIR/'"
