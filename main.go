package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/graphql-go/handler"
	"github.com/joho/godotenv"
	"github.com/mikiiiiihara/go-graphql-tutorial/resolver"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var users []resolver.User
var db *gorm.DB

func main() {
	envImportErr := godotenv.Load()
	if envImportErr != nil {
		fmt.Println("Error loading .env file")
		return
	}
	dsn := "host=localhost user=myuser password=mypassword dbname=mydbname port=5432 sslmode=disable TimeZone=Asia/Tokyo"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&resolver.User{})

	resolver := resolver.NewResolver(db)
	schema, err := resolver.CreateSchema()
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	// エンドポイントを定義
	http.Handle("/graphql", h)
	fmt.Println("Listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
