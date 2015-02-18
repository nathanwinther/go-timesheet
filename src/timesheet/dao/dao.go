package dao

import (
    _ "github.com/mattn/go-sqlite3"
    "database/sql"
    "timesheet/config"
)

func Exec(q string, params []interface{}) (sql.Result, error) {
    db, err := sql.Open("sqlite3", config.Get("dbf"))
    if err != nil {
        return nil, err
    }
    defer db.Close()

    stmt, err := db.Prepare(q)
    if err != nil {
        return nil, err
    }
    defer stmt.Close()

    return stmt.Exec(params...)
}

func Row(q string, params []interface{}, bind []interface{}) error {
    db, err := sql.Open("sqlite3", config.Get("dbf"))
    if err != nil {
        return err
    }
    defer db.Close()

    stmt, err := db.Prepare(q)
    if err != nil {
        return err
    }
    defer stmt.Close()

    row := stmt.QueryRow(params...)
    if err != nil {
        return err
    }

    err = row.Scan(bind...)
    if err != nil {
        return err
    }

    return nil
}

