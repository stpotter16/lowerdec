#!/bin/bash
set -eu

echo "Attempting to restore DB if it does not exist"
litestream restore -if-db-not-exists -if-replica-exists "$DB_PATH"

exec litestream replicate -exec "/usr/local/bin/lowerdec"
