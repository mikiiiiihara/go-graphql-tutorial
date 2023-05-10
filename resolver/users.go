package resolver

import (
	"fmt"

	"github.com/graphql-go/graphql"
)
type User struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}
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

// ユーザー登録
func (r *Resolver) CreateUser(p graphql.ResolveParams) (interface{}, error) {
	user := User{
		Name:     p.Args["name"].(string),
		Email:    p.Args["email"].(string),
		Password: p.Args["password"].(string),
	}
	r.db.Create(&user)
	return user, nil
}

// ユーザー更新
func (r *Resolver) UpdateUser(p graphql.ResolveParams) (interface{}, error) {
	id := p.Args["id"].(int)
	var user User
	r.db.First(&user, id)

	if name, ok := p.Args["name"].(string); ok {
		user.Name = name
	}

	if email, ok := p.Args["email"].(string); ok {
		user.Email = email
	}

	if password, ok := p.Args["password"].(string); ok {
		user.Password = password
	}

	r.db.Save(&user)
	return user, nil
}

// ユーザー削除
func (r *Resolver) DeleteUser(p graphql.ResolveParams) (interface{}, error) {
	id := p.Args["id"].(int)
	var user User
	r.db.First(&user, id)
	r.db.Delete(&user)
	return user, nil
}