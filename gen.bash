#!/bin/bash

ROOT_DIR=$(pwd)

echo "Generating sqlc..."

for file in `find . -type f -name "sqlc.yml"`; do
    echo `dirname $file`
    (cd `dirname $file`; $ROOT_DIR/bin/sqlc generate)
done
