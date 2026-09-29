# Docker

Build and run the server from the repository root:

```powershell
docker build -t travelraft:latest .
docker run --rm -p 8080:8080 -e NODE_ID=node1 travelraft:latest
```

The image uses a Go build stage and a non-root distroless runtime. The server listens on port 8080 by default and accepts `NODE_ID` and `PORT` environment variables.
