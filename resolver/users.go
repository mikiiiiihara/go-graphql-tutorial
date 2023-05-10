package resolver

import (
	"fmt"

	"github.com/graphql-go/graphql"
)

// ユーザー全取得
func (r *Resolver) GetAllUsers(p graphql.ResolveParams) (interface{}, error) {
	var users []User
	r.db.Find(&users)

	// コンソールにユーザーを出力
	for _, user := range users {
		fmt.Printf("ID: %d, Name: %s, Email: %s\n", user.ID, user.Name, user.Email)
	}

	return users, nil
}

// IDでユーザー検索
func (r *Resolver) GetUserByID(p graphql.ResolveParams) (interface{}, error) {
	id := p.Args["id"].(int)
	var user User
	r.db.First(&user, id)
	return user, nil
}