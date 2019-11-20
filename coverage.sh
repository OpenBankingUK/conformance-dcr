MIN_COVERAGE_PERCENT=69

calc(){ awk "BEGIN { print "$*" }"; }

TOTAL_COVERAGE=$(go test ./... --cover | awk '{if ($1 != "?") print $5; else print "0.0";}' | sed 's/\%//g' | awk '{s+=$1} END {printf "%.2f\n", s}')
NR_PACKAGES=$(go test ./... --cover | wc -l)

COVERAGE_PERCENT=$(calc "$TOTAL_COVERAGE"/"$NR_PACKAGES")

if (( $(echo "$COVERAGE_PERCENT < $MIN_COVERAGE_PERCENT" | bc -l) )); then
  echo "Code coverage bellow minimum $MIN_COVERAGE_PERCENT%: $COVERAGE_PERCENT%"
  exit 1
fi

echo "Code coverage $COVERAGE_PERCENT%"
