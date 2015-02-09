package flashdata

import (
    "encoding/base64"
    "fmt"
    "net/http"
    "strconv"
    "strings"
    "timesheet/config"
    "timesheet/logger"
)

func Get(w http.ResponseWriter, r *http.Request) (string, bool) {
    c, err := r.Cookie(fmt.Sprintf("%s-flash",
        config.Get("session_cookie_name")))
    if err != nil {
        return "", false
    }

    b, err := base64.URLEncoding.DecodeString(c.Value)
    if err != nil {
        return "", false
    }

    Set(w, "")

    return string(b), true
}

func Set(w http.ResponseWriter, s string) {
    secure, _ := strconv.ParseBool(config.Get("session_cookie_secure"))

    c := new(http.Cookie)
    c.Name = fmt.Sprintf("%s-flash", config.Get("session_cookie_name"))
    c.Path = config.Get("session_cookie_path")
    c.Value = base64.URLEncoding.EncodeToString([]byte(strings.TrimSpace(s)))
    if c.Value != "" {
        c.MaxAge = 0
    } else {
        c.MaxAge = -1
    }
    c.Secure = secure

    http.SetCookie(w, c)
    logger.Log(w, "SET-COOKIE", c.String())
}

