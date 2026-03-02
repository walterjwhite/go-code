package jobworx

import (
	"fmt"
	"log"
	"strings"

	atsclient "github.com/walterjwhite/go-code/lib/utils/ats-client"
)

type Account = atsclient.Account

type JobworxATS struct{}

func (j *JobworxATS) GetName() string {
	return "jobworx"
}

func (j *JobworxATS) RegisterAccount(executor *atsclient.Executor, account *Account) error {
	log.Printf("Starting Jobworx account registration for %s %s", account.FirstName, account.LastName)

	err := executor.Navigate("https://www.jobworx.com/")
	if err != nil {
		return fmt.Errorf("failed to navigate to Jobworx: %w", err)
	}


	registerSelectors := []string{
		"a[href*='register'], a[href*='Register'], a[href*='create'], a[href*='Create']",
		"a[href*='signup'], a[href*='SignUp'], a[href*='profile']",
		"button.register, button.create-account, button.signup",
		".register-link, .create-account, #register-btn",
	}

	registerFound := false
	for _, selector := range registerSelectors {
		err := executor.Click(selector)
		if err == nil {
			log.Println("Clicked register link/button")
			registerFound = true
			break
		}
	}

	if !registerFound {
		log.Println("No explicit register button found, proceeding with form fields")
	}

	fields := map[string]string{
		"input[name*='firstName'], input[name*='first_name'], input[id*='firstName'], input[id*='first_name']": account.FirstName,
		"input[name*='lastName'], input[name*='last_name'], input[id*='lastName'], input[id*='last_name']":     account.LastName,
		"input[name*='email'], input[id*='email'], input[name*='userName'], input[name*='username']":           account.Email,
		"input[name*='password'], input[id*='password']":                                                       account.Password,
		"input[name*='phone'], input[id*='phone'], input[name*='telephone']":                                   account.Phone,
		"input[name*='address'], input[id*='address']":                                                         account.Address,
		"input[name*='city'], input[id*='city']":                                                               account.City,
		"input[name*='state'], input[id*='state'], input[name*='region']":                                      account.State,
		"input[name*='zip'], input[name*='postal'], input[id*='zip'], input[id*='postal']":                     account.ZipCode,
		"input[name*='country'], input[id*='country']":                                                         account.Country,
	}

	for selector, value := range fields {
		if value != "" { // Only fill fields that have values
			err := executor.SetValue(selector, value)
			if err != nil {
				log.Printf("Could not set value for selector '%s': %v", selector, err)
			} else {
				log.Printf("Set value for field with selector '%s'", selector)
			}
		}
	}

	submitSelectors := []string{
		"input[type='submit'][value*='Register'], input[type='submit'][value*='Create']",
		"input[type='submit'][value*='register'], input[type='submit'][value*='create']",
		"button.register, button.submit, button[type='submit']",
		".submit-btn, .register-btn, #register-submit, #create-account",
	}

	submitted := false
	for _, selector := range submitSelectors {
		err := executor.Click(selector)
		if err != nil {
			log.Printf("Failed to click submit button with selector '%s': %v", selector, err)
		} else {
			log.Println("Registration form submitted successfully")
			submitted = true
			break
		}
	}

	if !submitted {
		return fmt.Errorf("could not find and click submit button for registration")
	}

	log.Println("Jobworx account registration completed")
	return nil
}

func (j *JobworxATS) LoginAccount(executor *atsclient.Executor, email, password string) error {
	log.Println("Starting Jobworx login")

	err := executor.Navigate("https://www.jobworx.com/")
	if err != nil {
		return fmt.Errorf("failed to navigate to Jobworx login page: %w", err)
	}

	loginSelectors := []string{
		"a[href*='login'], a[href*='Login'], a[href*='signin'], a[href*='SignIn']",
		"button.login, button.signin",
		".login-link, .signin-link, #login-btn",
	}

	loginFound := false
	for _, selector := range loginSelectors {
		err := executor.Click(selector)
		if err == nil {
			log.Println("Clicked login link/button")
			loginFound = true
			break
		}
	}

	if !loginFound {
		log.Println("No explicit login button found, proceeding with form fields")
	}

	loginFields := map[string]string{
		"input[name*='email'], input[id*='email'], input[name*='userName'], input[name*='username'], input[name*='login']": email,
		"input[name*='password'], input[id*='password']":                                                                   password,
	}

	for selector, value := range loginFields {
		if value != "" {
			err := executor.SetValue(selector, value)
			if err != nil {
				log.Printf("Could not set value for login field with selector '%s': %v", selector, err)
			} else {
				log.Printf("Set value for login field with selector '%s'", selector)
			}
		}
	}

	loginSubmitSelectors := []string{
		"input[type='submit'][value*='Login'], input[type='submit'][value*='Sign'], input[type='submit'][value*='Enter']",
		"input[type='submit'][value*='login'], input[type='submit'][value*='sign'], input[type='submit'][value*='enter']",
		"button.login, button.signin, button[type='submit']",
		".login-btn, .signin-btn, #login-submit, #signin-submit",
	}

	submitted := false
	for _, selector := range loginSubmitSelectors {
		err := executor.Click(selector)
		if err != nil {
			log.Printf("Failed to click login submit button with selector '%s': %v", selector, err)
		} else {
			log.Println("Login form submitted successfully")
			submitted = true
			break
		}
	}

	if !submitted {
		return fmt.Errorf("could not find and click submit button for login")
	}

	waitSelectors := []string{
		".dashboard, .profile, .welcome, #dashboard, #profile",
		".header, .navbar, .navigation", // Common elements that appear after login
	}

	loginCompleted := false
	for _, selector := range waitSelectors {
		err := executor.WaitForElement(selector)
		if err == nil {
			log.Println("Login appears to be successful - found post-login element")
			loginCompleted = true
			break
		}
	}

	if !loginCompleted {
		log.Println("Could not verify login success, but form was submitted")
	}

	log.Println("Jobworx login completed")
	return nil
}

