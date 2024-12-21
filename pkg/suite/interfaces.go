package suite

import (
	"testing"
	"time"
)

// *testing.T interface.
type testingT interface {
	testing.TB

	Deadline() (deadline time.Time, ok bool)
	Parallel()
	Run(name string, f func(t *testing.T)) bool
}

// SetupAllSuite has a SetupSuite method, which will run before the
// tests in the suite are run.
type SetupAllSuite interface {
	SetupSuite(t *T)
}

// SetupTestSuite has a SetupTest method, which will run before each
// test in the suite.
type SetupTestSuite interface {
	SetupTest(t *T)
}

// TearDownAllSuite has a TearDownSuite method, which will run after
// all the tests in the suite have been run.
type TearDownAllSuite interface {
	TearDownSuite(t *T)
}

// TearDownTestSuite has a TearDownTest method, which will run after
// each test in the suite.
type TearDownTestSuite interface {
	TearDownTest(t *T)
}
