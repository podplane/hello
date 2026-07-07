# Podplane Hello World Image

This repository is responsible for the [Podplane](https://podplane.dev) "Hello World" container image published at:

```
ghcr.io/podplane/hello:latest
```

which is used as an example image for the `web` app [template](https://github.com/podplane/templates).

## Configuration

Set `HELLO_MESSAGE` env var to change the message shown on the page:

```
HELLO_MESSAGE="Hello from my app!" make run
```

If `HELLO_MESSAGE` starts with `/`, it is treated as an absolute file path. The page shows the file contents, or the read error if the file cannot be opened or read. File messages must be UTF-8 text, must not contain binary control bytes, and are limited to 200 KiB. Messages are HTML-escaped before rendering.

If `HELLO_MESSAGE` is not set, the page shows `Hello, World!` by default.

The image listens on port 8080 by default; set `PORT` to choose a different port.

## Learn More

Read more about how templates work in the Podplane [templates documentation](https://podplane.dev/docs/templates).

Learn more about Podplane at the official project website: [podplane.dev](https://podplane.dev)

## License

Podplane is licensed under the Apache License, Version 2.0.
Copyright 2026 Nadrama Pty Ltd.

See the [LICENSE](./LICENSE) file for details.
