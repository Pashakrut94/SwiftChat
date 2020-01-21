package users

import (
	"database/sql"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Pashakrut94/SwiftChat/utility"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// func RouterCreate(db *sql.DB) *mux.Router {
// 	r := mux.NewRouter()
// 	r.HandleFunc("/api/users", CreateUser(*NewUserRepo(db)))

// 	return r
// }

//dont need test for unnecessary handler(delete CreateUser handler)
// func TestCreateUser(t *testing.T) {
// 	db, dropschema := utility.SetupSchema(t)
// 	defer dropschema()

// 	testUser := User{Name: "Pasha", Password: "q1w2e3r4", Phone: "375291112233"}

// 	data, err := json.Marshal(testUser)
// 	require.NoError(t, err)
// 	body := bytes.NewReader(data)

// 	req := httptest.NewRequest("POST", "/api/signin", body)
// 	req.Header.Add("Content-Type", "application/json")

// 	rec := httptest.NewRecorder()
// 	RouterCreate(db).ServeHTTP(rec, req)

// 	var user User
// 	err = json.Unmarshal(rec.Body.Bytes(), &user)
// 	require.NoError(t, err)

// 	assert.Equal(t, http.StatusOK, rec.Code)
// 	AssertUsers(t, testUser, user)

// }

// repo := *NewUserRepo(db)

// testUsers := []User{
// 	{Name: "Pasha", Password: "q1w2e3r4", Phone: "375291112233"},
// 	{Name: "Roman", Password: "z1x2c3v4", Phone: "375331112233"},
// 	{Name: "Denis", Password: "a1s2d3f4", Phone: "375441112233"},
// }

// for i := 0; i < len(testUsers); i++ {

// 	err = repo.Create(&testUsers[i])
// 	assert.NoError(t, err)
// }

func Router(db *sql.DB) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/users/{UserID:[0-9]+}", GetUser(*NewUserRepo(db)))

	return r
}

func GetUserCase(t *testing.T) {
	db, dropSchema := utility.SetupSchema(t)
	defer dropSchema()

	testCases := []struct {
		Name string
		ID   int
		Want []byte
	}{
		{Name: "first", ID: 1, Want: []byte(`{"data":{"id":1,"name":"Pasha","password":"qwerty123456","phone":"123456789012"}}`)},
		{Name: "second", ID: 2, Want: []byte(`{"data":{"id":2,"name":"Masha","password":"zxcvbn123","phone":"09876543212"}}`)},
		{Name: "third", ID: 3, Want: []byte(`{"data":{"id":3,"name":"Nikolay","password":"password123","phone":"375291234599"}}`)},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/api/users/%d", tc.ID), nil)

			// assert.NoError(t, err)
			// handler.ServeHTTP(rec, req)

			Router(db).ServeHTTP(rec, req)
			assert.Equal(t, tc.Want, rec.Body.Bytes())
		})
	}
}
