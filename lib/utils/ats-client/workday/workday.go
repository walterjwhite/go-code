package workdayats

import (
	"context"
	"fmt"
	"strings"

	atsclient "github.com/walterjwhite/go-code/lib/utils/ats-client"
	"github.com/walterjwhite/go-code/lib/utils/ats-client/ai"
)

type Account = atsclient.Account

const (
	loginButtonSelector    = `[data-automation-id="loginButton"]`
	registerButtonSelector = `[data-automation-id="registerButton"]`
	emailInputSelector     = `[data-automation-id="emailInput"]`
	passwordInputSelector     = `[data-automation-id="passwordInput"]`
	firstNameInputSelector    = `[data-automation-id="firstNameInput"]`
	lastNameInputSelector     = `[data-automation-id="lastNameInput"]`
	phoneInputSelector        = `[data-automation-id="phoneInput"]`
	countrySelector           = `[data-automation-id="countrySelect"]`
	addressInputSelector      = `[data-automation-id="addressInput"]`
	cityInputSelector         = `[data-automation-id="cityInput"]`
	stateSelector             = `[data-automation-id="stateSelect"]`
	zipCodeInputSelector      = `[data-automation-id="zipCodeInput"]`
	resumeUploadSelector      = `[data-automation-id="resumeUpload"]`
	nextButtonSelector        = `[data-automation-id="nextButton"]`
	submitButtonSelector      = `[data-automation-id="submitButton"]`
	questionContainerSelector = `[data-automation-id="questionContainer"]`
	applicationFormSelector   = `[data-automation-id="applicationForm"]`
)

type WorkdayATS struct{}

func (w *WorkdayATS) GetName() string {
	return "workday"
}

func (w *WorkdayATS) RegisterAccount(executor *atsclient.Executor, account *Account) error {
	err := executor.WaitAndClick(registerButtonSelector)
	if err != nil {
		return fmt.Errorf("failed to click register button: %w", err)
	}

	err = executor.SetValue(emailInputSelector, account.Email)
	if err != nil {
		return fmt.Errorf("failed to set email: %w", err)
	}

	err = executor.SetValue(passwordInputSelector, account.Password)
	if err != nil {
		return fmt.Errorf("failed to set password: %w", err)
	}

	err = executor.SetValue(firstNameInputSelector, account.FirstName)
	if err != nil {
		return fmt.Errorf("failed to set first name: %w", err)
	}

	err = executor.SetValue(lastNameInputSelector, account.LastName)
	if err != nil {
		return fmt.Errorf("failed to set last name: %w", err)
	}

	err = executor.SetValue(phoneInputSelector, account.Phone)
	if err != nil {
		return fmt.Errorf("failed to set phone: %w", err)
	}

	err = executor.SetValue(addressInputSelector, account.Address)
	if err != nil {
		return fmt.Errorf("failed to set address: %w", err)
	}

	err = executor.SetValue(cityInputSelector, account.City)
	if err != nil {
		return fmt.Errorf("failed to set city: %w", err)
	}

	err = executor.SetValue(countrySelector, account.Country)
	if err != nil {
		return fmt.Errorf("failed to set country: %w", err)
	}

	err = executor.SetValue(stateSelector, account.State)
	if err != nil {
		return fmt.Errorf("failed to set state: %w", err)
	}

	err = executor.SetValue(zipCodeInputSelector, account.ZipCode)
	if err != nil {
		return fmt.Errorf("failed to set zip code: %w", err)
	}

	err = executor.WaitAndClick(nextButtonSelector)
	if err != nil {
		return fmt.Errorf("failed to click next button: %w", err)
	}

	return nil
}

func (w *WorkdayATS) LoginAccount(executor *atsclient.Executor, email, password string) error {
	err := executor.SetValue(emailInputSelector, email)
	if err != nil {
		return fmt.Errorf("failed to set email: %w", err)
	}

	err = executor.SetValue(passwordInputSelector, password)
	if err != nil {
		return fmt.Errorf("failed to set password: %w", err)
	}

	err = executor.WaitAndClick(loginButtonSelector)
	if err != nil {
		return fmt.Errorf("failed to click login button: %w", err)
	}

	return nil
}

