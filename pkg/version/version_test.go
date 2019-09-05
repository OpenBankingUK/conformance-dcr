package version

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVersion_Print(t *testing.T) {
	require := require.New(t)

	expected := ""
	expected += fmt.Sprintf("  Version:    %+v\n", version)
	expected += fmt.Sprintf("  Git commit: %+v\n", commitHash)
	expected += fmt.Sprintf("  Built:      %+v\n", buildTime)

	buff := bytes.NewBufferString("")
	actual := bufio.NewWriter(buff)

	require.NoError(Print(actual))
	require.Equal(expected, buff.String())
}
