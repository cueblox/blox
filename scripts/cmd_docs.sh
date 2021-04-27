#!/bin/sh
set -e

SED="sed"
if which gsed >/dev/null 2>&1; then
	SED="gsed"
fi

rm -rf dogfood/sites/docs/docs/cmd/
mkdir -p dogfood/sites/docs/docs/cmd/
go run ./cmd/blox docs
"$SED" \
	-i'' \
	-e 's/SEE ALSO/See also/g' \
	-e 's/^## /# /g' \
	-e 's/^### /## /g' \
	-e 's/^#### /### /g' \
	-e 's/^##### /#### /g' \
	./dogfood/sites/docs/docs/docs/cmd/*.md
