package web

import (
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"strings"

	"github.com/alexedwards/scs/v2"
	"github.com/garnizeH/dimdim/service/user"
	"github.com/labstack/echo/v4"
	"github.com/microcosm-cc/bluemonday"
)

const (
	minPasswordLength = 1
	contextKeyEmail   = "email"
)

var (
	ErrInvalidEmail      = errors.New("invalid email")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrInvalidName       = errors.New("invalid name")
	ErrNotMatchPasswords = errors.New("passwords do not match")
	ErrInvalidToken      = errors.New("invalid token")
)

type signinRequest struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func (r *signinRequest) validate(c echo.Context, input *bluemonday.Policy) error {
	r.Email = input.Sanitize(strings.TrimSpace(r.Email))
	if _, e := mail.ParseAddress(r.Email); e != nil {
		return ErrInvalidEmail
	}

	r.Password = strings.TrimSpace(r.Password)
	if r.Password == "" || len(r.Password) < minPasswordLength {
		return ErrInvalidPassword
	}

	return nil
}

func (h *Handler) Signin(c echo.Context) error {
	r := signinRequest{}

	setFields := func() {
		setSessionDataFields(c, struct {
			Email    string
			Password string
		}{
			Email:    r.Email,
			Password: r.Password,
		})
	}

	if err := h.validateRequest(c, &r, "signin"); err != nil {
		setFields()
		return err
	}

	ctx := c.Request().Context()
	u, err := h.service.User().Signin(ctx, r.Email, r.Password)
	if err != nil {
		setFields()
		return h.errTmpl("signin", err.Error())
	}

	h.sess.Put(ctx, contextKeyEmail, u.Email)
	return c.Redirect(http.StatusSeeOther, "/")
}

type signupRequest struct {
	Email    string `form:"email"`
	Name     string `form:"name"`
	Password string `form:"password"`
	Confirm  string `form:"confirm"`
}

func (r *signupRequest) validate(c echo.Context, input *bluemonday.Policy) error {
	r.Name = input.Sanitize(strings.TrimSpace(r.Name))
	if r.Name == "" {
		return ErrInvalidName
	}

	r.Email = input.Sanitize(strings.TrimSpace(r.Email))
	email, err := mail.ParseAddress(r.Email)
	if err != nil {
		return ErrInvalidEmail
	}

	r.Email = email.Address

	r.Password = strings.TrimSpace(r.Password)
	r.Confirm = strings.TrimSpace(r.Confirm)
	if r.Password == "" || len(r.Password) < minPasswordLength {
		return ErrInvalidPassword
	}
	if r.Password != r.Confirm {
		return ErrNotMatchPasswords
	}

	return nil
}

func (h *Handler) Signup(c echo.Context) error {
	r := signupRequest{}

	setFields := func() {
		setSessionDataFields(c, struct {
			Email    string
			Name     string
			Password string
			Confirm  string
		}{
			Email:    r.Email,
			Name:     r.Name,
			Password: r.Password,
			Confirm:  r.Confirm,
		})
	}

	if err := h.validateRequest(c, &r, "signup"); err != nil {
		setFields()
		return err
	}

	ctx := c.Request().Context()
	err := h.service.User().Signup(ctx, h.baseURL, r.Email, r.Name, r.Password)
	if err != nil {
		setFields()
		return h.errTmpl("signup", err.Error())
	}

	return pageRendererWithFlashMsg(c, "signin", "check your mailbox")
}

func (h *Handler) SignupToken(c echo.Context) error {
	token := c.Param("token")

	ctx := c.Request().Context()
	u, err := h.service.User().ValidateSignupToken(ctx, token)
	if err != nil {
		return h.errTmpl("resend-signup-token", err.Error())
	}

	h.sess.Put(ctx, contextKeyEmail, u.Email)
	return c.Redirect(http.StatusSeeOther, "/")
}

func (h *Handler) Signout(c echo.Context) error {
	ctx := c.Request().Context()
	if err := h.sess.Destroy(ctx); err != nil {
		destroyCSRFCookie(c)
		return h.errMsg(err.Error())
	}

	return c.Redirect(http.StatusFound, "/auth/signin")
}

type resendSignupTokenRequest struct {
	Email string `form:"email"`
}

func (r *resendSignupTokenRequest) validate(c echo.Context, input *bluemonday.Policy) error {
	r.Email = input.Sanitize(strings.TrimSpace(r.Email))
	email, err := mail.ParseAddress(r.Email)
	if err != nil {
		return ErrInvalidEmail
	}

	r.Email = email.Address

	return nil
}

func (h *Handler) ResendSignupToken(c echo.Context) error {
	r := resendSignupTokenRequest{}

	setFields := func() {
		setSessionDataFields(c, struct {
			Email    string
			Password string
		}{
			Email:    r.Email,
			Password: "",
		})
	}

	if err := h.validateRequest(c, &r, "resend-signup-token"); err != nil {
		setFields()
		return err
	}

	ctx := c.Request().Context()
	err := h.service.User().ResendSignupToken(ctx, h.baseURL, r.Email)
	if err != nil {
		setFields()

		if errors.Is(err, user.ErrUserAlreadyVerified) {
			return h.errTmpl("signin", err.Error())
		}

		return h.errTmpl("resend-signup-token", err.Error())
	}

	return pageRendererWithFlashMsg(c, "signin", "check your mailbox")
}

type resetPasswordRequest struct {
	Email string `form:"email"`
}

