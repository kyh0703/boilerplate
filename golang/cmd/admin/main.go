package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/kyh0703/template/internal/core/domain/model"
	_ "modernc.org/sqlite"
)

func main() {
	email := flag.String("email", "", "User email to promote to admin")
	demote := flag.Bool("demote", false, "Demote admin to regular user")
	dbPath := flag.String("db", "flow.db", "Path to SQLite database file")
	flag.Parse()

	if *email == "" {
		log.Fatal("Email is required. Usage: admin -email=user@example.com")
	}

	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?cache=shared&mode=rwc", *dbPath))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	queries := model.New(db)

	user, err := queries.GetUserByEmail(context.Background(), *email)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Fatalf("User not found with email: %s", *email)
		}
		log.Fatalf("Failed to get user: %v", err)
	}

	isAdmin := 1
	if *demote {
		isAdmin = 0
	}

	err = queries.PatchUser(context.Background(), model.PatchUserParams{
		ID:      user.ID,
		IsAdmin: sql.NullInt64{Int64: int64(isAdmin), Valid: true},
	})
	if err != nil {
		log.Fatalf("Failed to update user: %v", err)
	}

	action := "promoted to"
	if *demote {
		action = "demoted from"
	}
	fmt.Printf("User %s has been %s admin\n", *email, action)
}
