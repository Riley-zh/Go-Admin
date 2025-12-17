# 部署说明

## 系统要求

- Go 1.16 或更高版本
- MySQL 5.7 或更高版本
- 操作系统：Linux, macOS, 或 Windows

## 部署步骤

### 1. 克隆代码库

```bash
git clone <repository-url>
cd go-admin
```

### 2. 配置环境

#### 2.1 数据库配置
1. 创建MySQL数据库：
```sql
CREATE DATABASE go_admin DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

2. 执行初始化脚本：
```bash
mysql -u root -p go_admin < init.sql
```

#### 2.2 环境变量配置
复制 `.env.example` 文件并根据实际情况修改配置：

```bash
cp .env.example .env
```

配置项说明：
- `APP_NAME`: 应用名称
- `APP_ENV`: 运行环境 (local, dev, test, prod)
- `APP_PORT`: 应用端口
- `DB_HOST`: 数据库主机地址
- `DB_PORT`: 数据库端口
- `DB_USER`: 数据库用户名
- `DB_PASSWORD`: 数据库密码
- `DB_NAME`: 数据库名称
- `LOG_LEVEL`: 日志级别 (debug, info, warn, error)
- `LOG_OUTPUT`: 日志输出方式 (console, file, both)
- `JWT_SECRET`: JWT密钥
- `JWT_EXPIRE`: JWT过期时间
- `CACHE_MAXSIZE`: 缓存最大大小
- `CACHE_GCINTERVAL`: 缓存垃圾回收间隔

### 3. 构建应用

```bash
go mod tidy
go build -o go-admin main.go
```

### 4. 运行应用

```bash
./go-admin
```

或者在Windows系统上：

```cmd
go-admin.exe
```

### 5. 验证部署

应用启动后，可以通过以下URL验证：

- 健康检查: `http://localhost:8080/health`
- API文档: `http://localhost:8080/swagger/index.html` (如果已集成Swagger)

默认管理员账户：
- 用户名: admin
- 密码: admin123

## 监控和维护

### 日志查看
日志文件位于 `logs/` 目录下，按日期分割。

### 性能监控
可以通过以下API端点监控系统性能：
- `/api/v1/metrics`: 系统指标
- `/api/v1/health`: 健康检查
- `/api/v1/cache/stats`: 缓存统计
- `/api/v1/db/stats`: 数据库统计

### 配置热重载
修改 `.env` 文件后，应用会自动重新加载配置。

## 故障排除

### 常见问题

1. **数据库连接失败**
   - 检查数据库服务是否启动
   - 验证数据库连接配置是否正确
   - 确认数据库用户权限

2. **端口被占用**
   - 修改 `APP_PORT` 配置项
   - 或者终止占用端口的进程

3. **JWT密钥问题**
   - 确保 `JWT_SECRET` 配置项不为空
   - 生产环境建议使用强随机密钥

### 日志分析
查看日志文件以诊断问题：
```bash
tail -f logs/*.log
```

## 安全建议

1. 生产环境务必修改默认管理员密码
2. 使用强随机JWT密钥
3. 配置HTTPS证书
4. 定期备份数据库
5. 限制数据库用户权限
6. 配置防火墙规则

## 升级指南

1. 备份当前数据库
2. 拉取最新代码
3. 执行数据库迁移脚本（如果有）
4. 重启应用服务