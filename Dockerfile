# 基础镜像
FROM ubuntu:10.04

COPY webook /app/webook
WORKDIR /app

ENTRYPOINT ["/app/webook"]