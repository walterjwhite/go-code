package greenhouse

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	atsclient "github.com/walterjwhite/go-code/lib/utils/ats-client"
)

type Account = atsclient.Account

type GreenhouseATS struct{}

func (g *GreenhouseATS) GetName() string {
	return "greenhouse"
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
		"\\", "",
		"&", "",
		"|", "",
		";", "",
		"<", "",
	)
	return replacer.Replace(strings.ToLower(input))
}

func (g *GreenhouseATS) RegisterAccount(executor *atsclient.Executor, account *Account) error {
	log.Println("Starting Greenhouse account registration")

	err := executor.Navigate("https://boards.greenhouse.io/")
	if err != nil {
		return fmt.Errorf("failed to navigate to Greenhouse: %w", err)
	}



	log.Println("Greenhouse typically handles job applications rather than account registration.")
	log.Println("Most Greenhouse implementations focus on the job application process.")

	return nil
}

func (g *GreenhouseATS) LoginAccount(executor *atsclient.Executor, email, password string) error {
	log.Println("Attempting to log into Greenhouse")


	err := executor.Navigate("https://boards.greenhouse.io/")
	if err != nil {
		return fmt.Errorf("failed to navigate to Greenhouse: %w", err)
	}

	log.Println("Greenhouse typically handles job applications rather than centralized account login.")
	log.Println("Each company using Greenhouse manages its own candidate accounts if any.")

	return nil
}

func (g *GreenhouseATS) ApplyForJob(executor *atsclient.Executor, resumePath string, qaMap map[string]string, aiEnabled bool) error {
	log.Println("Starting Greenhouse job application process")

	if resumePath != "" {
		if err := g.validateResumePath(resumePath); err != nil {
			return err
		}
		log.Println("Resume file validation passed")

		resumeSelectors := []string{
			"input[type='file'][name*='resume'], input[type='file'][data-resume], input[type='file'][id*='resume']",
			"input[type='file'][name*='attachment'], input[type='file'][name*='cv']",
			".file-input input[type='file'], .resume-upload input[type='file']",
		}

		foundResumeUpload := false
		for _, selector := range resumeSelectors {
			err := g.uploadFile(executor, selector, resumePath)
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
		if err := g.validateInput(questionPattern, answer); err != nil {
			log.Printf("Invalid input detected: %v, skipping question", err)
			continue
		}

		sanitizedPattern := sanitizeForSelector(questionPattern)
		sanitizedAnswer := sanitizeForSelector(answer)

		selector := fmt.Sprintf("input[type='text'][name*='%s'], textarea[name*='%s'], input[type='radio'][value*='%s'], input[type='checkbox'][value*='%s']",
			sanitizedPattern, sanitizedPattern, sanitizedAnswer, sanitizedAnswer)


		err := executor.SetValue(selector, answer)
		if err != nil {
			altSelector := fmt.Sprintf("[placeholder*='%s'], [title*='%s']", sanitizedPattern, sanitizedPattern)
			err = executor.SetValue(altSelector, answer)
			if err != nil {
				log.Println("Could not find field for question pattern, continuing")
			} else {
				log.Println("Successfully filled form field")
			}
		} else {
			log.Println("Successfully filled form field")
		}
	}

	for _, answer := range qaMap {
		sanitizedAnswer := sanitizeForSelector(answer)

		radioSelector := fmt.Sprintf("input[type='radio'][value*='%s'], input[type='checkbox'][value*='%s']",
			sanitizedAnswer, sanitizedAnswer)

		err := executor.Click(radioSelector)
		if err != nil {
			log.Println("Could not select radio/checkbox option")
		} else {
			log.Println("Successfully selected radio/checkbox option")
		}
	}

	submitSelectors := []string{
		"input[type='submit'][value*='Apply'], input[type='submit'][value*='Submit']",
		"button[type='submit'], button:contains('Apply'), button:contains('Submit')",
		".submit-btn, .apply-btn, #submit, #apply",
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

	log.Println("Greenhouse job application completed")
	return nil
}

func (g *GreenhouseATS) uploadFile(executor *atsclient.Executor, selector, filePath string) error {
	return executor.SetValue(selector, filePath)
}

func (g *GreenhouseATS) validateResumePath(resumePath string) error {
	absPath, err := filepath.Abs(resumePath)
	if err != nil {
		return fmt.Errorf("invalid resume path: %w", err)
	}

	fileInfo, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("resume file not accessible: %w", err)
	}

	if !fileInfo.Mode().IsRegular() {
		return fmt.Errorf("resume path is not a regular file")
	}

	return nil
}

func (g *GreenhouseATS) validateInput(pattern, value string) error {
	const maxInputLength = 1024
	if len(value) > maxInputLength {
		return fmt.Errorf("input exceeds maximum length of %d", maxInputLength)
	}

	if len(pattern) > maxInputLength {
		return fmt.Errorf("pattern exceeds maximum length of %d", maxInputLength)
	}

	dangerousChars := []string{"<", ">", "\"", "'", "`", "\\"}
	for _, char := range dangerousChars {
		if strings.Contains(value, char) || strings.Contains(pattern, char) {
			return fmt.Errorf("input contains potentially dangerous characters")
		}
	}

	return nil
}
