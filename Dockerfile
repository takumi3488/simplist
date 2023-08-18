FROM golang:1.21.0-alpine3.18 AS dev
WORKDIR /app
RUN apk update && apk add --no-cache git postgresql-client
COPY . .
RUN go mod download
RUN go install github.com/cweill/gotests/gotests@latest && \
    go install github.com/fatih/gomodifytags@latest && \
    go install github.com/josharian/impl@latest && \
    go install github.com/haya14busa/goplay/cmd/goplay@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest && \
    go install honnef.co/go/tools/cmd/staticcheck@latest && \
    go install golang.org/x/tools/gopls@latest && \
    go install github.com/go-task/task/v3/cmd/task@latest
CMD [ "go", "run", "server.go" ]


FROM golang:1.21.0-alpine3.18 AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o server .


FROM scratch AS prod
WORKDIR /app
COPY --from=build /app/server .

CMD [ "./server" ]
