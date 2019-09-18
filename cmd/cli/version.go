package main

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

type VersionInfo struct {
	version    string
	commitHash string
	buildTime  string
}

// Print - print the version, for example:
//
//   Version:    0.0.1
//   Git commit: 689c4d6eecd88591d303274d8707d25ec53e1771
//   Built:      Thu Sep  5 08:19:53 UTC 2019
func (v *VersionInfo) Print(output *bufio.Writer) error {
	message := "  Version:    %+v\n" +
		"  Git commit: %+v\n" +
		"  Built:      %+v\n"
	if _, err := fmt.Fprintf(output, message, v.version, v.commitHash, v.buildTime); err != nil {
		return err
	}
	if err := output.Flush(); err != nil {
		return err
	}

	return nil
}
