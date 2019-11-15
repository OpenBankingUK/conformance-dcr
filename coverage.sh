#!/usr/bin/env bash
set -eu -o pipefail

set -o noclobber    # Avoid overlay files (echo "hi" > foo)
set -o errexit      # Used to exit upon error, avoiding cascading errors
set -o pipefail     # Unveils hidden failures
set -o nounset      # Exposes unset variables

# set -o nullglob     # Non-matching globs are removed  ('*.foo' => '')
shopt -s nullglob   # Non-matching globs are removed  ('*.foo' => '')
# set -o failglob     # Non-matching globs throw errors
shopt -s failglob   # Non-matching globs throw errors
# set -o nocaseglob   # Case insensitive globs
shopt -s nocaseglob # Case insensitive globs
# set -o globstar     # Allow ** for recursive matches ('lib/**/*.rb' => 'lib/a/b/c.rb')
shopt -s globstar # Allow ** for recursive matches ('lib/**/*.rb' => 'lib/a/b/c.rb')

MIN_COVERAGE_PERCENT=71

calc(){ awk "BEGIN { print "$*" }"; }

TOTAL_COVERAGE=$(go test ./... --cover | awk '{if ($1 != "?") print $5; else print "0.0";}' | sed 's/\%//g' | awk '{s+=$1} END {printf "%.2f\n", s}')
NR_PACKAGES=$(go test ./... --cover | wc -l)
COVERAGE_PERCENT=$(calc "$TOTAL_COVERAGE"/"$NR_PACKAGES")

if (( $(echo "$COVERAGE_PERCENT < $MIN_COVERAGE_PERCENT" | bc -l) )); then
  echo "Code coverage bellow minimum $MIN_COVERAGE_PERCENT%: $COVERAGE_PERCENT%"
  exit 1
fi

echo "Code coverage $COVERAGE_PERCENT%"
