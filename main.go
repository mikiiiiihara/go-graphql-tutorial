package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

var users []User
var db *gorm.DB

func main() {
	dsn := "host=localhost user=myuser password=mypassword dbname=mydbname port=5432 sslmode=disable TimeZone=Asia/Tokyo"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&User{})

	userType := graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"email": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			// Query to get all users
			"users": &graphql.Field{
				Type: graphql.NewList(userType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					db.Find(&users)
					return users, nil
				},
			},
			// Query to get a user by ID
			"user": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					var user User
					db.First(&user, id)
					return user, nil
				},
			},
		},
	})
	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createUser": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					user := User{
						Name:     params.Args["name"].(string),
						Email:    params.Args["email"].(string),
						Password: params.Args["password"].(string),
					}
					db.Create(&user)
					return user, nil
				},
			},
			"updateUser": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"email": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id := params.Args["id"].(int)
					var user User
					db.First(&user, id)
	
					if name, ok := params.Args["name"].(string); ok {
						user.Name = name
					}
	
					if email, ok := params.Args["email"].(string); ok {
						user.Email = email
					}
	
					if password, ok := params.Args["password"].(string); ok {
						user.Password = password
					}
	
					db.Save(&user)
					return user, nil
				},
			},
			"deleteUser": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id := params.Args["id"].(int)
					var user User
					db.First(&user, id)
					db.Delete(&user)
					return user, nil
				},
			},
		},
	})	

	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery,
		Mutation: rootMutation,
	})
	// GraphQLハンドラーを設定
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
