package version

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

// AddFlag adds -version flags to the FlagSet.
// If triggered, the flags print version information and call os.Exit(0).
// If FlagSet is nil, it adds the flags to flag.CommandLine.
func AddFlag(f *flag.FlagSet) {
	if f == nil {
		f = flag.CommandLine
	}
	f.Var(boolFunc(version), "version", "print version information and exit")
}

// boolFunc implements flag.Value.
type boolFunc func(bool) error

func (f boolFunc) IsBoolFlag() bool {
	return true
}

func (f boolFunc) String() string {
	return ""
}

func (f boolFunc) Set(s string) error {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	return f(b)
}

func version(b bool) error {
	if !b {
		return nil
	}

	fmt.Fprintln(os.Stdout, String())

	os.Exit(0)

	return nil
}
