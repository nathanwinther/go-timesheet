package handler

import(
    "fmt"
    "html/template"
    "net/http"
    "path/filepath"
    "regexp"
    "strings"
    "github.com/nathanwinther/go-uuid4"
    "timesheet/config"
    "timesheet/flashdata"
    "timesheet/logger"
    "timesheet/session"
    "timesheet/user"
    "timesheet/validation"
)

type Handler struct {
    Rules []*Rule
    Templates *template.Template
    header string
}

type Rule struct {
    Pattern string
    Compiled *regexp.Regexp
    Handler func(http.ResponseWriter, *http.Request)
}

func New() (*Handler, error) {
    h := new(Handler)

    h.header = config.Get("response_header")

    h.Rules = []*Rule {
        &Rule{"GET:/timesheet", nil, h.handleHome},
        &Rule{"GET:/timesheet/forgot", nil, h.handleForgot},
        &Rule{"GET:/timesheet/logout", nil, h.handleLogout},
        &Rule{"GET:/timesheet/message", nil, h.handleMessage},
        &Rule{"GET:/timesheet/new", nil, h.handleNew},
        &Rule{"GET:/timesheet/purge", nil, h.handlePurge},
        &Rule{"GET:/timesheet/u", nil, h.handleLogin},
        &Rule{fmt.Sprintf("GET:/timesheet/u/%s", user.USERNAME_PATTERN),
            nil, h.handleUser},
        &Rule{fmt.Sprintf("GET:/timesheet/u/%s/client/new",
            user.USERNAME_PATTERN), nil, h.handleClientNew},
        &Rule{fmt.Sprintf("GET:/timesheet/u/%s/password",
            user.USERNAME_PATTERN), nil, h.handleUserPassword},
        &Rule{fmt.Sprintf("GET:/timesheet/u/%s/update", user.USERNAME_PATTERN),
            nil, h.handleUserUpdate},
        &Rule{"GET:/timesheet/verify/[A-Fa-f0-9][A-Fa-f0-9-]*", nil, h.handleVerify},
        &Rule{"POST:/timesheet/forgot", nil, h.handleForgotPost},
        &Rule{"POST:/timesheet/new", nil, h.handleNewPost},
        &Rule{"POST:/timesheet/u", nil, h.handleLoginPost},
        &Rule{fmt.Sprintf("POST:/timesheet/u/%s/password", user.USERNAME_PATTERN),
            nil, h.handleUserPasswordPost},
        &Rule{fmt.Sprintf("POST:/timesheet/u/%s/update", user.USERNAME_PATTERN),
            nil, h.handleUserUpdatePost},
    }

    // Compile rules
    for _, rule := range h.Rules {
        re, err := regexp.Compile(fmt.Sprintf("^%s$",
            strings.TrimRight(rule.Pattern, "/")))
        if err != nil {
            panic(err)
        }
        rule.Compiled = re
    }

    err := h.loadTemplates()
    if err != nil {
        return nil, err
    }

    return h, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    b := []byte(fmt.Sprintf("%s:/%s", r.Method, strings.Trim(r.URL.Path, "/")))

    rid, _ := uuid4.New()
    w.Header().Add(h.header, rid)

    logger.Log(w, "INCOMING", string(b))

    for _, rule := range h.Rules {
        if rule.Compiled.Match(b) {
            rule.Handler(w, r)
            return
        }
    }

    h.serveNotFound(w, r)
}

func (h *Handler) handleClientNew(w http.ResponseWriter, r *http.Request) {
    segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
    username := segments[2]

    s, err := session.Parse(r)
    if err != nil {
        http.Redirect(w, r, fmt.Sprintf("%s/u/%s",
            config.Get("baseurl"), username), http.StatusFound)
        return
    }

    if s.User.Username != username {
        http.Redirect(w, r, fmt.Sprintf("%s/u/%s",
            config.Get("baseurl"), username), http.StatusFound)
        return
    }

    m := map[string] interface{} {
        "baseurl": config.Get("baseurl"),
        "session": s,
        "form": map[string] string {
            "email": s.User.Email,
            "fullname": s.User.Fullname,
        },
    }

    h.Templates.ExecuteTemplate(w, "user_new_client.html", m)
}

