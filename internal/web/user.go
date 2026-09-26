package web

import (
	"net/http"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/johsqn/weboqk/internal/domain"
	"github.com/johsqn/weboqk/internal/service"
)

const (
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	// 和上面比起来，用 ` 看起来就比较清爽
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

type UserHandler struct {
	emailRexExp    *regexp.Regexp
	passwordRexExp *regexp.Regexp
	svc            *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		emailRexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:            svc,
	}
}

func (userHdl *UserHandler) RegisterRoutes(server *gin.Engine) {
	// REST 风格
	//server.POST("/user", userHdl.SignUp)
	//server.PUT("/user", userHdl.SignUp)
	//server.GET("/users/:username", userHdl.Profile)
	ug := server.Group("/users")
	// POST /users/signup
	ug.POST("/signup", userHdl.SignUp)
	// POST /users/login
	ug.POST("/login", userHdl.Login)
	// POST /users/edit
	ug.POST("/edit", userHdl.Edit)
	// GET /users/profile
	ug.GET("/profile", userHdl.Profile)
}

func (userHdl *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	isEmail, err := userHdl.emailRexExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isEmail {
		ctx.String(http.StatusOK, "非法邮箱格式")
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入密码不对")
		return
	}

	isPassword, err := userHdl.passwordRexExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "密码必须包含字母、数字、特殊字符，并且不少于八位")
		return
	}

	err = userHdl.svc.Signup(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	switch err {
	case nil:
		ctx.String(http.StatusOK, "注册成功")
	case service.ErrDuplicateEmail:
		ctx.String(http.StatusOK, "邮箱冲突，请换一个")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

func (userHdl *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	u, err := userHdl.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		sess := sessions.Default(ctx)
		sess.Set("userId", u.Id)
		sess.Options(sessions.Options{
			// 十五分钟
			MaxAge: 900,
		})
		err = sess.Save()
		if err != nil {
			ctx.String(http.StatusOK, "系统错误")
			return
		}
		ctx.String(http.StatusOK, "登录成功")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "用户名或者密码不对")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

func (userHdl *UserHandler) Edit(ctx *gin.Context) {

}

func (userHdl *UserHandler) Profile(ctx *gin.Context) {
	ctx.String(http.StatusOK, "这是 profile")
}
