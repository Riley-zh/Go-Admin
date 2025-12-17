# API 接口文档

## 基础信息

- Base URL: `http://localhost:8080/api/v1`
- 认证方式: JWT Token
- 数据格式: JSON

## 认证接口

### 用户注册
```
POST /register
```

**请求参数:**
```json
{
  "username": "string",
  "password": "string",
  "email": "string",
  "nickname": "string"
}
```

**响应:**
```json
{
  "code": 200,
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "username": "string",
    "email": "string",
    "nickname": "string",
    "avatar": "string",
    "status": 1,
    "created_at": "2023-01-01T00:00:00Z"
  }
}
```

### 用户登录
```
POST /login
```

**请求参数:**
```json
{
  "username": "string",
  "password": "string"
}
```

**响应:**
```json
{
  "code": 200,
  "message": "Login successful",
  "data": {
    "token": "jwt_token",
    "user": {
      "id": 1,
      "username": "string",
      "email": "string",
      "nickname": "string",
      "avatar": "string",
      "status": 1,
      "created_at": "2023-01-01T00:00:00Z"
    }
  }
}
```

### 用户登出
```
POST /logout
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Logout successful",
  "data": null
}
```

### 刷新令牌
```
POST /refresh
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Token refreshed successfully",
  "data": {
    "token": "new_jwt_token"
  }
}
```

## 用户管理接口

### 获取用户信息
```
GET /users/{id}
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "User retrieved successfully",
  "data": {
    "id": 1,
    "username": "string",
    "email": "string",
    "nickname": "string",
    "avatar": "string",
    "status": 1,
    "created_at": "2023-01-01T00:00:00Z"
  }
}
```

### 更新用户信息
```
PUT /users/{id}
```

**请求头:**
```
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "nickname": "string",
  "avatar": "string"
}
```

**响应:**
```json
{
  "code": 200,
  "message": "User updated successfully",
  "data": {
    "id": 1,
    "username": "string",
    "email": "string",
    "nickname": "string",
    "avatar": "string",
    "status": 1,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
}
```

### 删除用户
```
DELETE /users/{id}
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "User deleted successfully",
  "data": null
}
```

### 获取用户列表
```
GET /users
```

**请求参数:**
- page: 页码 (默认: 1)
- page_size: 每页数量 (默认: 10)

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Users retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "username": "string",
        "email": "string",
        "nickname": "string",
        "avatar": "string",
        "status": 1,
        "created_at": "2023-01-01T00:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

### 修改密码
```
PUT /users/change-password
```

**请求头:**
```
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "old_password": "string",
  "new_password": "string"
}
```

**响应:**
```json
{
  "code": 200,
  "message": "Password changed successfully",
  "data": null
}
```

## 角色管理接口

### 创建角色
```
POST /roles
```

**请求头:**
```
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "name": "string",
  "description": "string"
}
```

**响应:**
```json
{
  "code": 200,
  "message": "Role created successfully",
  "data": {
    "id": 1,
    "name": "string",
    "description": "string",
    "status": 1,
    "created_at": "2023-01-01T00:00:00Z"
  }
}
```

### 获取角色信息
```
GET /roles/{id}
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Role retrieved successfully",
  "data": {
    "id": 1,
    "name": "string",
    "description": "string",
    "status": 1,
    "created_at": "2023-01-01T00:00:00Z"
  }
}
```

### 更新角色信息
```
PUT /roles/{id}
```

**请求头:**
```
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "name": "string",
  "description": "string"
}
```

**响应:**
```json
{
  "code": 200,
  "message": "Role updated successfully",
  "data": {
    "id": 1,
    "name": "string",
    "description": "string",
    "status": 1,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
}
```

### 删除角色
```
DELETE /roles/{id}
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Role deleted successfully",
  "data": null
}
```

### 获取角色列表
```
GET /roles
```

