package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/AirflowSteve/blog-aggregator/internal/config"
	"github.com/AirflowSteve/blog-aggregator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("There's an error reading the config: %v\n", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	dbQueries := database.New(db)

	state := &State{
		db:     dbQueries,
		config: &cfg,
	}

	cmds := commands{
		handlerFunctions: make(map[string]func(*State, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)

	arguments := os.Args
	if len(arguments) < 2 {
		fmt.Println("not enough arguments")
		os.Exit(1)
	}

	cmd := command{
		Name: arguments[1],
		Args: arguments[2:],
	}

	err = cmds.run(state, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", cfg)
}
