package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/pressly/goose"
	"github.com/spf13/cobra"
	_ "gorm.io/driver/mysql"
	"messaging/utils/config"
)

var usageCommands = `
Run database migrations & seeder

Usage:
    [command]

Available Migration Commands:
    up                   Migrate the DB to the most recent version available
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with next version

`

func main() {
	var rootCmd = &cobra.Command{
		Use:   "migrate",
		Short: "MySql Migration Service",

		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				err := cmd.Help()
				if err != nil {
					log.Fatal(err)
				}
				os.Exit(0)
			}
			mysql := config.LoadDBConfig("mongodb")
			err := goose.SetDialect("mongodb")
			if err != nil {
				log.Fatal(err)
			}
			db, err := sql.Open("mongodb", mysql.User+":"+mysql.Password+"@tcp("+mysql.Host+":"+strconv.Itoa(mysql.Port)+")/"+mysql.DBName+"?charset=utf8&parseTime=True&loc=Local")
			if err != nil {
				log.Fatal(err)
			}

			appPath, _ := os.Getwd()
			dir := appPath + "/app/migration/mongodb/migration"
			if len(args) == 0 {
				err := cmd.Help()
				if err != nil {
					log.Fatal(err)
				}
				os.Exit(0)
			}
			command := args[0]
			cmdArgs := args[1:]
			if err := goose.Run(command, db, dir, cmdArgs...); err != nil {
				log.Fatalf("goose run: %v", err)
			}
		},
	}

	rootCmd.SetUsageTemplate(usageCommands)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
