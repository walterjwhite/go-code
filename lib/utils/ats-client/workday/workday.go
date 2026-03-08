package workdayats

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

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

var (
	emailRegex   = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	phoneRegex   = regexp.MustCompile(`^\+?[0-9\s\-\(\)]{7,20}$`)
	nameRegex    = regexp.MustCompile(`^[a-zA-Z\s'-]{1,50}$`)
	zipCodeRegex = regexp.MustCompile(`^[0-9]{5}(-[0-9]{4})?$`)
)

type WorkdayATS struct {
	lastOperationTime time.Time
	rateLimitDelay    time.Duration
}

func NewWorkdayATS() *WorkdayATS {
	return &WorkdayATS{
		rateLimitDelay: 500 * time.Millisecond, // Minimum delay between operations
	}
}

func (w *WorkdayATS) GetName() string {
	return "workday"
}

func (w *WorkdayATS) applyRateLimit() {
	if time.Since(w.lastOperationTime) < w.rateLimitDelay {
		time.Sleep(w.rateLimitDelay - time.Since(w.lastOperationTime))
	}
	w.lastOperationTime = time.Now()
}

func (w *WorkdayATS) validateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	if len(email) > 254 {
		return fmt.Errorf("email exceeds maximum length")
	}
	return nil
}

func (w *WorkdayATS) validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	if len(password) > 128 {
		return fmt.Errorf("password exceeds maximum length")
	}
	return nil
}

func (w *WorkdayATS) validateName(name string) error {
	if !nameRegex.MatchString(name) {
		return fmt.Errorf("invalid name format")
	}
	return nil
}

func (w *WorkdayATS) validatePhone(phone string) error {
	if phone == "" {
		return nil // Phone is optional
	}
	if !phoneRegex.MatchString(phone) {
		return fmt.Errorf("invalid phone format")
	}
	return nil
}

func (w *WorkdayATS) validateZipCode(zipCode string) error {
	if zipCode == "" {
		return nil // Zip code is optional
	}
	if !zipCodeRegex.MatchString(zipCode) {
		return fmt.Errorf("invalid zip code format")
	}
	return nil
}

func (w *WorkdayATS) sanitizeInput(input string) string {
	input = strings.Map(func(r rune) rune {
		if r < 32 && r != 0 {
			return -1
		}
		return r
	}, input)
	return strings.TrimSpace(input)
}

func (w *WorkdayATS) validateResumePath(resumePath string) error {
	if resumePath == "" {
		return nil
	}

	if strings.Contains(resumePath, "..") {
		return fmt.Errorf("invalid resume path: path traversal detected")
	}

	ext := strings.ToLower(filepath.Ext(resumePath))
	allowedExtensions := []string{".pdf", ".doc", ".docx", ".txt", ".rtf"}
	allowed := slices.Contains(allowedExtensions, ext)
	if !allowed {
		return fmt.Errorf("invalid resume file type: only PDF, DOC, DOCX, TXT, RTF allowed")
	}

	if strings.Contains(resumePath, "://") {
		return fmt.Errorf("invalid resume path: URL schemes not allowed")
	}

	if strings.HasPrefix(resumePath, "http") {
		_, err := url.ParseRequestURI(resumePath)
		if err != nil {
			return fmt.Errorf("invalid resume URL format")
		}
	}

	return nil
}

func (w *WorkdayATS) RegisterAccount(executor *atsclient.Executor, account *Account) error {
	w.applyRateLimit()

	if err := w.validateEmail(account.Email); err != nil {
		return fmt.Errorf("email validation failed: %w", err)
	}

	if err := w.validatePassword(account.Password); err != nil {
		return fmt.Errorf("password validation failed: %w", err)
	}

	if err := w.validateName(account.FirstName); err != nil {
		return fmt.Errorf("first name validation failed: %w", err)
	}

	if err := w.validateName(account.LastName); err != nil {
		return fmt.Errorf("last name validation failed: %w", err)
	}

	if err := w.validatePhone(account.Phone); err != nil {
		return fmt.Errorf("phone validation failed: %w", err)
	}

	if err := w.validateZipCode(account.ZipCode); err != nil {
		return fmt.Errorf("zip code validation failed: %w", err)
	}

	account.Email = w.sanitizeInput(account.Email)
	account.FirstName = w.sanitizeInput(account.FirstName)
	account.LastName = w.sanitizeInput(account.LastName)
	account.Phone = w.sanitizeInput(account.Phone)
	account.Address = w.sanitizeInput(account.Address)
	account.City = w.sanitizeInput(account.City)

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
	w.applyRateLimit()

	if err := w.validateEmail(email); err != nil {
		return fmt.Errorf("email validation failed: %w", err)
	}

	if err := w.validatePassword(password); err != nil {
		return fmt.Errorf("password validation failed: %w", err)
	}

	email = w.sanitizeInput(email)

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
	if err := w.validateResumePath(resumePath); err != nil {
		return fmt.Errorf("resume path validation failed: %w", err)
	}

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
					answer = w.sanitizeInput(predefinedAnswer)
					foundAnswer = true
					break
				}
			}

			if !foundAnswer && aiEnabled {
				answer, err = w.generateAnswerWithAI(questionText)
				if err != nil {
					answer = "Appropriate response to: " + questionText
				}
				answer = w.sanitizeInput(answer)
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
				fmt.Printf("Failed to process question: %v\n", err)
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
