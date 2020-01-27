package compliant

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/step"
	"github.com/sirupsen/logrus"
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

// Report marshals the result into json, then starts a server to host the generated report file.
func (r reporter) Report(result ManifestResult) error {
	file, err := json.MarshalIndent(r.mapToReport(result), "", " ")
	if err != nil {
		return err
	}

	reportFile := "report.json"
	files := []string{reportFile}
	err = ioutil.WriteFile(reportFile, file, os.ModePerm)
	if err != nil {
		return err
	}

	if r.debug {
		var debugFile *os.File
		debugFile, err = r.GetDebugFile(result)
		if err != nil {
			return err
		}

		files = append(files, debugFile.Name())
	}

	zip, err := ZipFiles(files)
	if err != nil {
		return err
	}

	r.startServer(zip)

	return nil
}

func (r reporter) startServer(report *os.File) {
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

func (r reporter) GetDebugFile(result ManifestResult) (*os.File, error) {
	file, err := os.Create("debug.json")
	if err != nil {
		return nil, err
	}
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(file)
	log.SetLevel(logrus.DebugLevel)

	for _, scenario := range result.Results {
		for _, testcase := range scenario.TestCaseResults {
			for _, result := range testcase.Results {
				for _, message := range result.Debug.Item {
					log.WithTime(message.Time).
						WithFields(logrus.Fields{
							"scenario_id":        scenario.Id,
							"scenario_name":      scenario.Name,
							"scenario_spec":      scenario.Spec,
							"scenario_pass":      !scenario.Fail(),
							"testcase_name":      testcase.Name,
							"testcase_pass":      !testcase.Fail(),
							"result_name":        result.Name,
							"result_pass":        result.Pass,
							"result_fail_reason": result.FailReason,
						}).
						Debug(message.Message)
				}
			}
		}
	}
	return file, nil
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
	report         *os.File
}

func (h downloadHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("download") != "" {
		b, err := ioutil.ReadFile(h.report.Name())
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Printf("Server error: %s", err.Error())
		}
		_, err = rw.Write(b)
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
