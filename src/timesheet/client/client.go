package client

import (
    "encoding/json"
    "timesheet/invoice"
)

type Client struct {
    Id string
    ClientName string
    ClientDescription string
    ClientAddress string
    ClientContact string
    InvoiceAddress string
    InvoiceContact string
    InvoiceRate float64
    Invoice *invoice.Invoice
    Fields map[string] string
}

func (c *Client) String() (string, error) {
    b, err := json.MarshalIndent(c, "", "    ")
    if err != nil {
        return "", err
    }

    return string(b), nil
}

