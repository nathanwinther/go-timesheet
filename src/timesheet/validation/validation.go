package validation

import (
    "fmt"
    "regexp"
    "timesheet/user"
)

type Validation struct {
    Errors map[string] string
    reEmail *regexp.Regexp
    reUsername *regexp.Regexp
}

func New() *Validation {
    reEmail, _ := regexp.Compile("..*@..*")
    reUsername, _ := regexp.Compile(fmt.Sprintf("^%s$", user.USERNAME_PATTERN))
    return &Validation{
        Errors: map[string] string {},
        reEmail: reEmail,
        reUsername: reUsername,
    }
}

func (v *Validation) Required(key string, val string, msg string) bool {
    if val != "" {
        return true
    }

    v.Errors[key] = msg
    return false
}

func (v *Validation) Email(key string, val string, msg string) bool {
    if v.reEmail.Match([]byte(val)) {
        return true
    }

    v.Errors[key] = msg
    return false
}

func (v *Validation) Username(key string, val string, msg string) bool {
    if v.reUsername.Match([]byte(val)) {
        return true
    }

    v.Errors[key] = msg
    return false
}

