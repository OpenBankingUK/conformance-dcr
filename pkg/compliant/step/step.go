package step

import (
	"fmt"
	"time"
)

type Step interface {
	Run(ctx Context) Result
}

type Result struct {
	Name       string
	Pass       bool
	FailReason string
	Debug      DebugMessages
}

type Results []Result

func (r Results) Fail() bool {
	for _, result := range r {
		if !result.Pass {
			return true
		}
	}
	return false
}

type DebugMessage struct {
	Time    time.Time
	Message string
}

type DebugMessages struct {
	Item []DebugMessage
}

func NewDebug() *DebugMessages {
	return &DebugMessages{}
}

func (d *DebugMessages) Log(msg string) {
	d.Item = append(d.Item, DebugMessage{
		Time:    time.Now(),
		Message: msg,
	})
}

func (d *DebugMessages) Logf(format string, a ...interface{}) {
	d.Item = append(d.Item, DebugMessage{
		Time:    time.Now(),
		Message: fmt.Sprintf(format, a...),
	})
}

func NewPassResult(name string) Result {
	return Result{Name: name, Pass: true}
}

func NewPassResultWithDebug(name string, log *DebugMessages) Result {
	return Result{Name: name, Pass: true, Debug: *log}
}

func NewFailResult(name, reason string) Result {
	return Result{Name: name, Pass: false, FailReason: reason}
}

func NewFailResultWithDebug(name, reason string, log *DebugMessages) Result {
	return Result{Name: name, Pass: false, FailReason: reason, Debug: *log}
}
