package utils

import (
	"bytes"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

func ExpandUriTemplate(tpl string, params map[string]string) (string, error) {
	re := regexp.MustCompile(`\{([?]?)(.+?)\}`)
	matches := re.FindAllStringSubmatch(tpl, -1)
	i := 0
	var err error
	result := re.ReplaceAllStringFunc(tpl, func(string) string {
		modifier := matches[i][1]
		match := matches[i][2]
		i += 1
		joinSep := ","
		var buf bytes.Buffer
		if modifier == "?" {
			joinSep = "&"
		}

		numWritten := 0
		for _, name := range strings.Split(match, ",") {
			val, defined := params[name]
			if !defined {
				if modifier == "?" {
					continue
				} else {
					err = fmt.Errorf("key %q required but was not found in params", name)
					break
				}
			}

			if numWritten == 0 {
				buf.WriteString(modifier)
			} else {
				buf.WriteString(joinSep)
			}
			if modifier == "?" {
				buf.WriteString(url.QueryEscape(name) + "=")
			}
			buf.WriteString(url.QueryEscape(val))
			numWritten += 1
		}

		return buf.String()
	})

	return result, err
}
