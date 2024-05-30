#!/bin/bash

show_help() {
    echo "Usage: ./get-parameter.sh [OPTIONS]"
    echo "Get parameters by API token."
    echo "Options:"
    echo "    -h, --help                          Display this help message."
    echo "    -o, --output-file-name <file-name>  Name of the output file. Default is 'parameters.txt'."
    echo "Requirements:"
    echo "    - 'jq' must be installed."
    echo "    - Environment variable \$PARAMETER_STORE_TOKEN must be set in user profile."
    echo ""
    echo "By Parameter Store - HUST - 20205059"
}

output_file="parameters.txt"
log_file="./get-parameters.log"

while [[ $# -gt 0 ]]; do
    key="$1"
    case $key in
        -h|--help)
        show_help
        exit 0
        ;;
        -o|--output-file-name)
        output_file="$2"
        shift
        shift
        ;;
        *)
        echo "Unknown option: $1"
        show_help
        exit 1
        ;;
    esac
done

if [ -z "$PARAMETER_STORE_TOKEN" ]; then
    echo "Environment variable \$PARAMETER_STORE_TOKEN is not set."
    exit 1
fi

response=$(curl -s  -D - -o $output_file  POST https://param-store-be.datn.live/api/v1/agents/auth-parameters \
    -H "Content-Type: application/json" \
    -d "{\"api_token\":\"$PARAMETER_STORE_TOKEN\"}")

# Separate the headers, body, and status code
headers=$(echo "$response" | sed -n '/^\r$/q;p')
body=$(echo "$response" | sed -n '/^\r$/,$p' | sed '1d')

# Extract the status code from the headers
status_code=$(echo "$headers" | grep HTTP | awk '{print $2}')

timestamp=$(date +"%Y-%m-%d %H:%M:%S")

if [ "$status_code" == "200" ]; then
    echo "Parameters retrieved successfully."
    echo "$timestamp $status_code Parameters retrieved successfully." >> "$log_file"
    exit 1
fi
if [ "$status_code" == "401" ]; then
    echo "Unauthorized. Check the API token."
    echo "$timestamp $status_code Unauthorized. Check the API token." >> "$log_file"
    exit 1
fi
if [ "$status_code" == "404" ]; then
    echo "No parameters found."
    echo "$timestamp $status_code No parameters found." >> "$log_file"
    exit 1
fi
if [ "$status_code" == "500" ]; then
    echo "Internal server error."
    echo "$timestamp $status_code Internal server error." >> "$log_file"
    exit 1
fi
if [ "$status_code" != "200" ]; then
    echo "Failed to get parameters: $body. Check $log_file for more details."
    echo "$timestamp $status_code Failed to get parameters: $body." >> "$log_file"
    exit 1
fi
parameters=$(echo "$body")

if [ -z "$parameters" ]; then
    echo "No parameters found in the response."
    echo "$timestamp $status_code No parameters found in the response." >> "$log_file"
    exit 1
fi

echo "$parameters" > "$output_file"
echo "Parameters written to $output_file successfully."
