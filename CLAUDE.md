# CLAUDE.md

本文档用于指导 Claude Code 在当前仓库中协作开发。

## 交流与输出要求

- 全程使用中文回复，包括代码注释、问题说明、报错分析与变更总结。
- 优先基于仓库现状开展工作，不要沿用“空项目模板”式假设。
- 修改代码前，先确认目标服务、模块边界和配置文件位置。

## 项目概览

当前仓库并非空目录，核心项目位于 `Online ride-hailing/`，是一个基于 Go Workspace 的网约车业务示例，包含多个独立模块与前端静态页面。

仓库根目录主要内容：

- `Online ride-hailing/`：主项目目录
- `AGENTS.md`：面向 Codex 的协作说明
- `CLAUDE.md`：当前文件
- `.idea/`：IDE 配置目录

## 工作区结构

`Online ride-hailing/` 下当前包含：

- `go.work`：Go 工作区配置
- `user-api/`：对外 HTTP API 服务，使用 go-zero，包含 WebSocket 和定时任务
- `user-srv/`：用户侧 gRPC 服务
- `driver-srv/`：司机侧 gRPC 服务
- `admin-srv/`：管理端 gRPC 服务
- `user-web/`：静态前端页面与 JS
- `docs/`：项目文档、接口资料、数据字典
- `nacos/`：本地 nacos 相关目录

## Go Workspace 信息

工作区文件：`Online ride-hailing/go.work`

当前 `use` 的模块有：

- `./admin-srv`
- `./driver-srv`
- `./user-api`
- `./user-srv`

Go 版本当前声明为：

```bash
go 1.26.2
```

如果执行 Go 命令，优先在 `Online ride-hailing/` 目录下进行，以便使用 `go.work` 统一解析各模块依赖。

## 各模块说明

### 1. user-api

- 模块名：`user-api`
- 入口文件：`Online ride-hailing/user-api/user.go`
- 配置文件：`Online ride-hailing/user-api/etc/user-api.yaml`
- 主要特征：
  - 使用 `go-zero/rest`
  - 提供 HTTP 接口
  - 提供 WebSocket 路由 `/ws`
  - 启动了基于 `cron` 的定时任务

### 2. user-srv

- 模块名：`user-srv`
- 入口文件：`Online ride-hailing/user-srv/user.go`
- 配置文件：`Online ride-hailing/user-srv/etc/user.yaml`
- 主要特征：
  - gRPC 服务
  - 使用 `go-zero/zrpc`
  - 涉及 MySQL、Redis、JWT、Nacos、Viper、GORM 等依赖

### 3. driver-srv

- 模块名：`driver-srv`
- 入口文件：`Online ride-hailing/driver-srv/driver.go`
- 配置文件：`Online ride-hailing/driver-srv/etc/driver.yaml`
- 主要特征：
  - gRPC 服务
  - 使用 `go-zero/zrpc`
  - 结构与 `user-srv` 类似

### 4. admin-srv

- 模块名：`admin-srv`
- 入口文件：`Online ride-hailing/admin-srv/admin.go`
- 配置文件：`Online ride-hailing/admin-srv/etc/admin.yaml`
- 主要特征：
  - gRPC 服务
  - 使用 `go-zero/zrpc`
  - 结构与 `user-srv`、`driver-srv` 接近

## 常用命令

以下命令建议在 `D:\GoWork\src\2309A\ClaudeCode\fly-pig-solo\Online ride-hailing` 目录下执行。

### 工作区级别

```bash
go work sync
go build ./...
go test ./...
go fmt ./...
go vet ./...
```

### 单模块运行

```bash
cd user-api
go run user.go -f etc/user-api.yaml
```

```bash
cd user-srv
go run user.go -f etc/user.yaml
```

```bash
cd driver-srv
go run driver.go -f etc/driver.yaml
```

```bash
cd admin-srv
go run admin.go -f etc/admin.yaml
```

### 单模块测试

```bash
go test ./user-api/...
go test ./user-srv/...
go test ./driver-srv/...
go test ./admin-srv/...
```

如果本地安装了 `golangci-lint`，可额外执行：

```bash
golangci-lint run
```

## 协作建议

- 优先在 `Online ride-hailing/` 内操作，不要把仓库根目录误判为空项目。
- 新增依赖前先确认属于哪个模块，避免错误修改到其他服务的 `go.mod`。
- 修改接口时，同时检查：
  - handler 或 server 实现
  - `internal/` 下业务逻辑
  - `pb/`、`.proto` 或 `.api` 定义
  - 配置文件与调用方是否需要同步调整
- 涉及前端联调时，关注 `user-web/` 中的静态页面与 `api.js`、`ws.js`。
- 涉及注册发现或本地运行问题时，检查各模块 `etc/*.yaml` 以及根目录下的 `nacos/` 相关内容。

## 注意事项

- 仓库中存在部分中文注释显示乱码现象，编辑文件时注意编码一致性，优先保持 UTF-8。
- 根目录现有 `AGENTS.md` 仍包含“目录为空”的旧描述，与当前实际结构不一致；后续如有需要，可一并更新。
- 如果后续项目结构、启动命令或依赖发生变化，应及时同步更新本文件。
