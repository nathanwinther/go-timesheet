package validation

import (
    "fmt"
    "regexp"
    "strconv"
    "time"
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

func (v *Validation) Date(key string, val string, msg string) bool {
    _, err := time.Parse("2006-01-02", val)
    if err != nil {
        v.Errors[key] = msg
        return false
    }

    return true
}

func (v *Validation) Email(key string, val string, msg string) bool {
    if v.reEmail.Match([]byte(val)) {
        return true
    }

    v.Errors[key] = msg
    return false
}

func (v *Validation) Money(key string, val string, msg string) bool {
    f, err := strconv.ParseFloat(val, 64)
    if err != nil {
        v.Errors[key] = msg
        return false
    }

    if f <= 0 {
        v.Errors[key] = msg
        return false
    }

    return true
}

func (v *Validation) Positive(key string, val string, msg string) bool {
    i, err := strconv.ParseInt(val, 10, 64)
    if err != nil {
        v.Errors[key] = msg
        return false
    }

    if i <= 0 {
        v.Errors[key] = msg
        return false
    }

    return true
}

func (v *Validation) Require(key string, val string, msg string) bool {
    if val != "" {
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

