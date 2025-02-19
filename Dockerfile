FROM golang:1.22-bookworm

WORKDIR /app
RUN go install github.com/jackc/tern/v2@latest
RUN go install go.uber.org/mock/mockgen@latest
ENTRYPOINT ["tail", "-f", "/dev/null"]
