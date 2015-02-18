<form action="{{.baseurl}}/new" class="standard register" method="POST" name="register">
    <p>
        <input class="standard{{if .errors.username}} error{{end}}" id="register-username" maxlength="256" name="username" placeholder="username" required="required" type="text" value="{{.form.username}}">
        {{if .errors.username}}<span class="error">{{.errors.username}}</span>{{end}}
    </p>
    <p>
        <input class="standard{{if .errors.email}} error{{end}}" id="register-email" maxlength="256" name="email" placeholder="email" required="required" type="email" value="{{.form.email}}">
        {{if .errors.email}}<span class="error">{{.errors.email}}</span>{{end}}
    </p>
    <p>
        <input class="standard{{if .errors.password}} error{{end}}" id="register-password" maxlength="256" name="password" placeholder="password" required="required" type="password">
        {{if .errors.password}}<span class="error">{{.errors.password}}</span>{{end}}
    </p>
    <p>
        <button class="standard" id="register-submit" type="submit">Sign up</button>
    </p>
    <p>
        Or <a href="{{.baseurl}}/u">Sign in.</a>
    </p>
</form>
