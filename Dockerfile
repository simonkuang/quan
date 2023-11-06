FROM alpine:latest AS builder

RUN sed -i \
    's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' \
    /etc/apk/repositories && \
    apk update && apk add curl && \
    mkdir -p /data/soft && \
    curl -Lo /data/soft/go1.21.3.linux-amd64.tar.gz \
    https://go.dev/dl/go1.21.3.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf /data/soft/go1.21.3.linux-amd64.tar.gz && \
    eval "$(/usr/local/go/bin/go env)" && \
    # git clone and build
    mkdir -p /data/src && \
    git clone https://github.com/simonkuang/quan.git /data/src/quan && \
    cd /data/src/quan && \
    /usr/local/go/bin/go build -o /data/bin/quan src/main.go

FROM alpine:latest AS prod

ENV QUAN_ADMIN_USERNAME ""
ENV QUAN_ADMIN_PASSWORD ""

COPY --from=builder /data/bin/quan /usr/bin/quan
COPY ./docker-entrypoint.sh /usr/bin/docker-entrypoint.sh

RUN mkdir -p /data/quan && \
    chown +x /usr/bin/docker-entrypoint.sh

WORKDIR /data/quan

CMD [ "/usr/bin/quan" ]

ENTRYPOINT [ "/usr/bin/docker-entrypoint.sh" ]
