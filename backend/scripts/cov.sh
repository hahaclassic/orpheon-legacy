#!/bin/bash

mkdir -p coverage

# 1. Full project
go test -covermode=atomic -coverprofile=./coverage/coverage_total.out ./internal/... ./pkg/... > coverage/result_total.log 2>&1

# 2. Only domain
go test -covermode=atomic -coverprofile=./coverage/coverage_domain.out ./internal/domain/... > coverage/result_domain.log 2>&1

print_coverage() {
    file=$1
    log=$2
    if [ $? -eq 0 ]; then
        if [[ "$3" == "-v" ]]; then
            cat "$log"
        fi
        cov=$(go tool cover -func="$file" | grep "total:" | awk '{print $3}')
        echo "$4 coverage: $cov"
    else
        echo "$4 coverage: FAIL"
        cat "$log"
        exit 1
    fi
}

print_coverage "./coverage/coverage_total.out" "coverage/result_total.log" "$1" "Total project"
print_coverage "./coverage/coverage_domain.out" "coverage/result_domain.log" "$1" "Business layer"
