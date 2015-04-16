package invoice

import (
    "encoding/json"
    "strconv"
    "time"
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
    Rate float64
    Days int
    StartDate int
    EndDate int
    Selected int
    Entries []Entry
    Total int
}

func New(startdate string, days int, rate float64) (*Invoice, error) {
    t, err := time.Parse("2006-01-02", startdate)
    if err != nil {
        return nil, err
    }

    v := new(Invoice)
    v.Rate = rate
    v.Days = days
    v.EndDate = days - 1
    v.Entries = make([]Entry, days)

    for i := 0; i < days; i++ {
        e := new(Entry)

        e.Key = t.Format("2006-01-02")
        e.D, _ = strconv.Atoi(t.Format("02"))
        e.DD = t.Format("02")
        e.DDD = t.Format("Mon")
        e.DDDD = t.Format("Monday")
        e.M, _ = strconv.Atoi(t.Format("01"))
        e.MM = t.Format("01")
        e.MMM = t.Format("Jan")
        e.MMMM = t.Format("January")
        e.YY = t.Format("06")
        e.YYYY, _ = strconv.Atoi(t.Format("2006"))

        v.Entries[i] = *e

        t = t.Add(time.Hour * 24)
    }

    return v, nil
}

func (v *Invoice) String() (string, error) {
    b, err := json.MarshalIndent(v, "", "    ")
    if err != nil {
        return "", err
    }

    return string(b), nil
}