**请求参数:**
- page: 页码 (默认: 1)
- page_size: 每页数量 (默认: 10)

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Roles retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "name": "string",
        "description": "string",
        "status": 1,
        "created_at": "2023-01-01T00:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

### 分配角色给用户
```
POST /roles/assign
```

**请求头:**
```
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "user_id": 1,
  "role_id": 1
}
```

**响应:**
```json
{
  "code": 200,
  "message": "Role assigned successfully",
  "data": null
}
```

### 移除用户角色
```
POST /roles/remove
```

**请求头:**
```
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "user_id": 1,
  "role_id": 1
}
```

**响应:**
```json
{
  "code": 200,
  "message": "Role removed successfully",
  "data": null
}
```

### 获取用户的角色列表
```
GET /users/{id}/roles
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "User roles retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "string",
      "description": "string",
      "status": 1,
      "created_at": "2023-01-01T00:00:00Z"
    }
  ]
}
```

## 权限管理接口

### 创建权限
```
POST /permissions
```

**请求头:**
```
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "name": "string",
  "description": "string",
  "resource": "string",
  "action": "string"
}
```

**响应:**
```json
{
  "code": 200,
  "message": "Permission created successfully",
  "data": {
    "id": 1,
    "name": "string",
    "description": "string",
    "resource": "string",
    "action": "string",
    "created_at": "2023-01-01T00:00:00Z"
  }
}
```

### 获取权限信息
```
GET /permissions/{id}
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Permission retrieved successfully",
  "data": {
    "id": 1,
    "name": "string",
    "description": "string",
    "resource": "string",
    "action": "string",
    "created_at": "2023-01-01T00:00:00Z"
  }
}
```

### 更新权限信息
```
PUT /permissions/{id}
```

**请求头:**
```
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "name": "string",
  "description": "string",
  "resource": "string",
  "action": "string"
}
```

**响应:**
```json
{
  "code": 200,
  "message": "Permission updated successfully",
  "data": {
    "id": 1,
    "name": "string",
    "description": "string",
    "resource": "string",
    "action": "string",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
}
```

### 删除权限
```
DELETE /permissions/{id}
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Permission deleted successfully",
  "data": null
}
```

### 获取权限列表
```
GET /permissions
```

**请求参数:**
- page: 页码 (默认: 1)
- page_size: 每页数量 (默认: 10)

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Permissions retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "name": "string",
        "description": "string",
        "resource": "string",
        "action": "string",
        "created_at": "2023-01-01T00:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

### 分配权限给角色
```
POST /permissions/assign
```

**请求头:**
```
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "role_id": 1,
  "permission_id": 1
}
```

**响应:**
```json
{
  "code": 200,
  "message": "Permission assigned successfully",
  "data": null
}
```

### 移除角色权限
```
POST /permissions/remove
```

**请求头:**
```
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "role_id": 1,
  "permission_id": 1
}
```

**响应:**
```json
{
  "code": 200,
  "message": "Permission removed successfully",
  "data": null
}
```

### 获取角色的权限列表
```
GET /roles/{id}/permissions
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Role permissions retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "string",
      "description": "string",
      "resource": "string",
      "action": "string",
      "created_at": "2023-01-01T00:00:00Z"
    }
  ]
}
```

### 获取用户的权限列表
```
GET /users/{id}/permissions
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "User permissions retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "string",
      "description": "string",
      "resource": "string",
      "action": "string",
      "created_at": "2023-01-01T00:00:00Z"
    }
  ]
}
```

## 菜单管理接口

### 创建菜单
```
POST /menus
```

**请求头:**
```
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "name": "string",
  "title": "string",
  "icon": "string",
  "path": "string",
  "component": "string",
  "redirect": "string",
  "permission": "string",
  "parent_id": 0,
  "sort": 0,
  "status": 1,
  "hidden": 0
}
```

