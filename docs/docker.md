# Docker

Build the image and run the three-node system from the repository root:

```cmd
docker build -t travelraft:latest .
docker compose up -d
docker compose ps
```

The image uses a Go build stage and a minimal scratch runtime. Compose starts one MySQL service and three TravelRaft nodes on ports `50051`, `50052`, and `50053`. Node identity and port are passed as `--id` and `--port` arguments, while database settings use `MYSQL_*` environment variables.
