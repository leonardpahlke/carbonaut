#!/bin/bash

# Exit script if you try to use an uninitialized variable.
set -o nounset
# Exit script if a statement returns a non-true return value.
set -o errexit
# Use the error status of the first failure, rather than that of the last item in a pipeline.
set -o pipefail

# TODO: fix links. There is a prefix which needs to be added.
# - [type Config](<#Config>) -> - [type Config](</docs/reference/schema/#type-config>) 

gomarkdoc --output docs/content/docs/reference/schema.md ./pkg/provider/...
