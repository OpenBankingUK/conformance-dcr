package compliant

import "errors"

// resolves what response type to use on register claims based on a list of supported from .wellknown
// DCR 3.2 spec allows: "code", "code id_token" or nil (defaults to "code id_token)
func responseTypeResolve(types *[]string) ([]string, error) {
	// nil provided so we pass nil (defaults to "code id_token)
	if types == nil {
		return nil, nil
	}

	// find if "code" and "code id_token" is present and add to the output
	// other values will be ignored
	var responseTypes []string
	for _, value := range *types {
		if value == "code" {
			responseTypes = append(responseTypes, value)
		}
		if value == "code id_token" {
			responseTypes = append(responseTypes, value)
		}
	}

	if responseTypes == nil {
		return nil, errors.New("supported response types must contain `code` and/or `code id_token`")
	}

	return responseTypes, nil
}
