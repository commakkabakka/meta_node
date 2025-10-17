## 一、环境

* Ubuntu 24.04.3 LTS
* GO 1.25
* MySQL 8.0.4

## 二、文件

```
├── ReadME.md
├── error
│   └── error.go    // 错误类型
├── go.mod
├── go.sum
├── main.go         // 程序入口、Gin 配置及相关路由
├── middleware
│   ├── auth.go     // JWT认证中间件
│   └── error.go    // 错误处理中间件
└── mysql
    ├── comment.go  // 评论模型及相关接口
    ├── mysql.go    // MySQL 数据库
    ├── post.go     // 文章模型及相关接口
    └── user.go     // 用户模型及相关接口
```

## 三、运行

```bash
go mod tidy
go run .
```

## 四、测试

PostMan 测试分享连接：

https://hello-5d93e094-4383719.postman.co/workspace/LianJiu-Zhang's-Workspace~383bd33c-991a-48f6-8c52-dba18f9c7536/collection/49287401-cafe58de-95ba-4e58-9f38-9083e41241a3?action=share&source=copy-link&creator=49287401

> 用户登录后返回 token，后续需要配置 token 进行调用。



















