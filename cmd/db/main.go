package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hadv/generator"
	"github.com/hadv/mysql-db-index-benchmark/core/model"
	"github.com/hadv/mysql-db-index-benchmark/core/repo"
	"github.com/hadv/mysql-db-index-benchmark/core/service"
	"github.com/hadv/mysql-db-index-benchmark/utils"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"gopkg.in/urfave/cli.v1"
)

var app *cli.App

func init() {
	app = utils.NewApp()
	app.Commands = []cli.Command{
		commandBulkInsert,
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var commandBulkInsert = cli.Command{
	Name:        "bulk-insert",
	Usage:       "",
	ArgsUsage:   "",
	Description: "",
	Flags:       []cli.Flag{},
	Action: func(ctx *cli.Context) error {
		viper.AutomaticEnv()
		dbURL := viper.GetString("DB_URL")
		fmt.Println(dbURL)
		db, err := sqlx.Connect("mysql", dbURL)
		db.Exec("USE testindex")
		if err != nil {
			log.Fatalf("Cannot connect to MySQL at %v: %v", dbURL, err)
		}
		account := service.NewAccount(repo.NewUser(db))
		users := make([]*model.User, 0)
		for i := 0; i < 1000000; i++ {
			generator := generator.New(20, "usr", "_")
			id, _ := generator.Get()
			user := &model.User{
				ID:              id,
				Firstname:       "Ha" + strconv.Itoa(i),
				Lastname:        "Dang" + strconv.Itoa(i),
				Email:           "dvietha" + strconv.Itoa(i) + "@gmail.com",
				Password:        "123456",
				ConfirmPassword: "123456",
			}

			users = append(users, user)

		}
		if err := account.BulkInsert(context.Background(), users, 10000); err != nil {
			return err
		}
		return nil
	},
}
