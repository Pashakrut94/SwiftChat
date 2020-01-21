package utility

import (
	"database/sql"
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

var (
	testPGUser       = flag.String("pg_user", "Pasha", "PostgreSQL name")
	testPGPwd        = flag.String("pg_pwd", "pwd0123456789", "PostgreSQL password")
	testPGHost       = flag.String("pg_host", "localhost", "PostgreSQL host")
	testPGPort       = flag.String("pg_port", "54320", "PostgreSQL port")
	testPGDBname     = flag.String("pg_dbname", "test_mydb", "PostgreSQL name of DB")
	migrationsSource = flag.String("migration_source", "file:///home/anduser/Documents/SwiftChat/migrations", "The path to migrations directory") //"file:///home/anduser/Documents/SwiftChat/migrations"
)

func SetupSchema(t *testing.T) (*sql.DB, func()) {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	schemaName := "test" + strconv.FormatInt(rand.Int63(), 10)
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password='%s' sslmode=disable search_path=%s", *testPGHost, *testPGPort, *testPGDBname, *testPGUser, *testPGPwd, schemaName)
	db, err := sql.Open("postgres", connectionString)
	require.NoError(t, err)

	_, err = db.Exec("CREATE SCHEMA " + schemaName)
	require.NoError(t, err)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	require.NoError(t, err)

	m, err := migrate.NewWithDatabaseInstance(*migrationsSource, "postgres", driver)
	require.NoError(t, err)

	err = m.Up()
	require.NoError(t, err)

	return db, func() {
		_, err := db.Exec("DROP SCHEMA " + schemaName + " CASCADE")
		require.NoError(t, err)

		db.Close()
	}
}
