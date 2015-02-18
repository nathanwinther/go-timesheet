package user

import (
    "bytes"
    "crypto/sha1"
    "encoding/base64"
    "fmt"
    "html/template"
    "path/filepath"
    "time"
    "github.com/nathanwinther/go-awsses"
    "github.com/nathanwinther/go-uuid4"
    "timesheet/config"
    "timesheet/dao"
)

type User struct {
    Id int64
    Key string
    Active int
    Username string
    Email string
    Fullname string
}

var (
    USERNAME_PATTERN = "[A-Za-z0-9][A-Za-z0-9_-]*"
    VERIFY_OFFSET = int64(60 * 60 * 24 * 2)
    Templates *template.Template
)

func Add(username string, email string, password string) error {
    // Add user
    ukey, err := uuid4.New()
    if err != nil {
        return err
    }

    q := `
        INSERT INTO user
        VALUES(
            NULL
            , ?
            , 0
            , ?
            , ?
            , ?
            , ''
            , ?
            , ?
        );
    `

    params := []interface{} {
        ukey,
        username,
        email,
        hashpassword(password),
        time.Now().Unix(),
        time.Now().Unix(),
    }

    result, err := dao.Exec(q, params)
    if err != nil {
        return err
    }

    uid, err := result.LastInsertId()
    if err != nil {
        return err
    }

    return SendVerify(uid, email, true)
}

func Find(username string) (*User, error) {
    u := new(User)

    q := `
        SELECT
            u.id
            , u.key
            , u.active
            , u.username
            , u.email
            , u.fullname
        FROM user u
        WHERE u.active = 1
        AND (
            u.username = ?
            OR u.email = ?
        )
    `

    params := []interface{} {
        username,
        username,
    }

    bind := []interface{} {
        &u.Id,
        &u.Key,
        &u.Active,
        &u.Username,
        &u.Email,
        &u.Fullname,
    }

    err := dao.Row(q, params, bind)

    return u, err
}

func Load(ukey string) (*User, error) {
    u := new(User)

    q := `
        SELECT
            u.id
            , u.key
            , u.active
            , u.username
            , u.email
            , u.fullname
        FROM user u
        WHERE u.key = ?;
    `

    params := []interface{} {
        ukey,
    }

    bind := []interface{} {
        &u.Id,
        &u.Key,
        &u.Active,
        &u.Username,
        &u.Email,
        &u.Fullname,
    }

    err := dao.Row(q, params, bind)

    return u, err
}

func LoadByUsername(username string) (*User, error) {
    u := new(User)

    q := `
        SELECT
            u.id
            , u.key
            , u.active
            , u.username
            , u.email
            , u.fullname
        FROM user u
        WHERE u.username = ?;
    `

    params := []interface{} {
        username,
    }

    bind := []interface{} {
        &u.Id,
        &u.Key,
        &u.Active,
        &u.Username,
        &u.Email,
        &u.Fullname,
    }

    err := dao.Row(q, params, bind)

    return u, err
}

func Login(username string, password string) (*User, error) {
    u := new(User)

    q := `
        SELECT
            u.id
            , u.key
            , u.active
            , u.username
            , u.email
            , u.fullname
        FROM user u
        WHERE u.active = 1
        AND (
            u.username = ?
            OR u.email = ?
        )
        AND u.password = ?;
    `

    params := []interface{} {
        username,
        username,
        hashpassword(password),
    }

    bind := []interface{} {
        &u.Id,
        &u.Key,
        &u.Active,
        &u.Username,
        &u.Email,
        &u.Fullname,
    }

    err := dao.Row(q, params, bind)

    return u, err
}

func SendVerify(uid int64, email string, activate bool) error {
    // Add verify
    vkey, err := uuid4.New()
    if err != nil {
        return err
    }

    q := `
        INSERT INTO user_verify
        VALUES(
            NULL
            , ?
            , ?
            , ?
            , ?
            , ?
        );
    `

    params := []interface{} {
        vkey,
        uid,
        time.Now().Unix() + VERIFY_OFFSET,
        time.Now().Unix(),
        time.Now().Unix(),
    }

    _, err = dao.Exec(q, params)
    if err != nil {
        return err
    }

    // Send verify email
    url := fmt.Sprintf("%s/verify/%s", config.Get("baseurl"), vkey)

    loadTemplates()

    tpl := "reset"
    subject := "Timesheet - Password reset"
    if activate {
        tpl = "activate"
        subject = "Timesheet - Please active your account"
    }

    html := new(bytes.Buffer)
    Templates.ExecuteTemplate(html, fmt.Sprintf("%s.html", tpl), url)

    text := new(bytes.Buffer)
    Templates.ExecuteTemplate(text, fmt.Sprintf("%s.txt", tpl), url)

    m := awsses.New(
        config.Get("awsses_sender"),
        email,
        subject,
        html.String(),
        text.String())

    return m.Send(
        config.Get("awsses_baseurl"),
        config.Get("awsses_accesskey"),
        config.Get("awsses_secretkey"))
}

func Verify(vkey string) (*User, error) {
    // Get user key from verify
    q := `
        SELECT
            u.key
        FROM user_verify v, user u
        WHERE v.user_id = u.id
        AND v.key = ?
        AND v.valid_until > ?;
    `

    params := []interface{} {
        vkey,
        time.Now().Unix(),
    }

    var ukey string

    bind := []interface{} {
        &ukey,
    }

    err := dao.Row(q, params, bind)
    if err != nil {
        return nil, err
    }

    // Activate user
    q = `
        UPDATE user SET
            active = ?
            , modified_date = ?
        WHERE key = ?;
    `

    params = []interface{} {
        1,
        time.Now().Unix(),
        ukey,
    }

    _, err = dao.Exec(q, params)
    if err != nil {
        return nil, err
    }

    u, err := Load(ukey)
    if err != nil {
        return nil, err
    }

    q = `
        UPDATE user_verify SET
            valid_until = 0
            , modified_date = ?
        WHERE key = ?;
    `

    params = []interface{} {
        time.Now().Unix(),
        vkey,
    }

    _, err = dao.Exec(q, params)
    if err != nil {
        return nil, err
    }

    return u, nil
}

func (u *User) Update(email string, fullname string) error {
    q := `
        UPDATE user SET
            email = ?
            , fullname = ?
            , modified_date = ?
        WHERE id = ?;
    `

    params := []interface{} {
        email,
        fullname,
        time.Now().Unix(),
        u.Id,
    }

    _, err := dao.Exec(q, params)

    return err
}

func (u *User) UpdatePassword(password string) error {
    q := `
        UPDATE user SET
            password = ?
            , modified_date = ?
        WHERE id = ?;
    `

    params := []interface{} {
        hashpassword(password),
        time.Now().Unix(),
        u.Id,
    }

    _, err := dao.Exec(q, params)

    return err
}

func hashpassword(s string) string {
    h := sha1.New()
    h.Write([]byte(s))
    return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func loadTemplates() {
    if Templates == nil {
        t, _ := template.ParseGlob(
            filepath.Join(config.Get("email"), "*.*"))
        Templates = t
    }
}

