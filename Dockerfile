# Dockerfile - Ubuntu 方案
FROM ubuntu:22.04

# 设置环境变量避免交互式提示
ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=Asia/Shanghai

# 使用国内镜像源（阿里云）
RUN sed -i 's/archive.ubuntu.com/mirrors.aliyun.com/g' /etc/apt/sources.list && \
    sed -i 's/security.ubuntu.com/mirrors.aliyun.com/g' /etc/apt/sources.list

# 安装基础工具
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    sudo vim git curl wget make \
    gcc g++ automake autoconf libtool \
    bison flex libssl-dev pkg-config \
    tar ca-certificates \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# 安装 Go 1.25.4（使用国内镜像）
RUN curl -OL https://golang.google.cn/dl/go1.25.4.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.25.4.linux-amd64.tar.gz && \
    rm go1.25.4.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPROXY="https://goproxy.cn,direct"
ENV GO111MODULE="on"

# 安装 Thrift 编译器（从Ubuntu官方仓库安装，避免源码编译）
RUN apt-get update && \
    apt-get install -y --no-install-recommends thrift-compiler && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /work

# 复制初始化脚本
COPY scripts/init-dev.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/init-dev.sh

CMD ["/usr/local/bin/init-dev.sh"]