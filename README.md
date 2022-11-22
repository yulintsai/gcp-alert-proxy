# gcp-alert-proxy
send alert proxy from gcp webhook

# docker-compose
```yaml
version: "3"

services:
  gcp-alert-proxy:
    ports:
      - "8080:${ListenPort}"
    image: rain123473/gcp-alert-proxy:latest
    restart: always
    command: ["./gcp-alert-proxy", "run", "${Basic Auth Username}", "${Basic Auth Password}", ${ListenPort}]
```
