package cmdutil

import (
	"fmt"
	"io"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/cli/cli/v2/internal/ghinstance"
)

func NewPrompter(stdin io.Reader, stdout, stderr io.Writer) Prompter {
	return &surveyPrompter{
		stdin:  stdin.(terminal.FileReader),
		stdout: stdout.(terminal.FileWriter),
		stderr: stderr,
	}
}

type surveyPrompter struct {
	stdin  terminal.FileReader
	stdout terminal.FileWriter
	stderr io.Writer
}

func (p *surveyPrompter) Select(message, defaultValue string, options []string) (result int, err error) {
	q := &survey.Select{
		Message:  message,
		Default:  defaultValue,
		Options:  options,
		PageSize: 20,
	}

	err = p.ask(q, &result)

	return
}

func (p *surveyPrompter) MultiSelect(message, defaultValue string, options []string) (result int, err error) {
	q := &survey.MultiSelect{
		Message:  message,
		Default:  defaultValue,
		Options:  options,
		PageSize: 20,
	}

	err = p.ask(q, &result)

	return
}

func (p *surveyPrompter) ask(q survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
	opts = append(opts, survey.WithStdio(p.stdin, p.stdout, p.stderr))
	err := survey.AskOne(q, response, opts...)
	if err == nil {
		return nil
	}
	return fmt.Errorf("could not prompt: %w", err)
}

func (p *surveyPrompter) Input(prompt, defaultValue string) (result string, err error) {
	err = p.ask(&survey.Input{
		Message: prompt,
		Default: defaultValue,
	}, &result)

	return
}

func (p *surveyPrompter) InputHostname() (result string, err error) {
	err = p.ask(
		&survey.Input{
			Message: "GHE hostname:",
		}, &result, survey.WithValidator(func(v interface{}) error {
			return ghinstance.HostnameValidator(v.(string))
		}))

	return
}

func (p *surveyPrompter) Password(prompt string) (result string, err error) {
	err = p.ask(&survey.Password{
		Message: prompt,
	}, &result)

	return
}

func (p *surveyPrompter) Confirm(prompt string, defaultValue bool) (result bool, err error) {
	p.ask(&survey.Confirm{
		Message: prompt,
		Default: defaultValue,
	}, &result)

	return
}
