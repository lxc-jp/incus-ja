#!/bin/bash
if [ $# -ne 1 ]; then
  >&2 echo Usage: $0 upstream_commit
  exit 1
fi
git checkout $1 -- AUTHORS client/ cmd/ go.mod go.sum grafana/ internal/ po/ scripts/ shared/ test/ \
	           .deepsource.toml
