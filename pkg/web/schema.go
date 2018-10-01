package web

type loginParam struct {
	Username string `form:"username" json:"username" binding:"required" description:"用户名"`
	Password string `form:"password" json:"password" binding:"required" description:"密码"`
	Service  string `form:"service" json:"service,omitempty" `
	Referer  string `form:"referer" json:"referer,omitempty" `
	Remember string `form:"remember" json:"remember,omitempty" `
}

type forgotParam struct {
	Username string `form:"username" json:"username" binding:"required" description:"用户名"`
	Mobile   string `form:"mobile" json:"mobile" binding:"required" description:"手机号"`
	Email    string `form:"email" json:"email" binding:"required" description:"邮箱"`
}

type resetParam struct {
	Username  string `form:"username" json:"username" binding:"required" description:"用户名"`
	Password  string `form:"password" json:"password" binding:"required" description:"密码"`
	Password2 string `form:"password_confirm" json:"password2" binding:"required" description:"密码"`
	Token     string `form:"rt" json:"token" binding:"required" description:"token"`
}

// 修改密码，需要在登录后
type passwordParam struct {
	OldPassword  string `form:"old_password" json:"old_password" binding:"required" description:"旧密码"`
	NewPassword  string `form:"password" json:"password" binding:"required" description:"密码"`
	NewPassword2 string `form:"password_confirm" json:"password2" binding:"required" description:"密码确认"`
}

type idParam struct {
	ID int `form:"id" json:"id" binding:"required"`
}
