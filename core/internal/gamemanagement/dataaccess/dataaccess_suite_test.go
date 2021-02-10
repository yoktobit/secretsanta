package dataaccess_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gda "github.com/yoktobit/secretsanta/internal/general/dataaccess"
)

var config = gda.Config{
	User:     "test",
	Password: "test",
	DB:       "test",
	Host:     "localhost",
	Port:     "5432",
}

func TestDataaccess(t *testing.T) {
	RegisterFailHandler(Fail)
	gda.InitDatabaseContainer(&config)
	RunSpecs(t, "Dataaccess Suite")
}
