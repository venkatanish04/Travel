FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/travelraft-server ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/travelraft-server /travelraft-server
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/travelraft-server"]
