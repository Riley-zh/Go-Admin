# Go Admin - 后台管理系统

一个使用Go语言开发的企业级后台管理系统，基于Gin框架构建，提供了完整的用户管理、权限控制、系统监控等功能。

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
- **缓存机制**: 自研缓存实现，支持统计和监控
- **数据库优化**: 连接池管理和监控
- **日志系统**: 结构化日志和动态日志级别调整
- **安全增强**: CSRF防护和请求频率限制
- **系统监控**: 性能指标收集和健康检查
- **配置管理**: 多环境配置和热重载

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
- **缓存**: 自研缓存实现 (基于sync.Map)
- **日志**: Zap
- **配置**: Viper
- **认证**: JWT
- **其他**: 限流、CSRF防护等中间件

## 项目结构

```
go-admin/
├── config/                 # 配置管理
├── internal/
│   ├── app/               # 应用启动和路由配置
│   ├── cache/             # 自研缓存实现
│   ├── database/          # 数据库连接管理
│   ├── handler/           # HTTP请求处理器
│   ├── logger/            # 日志系统
│   ├── middleware/        # 中间件
│   ├── model/             # 数据模型
│   ├── repository/        # 数据访问层
│   └── service/           # 业务逻辑层
├── pkg/
│   ├── errors/            # 自定义错误处理
│   ├── response/          # HTTP响应处理
│   └── utils/             # 工具函数
├── logs/                  # 日志文件目录
├── init.sql              # 数据库初始化脚本
├── DEPLOYMENT.md         # 部署说明
├── API.md                # API接口文档
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

详细的API接口文档请参考 [API.md](API.md) 文件。

## 部署说明

详细的部署说明请参考 [DEPLOYMENT.md](DEPLOYMENT.md) 文件。

## 开发指南

### 代码结构说明

本项目采用分层架构设计，各层职责如下：

1. **handler层**: 处理HTTP请求，负责参数解析和响应返回
2. **service层**: 实现业务逻辑
3. **repository层**: 数据访问层，负责与数据库交互
4. **model层**: 数据模型定义
5. **middleware层**: 中间件，处理跨域、认证、日志等通用功能
6. **cache层**: 缓存实现
7. **logger层**: 日志系统
8. **database层**: 数据库连接管理
9. **config层**: 配置管理

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

## 贡献指南

欢迎提交Issue和Pull Request来改进项目。

## 许可证

本项目采用MIT许可证，详情请见[LICENSE](LICENSE)文件。