**响应:**
```json
{
  "code": 200,
  "message": "Menu created successfully",
  "data": {
    "id": 1,
    "name": "string",
    "title": "string",
    "icon": "string",
    "path": "string",
    "component": "string",
    "redirect": "string",
    "permission": "string",
    "parent_id": 0,
    "sort": 0,
    "status": 1,
    "hidden": 0,
    "created_at": "2023-01-01T00:00:00Z"
  }
}
```

### 获取菜单信息
```
GET /menus/{id}
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Menu retrieved successfully",
  "data": {
    "id": 1,
    "name": "string",
    "title": "string",
    "icon": "string",
    "path": "string",
    "component": "string",
    "redirect": "string",
    "permission": "string",
    "parent_id": 0,
    "sort": 0,
    "status": 1,
    "hidden": 0,
    "created_at": "2023-01-01T00:00:00Z"
  }
}
```

### 更新菜单信息
```
PUT /menus/{id}
```

**请求头:**
```
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "name": "string",
  "title": "string",
  "icon": "string",
  "path": "string",
  "component": "string",
  "redirect": "string",
  "permission": "string",
  "parent_id": 0,
  "sort": 0,
  "status": 1,
  "hidden": 0
}
```

**响应:**
```json
{
  "code": 200,
  "message": "Menu updated successfully",
  "data": {
    "id": 1,
    "name": "string",
    "title": "string",
    "icon": "string",
    "path": "string",
    "component": "string",
    "redirect": "string",
    "permission": "string",
    "parent_id": 0,
    "sort": 0,
    "status": 1,
    "hidden": 0,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
}
```

### 删除菜单
```
DELETE /menus/{id}
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Menu deleted successfully",
  "data": null
}
```

### 获取菜单列表
```
GET /menus
```

**请求参数:**
- page: 页码 (默认: 1)
- page_size: 每页数量 (默认: 10)

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Menus retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "name": "string",
        "title": "string",
        "icon": "string",
        "path": "string",
        "component": "string",
        "redirect": "string",
        "permission": "string",
        "parent_id": 0,
        "sort": 0,
        "status": 1,
        "hidden": 0,
        "created_at": "2023-01-01T00:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

### 获取菜单树
```
GET /menus/tree
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Menu tree retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "string",
      "title": "string",
      "icon": "string",
      "path": "string",
      "component": "string",
      "redirect": "string",
      "permission": "string",
      "parent_id": 0,
      "sort": 0,
      "status": 1,
      "hidden": 0,
      "children": []
    }
  ]
}
```

## 系统日志接口

### 获取日志信息
```
GET /logs/{id}
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Log retrieved successfully",
  "data": {
    "id": 1,
    "level": "string",
    "method": "string",
    "path": "string",
    "status_code": 200,
    "client_ip": "string",
    "user_agent": "string",
    "request_id": "string",
    "user_id": 1,
    "username": "string",
    "message": "string",
    "request_body": "string",
    "error_detail": "string",
    "response": "string",
    "latency": 100,
    "created_at": "2023-01-01T00:00:00Z"
  }
}
```

### 获取日志列表
```
GET /logs
```

**请求参数:**
- page: 页码 (默认: 1)
- page_size: 每页数量 (默认: 10)
- level: 日志级别
- method: HTTP方法
- path: 请求路径
- username: 用户名

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Logs retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "level": "string",
        "method": "string",
        "path": "string",
        "status_code": 200,
        "client_ip": "string",
        "user_agent": "string",
        "request_id": "string",
        "user_id": 1,
        "username": "string",
        "message": "string",
        "request_body": "string",
        "error_detail": "string",
        "response": "string",
        "latency": 100,
        "created_at": "2023-01-01T00:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

### 删除日志
```
DELETE /logs/{id}
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Log deleted successfully",
  "data": null
}
```

### 清空日志
```
POST /logs/clear
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Logs cleared successfully",
  "data": null
}
```

## 字典管理接口

### 创建字典
```
POST /dictionaries
```

