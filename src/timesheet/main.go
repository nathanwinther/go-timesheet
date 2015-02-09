package main

import (
    "io/ioutil"
    "net/http"
    "os"
    "strconv"
    "timesheet/config"
    "timesheet/handler"
)

func main() {
    switch len(os.Args) {
        case 1:
            break
        case 2:
            err := os.Chdir(os.Args[1])
            if err != nil {
                panic(err)
            }
            break
        default:
            return
    }

    err := config.Load("./data/data.db")
    if err != nil {
        panic(err)
    }

    config.Dump(os.Stdout)

    h, err := handler.New()
    if err != nil {
        panic(err)
    }

    err = ioutil.WriteFile("./pid", []byte(strconv.Itoa(os.Getpid())), 0644)
    if err != nil {
        panic(err)
    }

    http.Handle("/", h)
    err = http.ListenAndServe(config.Get("bind"), nil)
    if err != nil {
        panic(err)
    }
}

