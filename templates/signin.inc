<form class="login standard" method="POST" name="login">
    <p>
        <input class="standard" id="login-username" maxlength="256" name="username" placeholder="username or email" required="required" type="text" value="{{.form.username}}">
        {{if .errors.username}}<span class="error">{{.errors.username}}</span>{{end}}
    </p>
    <p>
        <input class="standard" id="login-password" maxlength="256" name="password" placeholder="password" required="required" type="password">
        {{if .errors.password}}<span class="error">{{.errors.password}}</span>{{end}}
    </p>
    <p>
        <button class="standard" id="login-submit" type="submit">Sign in</button>
    </p>
    <p>
        <a href="{{.baseurl}}/new">Need an account? Sign up.</a>
    </p>
    <p>
        <a href="{{.baseurl}}/forgot">Forgot your password?</a>
    </p>
</form>
