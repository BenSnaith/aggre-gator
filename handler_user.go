package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/BenSnaith/aggre-gator/internal/database"
	"github.com/google/uuid"
	"os"
	"time"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return errors.New("usage: login <name>")
	}

	name := cmd.Args[0]

	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return errors.New("name is not present is database")
	}

	err = s.conf.SetUser(name)
	if err != nil {
		return err
	}

	fmt.Println("Username successfully changed!")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return errors.New("usage: register <name>")
	}

	name := cmd.Args[0]

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
	})
	if err != nil {
		os.Exit(1)
	}

	err = s.conf.SetUser(name)
	if err != nil {
		return err
	}

	fmt.Println("User was created successfully!")
	printUser(user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		os.Exit(1)
		return err
	}

	fmt.Println("User table reset successfully!")
	os.Exit(0)
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetAllUsers(context.Background())
	if err != nil {
		os.Exit(1)
		return err
	}

	for _, user := range users {
		if user.Name == s.conf.CurrentUserName {
			fmt.Printf(" * %s (current)\n", user.Name)
		} else {
			fmt.Printf(" * %s\n", user.Name)
		}
	}

	os.Exit(0)
	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}
