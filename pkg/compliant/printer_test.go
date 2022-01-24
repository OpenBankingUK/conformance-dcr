package compliant

import (
	"bytes"
	"flag"
	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/step"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"path/filepath"
	"testing"
)

// nolint:gochecknoglobals
var update = flag.Bool("update", false, "update .golden files")

func TestNewPrinter(t *testing.T) {
	result := ManifestResult{
		Results: []ScenarioResult{
			{
				Id:   "1",
				Name: "scenario one",
				Spec: "spec link",
				TestCaseResults: TestCaseResults{
					{
						Name: "tc one",
						Results: []step.Result{
							{
								Name:       "step one",
								Pass:       false,
								FailReason: "reasons",
								Debug: step.DebugMessages{
									Item: []step.DebugMessage{
										{
											Message: "debug",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		Name:    "manifest test test result",
		Version: "0.0",
	}
	w := &bytes.Buffer{}
	printer := NewPrinterWithOptions(true, w)

	err := printer.Print(result)
	require.NoError(t, err)

	gp := filepath.Join("testdata", t.Name()+".golden.json")

	if *update {
		t.Log("update golden file")
		err = ioutil.WriteFile(gp, w.Bytes(), 0644)
		require.NoError(t, err)
	}

	g, err := ioutil.ReadFile(gp)
	require.NoError(t, err)

	assert.Equal(t, g, w.Bytes())
}
