package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"test-post-article/handler"
	"test-post-article/repositories"
	"test-post-article/router"
	"test-post-article/services"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	dbInit, err := sql.Open("mysql", "root:@tcp(localhost:3306)/")

	if err != nil {
		panic(fmt.Sprintf("Failed to connect MySQL: %v", err))
	}

	_, err = dbInit.Exec("CREATE DATABASE IF NOT EXISTS article")

	if err != nil {
		panic(fmt.Sprintf("Failed to create database: %v", err))
	}
	dbInit.Close()

	m, err := migrate.New(
		"file://db/migrations",
		"mysql://root:@tcp(localhost:3306)/article",
	)

	if err != nil {
		panic(fmt.Sprintf("Failed to init the migration: %v", err))
	}

	err = m.Up()

	if err != nil && err != migrate.ErrNoChange {
		panic(fmt.Sprintf("Failed to run the migration: %v", err))
	}

	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/article?charset=utf8mb4&parseTime=true&loc=Local")

	if err != nil {
		panic(fmt.Sprintf("Failed to connect database: %v", err))
	}

	defer db.Close()

	repo := repositories.NewPostRepository(db)
	svc := services.NewPostService(repo)
	h := handler.NewPostHandler(svc)
	r := router.NewRouter()

	r.GET("/api/v1/articles/{limit}/{offset}", h.GetAll)
	r.POST("/api/v1/articles", h.Create)
	r.GET("/api/v1/articles/{id}", h.GetById)
	r.PUT("/api/v1/articles/{id}", h.Update)
	r.DELETE("/api/v1/articles/{id}", h.Delete)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}
