package suite

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"runtime/debug"
	"testing"
)

var matchMethod = flag.String("testify.m", "", "regular expression to select tests of the testify suite to run")

func failOnPanic(t *testing.T) {
	t.Helper()

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
func Run(t *testing.T, suite interface{}) {
	t.Helper()

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
			F: func(t *testing.T) {
				t.Helper()

				defer failOnPanic(t)

				if tearDownTestSuite, ok := suite.(TearDownTestSuite); ok {
					t.Cleanup(func() {
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
		t.Log("warning: no tests to run")
		return
	}

	if tearDownAllSuite, ok := suite.(TearDownAllSuite); ok {
		t.Cleanup(func() {
			tearDownAllSuite.TearDownSuite(t)
		})
	}

	for _, test := range tests {
		t.Run(test.Name, test.F)
	}
}
