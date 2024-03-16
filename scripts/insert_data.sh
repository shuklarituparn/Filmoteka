#!/bin/bash

SQL_FILE="./init/db_data.sql"

parse_env() {
    while read -r line; do
        if [[ "$line" =~ ^[[:space:]]*([^#][^=]+)[[:space:]]*=[[:space:]]*(.*)[[:space:]]*$ ]]; then
            echo $BASH_REMATCH
            key="${BASH_REMATCH[1]}"
            value="${BASH_REMATCH[2]}"
            export "$key"="$value"
        fi
    done < .env
}

execute_sql() {
    echo "Variables: "
    echo "DB_USER: $POSTGRES_USER"
    echo "DB_PORT: $POSTGRES_PORT"
    echo "DB_NAME: $POSTGRES_DB"
    echo "SQL_FILE: $SQL_FILE"
    echo "DB_PASSWORD: $POSTGRES_PASSWORD"
    echo "DB_HOST: $POSTGRES_HOST"
    psql -U "$POSTGRES_USER" -p "$POSTGRES_PORT" -d "$POSTGRES_DB" -f "$SQL_FILE" -W
}

parse_env

execute_sql
