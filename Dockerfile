FROM golang:1.24-alpine AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY LICENSE .
COPY README.md .
COPY cmd ./cmd
COPY server ./server
COPY utils ./utils
RUN CGO_ENABLED=0 GOOS=linux go build -a -o goseek ./cmd

FROM scratch
COPY --from=build /app/goseek /app/goseek
ENTRYPOINT ["/app/goseek"]