func (j *JobworxATS) ApplyForJob(executor *atsclient.Executor, resumePath string, qaMap map[string]string, aiEnabled bool) error {
	log.Println("Starting Jobworx job application process")


	if resumePath != "" {
		log.Printf("Uploading resume from: %s", resumePath)

		resumeSelectors := []string{
			"input[type='file'][name*='resume'], input[type='file'][data-resume], input[type='file'][id*='resume']",
			"input[type='file'][name*='attachment'], input[type='file'][name*='cv'], input[type='file'][name*='document']",
			".file-input input[type='file'], .resume-upload input[type='file'], .attachment-field input[type='file']",
			"input[type='file'][name*='upload'], input[type='file'][id*='upload']",
		}

		foundResumeUpload := false
		for _, selector := range resumeSelectors {
			err := j.uploadFile(executor, selector, resumePath)
			if err == nil {
				log.Println("Resume uploaded successfully")
				foundResumeUpload = true
				break
			} else {
				log.Printf("Failed to upload resume using selector '%s': %v", selector, err)
			}
		}

		if !foundResumeUpload {
			log.Println("Could not find resume upload field, continuing without resume")
		}
	}

	for questionPattern, answer := range qaMap {
		sanitizedPattern := sanitizeForSelector(questionPattern)
		sanitizedAnswer := sanitizeForSelector(answer)

		textFieldSelector := fmt.Sprintf("input[type='text'][name*='%s'], input[type='text'][id*='%s'], textarea[name*='%s'], textarea[id*='%s']",
			sanitizedPattern, sanitizedPattern, sanitizedPattern, sanitizedPattern)

		err := executor.SetValue(textFieldSelector, answer)
		if err != nil {
			altTextFieldSelector := fmt.Sprintf("input[name*='%s'], input[id*='%s'], textarea[name*='%s'], textarea[id*='%s']",
				sanitizedPattern, sanitizedPattern, sanitizedPattern, sanitizedPattern)

			err = executor.SetValue(altTextFieldSelector, answer)
			if err != nil {
				log.Println("Could not find text field for question")
			} else {
				log.Println("Filled text field for question")
			}
		} else {
			log.Println("Filled text field for question")
		}

		radioCheckboxSelector := fmt.Sprintf("input[type='radio'][value*='%s'], input[type='checkbox'][value*='%s'], input[type='radio'][id*='%s'], input[type='checkbox'][id*='%s']",
			sanitizedAnswer, sanitizedAnswer, sanitizedAnswer, sanitizedAnswer)

		err = executor.Click(radioCheckboxSelector)
		if err != nil {
			log.Println("Could not select radio/checkbox")
		} else {
			log.Println("Selected radio/checkbox")
		}

		dropdownSelector := fmt.Sprintf("select[name*='%s'], select[id*='%s']",
			sanitizedPattern, sanitizedPattern)

		err = j.selectDropdownOption(executor, dropdownSelector, answer)
		if err != nil {
			log.Println("Could not select dropdown option")
		} else {
			log.Println("Selected dropdown option")
		}
	}

	if aiEnabled {
		log.Println("AI enhancement enabled - would process responses with AI if implemented")
	}

	submitSelectors := []string{
		"input[type='submit'][value*='Apply'], input[type='submit'][value*='Submit'], input[type='submit'][value*='Next']",
		"button.apply, button.submit, button.next, button[type='submit']",
		".apply-btn, .submit-btn, .next-btn, #apply, #submit, #next",
	}

	submitted := false
	for _, selector := range submitSelectors {
		err := executor.Click(selector)
		if err != nil {
			log.Printf("Failed to click submit button with selector '%s': %v", selector, err)
		} else {
			log.Println("Application submitted successfully")
			submitted = true
			break
		}
	}

	if !submitted {
		return fmt.Errorf("could not find and click submit button")
	}

	log.Println("Jobworx job application completed")
	return nil
}

func (j *JobworxATS) uploadFile(executor *atsclient.Executor, selector, filePath string) error {
	return executor.SetValue(selector, filePath)
}

func (j *JobworxATS) selectDropdownOption(executor *atsclient.Executor, selector, optionValue string) error {
	err := executor.SetValue(selector, optionValue)
	if err != nil {
		return err
	}
	return nil
}

func sanitizeForSelector(input string) string {
	replacer := strings.NewReplacer(
		"'", "",
		"\"", "",
		"[", "",
		"]", "",
		"(", "",
		")", "",
		":", "",
		",", "",
		".", "",
		"#", "",
		">", "",
		"+", "",
		"~", "",
		" ", "_",
	)
	return replacer.Replace(strings.ToLower(input))
}
