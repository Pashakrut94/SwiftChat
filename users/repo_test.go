package users

import (
	"testing"

	"github.com/Pashakrut94/SwiftChat/utility"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func AssertUsers(t *testing.T, u1, u2 User) {
	assert.Equal(t, u1.Name, u2.Name)
	assert.Equal(t, u1.Password, u2.Password)
	assert.Equal(t, u1.Phone, u2.Phone)
}

func TestRepo(t *testing.T) {
	db, dropSchema := utility.SetupSchema(t)
	defer dropSchema()

	repo := NewUserRepo(db)

	testUsers := []User{
		{Name: "Pasha", Password: "q1w2e3r4", Phone: "375291112233"},
		{Name: "Roman", Password: "z1x2c3v4", Phone: "375331112233"},
		{Name: "Denis", Password: "a1s2d3f4", Phone: "375441112233"},
	}

	for i := 0; i < len(testUsers); i++ {

		err := repo.Create(&testUsers[i])
		assert.NoError(t, err)
	}
	userByID, err := repo.Get(1)
	assert.NoError(t, err)
	AssertUsers(t, testUsers[0], *userByID)

	userByPhone, err := repo.GetByPhone(testUsers[0].Phone)
	assert.NoError(t, err)
	AssertUsers(t, testUsers[0], *userByPhone)

	users, err := repo.List()
	assert.NoError(t, err)
	for i, user := range users {
		AssertUsers(t, testUsers[i], user)
	}
}
