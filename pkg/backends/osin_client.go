package backends

import (
	"time"

	"github.com/sethvargo/go-password/password"

	"github.com/liut/staffio/pkg/models/oauth"
	"github.com/liut/staffio/pkg/models/types"
)

// GenClientID ...
func GenClientID() string {
	now := time.Now()
	iid := types.IID(now.UnixNano())
	return iid.String()
}

// GenNewClient ...
func GenNewClient(name, redirectURI string) *oauth.Client {
	id := GenClientID()
	secret, err := password.Generate(28, 10, 0, false, false)
	if err != nil {
		logger().Infow("password generate fail", "err", err)
		return nil
	}
	client := oauth.NewClient(id, secret, redirectURI)
	client.Meta.Name = name
	logger().Debugw("new", "client", client, "secret", secret)

	return client
}
