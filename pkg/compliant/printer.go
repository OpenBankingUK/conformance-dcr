package compliant

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/step"
	"fmt"
	"github.com/logrusorgru/aurora"
	"io"
	"os"
)

func NewPrinter(debug bool) printer {
	return printer{
		debug:  debug,
		output: os.Stdout,
	}
}

func NewPrinterWithOptions(debug bool, w io.Writer) printer {
	return printer{
		debug:  debug,
		output: w,
	}
}

type printer struct {
	debug  bool
	output io.Writer
}

func (p printer) Print(result ManifestResult) error {
	for _, scenarioResult := range result.Results {
		_, err := fmt.Fprintf(p.output, "=== Scenario: %s - %s\n", scenarioResult.Id, scenarioResult.Name)
		if err != nil {
			return err
		}
		for _, testCasesResult := range scenarioResult.TestCaseResults {
			_, err := fmt.Fprintf(p.output, "\tTest case: %s\n", testCasesResult.Name)
			if err != nil {
				return err
			}
			for _, stepResult := range testCasesResult.Results {
				err := p.printColourTestResult(stepResult)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (p printer) printColourTestResult(result step.Result) error {
	if result.Pass {
		_, err := fmt.Fprintf(p.output, "\t\t%s %s\n", aurora.Green("PASS"), result.Name)
		if err != nil {
			return err
		}
	} else {
		_, err := fmt.Fprintf(p.output,
			"\t\t%s %s: %s\n",
			aurora.Red("FAIL"),
			result.Name,
			result.FailReason,
		)
		if err != nil {
			return err
		}
	}
	if p.debug {
		return p.printColourDebugMessages(result.Debug)
	}
	return nil
}

func (p printer) printColourDebugMessages(log step.DebugMessages) error {
	for _, msg := range log.Item {
		_, err := fmt.Fprintf(p.output,
			"%s %s\n",
			msg.Time.Format("2006/01/02 15:04:05"),
			aurora.Gray(15, msg.Message),
		)
		if err != nil {
			return err
		}
	}
	return nil
}
