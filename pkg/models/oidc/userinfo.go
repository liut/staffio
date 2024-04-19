package oidc

// UserInfo defined from https://openid.net/specs/openid-connect-core-1_0.html#UserInfo
// example:
//
//	{
//	  "sub": "248289761001",
//	  "name": "Jane Doe",
//	  "given_name": "Jane",
//	  "family_name": "Doe",
//	  "preferred_username": "j.doe",
//	  "email": "janedoe@example.com",
//	  "picture": "http://example.com/janedoe/me.jpg"
//	 }
type UserInfo struct {
	// basic of Standard
	// Subject - Identifier for the End-User at the Issuer.
	// 主题 - 发行方对最终用户的标识符。
	Sub string `json:"sub"`
	// Full name in displayable form including all name parts, possibly including titles and suffixes, ordered according to the End-User's locale and preferences.
	// 用于显示的全名，包括完整姓名，可能包括头衔和后缀，按照最终用户的语言环境和偏好排序。
	Name string `json:"name,omitempty"`
	// Given name(s) or first name(s) of the End-User. Note that in some cultures, people can have multiple given names; all can be present, with the names being separated by space characters.
	// 终端用户的名字或姓氏。请注意，在某些文化中，人们可能有多个名字；所有名字都可以出现，并用空格字符分隔。
	GivenName string `json:"given_name,omitempty"`
	// Surname(s) or last name(s) of the End-User. Note that in some cultures, people can have multiple family names or no family name; all can be present, with the names being separated by space characters.
	// 终端用户的姓氏。请注意，在某些文化中，人们可能有多个姓氏或无姓氏；这是完全可能出现的情况，并且各名字之间用空格分隔。
	FamilyName string `json:"family_name,omitempty"`
	// Casual name
	// 昵称
	Nickname string `json:"nickname,omitempty"`
	// Shorthand name by which the End-User wishes to be referred to at the RP.
	// The RP MUST NOT rely upon this value being unique
	// 终端用户希望在`RP`中被称为的简称，例如janedoe或j.doe。该值可以是任何有效的JSON字符串，包括特殊字符如@、/或空格。`RP`不得依赖于此值的唯一性。
	PreferredUsername string `json:"preferred_username,omitempty"`
	// preferred e-mail address
	Email         string `json:"email,omitempty"`
	EmailVerified bool   `json:"email_verified,omitempty"`
	// URL of the End-User's profile picture. This URL MUST refer to an image file (for example, a PNG, JPEG, or GIF image file), rather than to a Web page containing an image. Note that this URL SHOULD specifically reference a profile photo of the End-User suitable for displaying when describing the End-User, rather than an arbitrary photo taken by the End-User.
	Picture string `json:"picture,omitempty"`
}
