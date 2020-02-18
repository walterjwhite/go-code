package email

// email plugin
type EmailMessageField int

const (
	From EmailMessageField = iota
	To
	Cc
	Bcc
	Subject
	Message
	Attachments
)

type AttributeRule struct {
	EmailMessageField EmailMessageField
	Values            []string
	
	Invert bool
	CriteriaType CriteriaType
}

func (r *AttributeRule) Matches(data interface{}) bool {
	emailMessage, ok := email.(EmailMessage)
    if !ok {
        return false
    }

	switch r.EmailMessageField {
		case From:
			return emailMessage.From
		case To:
			return emailMessage.To
		case Cc:
			return emailMessage.Cc
		case Bcc:
			return emailMessage.Bcc
		case Subject:
			return emailMessage.Subject
		case Message:
			return emailMessage.Message
		case Attachments:
			return emailMessage.Attachments
	}
	emailMessage.
}

rule:
 ordering: int
 name: string
 rule: rule/group rule
 actions: []action

attribute rule:
 field: enum
 values: []string

group rule:
 criteria type:
 rules: []rule
