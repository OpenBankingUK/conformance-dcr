package compliant

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/step"
)

func NewReporter(debug bool, doneSignal chan<- bool, serverAddr string) reporter {
	return reporter{
		debug:          debug,
		doneSignalChan: doneSignal,
		serverAddr:     serverAddr,
	}
}

type reporter struct {
	debug          bool
	doneSignalChan chan<- bool
	serverAddr     string
}

// Report marshals the result and debug into json, zips them, then starts a server to host the generated zip file.
func (r reporter) Report(result ManifestResult) error {
	reportJson, err := json.MarshalIndent(r.mapToReport(result), "", " ")
	if err != nil {
		return err
	}
	var files = []ReportFile{
		{"report.json", string(reportJson)},
	}

	if r.debug {
		var debugJson []byte
		debugJson, err = json.MarshalIndent(r.GetDebugLog(result), "", " ")
		if err != nil {
			return err
		}

		files = append(files, ReportFile{"debug.json", string(debugJson)})
	}

	b, err := ZipReportFiles(files)
	if err != nil {
		return err
	}

	r.startServer(b)

	return nil
}

func (r reporter) startServer(report []byte) {
	go func() {
		handler := downloadHandler{
			doneSignalChan: r.doneSignalChan,
			report:         report,
		}

		server := &http.Server{Addr: r.serverAddr, Handler: handler}
		err := server.ListenAndServe()
		fmt.Printf("Error starting embedded webserver: %s", err)
	}()
}

type DebugLine struct {
	Time     string         `json:"time,omitempty"`
	Message  string         `json:"message,omitempty"`
	Scenario ReportScenario `json:"scenario,omitempty"`
	Testcase ReportTestcase `json:"testcase,omitempty"`
	Result   ReportStep     `json:"result,omitempty"`
}

func (r reporter) GetDebugLog(result ManifestResult) []DebugLine {

	var log []DebugLine

	for _, scenario := range result.Results {
		for _, testcase := range scenario.TestCaseResults {
			for _, result := range testcase.Results {
				for _, message := range result.Debug.Item {
					log = append(log, DebugLine{
						Message: message.Message,
						Time:    message.Time.Format(time.RFC3339),
						Scenario: ReportScenario{
							Id:   scenario.Id,
							Name: scenario.Name,
							Spec: scenario.Spec,
							Pass: !scenario.Fail(),
						},
						Testcase: ReportTestcase{
							Name: testcase.Name,
							Pass: !testcase.Fail(),
						},
						Result: ReportStep{
							Name:   result.Name,
							Pass:   result.Pass,
							Reason: result.FailReason,
						},
					})
				}
			}
		}
	}
	return log
}

func (r reporter) mapToReport(result ManifestResult) Report {
	results := make([]ReportScenario, len(result.Results))
	for key, scenario := range result.Results {
		results[key] = ReportScenario{
			Id:        scenario.Id,
			Name:      scenario.Name,
			Spec:      scenario.Spec,
			Pass:      !scenario.Fail(),
			TestCases: r.mapTCSToReport(scenario.TestCaseResults),
		}
	}
	return Report{
		Name:      result.Name,
		Version:   result.Version,
		Pass:      !result.Fail(),
		Scenarios: results,
	}
}

func (r reporter) mapTCSToReport(results TestCaseResults) []ReportTestcase {
	reportResults := make([]ReportTestcase, len(results))
	for key, result := range results {
		reportResults[key] = ReportTestcase{
			Name:  result.Name,
			Pass:  !result.Fail(),
			Steps: r.mapStepsToReport(result.Results),
		}
	}
	return reportResults
}

func (r reporter) mapStepsToReport(results step.Results) []ReportStep {
	stepResults := make([]ReportStep, len(results))
	for key, result := range results {
		stepResults[key] = ReportStep{
			Name:   result.Name,
			Pass:   result.Pass,
			Reason: result.FailReason,
		}
	}
	return stepResults
}

type Report struct {
	Name      string           `json:"name"`
	Version   string           `json:"version"`
	Pass      bool             `json:"pass"`
	Scenarios []ReportScenario `json:"scenarios,omitempty"`
}

type ReportScenario struct {
	Id        string           `json:"id"`
	Name      string           `json:"name"`
	Spec      string           `json:"spec"`
	Pass      bool             `json:"pass"`
	TestCases []ReportTestcase `json:"test_cases,omitempty"`
}

type ReportTestcase struct {
	Name  string       `json:"name"`
	Pass  bool         `json:"pass"`
	Steps []ReportStep `json:"steps,omitempty"`
}

type ReportStep struct {
	Name   string   `json:"name"`
	Pass   bool     `json:"pass"`
	Reason string   `json:"reason,omitempty"`
	Debug  []string `json:"debug,omitempty"`
}

type downloadHandler struct {
	doneSignalChan chan<- bool
	report         []byte
}

func (h downloadHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("download") != "" {
		_, err := rw.Write(h.report)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Printf("Server error: %s", err.Error())
		}
		rw.Header().Add("Content-Type", "application/zip")
		h.doneSignalChan <- true
		return
	}

	_, err := rw.Write([]byte(`<html><body><a href="?download=report">Click here to download report.</a></body></html>`))
	if err != nil {
		fmt.Printf("Server error: %s", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

type ReportFile struct {
	Name string
	Body string
}

func ZipReportFiles(files []ReportFile) ([]byte, error) {
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	for _, file := range files {
		f, err := w.Create(file.Name)
		if err != nil {
			return nil, err
		}
		_, err = f.Write([]byte(file.Body))
		if err != nil {
			return nil, err
		}
	}
	w.Close()
	return buf.Bytes(), nil
}
