FROM golang:1.16 as build

WORKDIR /zcx

ADD . /zcx

RUN go env -w GO111MODULE=on && \
    go env -w GOPROXY=https://goproxy.cn,direct && \
    GOOS=linux go build -o ./app ./cmd

FROM scratch as final

COPY --from=build /zcx/app /zcx/app.yaml /
ENV TZ=Asia/Shanghai

ENTRYPOINT [ "/app"]