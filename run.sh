#!/usr/bin/env bash

if [[ -x "$(command -v go)" ]]; then
  if ! go build ./main.go; then
    echo -e "\e[31m[FAILED] go build\e[0m"
    exit 1
  else
    echo -e "\e[32m[SUCCESS] go build\e[0m"
  fi
fi

echo -e "\n\e[33m[INFO] Running tests that should pass\e[0m\n"
# get all files that dont contain the word 'error.yail'
files=$(find ./examples -type f -name '*.yail' -not -name '*error.yail')
for file in $files; do
  if ! ./main "ast $file" 2>/dev/null; then
    echo -e "\e[31m[ERROR] $file\e[0m -- (It's supposed to fail)"
  else
    echo -e "\e[32m[SUCCESS] $file\e[0m"
  fi
done

echo -e "\n\e[33m[INFO] Running tests that should error\e[0m\n"
files=$(find ./examples -type f -name '*error.yail')
for file in $files; do
  if ! ./main "ast $file" 2>/dev/null; then
    echo -e "\e[31m[ERROR] $file\e[0m -- (It's supposed to fail)"
  else
    echo -e "\e[32m[SUCCESS] $file\e[0m"
  fi
done
