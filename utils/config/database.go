package config

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var (
	//DB DB
	DB *mongo.Database
)

//Database Database
type Database struct {
	Host      string
	User      string
	Password  string
	DBName    string
	DBNumber  int
	Port      int
	DebugMode bool
}

// LoadDBConfig load database configuration
func LoadDBConfig(name string) Database {
	db := viper.Sub("database." + name)
	conf := Database{
		Host:      db.GetString("host"),
		User:      db.GetString("user"),
		Password:  db.GetString("password"),
		DBName:    db.GetString("db_name"),
		Port:      db.GetInt("port"),
		DebugMode: db.GetBool("debug"),
	}
	return conf
}

func MongoConnect() *mongo.Database {
	if viper.Get("env") != "testing" {
		config := LoadDBConfig("mongo")
		mongoURL := fmt.Sprintf("mongodb://%s:%d", config.Host, config.Port)
		cmdMonitor := &event.CommandMonitor{
			Started: func(_ context.Context, evt *event.CommandStartedEvent) {
				log.Print(evt.Command)
			},
		}
		mongoOptions := options.Client().ApplyURI(mongoURL).SetMonitor(cmdMonitor)

		client, err := mongo.Connect(context.Background(), mongoOptions)
		if err != nil {
			panic(fmt.Sprintf("Unable to connect to database : %v", err))
		}

		DB = client.Database(config.DBName)
		collections, err := DB.ListCollectionNames(context.Background(), bson.D{})
		if err != nil {
			return nil
		}

		fmt.Println(fmt.Sprintf("Collections : %s", collections))
		return DB
	}
	return nil
}
