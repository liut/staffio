module github.com/liut/staffio

go 1.24

toolchain go1.24.2

require (
	daxv.cn/gopak/tencent-api-go v0.0.0-00010101000000-000000000000
	fhyx.online/lark-api-go v0.0.0-00010101000000-000000000000
	fhyx.online/welink-api-go v0.0.0-00010101000000-000000000000
	github.com/dchest/passwordreset v0.0.0-20190826080013-4518b1f41006
	github.com/getsentry/sentry-go v0.32.0
	github.com/getsentry/sentry-go/gin v0.32.0
	github.com/gin-gonic/gin v1.10.0
	github.com/go-jose/go-jose/v4 v4.1.0
	github.com/go-osin/osin v0.0.0-20240229065344-f0845653461e
	github.com/go-osin/session v1.3.4
	github.com/golang-jwt/jwt/v4 v4.5.2
	github.com/jmoiron/sqlx v1.4.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/lib/pq v1.10.9
	github.com/liut/keeper v0.0.0-20230310035549-ee21cc0ffcdd
	github.com/liut/simpauth v0.1.15
	github.com/liut/staffio-backend v0.3.1
	github.com/microcosm-cc/bluemonday v1.0.27
	github.com/mozillazg/go-slugify v0.2.0
	github.com/russross/blackfriday v1.6.0
	github.com/sethvargo/go-password v0.3.1
	github.com/spf13/cast v1.7.1
	github.com/spf13/cobra v1.9.1
	github.com/stretchr/testify v1.10.0
	github.com/ugorji/go/codec v1.2.12
	go.uber.org/zap v1.27.0
	golang.org/x/text v0.24.0
	gopkg.in/mail.v2 v2.3.1
)

replace (
	daxv.cn/gopak/tencent-api-go => github.com/fhyx/tencent-api-go v0.0.0-20250414082502-b24302591099
	fhyx.online/lark-api-go => github.com/fhyx/lark-api-go v0.0.0-20200706120648-5b8e3ff1d8ec
	fhyx.online/welink-api-go => github.com/fhyx/welink-api-go v0.0.0-20200606142605-e82b7908acf3
)

require (
	github.com/Azure/go-ntlmssp v0.0.0-20221128193559-754e69321358 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/bytedance/sonic v1.11.6 // indirect
	github.com/bytedance/sonic/loader v0.1.1 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dchest/authcookie v0.0.0-20190824115100-f900d2294c8e // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.8-0.20250403174932-29230038a667 // indirect
	github.com/go-ldap/ldap/v3 v3.4.11 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.20.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/css v1.0.1 // indirect
	github.com/icza/mighty v0.0.0-20180919140131-cfd07d671de6 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.7 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mozillazg/go-unidecode v0.1.1 // indirect
	github.com/pborman/uuid v1.2.1 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/philhofer/fwd v1.1.3-0.20240916144458-20a13a1f6b7c // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/tinylib/msgp v1.2.5 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/arch v0.8.0 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract [v0.0.0-0, v0.11.2]
