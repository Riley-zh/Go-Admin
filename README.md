# Go Admin - 企业级后台管理系统

一个使用Go语言开发的现代化企业级后台管理系统，基于Gin框架构建，提供了完整的用户管理、权限控制、系统监控、API优化等功能。系统经过深度优化，具备高性能、高可用性和强大的安全特性。

## 功能特性

### 核心功能
- **用户管理**: 用户注册、登录、信息管理、密码修改
- **权限控制**: 基于RBAC的权限模型，支持角色和权限分配
- **菜单管理**: 动态菜单配置，支持多级菜单
- **系统日志**: 完整的操作日志记录和查询
- **数据字典**: 系统数据字典管理
- **文件管理**: 文件上传、下载、删除
- **通知公告**: 系统通知和公告管理
- **定时任务**: 基于cron表达式的定时任务管理

### 系统优化
- **缓存机制**: 基于Redis的分布式缓存实现，支持统计和监控
- **数据库优化**: 连接池管理和查询优化，支持索引分析和慢查询监控
- **日志系统**: 基于Zap的高性能异步日志系统，支持采样和上下文信息注入
- **安全增强**: JWT认证、CSRF防护、请求签名验证、频率限制等多重安全保障
- **系统监控**: 基于Prometheus的性能指标收集和健康检查
- **配置管理**: 多环境配置和热重载
- **API优化**: 批量处理、流式处理、响应压缩等API性能优化

### 技术特性
- **现代化架构**: 采用分层架构设计，代码结构清晰
- **高性能**: 基于Gin框架，轻量级且高性能
- **易扩展**: 模块化设计，易于功能扩展
- **高可用**: 完善的监控和健康检查机制
- **安全性**: JWT认证、CSRF防护、频率限制等多重安全保障

## 技术栈

- **语言**: Go 1.16+
- **Web框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL
- **缓存**: Redis分布式缓存
- **日志**: Zap高性能日志库
- **配置**: Viper
- **认证**: JWT (支持IP绑定和User-Agent验证)
- **API文档**: Swagger/OpenAPI
- **安全**: 限流、CSRF防护、请求签名验证
- **监控**: Prometheus指标收集

## 项目结构

```
```go-admin/
├── config/                 # 配置管理
├── docs/                   # 文档和Swagger自动生成的API文档
├── internal/
│   ├── app/               # 应用启动和路由配置
│   ├── cache/             # 缓存实现
│   ├── database/          # 数据库连接管理和优化
│   ├── handler/           # HTTP请求处理器
│   ├── logger/            # 高性能日志系统
│   ├── middleware/        # 中间件
│   ├── metrics/           # 监控指标收集
│   ├── model/             # 数据模型
│   ├── repository/        # 数据访问层
│   └── service/           # 业务逻辑层
├── pkg/
│   ├── api/               # API客户端和优化工具
│   ├── errors/            # 自定义错误处理
│   ├── httpclient/        # HTTP客户端
│   ├── jsonutils/         # JSON处理工具
│   ├── middleware/        # 包级中间件
│   ├── response/          # HTTP响应处理
│   ├── utils/             # 工具函数
│   └── validation/        # 验证工具
├── logs/                  # 日志文件目录
├── init.sql              # 数据库初始化脚本
├── DEPLOYMENT.md         # 部署说明
└── main.go               # 程序入口
```

## 快速开始

### 环境要求
- Go 1.16 或更高版本
- MySQL 5.7 或更高版本

### 安装步骤

1. 克隆项目
```bash
git clone <repository-url>
cd go-admin
```

2. 初始化数据库
```sql
CREATE DATABASE go_admin DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

3. 执行初始化脚本
```bash
mysql -u root -p go_admin < init.sql
```

4. 配置环境变量
```bash
cp .env.example .env
# 修改.env文件中的配置项
```

5. 安装依赖
```bash
go mod tidy
```

6. 运行应用
```bash
go run main.go
```

或者构建后运行：
```bash
go build -o go-admin main.go
./go-admin
```

### 默认账户
- 用户名: admin
- 密码: admin123

## API文档

本项目集成了Swagger/OpenAPI自动生成API文档功能，提供了两种方式查看API文档：

### 1. 在线Swagger UI (推荐)
启动应用后，访问以下地址查看交互式API文档：
```
http://localhost:8080/swagger/index.html
```

Swagger UI提供了以下功能：
- 交互式API文档浏览
- 在线API测试
- 请求/响应示例
- 认证配置（JWT Token）

### 2. 静态API文档
详细的API接口文档可以通过在线Swagger UI查看，或者参考 [docs/swagger.json](docs/swagger.json) 文件。

### 3. 更新API文档
当添加新的API接口时，只需在handler函数中添加Swagger注释，然后运行以下命令更新文档：
```bash
go run github.com/swaggo/swag/cmd/swag@latest init
```

Swagger注释示例：
```go
// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with username, password and email
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateUserRequest true "User details"
// @Success 201 {object} map[string]interface{} "User created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
    // 实现代码
}
```

## 部署说明

详细的部署说明请参考 [DEPLOYMENT.md](DEPLOYMENT.md) 文件。

## 开发指南

### 代码结构说明

本项目采用分层架构设计，各层职责如下：

1. **handler层**: 处理HTTP请求，负责参数解析和响应返回
2. **service层**: 实现业务逻辑
3. **repository层**: 数据访问层，负责与数据库交互
4. **model层**: 数据模型定义
5. **middleware层**: 中间件，处理跨域、认证、日志、限流等通用功能
6. **cache层**: 缓存实现
7. **logger层**: 高性能异步日志系统
8. **database层**: 数据库连接管理和查询优化
9. **metrics层**: 监控指标收集
10. **config层**: 配置管理
11. **pkg层**: 通用工具包，包括API客户端、验证工具等

### 添加新功能

1. 在`model`目录下创建数据模型
2. 在`repository`目录下创建数据访问方法
3. 在`service`目录下实现业务逻辑
4. 在`handler`目录下创建HTTP处理方法
5. 在`internal/app/app.go`中注册路由

### 代码规范

1. 遵循Go语言官方编码规范
2. 使用`gofmt`格式化代码
3. 添加必要的注释和文档
4. 编写单元测试

## 测试

运行所有测试：
```bash
go test ./...
```

运行特定模块测试：
```bash
go test ./internal/service/...
```

## 监控和维护

### 健康检查
访问 `http://localhost:8080/health` 查看系统健康状态。

### 性能监控
通过以下端点监控系统性能：
- `/api/v1/metrics`: 系统指标
- `/api/v1/cache/stats`: 缓存统计
- `/api/v1/db/stats`: 数据库统计

### 日志查看
日志文件位于 `logs/` 目录下，按日期分割。

## 安全建议

1. 生产环境务必修改默认管理员密码
2. 使用强随机JWT密钥
3. 配置HTTPS证书
4. 定期备份数据库
5. 限制数据库用户权限
6. 配置防火墙规则
7. 启用请求签名验证防止重放攻击
8. 配置适当的频率限制策略
9. 定期审查审计日志

## 贡献指南

欢迎提交Issue和Pull Request来改进项目。

## 许可证

本项目采用MIT许可证，详情请见[LICENSE](LICENSE)文件。