func (h *Handler) handleForgot(w http.ResponseWriter, r *http.Request) {
    m := map[string] interface{} {
        "baseurl": config.Get("baseurl"),
    }

    h.Templates.ExecuteTemplate(w, "forgot.html", m)
}

func (h *Handler) handleForgotPost(w http.ResponseWriter, r *http.Request) {
    username := strings.TrimSpace(r.FormValue("username"))

    // Validate
    v := validation.New()
    v.Required("username", username, "username is required")

    if len(v.Errors) == 0 {
        u, err := user.Find(username)
        if err != nil {
            logger.Error(w, err)
            h.serveServerError(w, r)
            return
        }
        err = user.SendVerify(u.Id, u.Email, false)
        if err != nil {
            logger.Error(w, err)
            h.serveServerError(w, r)
            return
        }

        msg := `
            Password reset link sent
        `

        flashdata.Set(w, msg)

        url := fmt.Sprintf("%s/message", config.Get("baseurl"))
        http.Redirect(w, r, url, http.StatusFound)

        return
    }

    m := map[string] interface{} {
        "baseurl": config.Get("baseurl"),
        "form": map[string] string {
            "username": username,
        },
        "errors": v.Errors,
    }

    h.Templates.ExecuteTemplate(w, "forgot.html", m)
}

func (h *Handler) handleHome(w http.ResponseWriter, r *http.Request) {
    s, err := session.Parse(r)
    if err != nil {
        logger.Error(w, err)
    } else {
        err = s.Save(w, true)
        if err != nil {
            logger.Error(w, err)
        }
    }

    if s != nil {
        url := fmt.Sprintf("%s/u/%s", config.Get("baseurl"), s.User.Username)
        http.Redirect(w, r, url, http.StatusFound)
        return
    }

    m := map[string] interface{} {
        "baseurl": config.Get("baseurl"),
    }

    h.Templates.ExecuteTemplate(w, "home.html", m)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
    m := map[string] interface{} {
        "baseurl": config.Get("baseurl"),
    }

    h.Templates.ExecuteTemplate(w, "login.html", m)
}

func (h *Handler) handleLoginPost(w http.ResponseWriter, r *http.Request) {
    username := strings.TrimSpace(r.FormValue("username"))
    password := strings.TrimSpace(r.FormValue("password"))

    // Validate
    v := validation.New()
    v.Required("username", username, "username is required")
    v.Required("password", password, "password is required")

    if len(v.Errors) == 0 {
        u, err := user.Login(username, password)
        if err == nil {
            s, err := session.New(u)
            if err != nil {
                logger.Error(w, err)
                h.serveServerError(w, r)
                return
            }

            err = s.Save(w, true)
            if err != nil {
                logger.Error(w, err)
                h.serveServerError(w, r)
                return
            }

            url := fmt.Sprintf("%s/u/%s", config.Get("baseurl"), u.Username)
            http.Redirect(w, r, url, http.StatusFound)
            return
        } else {
            if err.Error() == "sql: no rows in result set" {
                v.Errors["username"] = "invalid username or password"
            }
        }
    }

    m := map[string] interface{} {
        "baseurl": config.Get("baseurl"),
        "form": map[string] string {
            "username": username,
        },
        "errors": v.Errors,
    }

    h.Templates.ExecuteTemplate(w, "login.html", m)
}

func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
    s, _ := session.Parse(r)
    if s != nil {
        s.Save(w, false)
    }
    http.Redirect(w, r, config.Get("baseurl"), http.StatusFound)
}

func (h *Handler) handleMessage(w http.ResponseWriter, r *http.Request) {
    s, ok := flashdata.Get(w, r)
    if !ok {
        http.Redirect(w, r, config.Get("baseurl"), http.StatusFound)
        return
    }

    s = strings.Replace(s, "\n", "<br>\n", -1)
    logger.Info(w, s)

    m := map[string] interface{} {
        "baseurl": config.Get("baseurl"),
        "message": template.HTML(s),
    }

    h.Templates.ExecuteTemplate(w, "message.html", m)
}

func (h *Handler) handleNew(w http.ResponseWriter, r *http.Request) {
    m := map[string] interface{} {
        "baseurl": config.Get("baseurl"),
    }

    h.Templates.ExecuteTemplate(w, "new.html", m)
}

