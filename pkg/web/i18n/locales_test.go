package i18n

import (
	"net/http"
	"testing"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/stretchr/testify/assert"
)

func TestLocale(t *testing.T) {
	assert.Equal(t, language.SimplifiedChinese, zhHans)
	assert.Equal(t, language.TraditionalChinese, zhHant)
	assert.Equal(t, language.AmericanEnglish, enUS)
}

func TestErrorMessage(t *testing.T) {
	ers := []ErrorCode{ErrSystemError, ErrSystemReadonly, ErrParamRequired, ErrParamInvalid,
		ErrLoginFailed, ErrAuthRequired, ErrVerifySend,
		ErrRegistFaild, ErrNoneMobile, ErrBadAlias, ErrBadEmail, ErrBadMobile,
		ErrAliasTaken, ErrMobileTaken, ErrBadVerifyCode, ErrTokenExpired, ErrTokenInvalid,
		ErrOldPassword, ErrEmptyPassword, ErrSimplePassword,
		ErrEqualOldMobile, ErrSNSInfoLost, ErrSNSBindFailed, ErrEnableTwoFactor, ErrTwoFactorCode}
	printers := []*Printer{
		message.NewPrinter(enUS),
		message.NewPrinter(zhHans),
		message.NewPrinter(zhHant),
	}
	for _, e := range ers {
		assert.NotZero(t, e.Code())
		assert.NotEmpty(t, e.Error())
		for _, p := range printers {
			assert.NotEqual(t, e.String(), e.ErrorString(p))
		}

	}
}

func TestTag(t *testing.T) {
	tag, _ := language.MatchStrings(matcher, "zh-hans", "zh-CN")
	ers := []ErrorCode{ErrLoginFailed, ErrRegistFaild}
	for _, e := range ers {
		assert.NotEqual(t, e.String(), e.ErrorString(message.NewPrinter(tag)))
	}
}

func TestRequest(t *testing.T) {
	req, _ := http.NewRequest("GET", "/?lang=zh-hans", nil)
	req.Header.Set("Accept-Language", "zh-TW,zh-CN;q=0.9,zh;q=0.8,en;q=0.7,en-US;q=0.6")

	// tag := GetTag(req)
	ers := []ErrorCode{ErrLoginFailed, ErrRegistFaild}
	for _, e := range ers {
		assert.NotEqual(t, e.String(), e.ErrorString(GetPrinter(req)))
	}
}
