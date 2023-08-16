//go:build integration

package store

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, &IntegrationTestSuite{})
}
