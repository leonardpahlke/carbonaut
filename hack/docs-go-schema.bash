#!/bin/bash

# TODO: fix links. There is a prefix which needs to be added.
# - [type Config](<#Config>) -> - [type Config](</docs/reference/schema/#type-config>) 
# TODO: update headers (align with the other docs)
# TODO: make use of drop downs (see server-api docs)
# {{< details title="`/static-data`" open=false >}}
# - **Method**: GET
# - **Description**: Provides static data retrieved via the connector. If available, data is served from the cache; otherwise, it is retrieved from the state and cached.
# {{< /details >}}

# gomarkdoc --output docs/content/docs/reference/schema.md ./pkg/provider/...

# Exit script if you try to use an uninitialized variable.
set -o nounset
# Exit script if a statement returns a non-true return value.
set -o errexit
# Use the error status of the first failure, rather than that of the last item in a pipeline.
set -o pipefail

gomarkdoc --output documentation/content/docs/reference/schema.md ./pkg/provider/...

temp_file=$(mktemp)

{
  echo "---"
  echo "weight: 2"
  echo "---"
  cat documentation/content/docs/reference/schema.md
} > "$temp_file"

mv "$temp_file" documentation/content/docs/reference/schema.md
