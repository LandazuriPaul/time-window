ARG GO_VERSION=1.23.2

FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /time-window

COPY . .

RUN go build -o /time-window/build/time-window ./cmd/cli/cli.go

FROM scratch AS runner

COPY --from=builder /time-window/build/time-window /time-window

ENTRYPOINT ["/time-window"]
