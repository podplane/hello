# Podplane <https://podplane.dev>
# Copyright 2026 Nadrama Pty Ltd
# SPDX-License-Identifier: Apache-2.0

IMAGE ?= ghcr.io/podplane/hello
TAG ?= latest
PORT ?= 8080

.PHONY: build run image push

build:
	go build ./...

run:
	PORT=$(PORT) go run .

image:
	docker build -f Containerfile -t $(IMAGE):$(TAG) .

push: image
	docker push $(IMAGE):$(TAG)