func (w *WorkdayATS) ApplyForJob(executor *atsclient.Executor, resumePath string, qaMap map[string]string, aiEnabled bool) error {
	err := executor.WaitForElement(applicationFormSelector)
	if err != nil {
		return fmt.Errorf("application form not found: %w", err)
	}

	if resumePath != "" {
		err = executor.SetValue(resumeUploadSelector, resumePath)
		if err != nil {
			return fmt.Errorf("failed to upload resume: %w", err)
		}
	}

	err = w.ProcessQuestionsWithAnswers(executor, qaMap, aiEnabled)
	if err != nil {
		return fmt.Errorf("failed to process questions: %w", err)
	}

	err = executor.WaitAndClick(submitButtonSelector)
	if err != nil {
		return fmt.Errorf("failed to submit application: %w", err)
	}

	return nil
}

type Question struct {
	ID      string
	Type    string // text, radio, checkbox, dropdown
	Text    string
	Options []string // For radio/checkbox/dropdown questions
}

type ApplicationData struct {
	Questions  map[string]string // Maps question IDs to answers
	ResumePath string
}

func (w *WorkdayATS) ProcessQuestionsWithAnswers(executor *atsclient.Executor, qaMap map[string]string, aiEnabled bool) error {
	questionSelectors := []string{".question", "[data-automation-id*='question']", ".application-question"}

	for _, selector := range questionSelectors {
		elements, err := w.getQuestionElements(executor, selector)
		if err != nil {
			continue // If selector doesn't match anything, try the next one
		}

		for _, element := range elements {
			questionText, err := w.extractQuestionText(executor, element)
			if err != nil {
				continue
			}

			inputSelector, inputType, err := w.findInputForQuestion(executor, element)
			if err != nil {
				continue
			}

			answer := ""
			foundAnswer := false

			for qPattern, predefinedAnswer := range qaMap {
				if strings.Contains(strings.ToLower(questionText), strings.ToLower(qPattern)) {
					answer = predefinedAnswer
					foundAnswer = true
					break
				}
			}

			if !foundAnswer && aiEnabled {
				answer, err = w.generateAnswerWithAI(questionText)
				if err != nil {
					answer = "Appropriate response to: " + questionText
				}
			}

			switch inputType {
			case "text", "textarea":
				err = executor.SetValue(inputSelector, answer)
			case "radio", "checkbox":
				err = w.selectOption(executor, inputSelector, answer)
			case "select-one":
				err = executor.SetValue(inputSelector, answer)
			default:
				err = executor.SetValue(inputSelector, answer)
			}

			if err != nil {
				fmt.Printf("Failed to answer question '%s': %v\n", questionText, err)
			}
		}
	}

	return nil
}

func (w *WorkdayATS) getQuestionElements(executor *atsclient.Executor, selector string) ([]string, error) {
	return []string{}, nil
}

func (w *WorkdayATS) extractQuestionText(executor *atsclient.Executor, element string) (string, error) {
	return "", nil
}

func (w *WorkdayATS) findInputForQuestion(executor *atsclient.Executor, questionElement string) (string, string, error) {
	return "", "", nil
}

func (w *WorkdayATS) selectOption(executor *atsclient.Executor, selector, answer string) error {
	return nil
}

func (w *WorkdayATS) generateAnswerWithAI(question string) (string, error) {
	client, err := ai.NewAIClient()
	if err != nil {
		return "", fmt.Errorf("failed to create AI client: %w", err)
	}

	ctx := context.Background()
	answer, err := client.GenerateAnswer(ctx, question)
	if err != nil {
		return "", fmt.Errorf("failed to generate answer with AI: %w", err)
	}

	return answer, nil
}
