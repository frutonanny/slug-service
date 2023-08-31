FROM golang:1.21.0 as build

ENV BIN_FILE /opt/slug-service/slug-app
ENV CODE_DIR /go/src/github.com/frutonanny/slug-service

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/service/*

LABEL SERVICE="slug-service"

ENV CONFIG_FILE /etc/slug-service/config.json
COPY /config/config.dev.json ${CONFIG_FILE}

CMD ${BIN_FILE} -config ${CONFIG_FILE}
