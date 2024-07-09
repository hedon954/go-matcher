#!/bin/bash

# Find all .go files excluding vendor directory, main.go, and *_test.go files
go_files=$(find . -name '*.go' -not -path "./vendor/*" -not -path "./merr/*" -not -path "./pto/*"  -not -path "./enum/*" -not -name 'main.go' -not -name '*_test.go' -not -name 'doc.go' -not -name 'interface.go')

missing_tests=0

for go_file in $go_files; do
  test_file="${go_file%.go}_test.go"
  if [[ ! -f $test_file ]]; then
    echo "Missing test file for $go_file"
    missing_tests=1
  fi
done

exit $missing_tests