func (h *Handler) handleNewPost(w http.ResponseWriter, r *http.Request) {
    username := strings.TrimSpace(r.FormValue("username"))
    email := strings.TrimSpace(r.FormValue("email"))
    password := strings.TrimSpace(r.FormValue("password"))

    // Validate
    v := validation.New()
    if v.Required("username", username, "username is required") {
        v.Username("username", username, "invalid username")
    }
    if v.Required("email", email, "email is required") {
        v.Email("email", email, "invalid email")
    }
    v.Required("password", password, "password is required")

    if len(v.Errors) == 0 {
        err := user.Add(username, email, password)
        if err == nil {
            msg := `
                Account created.
                Please check your email for your verification link
            `

            flashdata.Set(w, msg)

            url := fmt.Sprintf("%s/message", config.Get("baseurl"))
            http.Redirect(w, r, url, http.StatusFound)

            return
        } else {
            if err.Error() == "UNIQUE constraint failed: user.username" {
                v.Errors["username"] = "username already exists"
            } else if err.Error() == "UNIQUE constraint failed: user.email" {
                v.Errors["email"] = "email already exists"
            } else {
                logger.Error(w, err)
                h.serveServerError(w, r)
                return
            }
        }
    }

    m := map[string] interface{} {
        "baseurl": config.Get("baseurl"),
        "form": map[string] string {
            "username": username,
            "email": email,
        },
        "errors": v.Errors,
    }

    h.Templates.ExecuteTemplate(w, "new.html", m)
}

func (h *Handler) handlePurge(w http.ResponseWriter, r *http.Request) {
    err := h.loadTemplates()
    if err != nil {
        logger.Error(w, err)
        h.serveServerError(w, r)
        return
    }

    w.Header().Add("Content-Type", "text/plain")
    w.Write([]byte("Templates Reloaded"))
}

func (h *Handler) handleUser(w http.ResponseWriter, r *http.Request) {
    segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
    username := segments[2]

    s, err := session.Parse(r)
    if err != nil {
        logger.Error(w, err)
    } else {
        err = s.Save(w, true)
        if err != nil {
            logger.Error(w, err)
        }
    }

    var owner bool
    var u *user.User
    if s != nil && s.User.Username == username {
        owner = true
        u = s.User
    } else {
        owner = false
        u, err = user.LoadByUsername(username)
        if err != nil {
            logger.Error(w, err)
            h.serveNotFound(w, r)
            return
        }
    }

    msg, _ := flashdata.Get(w, r)

    m := map[string] interface{} {
        "baseurl": config.Get("baseurl"),
        "session": s,
        "message": msg,
        "user": u,
        "owner": owner,
    }

    h.Templates.ExecuteTemplate(w, "user.html", m)
}

func (h *Handler) handleUserPassword(w http.ResponseWriter, r *http.Request) {
    segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
    username := segments[2]

    s, err := session.Parse(r)
    if err != nil {
        http.Redirect(w, r, fmt.Sprintf("%s/u/%s",
            config.Get("baseurl"), username), http.StatusFound)
        return
    }

    if s.User.Username != username {
        http.Redirect(w, r, fmt.Sprintf("%s/u/%s",
            config.Get("baseurl"), username), http.StatusFound)
        return
    }

    msg, _ := flashdata.Get(w, r)

    m := map[string] interface{} {
        "baseurl": config.Get("baseurl"),
        "session": s,
        "message": msg,
    }

    h.Templates.ExecuteTemplate(w, "user_update_password.html", m)
}

func (h *Handler) handleUserPasswordPost(w http.ResponseWriter, r *http.Request) {
    segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
    username := segments[2]

    s, err := session.Parse(r)
    if err != nil {
        http.Redirect(w, r, fmt.Sprintf("%s/u/%s",
            config.Get("baseurl"), username), http.StatusFound)
        return
    }

    if s.User.Username != username {
        http.Redirect(w, r, fmt.Sprintf("%s/u/%s",
            config.Get("baseurl"), username), http.StatusFound)
        return
    }

    password := strings.TrimSpace(r.FormValue("password"))

    // Validate
    v := validation.New()
    v.Required("password", password, "new password is required")

    if len(v.Errors) == 0 {
        err := s.User.UpdatePassword(password)
        if err != nil {
            logger.Error(w, err)
            h.serveServerError(w, r)
            return
        }
        flashdata.Set(w, "Password updated")
        http.Redirect(w, r, fmt.Sprintf("%s/u/%s",
            config.Get("baseurl"), s.User.Username), http.StatusFound)
        return
    }

    m := map[string] interface{} {
        "baseurl": config.Get("baseurl"),
        "session": s,
        "errors": v.Errors,
    }

    h.Templates.ExecuteTemplate(w, "user_update_password.html", m)
}

