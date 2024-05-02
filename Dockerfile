FROM golang:1.22 as builder

COPY . /src/lowerdec
WORKDIR /src/lowerdec

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    go build -ldflags "-s -w -extldflags '-static'" -tags osusergo,netgo,sqlite_omit_load_extension -o /usr/local/bin/lowerdec .

FROM alpine

COPY --from=builder /usr/local/bin/lowerdec /usr/local/bin/lowerdec

RUN apk add bash

EXPOSE 8080

CMD [ "/usr/local/bin/lowerdec" ]
