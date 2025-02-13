FROM golang:1.22-bookworm

WORKDIR /app
RUN go install github.com/jackc/tern/v2@latest
ENTRYPOINT ["tail", "-f", "/dev/null"]