**请求头:**
```
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "name": "string",
  "title": "string",
  "description": "string"
}
```

**响应:**
```json
{
  "code": 200,
  "message": "Dictionary created successfully",
  "data": {
    "id": 1,
    "name": "string",
    "title": "string",
    "description": "string",
    "status": 1,
    "created_at": "2023-01-01T00:00:00Z"
  }
}
```

### 获取字典信息
```
GET /dictionaries/{id}
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Dictionary retrieved successfully",
  "data": {
    "id": 1,
    "name": "string",
    "title": "string",
    "description": "string",
    "status": 1,
    "created_at": "2023-01-01T00:00:00Z"
  }
}
```

### 更新字典信息
```
PUT /dictionaries/{id}
```

**请求头:**
```
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "name": "string",
  "title": "string",
  "description": "string"
}
```

**响应:**
```json
{
  "code": 200,
  "message": "Dictionary updated successfully",
  "data": {
    "id": 1,
    "name": "string",
    "title": "string",
    "description": "string",
    "status": 1,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
}
```

### 删除字典
```
DELETE /dictionaries/{id}
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Dictionary deleted successfully",
  "data": null
}
```

### 获取字典列表
```
GET /dictionaries
```

**请求参数:**
- page: 页码 (默认: 1)
- page_size: 每页数量 (默认: 10)

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Dictionaries retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "name": "string",
        "title": "string",
        "description": "string",
        "status": 1,
        "created_at": "2023-01-01T00:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

### 创建字典项
```
POST /dictionaries/{dict_id}/items
```

**请求头:**
```
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "label": "string",
  "value": "string",
  "sort": 0
}
```

**响应:**
```json
{
  "code": 200,
  "message": "Dictionary item created successfully",
  "data": {
    "id": 1,
    "dictionary_id": 1,
    "label": "string",
    "value": "string",
    "sort": 0,
    "status": 1,
    "created_at": "2023-01-01T00:00:00Z"
  }
}
```

### 获取字典项信息
```
GET /dictionaries/{dict_id}/items/{item_id}
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Dictionary item retrieved successfully",
  "data": {
    "id": 1,
    "dictionary_id": 1,
    "label": "string",
    "value": "string",
    "sort": 0,
    "status": 1,
    "created_at": "2023-01-01T00:00:00Z"
  }
}
```

### 更新字典项信息
```
PUT /dictionaries/{dict_id}/items/{item_id}
```

**请求头:**
```
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "label": "string",
  "value": "string",
  "sort": 0
}
```

**响应:**
```json
{
  "code": 200,
  "message": "Dictionary item updated successfully",
  "data": {
    "id": 1,
    "dictionary_id": 1,
    "label": "string",
    "value": "string",
    "sort": 0,
    "status": 1,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
}
```

### 删除字典项
```
DELETE /dictionaries/{dict_id}/items/{item_id}
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Dictionary item deleted successfully",
  "data": null
}
```

### 获取字典项列表
```
GET /dictionaries/{dict_id}/items
```

**请求参数:**
- page: 页码 (默认: 1)
- page_size: 每页数量 (默认: 10)

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Dictionary items retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "dictionary_id": 1,
        "label": "string",
        "value": "string",
        "sort": 0,
        "status": 1,
        "created_at": "2023-01-01T00:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

### 获取所有字典项
```
GET /dictionaries/{dict_id}/items-all
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "All dictionary items retrieved successfully",
  "data": [
    {
      "id": 1,
      "dictionary_id": 1,
      "label": "string",
      "value": "string",
      "sort": 0,
      "status": 1,
      "created_at": "2023-01-01T00:00:00Z"
    }
  ]
}
```

## 文件管理接口

### 上传文件
```
POST /files/upload
```

