package customError

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