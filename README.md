# dbconnector

一个基于 GORM 和 Viper 的 MySQL 数据库连接管理库，提供简单易用的数据库连接初始化和管理功能。

## 功能特性

- 支持多种初始化方式：配置对象、DSN 字符串、Viper 配置
- 自动管理数据库连接池
- 支持模型自动迁移
- 提供数据库连接有效性检查
- 支持无密码连接配置

## 安装

```bash
go get github.com/xuxinzhi007/dbconnector
```

## 使用示例

### 1. 使用配置对象初始化

```go
package main

import (
    "github.com/xuxinzhi007/dbconnector"
)

func main() {
    // 创建配置
    config := &dbconnector.MysqlConfig{
        DataSource: "user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local",
    }
    
    // 初始化数据库
    if err := dbconnector.InitDB(config); err != nil {
        panic(err)
    }
    
    // 获取数据库实例
    db := dbconnector.GetDB()
    
    // 使用 db 进行数据库操作
}
```

### 2. 使用 DSN 字符串初始化

```go
package main

import (
    "github.com/xuxinzhi007/dbconnector"
)

func main() {
    dsn := "user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
    
    // 初始化数据库
    if err := dbconnector.InitDBWithDSN(dsn); err != nil {
        panic(err)
    }
    
    // 获取数据库实例
    db := dbconnector.GetDB()
    
    // 使用 db 进行数据库操作
}
```

### 3. 使用 Viper 配置初始化

```go
package main

import (
    "github.com/spf13/viper"
    "github.com/xuxinzhi007/dbconnector"
)

func main() {
    // 配置 Viper
    viper.SetConfigFile("config.yaml")
    if err := viper.ReadInConfig(); err != nil {
        panic(err)
    }
    
    // 初始化数据库
    if err := dbconnector.InitDBWithViper(); err != nil {
        panic(err)
    }
    
    // 获取数据库实例
    db := dbconnector.GetDB()
    
    // 使用 db 进行数据库操作
}
```

### 4. 注册模型并自动迁移

```go
package main

import (
    "github.com/xuxinzhi007/dbconnector"
)

// 定义模型
type User struct {
    ID   uint   `gorm:"primaryKey"`
    Name string `gorm:"size:255"`
    Age  int
}

func main() {
    // 注册模型
    dbconnector.RegisterModels(&User{})
    
    // 初始化数据库（任意方式）
    dsn := "user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
    if err := dbconnector.InitDBWithDSN(dsn); err != nil {
        panic(err)
    }
    
    // 数据库连接成功后会自动执行模型迁移
    db := dbconnector.GetDB()
    
    // 使用 db 进行数据库操作
}
```

## 配置说明

### Viper 配置文件格式

```yaml
# config.yaml
database:
  host: localhost
  port: 3306
  user: root
  password: password
  dbname: test_db
  charset: utf8mb4
  parseTime: true
  loc: Local
```

### 配置项说明

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| host | string | - | 数据库主机地址 |
| port | int | 3306 | 数据库端口 |
| user | string | - | 数据库用户名 |
| password | string | - | 数据库密码（可选） |
| dbname | string | - | 数据库名称 |
| charset | string | utf8mb4 | 数据库字符集 |
| parseTime | bool | false | 是否解析时间类型 |
| loc | string | Local | 时区设置 |

## API 文档

### 1. 初始化函数

#### `InitDB(config *MysqlConfig) error`
通过配置对象初始化数据库连接

- **参数**：
  - `config`：MySQL 配置对象，包含 DataSource 或分拆的配置项

#### `InitDBWithDSN(dsn string) error`
通过完整 DSN 字符串初始化数据库连接

- **参数**：
  - `dsn`：完整的 MySQL DSN 连接字符串

#### `InitDBWithViper() error`
通过 Viper 配置初始化数据库连接

- **说明**：
  - 从 Viper 配置中读取 `database` 前缀的配置项

### 2. 数据库操作

#### `GetDB() *gorm.DB`
获取数据库实例，并确保连接有效

- **返回值**：
  - `*gorm.DB`：GORM 数据库实例
  - 如果连接无效或未初始化，会触发 `log.Fatalf`

### 3. 模型管理

#### `RegisterModels(models ...interface{})`
注册需要自动迁移的模型

- **参数**：
  - `models`：可变参数，需要自动迁移的模型结构体指针

## 依赖

- [github.com/spf13/viper](https://github.com/spf13/viper) - 配置管理
- [gorm.io/driver/mysql](https://gorm.io/driver/mysql) - MySQL 驱动
- [gorm.io/gorm](https://gorm.io/gorm) - ORM 框架

## 许可证

MIT License
