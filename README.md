通过 thrift 学习 rpc 的设计思想

## 流程
1. 容器搭建
2. thrift 客户端、服务端代码编写，调试通过
3. 学习源码
   1. client.TwoWay 调用是怎么把数据发到网络的？
   2. 客户端的分层
   3. 服务端的分层
   4. 服务端的 processor和业务回调怎么关联的？
4. 项目优化， 完成下面的 todo
5. 


## work log
* 2025-12-01，构建容器

## todo
1. 目录结构调整 “协议层 -> 传输层 -> 业务层”，是否有工具层，比如日志、监控、链路追踪这些
tthrift/
├── cmd/               # 可执行文件入口
│   ├── tserver/       # 服务端 main.go
│   └── tclient/       # 客户端 main.go
├── internal/          # 私有应用代码
│   ├── handler/       # 业务处理器（你的 register.go, handler.go）
│   ├── server/        # 服务端逻辑（你的 server.go）
│   └── client/        # 客户端连接池等（你的 client.go, selector.go）
├── pkg/               # 可公开的库代码
│   └── protocol/      # 协议相关（你的 Marshal 等）
└── ...                # 其他 (idl, gen-go, docker 等保持不变)

2. 连接池优化
3. 配置化
4. 错误分类，细化
5. 日志组件、切分（zap、）
6. 配置接入 etcd、Apollo
7. 接入 jaeger
8. 接入 prometheus
9. 序列化、接口化