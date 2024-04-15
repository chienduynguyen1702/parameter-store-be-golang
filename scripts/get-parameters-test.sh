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
log_file="get-parameters.log"

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

response=$(curl -s -X POST http://localhost:8086/api/v1/agents/auth-parameters \
    -H "Content-Type: application/json" \
    -d "{\"api_token\":\"$PARAMETER_STORE_TOKEN\"}")

if [ $? -ne 0 ]; then
    echo "Error executing curl command."
    exit 1
fi
parameters=$(echo "$response" | jq -r '.parameters[] | "\(.name)=\(.value)"')

if [ -z "$parameters" ]; then
    echo "No parameters found in the response."
    exit 1
fi

echo "$parameters" > "$output_file"
echo "Parameters written to $output_file successfully."
