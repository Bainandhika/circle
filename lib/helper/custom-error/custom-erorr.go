package customError

import "fmt"

func Success() string {
	return "Success"
}

func InternalServerError() string {
	return "Internal Server Error"
}

func NotFoundError(resource string) string {
	return "Resource " + resource + " not found"
}

func IncorrectPassword() string {
	return "Incorrect Password"
}

func NeedIDasPathParam() string {
	return "Need id as path param"
}

func NoRowAffected() string {
	return "No row affected"
}

func MustHaveAccount() string {
	return "You must have an account to create order"
}

func NeedQueryParams(jsonParams ...string) string {
	var params string
	for i, p := range jsonParams {
		if i == len(jsonParams)-1 {
			params += fmt.Sprintf("& %s", p)
		} else {
			params += fmt.Sprintf("%s, ", p)
		}
	}

    return fmt.Sprintf("Need query params: %s", params)
}
