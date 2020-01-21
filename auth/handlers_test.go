package auth

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Pashakrut94/SwiftChat/users"
	"github.com/Pashakrut94/SwiftChat/utility"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Router(db *sql.DB) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/signup", SignUp(users.NewUserRepo(db), *NewSessionRepo(db)))
	r.HandleFunc("/api/signin", SignIn(users.NewUserRepo(db), *NewSessionRepo(db)))
	r.HandleFunc("/api/logout", Logout(*users.NewUserRepo(db), *NewSessionRepo(db)))
	return r
}

func TestSignUp(t *testing.T) {
	db, dropschema := utility.SetupSchema(t)
	defer dropschema()

	signUpRequest := SignUpRequest{Username: "Pasha", Password: "q1w2e3r4", Phone: "375291112233"}
	data, err := json.Marshal(signUpRequest)
	require.NoError(t, err)

	body := bytes.NewReader(data)
	req := httptest.NewRequest("POST", "/api/signup", body)
	req.Header.Add("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	Router(db).ServeHTTP(rec, req)

	var resp struct {
		Data users.User `json:"data"`
	}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)

	// User data exists in signup responses
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, 1, resp.Data.ID)
	assert.NotEqual(t, signUpRequest.Password, resp.Data.Password)

	// Signed up user exists in db
	repo := users.NewUserRepo(db)
	userByID, err := repo.Get(resp.Data.ID)
	require.NoError(t, err)
	assert.Equal(t, userByID.Name, resp.Data.Name)
	assert.Equal(t, userByID.Phone, resp.Data.Phone)
	assert.Equal(t, userByID.Password, resp.Data.Password)

	// Valid session created
	sid := rec.Result().Cookies()[0].Value
	sessRepo := NewSessionRepo(db)
	sess, err := sessRepo.Get(sid)
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, sess.UserID)
	assert.Nil(t, sess.DeletedAt)

	rec = httptest.NewRecorder()
	body.Seek(0, 0)
	Router(db).ServeHTTP(rec, req)

	var errResp struct {
		Error string `json:"error"`
	}
	err = json.Unmarshal(rec.Body.Bytes(), &errResp)
	require.NoError(t, err)

	// User already exists
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "user already exists", errResp.Error)
}

func TestSignIn(t *testing.T) {
	db, dropschema := utility.SetupSchema(t)
	defer dropschema()

	signInRequest := users.User{Phone: "375441235588", Password: "q1w2e3r4", Name: "Pasha"}
	repo := users.NewUserRepo(db)

	user, err := HandleSignUp(repo, signInRequest.Name, signInRequest.Password, signInRequest.Phone)
	require.NoError(t, err)

	data, err := json.Marshal(signInRequest)
	require.NoError(t, err)
	body := bytes.NewReader(data)

	req := httptest.NewRequest("POST", "/api/signin", body)
	req.Header.Add("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	Router(db).ServeHTTP(rec, req)

	var resp struct {
		Data users.User `json:"data"`
	}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)

	// User data exists in signin response
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, user.ID, resp.Data.ID)
	assert.Equal(t, user.Password, resp.Data.Password)
	assert.Equal(t, user.Phone, resp.Data.Phone)

	// Signen in user exists in db
	userByID, err := repo.Get(resp.Data.ID)
	require.NoError(t, err)
	assert.Equal(t, userByID.Name, resp.Data.Name)
	assert.Equal(t, userByID.Phone, resp.Data.Phone)
	assert.Equal(t, userByID.Password, resp.Data.Password)

	// Valid session created
	sid := rec.Result().Cookies()[0].Value
	sessRepo := NewSessionRepo(db)
	sess, err := sessRepo.Get(sid)
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, sess.UserID)
	assert.Nil(t, sess.DeletedAt)
}

func TestLogout(t *testing.T) {
	// Going to signup handler for seed user and set cookies
	db, dropschema := utility.SetupSchema(t)
	defer dropschema()

	signUpRequest := SignUpRequest{Username: "Pasha", Password: "q1w2e3r4", Phone: "375291112233"}
	data, err := json.Marshal(signUpRequest)
	require.NoError(t, err)

	body := bytes.NewReader(data)
	req := httptest.NewRequest("POST", "/api/signup", body)
	req.Header.Add("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	Router(db).ServeHTTP(rec, req)

	sid := rec.Result().Cookies()[0].Value

	sessRepo := NewSessionRepo(db)
	sess, err := sessRepo.Get(sid)
	assert.NoError(t, err)
	assert.Nil(t, sess.DeletedAt)

	rec = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/api/logout", nil)
	req.AddCookie(&http.Cookie{Name: sessionName, Value: sid})
	req.Header.Add("Content-Type", "application/json")
	Router(db).ServeHTTP(rec, req)

	sess, err = sessRepo.Get(sid)
	assert.Error(t, err)

	cookie := rec.Result().Cookies()[0]
	assert.Empty(t, cookie.Value)
	assert.Equal(t, time.Unix(0, 0).UTC(), cookie.Expires)
}
