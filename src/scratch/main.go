package main

import (
    "fmt"
    "timesheet/client"
    "timesheet/invoice"
)

func main() {
    c := new(client.Client)
    c.ClientName = "Adverator"
    c.ClientDescription = "Adverator"
    c.InvoiceRate = 100.0

    v, err := invoice.New("2015-01-23", 16)
    if err != nil {
        panic(err)
    }

    c.Invoice = v
    c.Fields = map[string] string {}

    s, err := c.String()
    if err != nil {
        panic(err)
    }
    
    fmt.Println(s)
}