func (h *Handler) handleUserUpdate(w http.ResponseWriter, r *http.Request) {
    segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
    username := segments[2]

    s, err := session.Parse(r)
    if err != nil {
        http.Redirect(w, r, fmt.Sprintf("%s/u/%s",
            config.Get("baseurl"), username), http.StatusFound)
        return
    }

    if s.User.Username != username {
        http.Redirect(w, r, fmt.Sprintf("%s/u/%s",
            config.Get("baseurl"), username), http.StatusFound)
        return
    }

    msg, _ := flashdata.Get(w, r)

    m := map[string] interface{} {
        "baseurl": config.Get("baseurl"),
        "session": s,
        "message": msg,
        "form": map[string] string {
            "email": s.User.Email,
            "fullname": s.User.Fullname,
        },
    }

    h.Templates.ExecuteTemplate(w, "user_update.html", m)
}

func (h *Handler) handleUserUpdatePost(w http.ResponseWriter, r *http.Request) {
    segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
    username := segments[2]

    s, err := session.Parse(r)
    if err != nil {
        http.Redirect(w, r, fmt.Sprintf("%s/u/%s",
            config.Get("baseurl"), username), http.StatusFound)
        return
    }

    if s.User.Username != username {
        http.Redirect(w, r, fmt.Sprintf("%s/u/%s",
            config.Get("baseurl"), username), http.StatusFound)
        return
    }

    email := strings.TrimSpace(r.FormValue("email"))
    fullname := strings.TrimSpace(r.FormValue("fullname"))

    // Validate
    v := validation.New()
    if v.Required("email", email, "email is required") {
        v.Email("email", email, "invalid email")
    }

    if len(v.Errors) == 0 {
        err = s.User.Update(email, fullname)
        if err != nil {
            logger.Error(w, err)
            h.serveServerError(w, r)
            return
        }
        flashdata.Set(w, "Profile updated")
        http.Redirect(w, r, fmt.Sprintf("%s/u/%s", config.Get("baseurl"),
            s.User.Username), http.StatusFound)
        return
    }

    m := map[string] interface{} {
        "baseurl": config.Get("baseurl"),
        "session": s,
        "form": map[string] string {
            "email": email,
            "fullname": fullname,
        },
        "errors": v.Errors,
    }

    h.Templates.ExecuteTemplate(w, "user_update.html", m)
}

func (h *Handler) handleVerify(w http.ResponseWriter, r *http.Request) {
    segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
    vid := segments[2]

    u, err := user.Verify(vid)
    if err != nil {
        logger.Error(w, err)
        h.serveServerError(w, r)
        return
    }

    s, err := session.New(u)
    if err != nil {
        logger.Error(w, err)
        h.serveServerError(w, r)
        return
    }

    // Drop cookie
    err = s.Save(w, true)
    if err != nil {
        logger.Error(w, err)
        h.serveServerError(w, r)
        return
    }

    url := fmt.Sprintf("%s/u/%s", config.Get("baseurl"), u.Username)
    http.Redirect(w, r, url, http.StatusFound)
}

func (h *Handler) loadTemplates() error {
    t, err := template.ParseGlob(filepath.Join(config.Get("templates"), "*.*"))
    if err != nil {
        return err
    }

    h.Templates = t

    return nil
}

func (h *Handler) serveNotFound(w http.ResponseWriter, r *http.Request) {
    s, _ := session.Parse(r)
    m := map[string] interface{} {
        "baseurl": config.Get("baseurl"),
        "session": s,
    }
    w.WriteHeader(http.StatusNotFound)
    h.Templates.ExecuteTemplate(w, "error404.html", m)
}

func (h *Handler) serveServerError(w http.ResponseWriter, r *http.Request) {
    s, _ := session.Parse(r)
    m := map[string] interface{} {
        "baseurl": config.Get("baseurl"),
        "session": s,
    }
    w.WriteHeader(http.StatusInternalServerError)
    h.Templates.ExecuteTemplate(w, "error500.html", m)
}

