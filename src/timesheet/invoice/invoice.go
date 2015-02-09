package invoice

import (
)

type Entry struct {
    Key string
    Hours int
    Selected bool
    Today bool
    D int
    DD string
    DDD string
    DDDD string
    M int
    MM string
    MMM string
    MMMM string
    YY string
    YYYY int
}

type Invoice struct {
    StartDate int
    Days int
    EndDate int
    Selected int
    Entries []Entry
    Total int
}

