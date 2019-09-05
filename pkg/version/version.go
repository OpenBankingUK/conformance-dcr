package version

import (
	"bufio"
	"fmt"
)

// nolint:gochecknoglobals
var (
	version    string
	commitHash string
	buildTime  string
)

// Print - print the version, for example:
//
// Dynamic Client Registration Conformance Tool cli
//   Version:    0.0.1
//   Git commit: 689c4d6eecd88591d303274d8707d25ec53e1771
//   Built:      Thu Sep  5 08:19:53 UTC 2019
func Print(output *bufio.Writer) error {
	if _, err := fmt.Fprintf(output, "  Version:    %+v\n", version); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(output, "  Git commit: %+v\n", commitHash); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(output, "  Built:      %+v\n", buildTime); err != nil {
		return err
	}
	if err := output.Flush(); err != nil {
		return err
	}

	return nil
}
