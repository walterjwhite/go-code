package format

type LanguageType struct {
	Name      string
	Command   string
	Arguments []string
}

var Languages = [2]LanguageType{}

func init() {
	Languages[0] = LanguageType{"Java", "mvn", []string{"com.coveo:fmt-maven-plugin:format"}}
	Languages[1] = LanguageType{"XML", "mvn", []string{"au.com.acegi:xml-format-maven-plugin:xml-format"}}
}