**请求头:**
```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**请求参数:**
- file: 文件

**响应:**
```json
{
  "code": 200,
  "message": "File uploaded successfully",
  "data": {
    "id": 1,
    "name": "string",
    "path": "string",
    "size": 1024,
    "mime_type": "string",
    "created_by": 1,
    "created_at": "2023-01-01T00:00:00Z"
  }
}
```

### 获取文件信息
```
GET /files/{id}
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "File retrieved successfully",
  "data": {
    "id": 1,
    "name": "string",
    "path": "string",
    "size": 1024,
    "mime_type": "string",
    "created_by": 1,
    "created_at": "2023-01-01T00:00:00Z"
  }
}
```

### 获取文件列表
```
GET /files
```

**请求参数:**
- page: 页码 (默认: 1)
- page_size: 每页数量 (默认: 10)

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Files retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "name": "string",
        "path": "string",
        "size": 1024,
        "mime_type": "string",
        "created_by": 1,
        "created_at": "2023-01-01T00:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

### 删除文件
```
DELETE /files/{id}
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "File deleted successfully",
  "data": null
}
```

### 下载文件
```
GET /files/{id}/download
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
文件流

## 通知公告接口

### 创建通知
```
POST /notifications
```

**请求头:**
```
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "title": "string",
  "content": "string",
  "type": "string",
  "status": "string",
  "start_date": "2023-01-01",
  "end_date": "2023-01-31"
}
```

**响应:**
```json
{
  "code": 200,
  "message": "Notification created successfully",
  "data": {
    "id": 1,
    "title": "string",
    "content": "string",
    "type": "string",
    "status": "string",
    "start_date": "2023-01-01T00:00:00Z",
    "end_date": "2023-01-31T00:00:00Z",
    "created_by": 1,
    "created_at": "2023-01-01T00:00:00Z"
  }
}
```

### 获取通知信息
```
GET /notifications/{id}
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Notification retrieved successfully",
  "data": {
    "id": 1,
    "title": "string",
    "content": "string",
    "type": "string",
    "status": "string",
    "start_date": "2023-01-01T00:00:00Z",
    "end_date": "2023-01-31T00:00:00Z",
    "created_by": 1,
    "created_at": "2023-01-01T00:00:00Z"
  }
}
```

### 更新通知信息
```
PUT /notifications/{id}
```

**请求头:**
```
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "title": "string",
  "content": "string",
  "type": "string",
  "status": "string",
  "start_date": "2023-01-01",
  "end_date": "2023-01-31"
}
```

**响应:**
```json
{
  "code": 200,
  "message": "Notification updated successfully",
  "data": {
    "id": 1,
    "title": "string",
    "content": "string",
    "type": "string",
    "status": "string",
    "start_date": "2023-01-01T00:00:00Z",
    "end_date": "2023-01-31T00:00:00Z",
    "created_by": 1,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
}
```

### 删除通知
```
DELETE /notifications/{id}
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Notification deleted successfully",
  "data": null
}
```

### 获取通知列表
```
GET /notifications
```

**请求参数:**
- page: 页码 (默认: 1)
- page_size: 每页数量 (默认: 10)
- status: 状态
- type: 类型

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Notifications retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "title": "string",
        "content": "string",
        "type": "string",
        "status": "string",
        "start_date": "2023-01-01T00:00:00Z",
        "end_date": "2023-01-31T00:00:00Z",
        "created_by": 1,
        "created_at": "2023-01-01T00:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

### 获取活跃通知
```
GET /notifications/active
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Active notifications retrieved successfully",
  "data": [
    {
      "id": 1,
      "title": "string",
      "content": "string",
      "type": "string",
      "status": "string",
      "start_date": "2023-01-01T00:00:00Z",
      "end_date": "2023-01-31T00:00:00Z",
      "created_by": 1,
      "created_at": "2023-01-01T00:00:00Z"
    }
  ]
}
```

## 系统监控接口

### 获取系统信息
```
GET /monitor/info
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "System information retrieved successfully",
  "data": {
    "os": "string",
    "arch": "string",
    "go_version": "string",
    "app_version": "string",
    "uptime": "string"
  }
}
```

