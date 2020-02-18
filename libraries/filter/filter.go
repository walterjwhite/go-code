package filter

type Matcher interface {
	Matches(interface{}) bool
}

type Action interface {
	OnMatch(matcher Matcher, data interface{})
}

type Rule struct {
	Ordering       int
	Matcher        Matcher
	Actions        []Action
	FallbackAction Action
	MatchType MatchType
}

type CriteriaType int

const (
	Must CriteriaType = iota
	Should
)

type MatchType int

const (
	EqualsIgnoreCase MatchType = iota
	Equals
	ContainsIgnoreCase
	Contains
	Regex
	LessThan
	GreaterThan
)

type GroupRule struct {
	CriteriaType CriteriaType
	Invert       bool

	Rules []Matcher
}

func (r *GroupRule) Matches(data interface{}) bool {
	for i,rule := range r.Rules {
		
	}
	
}

func (r *Rule) Matches(data interface{}) bool {
	switch r.MatchType {
	case EqualsIgnoreCase:
		return true
	case Equals:
		return true
	case Contains:
		return true
	case Regex:
		return true
	case LessThan:
		return true
	case GreaterThan:
		return true
	}

	return false
}
