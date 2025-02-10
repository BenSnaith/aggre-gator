package main

// import the postgres driver
import (
	"database/sql"
	"github.com/BenSnaith/aggre-gator/internal/database"
	_ "github.com/lib/pq"
)

import (
	"log"
	"os"

	"github.com/BenSnaith/aggre-gator/internal/config"
)

type state struct {
	conf *config.Config
	db   *database.Queries
}

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Fatalf("%v", err)
	}

	programState := &state{
		conf: &conf,
	}
	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)

	db, err := sql.Open("postgres", programState.conf.DbURL)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	dbQueries := database.New(db)
	programState.db = dbQueries

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
		return
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}
