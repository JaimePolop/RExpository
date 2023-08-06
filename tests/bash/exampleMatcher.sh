#!/bin/bash

# check if yq is installed
if ! command -v yq &> /dev/null
then
    echo "yq not found. Please install it: https://github.com/mikefarah/yq"
    exit
fi

# check if jq is installed
if ! command -v jq &> /dev/null
then
    echo "jq not found. Please install it: https://github.com/jqlang/jq"
    exit
fi

# read the YAML file
yaml_file="../../regex.yaml"
verbose=""
failure=0
cont=0


while IFS=$'\t' read -r rex_name rex_regex rex_example _; do
  
  # Remove escaped "\\" & escaped "\n" from the regex
  rex_regex="$(echo -n $rex_regex | sed 's/\\\\/\\/g' | sed 's/\\n//g')"

  # Remove escaped "\n" from the example
  rex_example="$(echo -n $rex_example | sed 's/\\n//g')"  

  # Check if the regex finds the example
  if echo -n "$rex_example" | grep -q -E "$rex_regex"; then
    if [ "$verbose" ]; then
        echo "Success ($rex_name): Regex '$rex_regex' matched the example '$rex_example'"
    fi
  else
    echo "Failure ($rex_name): Regex '$rex_regex' did not match the example '$rex_example'"
    failure=$((failure+1))
  fi

  # Add cont
  cont=$((cont+1))
done < <(yq -o=json $yaml_file | jq -r '.regular_expresions[] | .regexes[] | [.name, .regex, .example] | @tsv')


echo "Checked $cont examples"

if [ $failure -eq 0 ]; then
  echo "All examples matched their regular expressions"
else
  echo "$failure examples did not match their regular expressions"
  exit $failure
fi