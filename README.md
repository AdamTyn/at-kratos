## at-kratos
嗯，一个更适合开箱即用的 ***kratos-layout***

0. 环境
   
   ```txt
   Golang v1.17+
   Apollo
   数据库
   Redis
   ```

1. 文件结构梳理

    ```txt
    api/... # 提供grpc服务的proto文件
    cmd/at-kratos/... # 程序入口文件
    internal/biz/...
    internal/conf/... # (解析apollo)配置的proto文件
    internal/data/...
    internal/data/dao/...
    internal/data/entity/...
    internal/pkg/cache/... # 缓存组件
    internal/pkg/databse/... # 数据库组件
    internal/pkg/grpc_client/... # 调用外部的grpc客户端
    internal/pkg/http_client/... # 调用外部的http客户端
    internal/pkg/util/... # 通用组件
    internal/server/... # 注册grpc服务
    internal/service/...
    migrations/v1/... # 数据库迁移
    tests/v1/... # 单元测试
    third_party/... # protobuf的官方依赖
    Dockerfile
    Makefile # 构建指令
    README.md
    ......
    ```

2. 构建指令详解
   - 生成依赖注入相关文件

     ```bash
     make generate
     ```
   - 编译 *internal/conf/conf.proto* 文件
     
     ```bash
     make config
     ```
   - 编译 **grpc服务** 的 *.proto* 文件
     
     ```bash
     make api
     ```
   - 编译调用外部的 **grpc客户端** 的 *.proto* 文件
     
     ```bash
     make grpc-client
     ```
   - 编译调用外部的 **http客户端** 的 *.proto* 文件
     
     ```bash
     make http-client
     ```
   - 编译打包
     
     ```bash
     make build
     ```

3. 感谢
   - [**kratos-doc**](https://go-kratos.dev/docs/)
   - [**proto3-doc**](https://developers.google.cn/protocol-buffers/docs/proto3)
   - [**insomnia**](https://github.com/Kong/insomnia/releases)
   - [**zorm**](https://gitee.com/chunanyong/zorm)

4. 联系方式
   - 邮箱：tynadam@foxmail.com
   - 闲鱼：欧布00