
## httpbin

simple http server with a health check at `/healthz` the server is configured to automically listen on `:8080`

### Tooling

The `Dockerfile` uses [`trevatk/cfgo`](https://github.com/trevatk/cfgo) as the builder base image 

### Security

Trivy scans are provided as part of CI/CD pipeline. This package is not recommended for production. 