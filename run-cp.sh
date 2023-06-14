#!/usr/bin/env bash

if [[ -x "$(command -v go)" ]]; then
  if ! make build; then
    echo -e "\e[31m[FAILED] go build\e[0m"
    exit 1
  else
    echo -e "\e[32m[SUCCESS] go build\e[0m"
  fi
fi

# get all files that dont contain the word 'error.yail'
files=$(find ./examples_2 -type f -name '*.yail')
for file in $files; do
  if ! ./yail vm "$file" 2>/dev/null; then
    echo -e "\e[31m[ERROR] $file\e[0m"
  else
    echo -e "\e[32m[SUCCESS] $file\e[0m"
  fi
done
