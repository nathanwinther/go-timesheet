package client

import (
    "encoding/json"
    "time"
    "timesheet/dao"
    "timesheet/invoice"
)

type Client struct {
    Id int64
    UserId int64
    Name string
    Description string
    Client *Company
    Company *Company
    Invoice *invoice.Invoice
    Fields map[string] string
}

type Company struct {
    Address string
    Contact string
}

func (c *Client) Save() error {
    q := `
        INSERT INTO user_client VALUES(
            NULL
            , ?
            , ?
            , ''
            , ?
            , ?
        );
    `

    params := []interface{} {
        c.UserId,
        c.Name,
        time.Now().Unix(),
        time.Now().Unix(),
    }

    result, err := dao.Exec(q, params)
    if err != nil {
        return err
    }

    cid, err := result.LastInsertId()
    if err != nil {
        return err
    }

    c.Id = cid

    return c.Update()
}

func (c *Client) String() (string, error) {
    b, err := json.MarshalIndent(c, "", "    ")
    if err != nil {
        return "", err
    }

    return string(b), nil
}

func (c *Client) Update() error {
    q := `
        UPDATE user_client SET
            data = ?
            , modified_date = ?
        WHERE id = ?;
    `

    data, err := c.String()
    if err != nil {
        return err
    }

    params := []interface{} {
        data,
        time.Now().Unix(),
        c.Id,
    }

    _, err = dao.Exec(q, params)

    return err
}