### 获取系统指标
```
GET /monitor/metrics
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "System metrics retrieved successfully",
  "data": {
    "id": 1,
    "timestamp": "2023-01-01T00:00:00Z",
    "cpu_usage": 25.5,
    "memory_usage": 30.2,
    "disk_usage": 45.7,
    "network_inbound": 1024.0,
    "network_outbound": 2048.0,
    "request_count": 1000,
    "error_count": 5,
    "created_at": "2023-01-01T00:00:00Z"
  }
}
```

### 获取最近指标
```
GET /monitor/recent
```

**请求参数:**
- count: 数量 (默认: 10, 最大: 100)

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Recent metrics retrieved successfully",
  "data": [
    {
      "id": 1,
      "timestamp": "2023-01-01T00:00:00Z",
      "cpu_usage": 25.5,
      "memory_usage": 30.2,
      "disk_usage": 45.7,
      "network_inbound": 1024.0,
      "network_outbound": 2048.0,
      "request_count": 1000,
      "error_count": 5,
      "created_at": "2023-01-01T00:00:00Z"
    }
  ]
}
```

## 定时任务接口

### 创建任务
```
POST /tasks
```

**请求头:**
```
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "name": "string",
  "description": "string",
  "cron_expr": "string",
  "handler": "string"
}
```

**响应:**
```json
{
  "code": 200,
  "message": "Task created successfully",
  "data": {
    "id": 1,
    "name": "string",
    "description": "string",
    "cron_expr": "string",
    "handler": "string",
    "status": "active",
    "last_run": "2023-01-01T00:00:00Z",
    "next_run": "2023-01-01T00:00:00Z",
    "created_by": 1,
    "created_at": "2023-01-01T00:00:00Z"
  }
}
```

### 获取任务信息
```
GET /tasks/{id}
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Task retrieved successfully",
  "data": {
    "id": 1,
    "name": "string",
    "description": "string",
    "cron_expr": "string",
    "handler": "string",
    "status": "active",
    "last_run": "2023-01-01T00:00:00Z",
    "next_run": "2023-01-01T00:00:00Z",
    "created_by": 1,
    "created_at": "2023-01-01T00:00:00Z"
  }
}
```

### 更新任务信息
```
PUT /tasks/{id}
```

**请求头:**
```
Authorization: Bearer <token>
```

**请求参数:**
```json
{
  "name": "string",
  "description": "string",
  "cron_expr": "string",
  "handler": "string",
  "status": "string"
}
```

**响应:**
```json
{
  "code": 200,
  "message": "Task updated successfully",
  "data": {
    "id": 1,
    "name": "string",
    "description": "string",
    "cron_expr": "string",
    "handler": "string",
    "status": "string",
    "last_run": "2023-01-01T00:00:00Z",
    "next_run": "2023-01-01T00:00:00Z",
    "created_by": 1,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
}
```

### 删除任务
```
DELETE /tasks/{id}
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Task deleted successfully",
  "data": null
}
```

### 获取任务列表
```
GET /tasks
```

**请求参数:**
- page: 页码 (默认: 1)
- page_size: 每页数量 (默认: 10)
- status: 状态

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Tasks retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "name": "string",
        "description": "string",
        "cron_expr": "string",
        "handler": "string",
        "status": "string",
        "last_run": "2023-01-01T00:00:00Z",
        "next_run": "2023-01-01T00:00:00Z",
        "created_by": 1,
        "created_at": "2023-01-01T00:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

### 立即执行任务
```
POST /tasks/{id}/run
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Task started successfully",
  "data": null
}
```

## 数据导入导出接口

### 导出用户数据
```
GET /export/users
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
Excel文件流

### 导入用户数据
```
POST /import/users
```

**请求参数:**
- file: Excel文件
- has_headers: 是否包含表头 (默认: true)

