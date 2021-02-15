package dataaccess // import github.com/yoktobit/secretsanta/internal/general/dataacces

import (
	"fmt"
	"net/url"
	"os"

	log "github.com/sirupsen/logrus"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config is the configuration for the connection
type Config struct {
	User     string
	Password string
	DB       string
	Host     string
	Port     string
}

// Connection encapsulates some DB connection
type Connection interface {
	Connection() *gorm.DB
	NewTransaction(f func(Connection) error)
}

type connection struct {
	db *gorm.DB
}

// NewConnectionWithEnvironment connects the database by using environment parameters
func NewConnectionWithEnvironment() Connection {

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	return connect(user, password, dbname, host, port)
}

// NewConnectionWithParameters connects the database by given parameters
func NewConnectionWithParameters(user string, password string, dbname string, host string, port string) Connection {

	return connect(user, password, dbname, host, port)
}

// NewConnectionWithConfig connects the database by given config
func NewConnectionWithConfig(config Config) Connection {

	return connect(config.User, config.Password, config.DB, config.Host, config.Port)
}

// Connection gets the Gorm-Connection from the Connection object
func (c *connection) Connection() *gorm.DB {

	return c.db
}

// NewTransaction defines a new Transaction
func (c connection) NewTransaction(f func(Connection) error) {
	c.Connection().Transaction(func(tx *gorm.DB) error {
		return f(&connection{db: tx})
	})
}

func connect(user string, password string, dbname string, host string, port string) Connection {

	dsn := connectionString(user, password, dbname, host, port)
	log.Debug("Connecting to " + dsn)
	if dsn == "" {
		dsn = "user=santa password=santa dbname=secretsanta port=5432"
	}
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	return &connection{db: database}
}

func connectionString(user string, password string, dbname string, host string, port string) string {
	dsn := url.URL{
		User:   url.UserPassword(user, password),
		Scheme: "postgres",
		Host:   fmt.Sprintf("%s:%s", host, port),
		Path:   dbname,
	}
	return dsn.String()
}
