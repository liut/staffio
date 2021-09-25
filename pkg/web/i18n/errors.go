package i18n

//go:generate stringer -type=ErrorValue -trimprefix=ErrorValue -output errors_string.go

// ErrorValue ...
type ErrorValue int

// consts of error value
const (
	ErrOK             ErrorValue = iota // ok
	ErrSystemError                      // system error
	ErrSystemFailed                     // system faild
	ErrSystemPause                      // system pause
	ErrSystemReadonly                   // system readonly

	_

	ErrNotFound      // 404
	ErrForbiddedn    // 403
	ErrParamRequired // need some param or value input
	ErrParamInvalid  // invalid param

	_

	ErrLoginFailed    // Incorrect username or password
	ErrAuthRequired   // 401 need login
	ErrVerifySend     // old error (1402, "😓发送验证码失败")
	ErrRegistFaild    // old error (1403, "😓注册失败")
	ErrNoneMobile     // old error (1404, "😓没有这个手机号！")
	ErrBadAlias       // old error (1405, "🤔貌似不像合法的名称?")
	ErrBadEmail       // old error (1408, "🤔您确认这是邮箱地址?")
	ErrBadMobile      // old error (1406, "🤔您确认这是手机号码?")
	ErrAliasTaken     // old error (1407, "😓此用户已存在，请直接登录")
	ErrMobileTaken    // old error (1416, "😓此手机号已存在")
	ErrBadVerifyCode  // old error (1409, "😓验证码不正确")
	ErrTokenExpired   // old error (1410, "😓登录身份已过期")
	ErrTokenInvalid   // old error (1411, "😓登录身份无效或已过期")
	ErrOldPassword    // old error (1412, "😓原密码不正确")
	ErrEmptyPassword  // old error (1413, "😓密码不能为空")
	ErrSimplePassword // old error (1415, "😓您的密码太简单了")
	ErrMultiOnline    // old error (1414, "🤔您似乎已经登录了")
	ErrEqualOldMobile // old error (1417, "😓新手机号和旧的一样唉")
	ErrAliasTooFew    // 有一些必需的别名不能解绑

	_

	ErrEnableTwoFactor // old error (1420, "😓两步认证未开启，请先开启两步认证")
	ErrTwoFactorCode   // old error (1421, "😓两步认证验证码输入有误")

	_

	ErrSNSInfoLost   // 第三方(绑定)信息因过期而丢失
	ErrSNSBindFailed // 绑定第三方信息失败

)

// ErrorString return locale string with message printer
func (ev ErrorValue) ErrorString(p *Printer) string {
	switch ev {
	case ErrSystemError:
		return p.Sprintf("System error")
	case ErrSystemReadonly:
		return p.Sprintf("The system is currently in read-only mode.")
	case ErrParamRequired:
		return p.Sprintf("Required parameters")
	case ErrParamInvalid:
		return p.Sprintf("Invalid parameters")
	case ErrLoginFailed:
		return p.Sprintf("Incorrect username or password")
	case ErrAuthRequired:
		return p.Sprintf("You must be authenticated to see this resource")
	case ErrVerifySend:
		return p.Sprintf("Failed to send verifaction code")
	case ErrRegistFaild:
		return p.Sprintf("Failed to Register")
	case ErrNoneMobile:
		return p.Sprintf("No such mobile number")
	case ErrBadAlias:
		return p.Sprintf("Doesn't seem like a valid name?")
	case ErrBadEmail:
		return p.Sprintf("Are you sure this is the email address?")
	case ErrBadMobile:
		return p.Sprintf("Are you sure this is a cell phone number?")
	case ErrAliasTaken:
		return p.Sprintf("This user already exists, please sign in.")
	case ErrAliasTooFew:
		return p.Sprintf("Unbinding failed, at least one phone and email address required")
	case ErrMobileTaken:
		return p.Sprintf("This mobile already used.")
	case ErrBadVerifyCode:
		return p.Sprintf("The verification code is incorrect.")
	case ErrTokenExpired:
		return p.Sprintf("Expired token")
	case ErrTokenInvalid:
		return p.Sprintf("Invalid token")
	case ErrOldPassword:
		return p.Sprintf("The original password is incorrect.")
	case ErrEmptyPassword:
		return p.Sprintf("Password cannot be empty")
	case ErrSimplePassword:
		return p.Sprintf("Your password is too simple.")
	case ErrEqualOldMobile:
		return p.Sprintf("The new phone number is the same as the old one.")
	case ErrSNSInfoLost:
		return p.Sprintf("Expired or lost third-party information")
	case ErrSNSBindFailed:
		return p.Sprintf("Failure to bind third party information")
	case ErrEnableTwoFactor:
		return p.Sprintf("Two-step authentication not enabled")
	case ErrTwoFactorCode:
		return p.Sprintf("Two-step authentication code entered incorrectly")
	}
	return ev.String()
}

// Code ...
func (ev ErrorValue) Code() int {
	return int(ev)
}

func (ev ErrorValue) Error() string {
	return ev.String()
}
