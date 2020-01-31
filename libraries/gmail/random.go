package gmail

import (
	"github.com/walterjwhite/go-application/libraries/random"
	"time"
)

func NewRandom(phonePreference *PhonePreference) *Account {
	account := &Account{PhonePreference: phonePreference}

	account.FirstName = randomString(firstNameMinimumLength)
	account.LastName = randomString(lastNameMinimumLength)
	account.Username = randomString(usernameMinimumLength)
	account.Password = randomString(passwordMinimumLength)

	account.Gender = Male
	//random.Of( /*int(Custom)*/ 4)
	account.BirthDate = &BirthDate{Month: 1 + random.Of(11),
		Day:  1 + random.Of(27),
		Year: birthYear()}

	// TODO: integrate secrets here
	// 1. secret label
	// 2. actual secret creation (email-address, first-name, last-name, password, birth-date ...)

	return account
}

func birthYear() int {
	age := randomAge()

	currentYear := time.Now().Year()

	return currentYear - age
}

func randomAge() int {
	return minimumAge + random.Of(maximumAge)
}

func randomString(minimumLength int) string {
	return random.String(minimumLength + random.Of(deviationFactor*minimumLength))
}
