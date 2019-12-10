package compliant

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/step"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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

func (r reporter) Report(result ManifestResult) error {
	file, err := json.MarshalIndent(r.mapToReport(result), "", " ")
	if err != nil {
		return err
	}

	r.startServer(file)

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

		var debug []string
		if r.debug {
			debug = r.mapDebugToReport(result.Debug)
		}

		stepResults[key] = ReportStep{
			Name:   result.Name,
			Pass:   result.Pass,
			Reason: result.FailReason,
			Debug:  debug,
		}
	}
	return stepResults
}

func (r reporter) mapDebugToReport(messages step.DebugMessages) []string {
	result := make([]string, len(messages.Item))
	for key, message := range messages.Item {
		result[key] = fmt.Sprintf(
			"%s %s",
			message.Time.Format(time.RFC3339),
			message.Message,
		)
	}
	return result
}

type Report struct {
	Name      string           `json:"name"`
	Version   string           `json:"version"`
	Pass      bool             `json:"pass"`
	Scenarios []ReportScenario `json:"scenarios"`
}

type ReportScenario struct {
	Id        string           `json:"id"`
	Name      string           `json:"name"`
	Spec      string           `json:"spec"`
	Pass      bool             `json:"pass"`
	TestCases []ReportTestcase `json:"test_cases"`
}

type ReportTestcase struct {
	Name  string       `json:"name"`
	Pass  bool         `json:"pass"`
	Steps []ReportStep `json:"steps"`
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
		h.doneSignalChan <- true
		return
	}

	_, err := rw.Write([]byte(`<html><body><a href="?download=report">Click here to download report.</a></body></html>`))
	if err != nil {
		fmt.Printf("Server error: %s", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
	}
}
