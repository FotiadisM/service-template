package authv1

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type UnitTestSuit struct {
	suite.Suite
}

func TestUnitTestSuit(t *testing.T) {
	suite.Run(t, &UnitTestSuit{})
}
