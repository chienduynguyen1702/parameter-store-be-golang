#!/bin/bash

# Load environment variables from .env file
dotenv_path=".env"
if [ -f "$dotenv_path" ]; then
    echo "Loading environment variables from $dotenv_path"
    source "$dotenv_path"
else
    echo "$dotenv_path not found. Make sure it exists in the same directory as this script."
    exit 1
fi

psql -h $DB_HOST -U $DB_USERNAME -d $DB_DATABASE -p $DB_PORT
# psql --help