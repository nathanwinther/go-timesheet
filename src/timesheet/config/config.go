package config

import (
    _ "github.com/mattn/go-sqlite3"
    "database/sql"
    "fmt"
    "io"
    "sort"
)

var (
    m = map[string] string {}
)

func Dump(w io.Writer) {
    keys := make([]string, len(m))

    i := 0
    for k, _ := range m {
        keys[i] = k
        i = i + 1
    }

    sort.Strings(keys)

    for _, k := range keys {
        fmt.Fprintf(w, "%s = %s\n", k, m[k])
    }
}

func Load(dbf string) error {
    m["dbf"] = dbf

    q := `
        SELECT
            key
            , value
        FROM config
        ORDER BY key ASC;
    `

    db, err := sql.Open("sqlite3", m["dbf"])
    if err != nil {
        return err
    }
    defer db.Close()

    rows, err := db.Query(q)
    if err != nil {
        return err
    }
    defer rows.Close()

    for rows.Next() {
        var key string
        var value string

        err = rows.Scan(&key, &value)
        if err != nil {
            return err
        }

        m[key] = value
    }

    return nil
}

func Get(key string) string {
    return m[key]
}

func Test(key string) (string, bool) {
    v, ok := m[key]
    return v, ok
}

