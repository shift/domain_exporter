#!/usr/bin/env bash

export GO_VERSION=$(go version | awk '{print $3}')
echo "$REG_PW" | docker login quay.io --username "$REG_USER" --password-stdin
curl -sL https://git.io/goreleaser | bash
