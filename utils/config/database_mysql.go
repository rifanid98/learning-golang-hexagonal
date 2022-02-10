package config

// import (
// 	"fmt"
// 	"time"

// 	"github.com/spf13/viper"
// 	"gorm.io/gorm"

// 	//driver
// 	"gorm.io/driver/mysql"
// )

// var (
// 	//DB DB
// 	DB *gorm.DB
// )

// //Database Database
// type Database struct {
// 	Host      string
// 	User      string
// 	Password  string
// 	DBName    string
// 	DBNumber  int
// 	Port      int
// 	DebugMode bool
// }

// // LoadDBConfig load database configuration
// func LoadDBConfig(name string) Database {
// 	db := viper.Sub("database." + name)
// 	conf := Database{
// 		Host:      db.GetString("host"),
// 		User:      db.GetString("user"),
// 		Password:  db.GetString("password"),
// 		DBName:    db.GetString("db_name"),
// 		Port:      db.GetInt("port"),
// 		DebugMode: db.GetBool("debug"),
// 	}
// 	return conf
// }

// func OpenDbPool() {
// 	DB = MysqlConnect("mysql")
// 	pool := viper.Sub("database.mysql.pool")
// 	db, err := DB.DB()
// 	if err != nil {
// 		panic(err)
// 	}
// 	db.SetMaxOpenConns(pool.GetInt("maxOpenConns"))
// 	db.SetMaxIdleConns(pool.GetInt("maxIdleConns"))
// 	db.SetConnMaxLifetime(pool.GetDuration("maxLifetime") * time.Second)
// }

// // MysqlConnect connect to mysql using setting name. return *gorm.DB incstance
// func MysqlConnect(configName string) *gorm.DB {
// 	mysqlConfig := LoadDBConfig(configName)

// 	dsn := fmt.Sprintf(
// 		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
// 		mysqlConfig.User,
// 		mysqlConfig.Password,
// 		mysqlConfig.Host,
// 		mysqlConfig.Port,
// 		mysqlConfig.DBName,
// 	)
// 	connection, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

// 	//connection, err := gorm.Open("mysql", mysqlConfig.User+":"+mysqlConfig.Password+"@tcp("+mysqlConfig.Host+":"+strconv.Itoa(mysqlConfig.Port)+")/"+mysqlConfig.DBName+"?charset=utf8&parseTime=True&loc=Local")
// 	// enable debug
// 	if err != nil {
// 		panic(err)
// 	}

// 	//connection.LogMode(true)

// 	if mysqlConfig.DebugMode {
// 		return connection.Debug()
// 	}

// 	return connection
// }
