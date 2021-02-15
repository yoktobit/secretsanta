package logic_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/yoktobit/secretsanta/internal/gamemanagement/logic/to"
)

func NewCreateGameTo() to.CreateGameTo {
	return to.CreateGameTo{Title: "ABC", Description: "DEF", AdminUser: "Martin", AdminPassword: "Test12345"}
}

func TestLogic(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Logic Suite")
}
