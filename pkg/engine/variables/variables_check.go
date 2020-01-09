package variables

import (
	"fmt"
	"regexp"
	"strconv"
)

//CheckVariables checks if the variable regex has been used
func CheckVariables(pattern interface{}, variables []string, path string) error {
	switch typedPattern := pattern.(type) {
	case map[string]interface{}:
		return checkMap(typedPattern, variables, path)
	case []interface{}:
		return checkArray(typedPattern, variables, path)
	case string:
		return checkValue(typedPattern, variables, path)
	default:
		return nil
	}
}

func checkMap(patternMap map[string]interface{}, variables []string, path string) error {
	for patternKey, patternElement := range patternMap {

		if err := CheckVariables(patternElement, variables, path+patternKey+"/"); err != nil {
			return err
		}
	}
	return nil
}

func checkArray(patternList []interface{}, variables []string, path string) error {
	for idx, patternElement := range patternList {
		if err := CheckVariables(patternElement, variables, path+strconv.Itoa(idx)+"/"); err != nil {
			return err
		}
	}
	return nil
}

func checkValue(valuePattern string, variables []string, path string) error {
	operatorVariable := getOperator(valuePattern)
	variable := valuePattern[len(operatorVariable):]
	if checkValueVariable(variable, variables) {
		return fmt.Errorf(path + valuePattern)
	}
	return nil
}

func checkValueVariable(valuePattern string, variables []string) bool {
	variableRegex := regexp.MustCompile(`\{\{([^{}]*)\}\}`)
	groups := variableRegex.FindAllStringSubmatch(valuePattern, -1)
	if len(groups) == 0 {
		// no variable defined
		return false
	}
	// if variables are defined, check against the list of variables to be filtered
	for _, group := range groups {
		if len(group) == 2 {
			// group[0] -> {{variable}}
			// group[1] -> variable
			if variablePatternSearch(group[1], variables) {
				return true
			}
		}
	}
	return false
}

func variablePatternSearch(pattern string, regexs []string) bool {
	for _, regex := range regexs {
		varRegex := regexp.MustCompile(regex)
		found := varRegex.FindString(pattern)
		if found != "" {
			return true
		}
	}
	return true
}
