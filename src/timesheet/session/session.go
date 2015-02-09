package session

import (
    "net/http"
    "strconv"
    "time"
    "github.com/nathanwinther/go-uuid4"
    "timesheet/config"
    "timesheet/dao"
    "timesheet/logger"
    "timesheet/user"
)

type Session struct {
    Id int64
    Key string
    User *user.User
}

var (
    SESSION_OFFSET = int64(60 * 60 * 24 * 14)
)

func New(u *user.User) (*Session, error) {
    // Create new session
    skey, err := uuid4.New()
    if err != nil {
        return nil, err
    }

    q := `
        INSERT INTO user_session VALUES(
            NULL
            , ?
            , ?
            , ?
            , ?
            , ?
        );
    `

    params := []interface{} {
        skey,
        u.Id,
        time.Now().Unix() + SESSION_OFFSET,
        time.Now().Unix(),
        time.Now().Unix(),
    }

    sid, err := dao.Exec(q, params)
    if err != nil {
        return nil, err
    }

    return &Session{
        Id: sid,
        Key: skey,
        User: u,
    }, nil
}

func Parse(r *http.Request) (*Session, error) {
    c, err := r.Cookie(config.Get("session_cookie_name"))
    if err != nil {
        return nil, err
    }

    q := `
        SELECT
            s.id
            , u.key
        FROM user u, user_session s
        WHERE u.id = s.user_id
        AND u.active = 1
        AND s.key = ?
        AND s.valid_until > ?;
    `

    params := []interface{} {
        c.Value,
        time.Now().Unix(),
    }

    var sid int64
    var ukey string

    bind := []interface{} {
        &sid,
        &ukey,
    }

    err = dao.Row(q, params, bind)
    if err != nil {
        return nil, err
    }

    u, err := user.Load(ukey)
    if err != nil {
        return nil, err
    }

    return &Session{
        Id: sid,
        Key: c.Value,
        User: u,
    }, nil
}

func (s *Session) Save(w http.ResponseWriter, keepalive bool) error {
    q := `
        UPDATE user_session SET
            valid_until = ?
            , modified_date = ?
        WHERE id = ?;
    `

    valid := int64(0)

    if keepalive {
        valid = time.Now().Unix() + SESSION_OFFSET
    }

    params := []interface{} {
        valid,
        time.Now().Unix(),
        s.Id,
    }

    _, err := dao.Exec(q, params)
    if err != nil {
        return err
    }

    expires := -1

    if keepalive {
        expires, _ = strconv.Atoi(config.Get("session_cookie_expires"))
    }

    secure, _ := strconv.ParseBool(config.Get("session_cookie_secure"))

    c := new(http.Cookie)
    c.Name = config.Get("session_cookie_name")
    if keepalive {
        c.Value = s.Key
    }
    c.Path = config.Get("session_cookie_path")
    c.MaxAge = expires
    c.Secure = secure
    
    http.SetCookie(w, c)
    logger.Log(w, "SET-COOKIE", c.String())
    return nil
}

