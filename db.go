package dbconnector

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库实例
var DB *gorm.DB

// registeredModels 已注册的模型列表
var registeredModels []interface{}

// InitDB 通过MysqlConfig初始化数据库连接
func InitDB(config *MysqlConfig) error {
	if config == nil {
		return fmt.Errorf("MysqlConfig is nil")
	}

	if config.DataSource == "" {
		return fmt.Errorf("DataSource is empty")
	}

	return initDBWithDSN(config.DataSource)
}

// InitDBWithDSN 通过完整DSN字符串初始化数据库连接
func InitDBWithDSN(dsn string) error {
	if dsn == "" {
		return fmt.Errorf("DSN is empty")
	}

	return initDBWithDSN(dsn)
}

// InitDBWithViper 通过viper分拆配置初始化数据库连接
func InitDBWithViper() error {
	// 构建MySQL连接DSN
	dsn, err := buildDSNFromViper()
	if err != nil {
		return fmt.Errorf("构建DSN失败: %v", err)
	}

	return initDBWithDSN(dsn)
}

// GetDB 获取数据库实例，并确保连接有效
func GetDB() *gorm.DB {
	if DB == nil {
		log.Fatalf("数据库未初始化，请先调用InitDB或相关初始化函数")
	}

	// 检查连接是否有效
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("获取底层sql.DB失败: %v", err)
	}

	// 使用Ping检查连接是否活跃
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("数据库连接已断开: %v", err)
	}

	return DB
}

// RegisterModels 注册需要自动迁移的模型
func RegisterModels(models ...interface{}) {
	registeredModels = append(registeredModels, models...)
}

// 内部方法：通过DSN初始化数据库连接
func initDBWithDSN(dsn string) error {
	// 配置gorm日志
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// 连接MySQL数据库
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("连接MySQL数据库失败: %v", err)
	}

	// 配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取底层sql.DB失败: %v", err)
	}

	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(10)
	// 设置最大打开连接数
	sqlDB.SetMaxOpenConns(100)
	// 设置连接最大生命周期
	sqlDB.SetConnMaxLifetime(1 * time.Hour)

	// 自动迁移数据库表
	if err := autoMigrate(); err != nil {
		return fmt.Errorf("自动迁移数据库表失败: %v", err)
	}

	log.Println("数据库连接成功")
	return nil
}

// 内部方法：从viper配置构建DSN
func buildDSNFromViper() (string, error) {
	// 从viper获取配置
	host := viper.GetString("database.host")
	port := viper.GetInt("database.port")
	user := viper.GetString("database.user")
	password := viper.GetString("database.password")
	dbname := viper.GetString("database.dbname")
	charset := viper.GetString("database.charset")
	parseTime := viper.GetBool("database.parseTime")
	loc := viper.GetString("database.loc")

	// 验证必要的配置项
	if host == "" {
		return "", fmt.Errorf("database.host is empty")
	}
	if user == "" {
		return "", fmt.Errorf("database.user is empty")
	}
	if dbname == "" {
		return "", fmt.Errorf("database.dbname is empty")
	}

	// 设置默认值
	if port == 0 {
		port = 3306
	}
	if charset == "" {
		charset = "utf8mb4"
	}
	if loc == "" {
		loc = "Local"
	}

	// 构建DSN，兼容无密码情况
	var dsn string
	if password != "" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
			user, password, host, port, dbname, charset, parseTime, loc)
	} else {
		dsn = fmt.Sprintf("%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
			user, host, port, dbname, charset, parseTime, loc)
	}

	return dsn, nil
}

// 内部方法：自动迁移数据库表
func autoMigrate() error {
	if len(registeredModels) > 0 {
		if err := DB.AutoMigrate(registeredModels...); err != nil {
			return err
		}
		log.Printf("自动迁移完成，共迁移 %d 个模型", len(registeredModels))
	}
	return nil
}
