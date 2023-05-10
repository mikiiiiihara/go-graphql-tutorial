package resolver

import (
	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
)


type Resolver struct {
	db *gorm.DB
}

type User struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

type Ticker struct {
	ProductCode      string  `json:"product_code"`
	State            string  `json:"state"`
	Timestamp        string  `json:"timestamp"`
	TickID           int     `json:"tick_id"`
	BestBid          float64 `json:"best_bid"`
	BestAsk          float64 `json:"best_ask"`
	BestBidSize      float64 `json:"best_bid_size"`
	BestAskSize      float64 `json:"best_ask_size"`
	TotalBidDepth    float64 `json:"total_bid_depth"`
	TotalAskDepth    float64 `json:"total_ask_depth"`
	MarketBidSize    float64 `json:"market_bid_size"`
	MarketAskSize    float64 `json:"market_ask_size"`
	Ltp              float64 `json:"ltp"`
	Volume           float64 `json:"volume"`
	VolumeByProduct  float64 `json:"volume_by_product"`
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
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					user := User{
						Name:     params.Args["name"].(string),
						Email:    params.Args["email"].(string),
						Password: params.Args["password"].(string),
					}
					r.db.Create(&user)
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
					r.db.First(&user, id)
	
					if name, ok := params.Args["name"].(string); ok {
						user.Name = name
					}
	
					if email, ok := params.Args["email"].(string); ok {
						user.Email = email
					}
	
					if password, ok := params.Args["password"].(string); ok {
						user.Password = password
					}
	
					r.db.Save(&user)
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
					r.db.First(&user, id)
					r.db.Delete(&user)
					return user, nil
				},
			},
		},
	})	

	schemaConfig := graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	}

	return graphql.NewSchema(schemaConfig)
}
