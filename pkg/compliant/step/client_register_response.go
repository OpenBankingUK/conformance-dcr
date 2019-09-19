package step

import (
	"fmt"
	"io/ioutil"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/auth"
)

type clientRegisterResponse struct {
	stepName          string
	responseCtxKey    string
	clientCtxKey      string
	debug             *DebugMessages
	authoriserBuilder auth.AuthoriserBuilder
}

func NewClientRegisterResponse(responseCtxKey, clientCtxKey string, authoriserBuilder auth.AuthoriserBuilder) Step {
	return clientRegisterResponse{
		stepName:          "Decode client register response",
		responseCtxKey:    responseCtxKey,
		clientCtxKey:      clientCtxKey,
		debug:             NewDebug(),
		authoriserBuilder: authoriserBuilder,
	}
}

func (s clientRegisterResponse) Run(ctx Context) Result {
	s.debug.Logf("get response object from ctx var: %s", s.responseCtxKey)
	response, err := ctx.GetResponse(s.responseCtxKey)
	if err != nil {
		return s.failResult(fmt.Sprintf("getting response object from context: %s", err.Error()))
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return s.failResult(fmt.Sprintf("client register: %s", err.Error()))
	}

	s.debug.Log("getting client")
	s.debug.Logf("register res: %+v", string(body))
	authoriser, err := s.authoriserBuilder.Build()
	if err != nil {
		return NewFailResultWithDebug(
			s.stepName,
			err.Error(),
			s.debug,
		)
	}
	client, err := authoriser.Client(body)
	if err != nil {
		return s.failResult(fmt.Sprintf("client register: %s", err.Error()))
	}

	s.debug.Logf("setting software client in context var: %s", s.clientCtxKey)
	ctx.SetClient(s.clientCtxKey, client)

	return NewPassResultWithDebug(s.stepName, s.debug)
}

func (s clientRegisterResponse) failResult(msg string) Result {
	return NewFailResultWithDebug(
		s.stepName,
		msg,
		s.debug,
	)
}
