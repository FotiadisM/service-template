package suite

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"runtime/debug"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var matchMethod = flag.String("testify.m", "", "regular expression to select tests of the testify suite to run")

// T is a basic testing suite with methods for accessing
// appropriate *testing.T objects in tests.
type T struct {
	*assert.Assertions
	require  *require.Assertions
	testingT testingT
}

// setT sets the current *testing.T context.
func (t *T) setT(testingT *testing.T) { //nolint:thelper
	if t.testingT != nil {
		panic("T.testingT already set, can't overwrite")
	}
	t.testingT = testingT
	t.Assertions = assert.New(testingT)
	t.require = require.New(testingT)
}

func (t *T) Cleanup(f func()) {
	t.testingT.Cleanup(f)
}

func (t *T) Deadline() (deadline time.Time, ok bool) {
	return t.testingT.Deadline()
}

func (t *T) Error(args ...interface{}) {
	t.testingT.Error(args...)
}

func (t *T) Errorf(format string, args ...interface{}) {
	t.testingT.Errorf(format, args...)
}

func (t *T) Fail() {
	t.testingT.Fail()
}

func (t *T) FailNow() {
	t.testingT.FailNow()
}

func (t *T) Failed() bool {
	return t.testingT.Failed()
}

func (t *T) Fatal(args ...interface{}) {
	t.testingT.Fatal(args...)
}

func (t *T) Fatalf(format string, args ...interface{}) {
	t.testingT.Fatalf(format, args...)
}

func (t *T) Helper() {
	t.testingT.Helper()
}

func (t *T) Log(args ...interface{}) {
	t.testingT.Log(args...)
}

func (t *T) Logf(format string, args ...interface{}) {
	t.testingT.Logf(format, args...)
}

func (t *T) Name() string {
	return t.testingT.Name()
}

// Run provides suite functionality around golang subtests.  It should be
// called in place of t.Run(name, func(t *testing.T)) in test suite code.
// The passed-in func will be executed as a subtest with a fresh instance of t.
// Provides compatibility with go test pkg -run TestSuite/TestName/SubTestName.
func (t *T) Run(name string, subtest func(t *T)) bool {
	return t.testingT.Run(name, func(testingT *testing.T) { //nolint:thelper
		t := &T{}
		t.setT(testingT)
		subtest(t)
	})
}

func (t *T) Setenv(key, value string) {
	t.testingT.Setenv(key, value)
}

func (t *T) Skip(args ...interface{}) {
	t.testingT.Skip(args...)
}

func (t *T) SkipNow() {
	t.testingT.SkipNow()
}

func (t *T) Skipf(format string, args ...interface{}) {
	t.testingT.Skipf(format, args...)
}

func (t *T) Skipped() bool {
	return t.testingT.Skipped()
}

func (t *T) TempDir() string {
	return t.testingT.TempDir()
}

func (t *T) Parallel() {
	t.testingT.Parallel()
}

// Require returns a require context for suite.
func (t *T) Require() *require.Assertions {
	if t.testingT == nil {
		panic("T.testingT not set, can't get Require object")
	}
	return t.require
}

func failOnPanic(t *T) {
	r := recover()
	if r != nil {
		t.Errorf("test panicked: %v\n%s", r, debug.Stack())
		t.FailNow()
	}
}

// Filtering method according to set regular expression
// specified command-line argument -m.
func methodFilter(name string) (bool, error) {
	if ok, _ := regexp.MatchString("^Test", name); !ok {
		return false, nil
	}
	return regexp.MatchString(*matchMethod, name)
}

// Run takes a testing suite and runs all of the tests attached
// to it.
func Run(testingT *testing.T, suite interface{}) { //nolint:thelper
	t := &T{}
	t.setT(testingT)

	defer failOnPanic(t)

	var isSetupFinished bool

	tests := []testing.InternalTest{}
	methodFinder := reflect.TypeOf(suite)

	for i := 0; i < methodFinder.NumMethod(); i++ {
		method := methodFinder.Method(i)

		ok, err := methodFilter(method.Name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "testify: invalid regexp for -m: %s\n", err)
			os.Exit(1)
		}

		if !ok {
			continue
		}

		if !isSetupFinished {
			if setupAllSuite, ok := suite.(SetupAllSuite); ok {
				setupAllSuite.SetupSuite(t)
			}
			isSetupFinished = true
		}

		test := testing.InternalTest{
			Name: method.Name,
			F: func(testingT *testing.T) { //nolint:thelper
				t := &T{}
				t.setT(testingT)

				defer failOnPanic(t)

				if tearDownTestSuite, ok := suite.(TearDownTestSuite); ok {
					testingT.Cleanup(func() {
						tearDownTestSuite.TearDownTest(t)
					})
				}

				if setupTestSuite, ok := suite.(SetupTestSuite); ok {
					setupTestSuite.SetupTest(t)
				}

				method.Func.Call([]reflect.Value{reflect.ValueOf(suite), reflect.ValueOf(t)})
			},
		}
		tests = append(tests, test)
	}

	if len(tests) == 0 {
		testingT.Log("warning: no tests to run")
		return
	}

	if tearDownAllSuite, ok := suite.(TearDownAllSuite); ok {
		testingT.Cleanup(func() {
			tearDownAllSuite.TearDownSuite(t)
		})
	}

	for _, test := range tests {
		testingT.Run(test.Name, test.F)
	}
}