**请求头:**
```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**响应:**
```json
{
  "code": 200,
  "message": "Users imported successfully",
  "data": {
    "row_count": 100,
    "data": []
  }
}
```

### 导出自定义数据
```
GET /export/data
```

**请求参数:**
- sheet_name: 工作表名称 (默认: Data)

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
Excel文件流

## 缓存管理接口

### 获取缓存统计
```
GET /cache/stats
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Cache statistics retrieved successfully",
  "data": {
    "hits": 1000,
    "misses": 100,
    "hit_rate": 0.909,
    "size": 50
  }
}
```

### 重置缓存统计
```
POST /cache/reset-stats
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Cache statistics reset successfully",
  "data": null
}
```

### 清空缓存
```
POST /cache/clear
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Cache cleared successfully",
  "data": null
}
```

## 数据库管理接口

### 获取数据库统计
```
GET /db/stats
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Database statistics retrieved successfully",
  "data": {
    "max_open_connections": 100,
    "open_connections": 10,
    "in_use": 5,
    "idle": 5,
    "wait_count": 0,
    "wait_duration": 0,
    "max_idle_closed": 0,
    "max_lifetime_closed": 0,
    "max_idle_time_closed": 0
  }
}
```

## 日志级别管理接口

### 获取日志级别
```
GET /log/level
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Current log level retrieved successfully",
  "data": {
    "level": "info"
  }
}
```

### 设置日志级别
```
POST /log/level
```

**请求参数:**
```json
{
  "level": "string"
}
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Log level updated successfully",
  "data": {
    "level": "debug"
  }
}
```

## 安全管理接口

### 获取CSRF令牌
```
GET /security/csrf-token
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "CSRF token generated successfully",
  "data": {
    "token": "csrf_token"
  }
}
```

### 获取频率限制配置
```
GET /security/rate-limit-config
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Rate limit configuration retrieved successfully",
  "data": {
    "requests_per_minute": 60,
    "burst": 10
  }
}
```

## 健康检查接口

### 获取健康状态
```
GET /health
```

**响应:**
```json
{
  "code": 200,
  "message": "Health status retrieved successfully",
  "data": {
    "status": "healthy",
    "timestamp": "2023-01-01T00:00:00Z",
    "components": {
      "database": {
        "status": true,
        "details": {}
      },
      "cache": {
        "status": true
      },
      "system": {
        "goroutines": 10,
        "memory": {
          "allocated": "N/A",
          "system": "N/A"
        }
      }
    }
  }
}
```

### 获取详细健康状态
```
GET /health/detailed
```

**响应:**
```json
{
  "code": 200,
  "message": "Detailed health status retrieved successfully",
  "data": {
    "status": "healthy",
    "timestamp": "2023-01-01T00:00:00Z",
    "components": {
      "database": {
        "status": true,
        "details": {}
      },
      "cache": {
        "status": true
      },
      "system": {
        "goroutines": 10,
        "memory": {
          "allocated": "N/A",
          "system": "N/A"
        }
      }
    },
    "uptime": "1h30m"
  }
}
```

## 配置管理接口

### 获取配置信息
```
GET /config
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Configuration retrieved successfully",
  "data": {
    "app": {
      "name": "go-admin",
      "env": "local",
      "port": "8080"
    },
    "db": {
      "host": "localhost",
      "port": "3306",
      "user": "root",
      "name": "go_admin"
    },
    "log": {
      "level": "info",
      "output": "console"
    },
    "cache": {
      "maxsize": 10000,
      "gcinterval": "10m0s"
    }
  }
}
```

### 重新加载配置
```
POST /config/reload
```

**请求头:**
```
Authorization: Bearer <token>
```

**响应:**
```json
{
  "code": 200,
  "message": "Configuration reloaded successfully",
  "data": {
    "app": {
      "name": "go-admin",
      "env": "local",
      "port": "8080"
    }
  }
}
```