func (r *resetPasswordRequest) validate(c echo.Context, input *bluemonday.Policy) error {
	r.Email = input.Sanitize(strings.TrimSpace(r.Email))
	email, err := mail.ParseAddress(r.Email)
	if err != nil {
		return ErrInvalidEmail
	}

	r.Email = email.Address

	return nil
}

func (h *Handler) ResetPassword(c echo.Context) error {
	r := resetPasswordRequest{}

	setFields := func() {
		setSessionDataFields(c, struct {
			Email string
		}{
			Email: r.Email,
		})
	}

	if err := h.validateRequest(c, &r, "reset-password"); err != nil {
		setFields()
		return err
	}

	ctx := c.Request().Context()
	err := h.service.User().ResetPassword(ctx, h.baseURL, r.Email)
	if err != nil {
		setFields()

		return h.errTmpl("reset-password", err.Error())
	}

	return pageRendererWithFlashMsg(c, "signin", "check your mailbox")
}

func (h *Handler) ResetPasswordToken(c echo.Context) error {
	token := c.Param("token")

	ctx := c.Request().Context()
	if err := h.service.User().ResetPasswordToken(ctx, token); err != nil {
		return h.errTmpl("reset-password", err.Error())
	}

	setSessionDataFields(c, struct {
		Token    string
		Password string
		Confirm  string
	}{
		Token:    token,
		Password: "",
		Confirm:  "",
	})
	return pageRendererWithFlashMsg(c, "reset-password-token", "")
}

type changePasswordRequest struct {
	Password string `form:"password"`
	Confirm  string `form:"confirm"`
}

func (r *changePasswordRequest) validate(c echo.Context, input *bluemonday.Policy) error {
	r.Password = strings.TrimSpace(r.Password)
	r.Confirm = strings.TrimSpace(r.Confirm)
	if r.Password == "" || len(r.Password) < minPasswordLength {
		return ErrInvalidPassword
	}
	if r.Password != r.Confirm {
		return ErrNotMatchPasswords
	}

	return nil
}

func (h *Handler) ChangePassword(c echo.Context) error {
	r := changePasswordRequest{}

	setFields := func() {
		setSessionDataFields(c, struct {
			Password string
			Confirm  string
		}{
			Password: r.Password,
			Confirm:  r.Confirm,
		})
	}

	if err := h.validateRequest(c, &r, "change-password"); err != nil {
		setFields()
		return err
	}

	ctx := c.Request().Context()
	email := h.sess.GetString(ctx, contextKeyEmail)
	if err := h.service.User().ChangePassword(ctx, email, r.Password); err != nil {
		setFields()

		return h.errTmpl("change-password", err.Error())
	}

	return pageRendererWithFlashMsg(c, "index", "password updated")
}

type changePasswordWithTokenRequest struct {
	Token    string `form:"token"`
	Password string `form:"password"`
	Confirm  string `form:"confirm"`
}

func (r *changePasswordWithTokenRequest) validate(c echo.Context, input *bluemonday.Policy) error {
	r.Token = strings.TrimSpace(r.Token)
	if r.Token == "" {
		return ErrInvalidToken
	}

	r.Password = strings.TrimSpace(r.Password)
	r.Confirm = strings.TrimSpace(r.Confirm)
	if r.Password == "" || len(r.Password) < minPasswordLength {
		return ErrInvalidPassword
	}
	if r.Password != r.Confirm {
		return ErrNotMatchPasswords
	}

	return nil
}

func (h *Handler) ChangePasswordWithToken(c echo.Context) error {
	r := changePasswordWithTokenRequest{}

	setFields := func() {
		setSessionDataFields(c, struct {
			Token    string
			Password string
			Confirm  string
		}{
			Token:    r.Token,
			Password: r.Password,
			Confirm:  r.Confirm,
		})
	}

	if err := h.validateRequest(c, &r, "reset-password-token"); err != nil {
		setFields()
		return err
	}

	ctx := c.Request().Context()
	u, err := h.service.User().ChangePasswordWithToken(ctx, r.Token, r.Password)
	if err != nil {
		setFields()

		return h.errTmpl("reset-password-token", err.Error())
	}

	h.sess.Put(ctx, contextKeyEmail, u.Email)
	return c.Redirect(http.StatusSeeOther, "/")
}

func sessionDataMiddleware(sessionManager *scs.SessionManager, users *user.Service, appName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			// Skip static endpoints.
			if strings.HasPrefix(req.URL.Path, "/static") {
				return next(c)
			}

			sessionData := SessionData{AppName: appName}

			ctx := req.Context()
			email := sessionManager.GetString(ctx, contextKeyEmail)
			if email != "" {
				user, err := users.GetUser(ctx, email)
				if err != nil {
					// TODO: need to clear the session/cookie and redirect to signin.
					panic(fmt.Sprintf("failed to get user with email %q: %v", email, err))
				}

				sessionData.Email = user.Email
				sessionData.Name = user.Name
			}
			tk, ok := c.Get("csc").(string)
			if ok {
				sessionData.CSRFToken = tk
			}

			c.Set("sessionData", sessionData)
			return next(c)
		}
	}
}

func signedInMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := c.Get("sessionData").(SessionData)
		if !sess.SignedIn() {
			return c.Redirect(http.StatusSeeOther, "/auth/signin")
		}

		return next(c)
	}
}

func signedOutMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := c.Get("sessionData").(SessionData)
		if sess.SignedIn() {
			return c.Redirect(http.StatusSeeOther, "/")
		}

		return next(c)
	}
}
