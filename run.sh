#!/usr/bin/env bash

if [[ -x "$(command -v go)" ]]; then
  if ! go build ./main.go; then
    echo -e "\e[31m[FAILED] go build\e[0m"
    exit 1
  else
    echo -e "\e[32m[SUCCESS] go build\e[0m"
  fi
fi

for file in ./examples/*.yail; do

  if ! ./main "$file"; then
    echo -e "\e[31m[FAILED] $file\e[0m"
    exit 1
  else
    echo -e "\e[32m[SUCCESS] $file\e[0m"
  fi
done
