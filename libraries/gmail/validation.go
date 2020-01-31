package gmail

import (
	"fmt"
	"github.com/walterjwhite/go-application/libraries/logging"
)

func (a *Account) validate() {
	if len(a.FirstName) < firstNameMinimumLength {
		validationError("First Name", ">=", firstNameMinimumLength, len(a.FirstName))
	}

	if len(a.LastName) < lastNameMinimumLength {
		validationError("Last Name", ">=", lastNameMinimumLength, len(a.LastName))
	}

	if len(a.Username) < usernameMinimumLength {
		validationError("Username", ">=", usernameMinimumLength, len(a.Username))
	}

	if len(a.Username) > usernameMaximumLength {
		validationError("Username", "<=", usernameMaximumLength, len(a.Username))
	}

	if len(a.Password) < passwordMinimumLength {
		validationError("Pasword", ">=", passwordMinimumLength, len(a.Password))
	}

	if len(a.PhonePreference.PhoneNumber) != phoneNumberLength {
		validationError("Phone Number", "=", phoneNumberLength, len(a.PhonePreference.PhoneNumber))
	}
}

func validationError(field string, operator string, length int, actualLength int) {
	logging.Panic(fmt.Errorf("%v should be %v %v characters (%v)", field, operator, length, actualLength))
}
