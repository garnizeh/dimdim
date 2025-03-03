package web

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/alexedwards/scs/v2"
	"github.com/garnizeH/dimdim/embeded"
	"github.com/garnizeH/dimdim/pkg/domain"
	"github.com/garnizeH/dimdim/service/user"
	"github.com/labstack/echo/v4"
	"github.com/microcosm-cc/bluemonday"
)

type SessionData struct {
	AppName   string
	Email     string
	Name      string
	ErrMsg    string
	FlashMsg  string
	CSRFToken string
	Fields    any
}

func (sd SessionData) SignedIn() bool {
	return sd.Email != ""
}

type Handler struct {
	baseURL string
	sess    *scs.SessionManager
	input   *bluemonday.Policy
	user    *user.Service
}

func NewHandler(
	domain domain.Domain,
	sess *scs.SessionManager,
	user *user.Service,
) *Handler {
	return &Handler{
		baseURL: domain.URL(""),
		sess:    sess,
		input:   bluemonday.StrictPolicy(),
		user:    user,
	}
}

func (h *Handler) LoadRoutes(e *echo.Echo, templates *embeded.Template) {
	// root
	templates.NewView("index", "base.tmpl", "menu.tmpl", "messages.tmpl", "index.tmpl")
	e.GET("/", pageRenderer("index"), signedInMiddleware)

	// auth
	auth := e.Group("/auth")
	h.loadRoutesAuth(auth, templates)
}

func (h *Handler) loadRoutesAuth(g *echo.Group, templates *embeded.Template) {
	// signin
	templates.NewView("signin", "base.tmpl", "messages.tmpl", "auth-links.tmpl", "auth/signin.tmpl")
	g.GET("/signin", pageRenderer("signin"), signedOutMiddleware)
	g.POST("/signin", h.Signin, signedOutMiddleware)

	// reset password
	templates.NewView("reset-password", "base.tmpl", "messages.tmpl", "auth-links.tmpl", "auth/reset-password.tmpl")
	templates.NewView("reset-password-token", "base.tmpl", "messages.tmpl", "auth/reset-password-token.tmpl")
	g.GET("/reset-password", pageRenderer("reset-password"), signedOutMiddleware)
	g.GET("/reset-password/:token", h.ResetPasswordToken)
	g.POST("/reset-password", h.ResetPassword, signedOutMiddleware)
	g.POST("/reset-password-token", h.ChangePasswordWithToken, signedOutMiddleware)

	// signup
	templates.NewView("signup", "base.tmpl", "messages.tmpl", "auth-links.tmpl", "auth/signup.tmpl")
	g.GET("/signup", pageRenderer("signup"), signedOutMiddleware)
	g.GET("/signup/:token", h.SignupToken)
	g.POST("/signup", h.Signup, signedOutMiddleware)
	g.GET("/signout", h.Signout, signedInMiddleware)

	// resend signup token
	templates.NewView("resend-signup-token", "base.tmpl", "messages.tmpl", "auth-links.tmpl", "auth/resend-signup-token.tmpl")
	g.GET("/resend-confirmation-email", pageRenderer("resend-signup-token"), signedOutMiddleware)
	g.POST("/resend-confirmation-email", h.ResendSignupToken, signedOutMiddleware)

	// change password
	templates.NewView("change-password", "base.tmpl", "messages.tmpl", "auth/change-password.tmpl")
	g.GET("/change-password", pageRenderer("change-password"), signedInMiddleware)
	g.POST("/change-password", h.ChangePassword, signedInMiddleware)
}

type validator interface {
	validate(echo.Context, *bluemonday.Policy) error
}

func (h *Handler) validateRequest(c echo.Context, req validator, tmpl ...string) error {
	var err error
	if err = c.Bind(req); err == nil {
		err = req.validate(c, h.input)
	}

	errTmpl := "index"
	if err != nil {
		if len(tmpl) > 0 {
			errTmpl = tmpl[0]
		}
		return h.errTmpl(errTmpl, err.Error())
	}

	return nil
}

func (h *Handler) errMsg(msg string) error {
	return h.errTmpl("index", msg)
}

type webError struct {
	msg  string
	tmpl string
}

func (we webError) Error() string {
	return we.msg
}

func (h *Handler) errTmpl(tmpl, msg string) error {
	err := webError{
		msg:  msg,
		tmpl: tmpl,
	}
	return echo.NewHTTPError(http.StatusOK).WithInternal(err)
}

func pageRenderer(page string) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess := getSessionData(c)

		return c.Render(http.StatusOK, page, sess)
	}
}

func pageRendererWithFlashMsg(c echo.Context, page, msg string) error {
	sess := getSessionData(c)
	sess.FlashMsg = msg

	return c.Render(http.StatusOK, page, sess)
}

func errorHandler(template *embeded.Template) func(error, echo.Context) {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		code := http.StatusInternalServerError
		msg := err.Error()
		tmpl := "index"
		he, ok := err.(*echo.HTTPError)
		if ok {
			code = he.Code
			out, ok := he.Internal.(webError)
			if ok {
				msg = out.msg
				if out.tmpl != "" {
					tmpl = out.tmpl
				}
			} else if m, _ := he.Message.(string); m != "" {
				msg = m
			}
		}

		if strings.HasPrefix(c.Request().URL.Path, "/static") {
			return
		}

		sess := getSessionData(c)
		sess.ErrMsg = msg

		buf := bytes.Buffer{}
		template.Render(&buf, tmpl, sess, c)
		m := buf.String()
		if err := c.HTML(code, m); err != nil {
			panic(fmt.Sprintf("failed to return html code: %v", err))
		}
	}
}

func getSessionData(c echo.Context) SessionData {
	common, _ := c.Get("sessionData").(SessionData)
	return common
}

func setSessionDataFields(c echo.Context, fields any) {
	common := getSessionData(c)
	common.Fields = fields
	c.Set("sessionData", common)
}

func destroyCSRFCookie(c echo.Context) {
	k, err := c.Cookie("_csc")
	if err != nil {
		return // nothing to do
	}

	k.Value = ""
	k.MaxAge = -1
	c.SetCookie(k)
}
