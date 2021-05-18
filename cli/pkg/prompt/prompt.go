package prompt

import "github.com/AlecAivazis/survey/v2"

func Confirm(question string) (bool, error) {
	var answer bool
	prompt := &survey.Confirm{
		Message: question,
	}
	err := survey.AskOne(prompt, &answer)
	if err != nil {
		return false, err
	}
	return answer, nil
}

func Select(question string, options []string) (string, error) {
	var answer string
	promptSelect := &survey.Select{
		Message: question,
		Options: options,
	}

	err := survey.AskOne(promptSelect, &answer)
	if err != nil {
		return answer, err
	}
	return answer, nil
}

func Input(question string) (string, error) {
	var answer string
	promptInput := &survey.Input{
		Message: question,
	}
	err := survey.AskOne(promptInput, &answer)
	if err != nil {
		return "", err
	}

	return answer, nil
}

func Password(question string) (string, error) {
	var answer string
	promptInput := &survey.Password{
		Message: question,
	}
	err := survey.AskOne(promptInput, &answer)
	if err != nil {
		return "", err
	}

	return answer, nil
}
