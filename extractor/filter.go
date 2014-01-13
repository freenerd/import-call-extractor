package extractor

import (
	"fmt"
	"strings"
)

var (
	NetworkCallSuspects = [...]string{
		"net.Dial",
		"net.DialTimeout",
		"net/http.Client",
		"crypto/tls.Client",
	}
)

func FilterForSuspectPackages(imports Imports) (Imports, error) {
	suspects := NetworkCallSuspects
	output := Imports{}

	for _, suspect := range suspects {
		splitSuspect := strings.Split(suspect, ".")

		if len(splitSuspect) != 2 {
			return nil, fmt.Errorf("incorrect suspect %s", suspect)
		}

		if pkg, ok := imports[splitSuspect[0]]; ok {
			if calls, ok := pkg[splitSuspect[1]]; ok {
				if _, present := output[splitSuspect[0]]; !present {
					output[splitSuspect[0]] = map[string][]string{}
				}

				output[splitSuspect[0]][splitSuspect[1]] = calls
			}
		}
	}

	return output, nil
}
