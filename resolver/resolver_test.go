package resolver

import (
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func prepareTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}
	db.AutoMigrate(&User{})
	return db
}

func TestGetAllUsers(t *testing.T) {
	testDB := prepareTestDB()
	testResolver := NewResolver(testDB)

	// Create a test user
	testUser := User{
		Name:     "Test User",
		Email:    "test.user@example.com",
		Password: "testpassword",
	}
	testDB.Create(&testUser)

	// Execute the query
	params := graphql.ResolveParams{}
	users, err := testResolver.GetAllUsers(params)

	assert.Nil(t, err)
	assert.NotNil(t, users)

	userList, ok := users.([]User)
	assert.True(t, ok)
	assert.Len(t, userList, 1)
	assert.Equal(t, testUser.Name, userList[0].Name)
	assert.Equal(t, testUser.Email, userList[0].Email)
}

func TestGetUserByID(t *testing.T) {
	testDB := prepareTestDB()
	testResolver := NewResolver(testDB)

	// Create a test user
	testUser := User{
		Name:     "Test User",
		Email:    "test.user@example.com",
		Password: "testpassword",
	}
	testDB.Create(&testUser)

	// Execute the query
	params := graphql.ResolveParams{
		Args: map[string]interface{}{
			"id": int(testUser.ID),
		},
	}
	user, err := testResolver.GetUserByID(params)

	assert.Nil(t, err)
	assert.NotNil(t, user)

	result, ok := user.(User)
	assert.True(t, ok)
	assert.Equal(t, testUser.Name, result.Name)
	assert.Equal(t, testUser.Email, result.Email)
}

// 他のテスト関数をここに追加してください。
