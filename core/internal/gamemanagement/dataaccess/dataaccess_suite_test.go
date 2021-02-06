package dataaccess_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/docker/go-connections/nat"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	testcontainers "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	gda "github.com/yoktobit/secretsanta/internal/general/dataaccess"
)

var config = gda.Config{
	User:     "test",
	Password: "test",
	DB:       "test",
	Host:     "localhost",
	Port:     "5432",
}

func initContainer(config *gda.Config) {
	ctx := context.Background()
	natPort := fmt.Sprintf("%s/tcp", config.Port)
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{natPort},
		Env: map[string]string{
			"POSTGRES_USER":     config.User,
			"POSTGRES_PASSWORD": config.Password,
		},
		WaitingFor: wait.ForListeningPort(nat.Port(natPort)),
	}
	pg, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}
	mp, err := pg.MappedPort(ctx, nat.Port(natPort))
	if err != nil {
		panic(err)
	}
	ma, err := pg.Host(ctx)
	if err != nil {
		panic(err)
	}
	config.Host = ma
	config.Port = mp.Port()
}

func TestDataaccess(t *testing.T) {
	RegisterFailHandler(Fail)
	initContainer(&config)
	RunSpecs(t, "Dataaccess Suite")
}
