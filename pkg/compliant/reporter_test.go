package compliant

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/step"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewReporter(t *testing.T) {
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

	doneSignal := make(chan bool, 1)
	serverAddr := "localhost:8001"
	reporter := NewReporter(true, doneSignal, serverAddr)

	err := reporter.Report(result)
	require.NoError(t, err)

	// wait for http server to start
	time.Sleep(time.Millisecond * 100)

	// download report
	out := filepath.Join("testdata", "temp.json")
	r, err := http.Get("http://" + serverAddr + "?download=report")
	require.NoError(t, err)
	file, err := os.Create(out)
	require.NoError(t, err)
	defer file.Close()
	_, err = io.Copy(file, r.Body)
	require.NoError(t, err)

	<-doneSignal

	gp := filepath.Join("testdata", t.Name()+".golden.json")

	if *update {
		err = Copy(out, gp)
		require.NoError(t, err)
	}

	g, err := ioutil.ReadFile(gp)
	require.NoError(t, err)

	report, err := ioutil.ReadFile(out)
	require.NoError(t, err)

	err = os.Remove(out)
	require.NoError(t, err)

	assert.Equal(t, g, report)
}

func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
