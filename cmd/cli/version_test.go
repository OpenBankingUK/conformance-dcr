package main

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVersion_Print(t *testing.T) {
	require := require.New(t)

	versionInfo := VersionInfo{
		version:    "0.1.2",
		commitHash: "commit-hash",
		buildTime:  "build-time",
	}

	expected := ""
	expected += fmt.Sprintf("  Version:    %+v\n", versionInfo.version)
	expected += fmt.Sprintf("  Git commit: %+v\n", versionInfo.commitHash)
	expected += fmt.Sprintf("  Built:      %+v\n", versionInfo.buildTime)

	buff := bytes.NewBufferString("")
	actual := bufio.NewWriter(buff)

	require.NoError(versionInfo.Print(actual))
	require.Equal(expected, buff.String())
}
