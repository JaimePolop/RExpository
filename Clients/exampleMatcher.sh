#!/bin/bash

# check if yq is installed
if ! command -v yq &> /dev/null
then
    echo "yq not found. Please install it: https://github.com/mikefarah/yq"
    exit
fi

# read the YAML file
REGEX_FILE="../regex.yaml"
REGEXES=$(yq eval '.regular_expresions[].regexes[]' "${REGEX_FILE}")
EXAMPLES=$(yq eval '.regular_expresions[].regexes[].example' "${REGEX_FILE}")
echo "{$REGEXES}"
#echo "{$EXAMPLES}"


# for REGEX in $REGEXES; do
#     for EXAMPLE in  $EXAMPLES; do
#       if grep -o -E "${REGEX}" <<< "${EXAMPLE}" > /dev/null; then
#         echo "${REGEX} regex matches example: ${EXAMPLE}"
#       fi
#     done
# done

# loop through the regexes
# for REGEX in $(echo "${REGEXES}" | yq eval '.regexes[].regex' -); do
#   #echo "{$REGEX}"
#   EXAMPLE=$(echo "${REGEXES}" | yq eval '.regexes[] | select(.regex == "'${REGEX}'").example' -)
#   NAME=$(echo "${REGEXES}" | yq eval '.name' -)
#   if grep -o -E "${REGEX}" <<< "${EXAMPLE}" > /dev/null; then
#     echo "  "
#     #"${NAME} regex matches example: ${EXAMPLE}"
#   else
#     echo " "
#     #"${NAME} regex does not match example: ${EXAMPLE}"
#   fi
# done