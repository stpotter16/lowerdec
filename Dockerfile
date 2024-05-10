FROM golang:1.22 as builder

COPY . /src/lowerdec
WORKDIR /src/lowerdec

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    go build -ldflags "-s -w -extldflags '-static'" -tags osusergo,netgo,sqlite_omit_load_extension -o /usr/local/bin/lowerdec .

ADD https://github.com/benbjohnson/litestream/releases/download/v0.3.13/litestream-v0.3.13-linux-amd64.tar.gz /tmp/litestream.tar.gz
RUN tar -C /usr/local/bin -xzf /tmp/litestream.tar.gz

FROM alpine

COPY --from=builder /usr/local/bin/lowerdec /usr/local/bin/lowerdec
COPY --from=builder /usr/local/bin/litestream /usr/local/bin/litestream

RUN apk add bash

COPY etc/litestream.yml /etc/litestream.yml
COPY scripts/run.sh /scripts/run.sh

EXPOSE 8080

CMD [ "/scripts/run.sh" ]
