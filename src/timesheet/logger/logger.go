package logger

import(
    "fmt"
    "net/http"
    "strings"
    "timesheet/config"
)

var (
    allow = false
    check = true
    header string
)

func allowed() bool {
    if check {
        ok := strings.ToLower(config.Get("logging"))
        allow = (ok == "true" || ok == "on" || ok == "1")
        header = config.Get("response_header")
        check = false
    }
    return allow
}

func Log(w http.ResponseWriter, messageType string, message string) {
    if allowed() {
        id := w.Header().Get(header)
        if id != "" {
            fmt.Printf("[%s] %s: %s\n", id, messageType, message)
        } else {
            fmt.Printf("%s: %s\n", messageType, message)
        }
    }
}

func Info(w http.ResponseWriter, message string) {
    Log(w, "INFO", message)
}

func Error(w http.ResponseWriter, err error) {
    Log(w, "ERROR", err.Error())
}

