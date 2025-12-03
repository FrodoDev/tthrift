namespace go trpc

service TService {
    TResponse twoWay(1: TRequest req);
    oneway void oneWay(1: TRequest req)
}

struct TRequest {
    1: optional i32 Type;
    2: optional binary Content;
    3: optional i64 ShardingId;
    4: optional binary Trace
}

struct TResponse {
    1: optional i64 ErrCode;
    2: optional i32 Type;
    3: optional binary Content;
    4: optional binary Trace
}