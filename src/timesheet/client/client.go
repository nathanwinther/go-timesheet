package client

import (
    "timesheet/user"
    "timesheet/invoice"
)

type Client struct {
    Id string
    Name string
    Invoice invoice.Invoice
}

