# Podplane <https://podplane.dev>
# Copyright 2026 Nadrama Pty Ltd
# SPDX-License-Identifier: Apache-2.0

# build image
FROM golang:alpine AS build
WORKDIR /src
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /hello .

# final image
FROM scratch
COPY --from=build /hello /hello
USER 65532:65532
EXPOSE 80
ENTRYPOINT ["/hello"]
