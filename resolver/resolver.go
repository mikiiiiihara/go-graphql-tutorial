package resolver

import (
	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
)


type Resolver struct {
	db *gorm.DB
}

func NewResolver(db *gorm.DB) *Resolver {
	return &Resolver{db: db}
}

func (r *Resolver) CreateSchema() (graphql.Schema, error) {
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
	},)

	tickerType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Ticker",
		Fields: graphql.Fields{
			"product_code":      &graphql.Field{Type: graphql.String},
			"state":             &graphql.Field{Type: graphql.String},
			"timestamp":         &graphql.Field{Type: graphql.String},
			"tick_id":           &graphql.Field{Type: graphql.Int},
			"best_bid":          &graphql.Field{Type: graphql.Float},
			"best_ask":          &graphql.Field{Type: graphql.Float},
			"best_bid_size":     &graphql.Field{Type: graphql.Float},
			"best_ask_size":     &graphql.Field{Type: graphql.Float},
			"total_bid_depth":   &graphql.Field{Type: graphql.Float},
			"total_ask_depth":   &graphql.Field{Type: graphql.Float},
			"market_bid_size":   &graphql.Field{Type: graphql.Float},
			"market_ask_size":   &graphql.Field{Type: graphql.Float},
			"ltp":               &graphql.Field{Type: graphql.Float},
			"volume":            &graphql.Field{Type: graphql.Float},
			"volume_by_product": &graphql.Field{Type: graphql.Float},
		},
	},)

	rootQuery := graphql.NewObject(
		graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"users": &graphql.Field{
				Type:    graphql.NewList(userType),
				Resolve: r.GetAllUsers,
			},
			"user": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: r.GetUserByID,
			},
			"cryptoTicker":&graphql.Field{
				Type: tickerType,
				Resolve: r.GetCryptoTickers,
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
				Resolve: r.CreateUser,
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
				Resolve: r.UpdateUser,
			},
			"deleteUser": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: r.DeleteUser,
			},
		},
	})	

	schemaConfig := graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	}

	return graphql.NewSchema(schemaConfig)
}
