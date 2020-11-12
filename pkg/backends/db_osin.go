package backends

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/openshift/osin"

	"github.com/liut/staffio/pkg/models/oauth"
)

// vars
var (
	clientsSortableFields = []string{"id", "created"}

	_ OSINStore = (*DbStorage)(nil)

	ToJSONKV = oauth.ToJSONKV
)

type JSONKV = oauth.JSONKV
type Client = oauth.Client
type ClientSpec = oauth.ClientSpec
type OSINStore = oauth.OSINStore

func NewClient(id, secret, redirectURI string) *Client {
	return oauth.NewClient(id, secret, redirectURI)
}

type DbStorage struct {
	pageSize int
	refresh  *sync.Map
	isDebug  bool
}

func NewStorage() *DbStorage {
	s := &DbStorage{
		pageSize: 20,
		refresh:  new(sync.Map),
	}

	return s
}

func (s *DbStorage) Clone() osin.Storage {
	return s
}

func (s *DbStorage) Close() {
}

func (s *DbStorage) SaveAuthorize(data *osin.AuthorizeData) error {
	extra, err := ToJSONKV(data.UserData)
	if err != nil {
		logger().Infow("userData to json fail", "data", data, "err", err)
		return err
	}
	qs := func(tx dbTxer) error {
		sql := `INSERT INTO
		 oauth_authorization_code(code, client_id, userdata, redirect_uri, expires_in, scopes, created)
		 VALUES($1, $2, $3, $4, $5, $6, $7);`
		r, err := tx.Exec(sql, data.Code, data.Client.GetId(), extra,
			data.RedirectUri, data.ExpiresIn, data.Scope, data.CreatedAt)
		if err != nil {
			logger().Infow("save authorizeData fail", "code", data.Code, "userData", data.UserData, "err", err)
			return err
		}

		logger().Debugw("save authorizeData OK", "code", data.Code, "result", r)

		return nil
	}
	return withTxQuery(qs)
}

func (s *DbStorage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	var (
		clientID string
		userdata JSONKV
		err      error
	)
	a := &osin.AuthorizeData{Code: code}
	qs := func(db dber) error {
		return db.QueryRow(`SELECT client_id, userdata, redirect_uri, expires_in, scopes, created
		 FROM oauth_authorization_code WHERE code = $1`,
			code).Scan(&clientID, &userdata, &a.RedirectUri, &a.ExpiresIn, &a.Scope, &a.CreatedAt)
	}
	err = withDbQuery(qs)
	if err == nil {
		a.UserData = userdata
		a.Client, err = s.GetClient(clientID)
	}
	if err != nil {
		logger().Infow("load authorization fail", "code", code, "err", err)
		return nil, err
	}
	logger().Debugw("loaded authorization OK", "code", code, "created", a.CreatedAt)
	return a, nil
}

func (s *DbStorage) RemoveAuthorize(code string) error {
	if code == "" {
		logger().Infow("empty code when remove authorize")
		return nil
	}
	qs := func(tx dbTxer) error {
		sql := `DELETE FROM oauth_authorization_code WHERE code = $1;`
		r, err := tx.Exec(sql, code)
		if err != nil {
			logger().Infow("delete authorization fail", "code", code, "err", err)
			return err
		}

		logger().Debugw("delete authorizeData OK", "code", code, "r", r)

		return nil
	}
	return withTxQuery(qs)
}

func (s *DbStorage) SaveAccess(data *osin.AccessData) (err error) {
	_, err = s.LoadAccess(data.AccessToken)
	if err == nil {
		return
	}
	if err != ErrNotFound {
		logger().Infow("load access fail", "accessToken", data.AccessToken, "err", err)
		return
	}
	prev := ""
	authorizeData := &osin.AuthorizeData{}

	if data.AccessData != nil {
		prev = data.AccessData.AccessToken
	}

	if data.AuthorizeData != nil {
		authorizeData = data.AuthorizeData
	}

	var (
		extra JSONKV
	)
	if extra, err = ToJSONKV(data.UserData); err != nil {
		logger().Infow("access.userdata fail", "userdata", data.UserData, "err", err)
		return
	}
	if data.Client == nil {
		return valueError
	}
	qs := func(tx dbTxer) error {
		r, err := tx.Exec(`INSERT INTO oauth_access_token (
			client_id, authorize_code, previous, access_token, refresh_token, expires_in, scopes, created, userdata)
			    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			data.Client.GetId(), authorizeData.Code, prev, data.AccessToken, data.RefreshToken,
			data.ExpiresIn, data.Scope, data.CreatedAt, extra)
		if err != nil {
			logger().Infow("save accessData fail", "data", data, "err", err)
			return err
		}
		logger().Infow("save accessData ok", "data", data, "r", r)

		// debug("save AccessData token %s OK %v", data.AccessToken, r)
		// str := `INSERT INTO
		//  oauth_access_token(client_id, userdata, access_token, refresh_token, expires_in, scopes, created)
		//  VALUES($1, $2, $3, $4, $5, $6, $7);`
		// _, err := tx.Exec(str, data.Client.GetId(), data.UserData,
		// 	data.AccessToken, data.RefreshToken, data.ExpiresIn, data.Scope, data.CreatedAt)
		// if err != nil {
		// 	logger().Infow("save accessData fail", "data", data, "err", err)
		// 	return err
		// }

		if data.RefreshToken != "" {
			s.refresh.Store(data.RefreshToken, data.AccessToken)
		}
		return nil
	}
	return withTxQuery(qs)
}

func (s *DbStorage) LoadAccess(code string) (*osin.AccessData, error) {
	var (
		clientID, authorizeCode, prevAccessToken string
		userdata                                 JSONKV
		err                                      error
		isFrozen                                 bool
		id                                       int
	)
	a := &osin.AccessData{AccessToken: code}
	qs := func(db dber) error {
		return db.QueryRow(`SELECT id, client_id, userdata, refresh_token, authorize_code, previous, expires_in, scopes, is_frozen, created
		 FROM oauth_access_token WHERE access_token = $1`,
			code).Scan(&id, &clientID, &userdata, &a.RefreshToken, &authorizeCode, &prevAccessToken, &a.ExpiresIn, &a.Scope, &isFrozen, &a.CreatedAt)
	}
	err = withDbQuery(qs)
	if err != nil {
		logger().Infow("loadAccess fail", "code", code, "err", err)
		return nil, err
	}
	a.UserData = userdata
	a.Client, err = s.GetClient(clientID)
	a.AuthorizeData, _ = s.LoadAuthorize(authorizeCode)
	prevAccess, _ := s.LoadAccess(prevAccessToken)
	a.AccessData = prevAccess
	if err != nil {
		logger().Warnw("load access fail", "code", code, "err", err)
		return nil, err
	}
	return a, nil
}

func (s *DbStorage) RemoveAccess(code string) error {
	qs := func(tx dbTxer) error {
		str := `DELETE FROM oauth_access_token WHERE access_token = $1;`
		_, err := tx.Exec(str, code)
		if err != nil {
			return err
		}

		return nil
	}
	return withTxQuery(qs)
}

func (s *DbStorage) LoadRefresh(code string) (*osin.AccessData, error) {
	if v, ok := s.refresh.Load(code); ok {
		return s.LoadAccess(v.(string))
	}
	return nil, ErrNotFound
}

func (s *DbStorage) RemoveRefresh(code string) error {
	s.refresh.Delete(code)
	return nil
}

func (s *DbStorage) GetClient(id string) (c osin.Client, err error) {
	c = new(oauth.Client)
	err = withDbQuery(func(db dber) error {
		return db.Get(c, "SELECT * FROM oauth_client WHERE id = $1", id)
	})
	return
}

func (s *DbStorage) LoadClient(id string) (*Client, error) {
	c, err := s.GetClient(id)
	if err != nil {
		return nil, err
	}
	return c.(*Client), nil
}

func (s *DbStorage) LoadClients(spec *oauth.ClientSpec) (clients []oauth.Client, err error) {
	if spec.Limit < 1 {
		spec.Limit = s.pageSize
	}
	if spec.Page < 1 {
		spec.Page = 1
	} else {
		withDbQuery(func(db dber) error {
			err := db.Get(&spec.Total, "SELECT COUNT(id) as total FROM oauth_client")
			if err != nil {
				logger().Infow("count oauth_client ERR ", "err", err)
			}
			return err
		})
		if spec.Total == 0 {
			return
		}
	}

	str := `SELECT * FROM oauth_client `

	if len(spec.Orders) > 0 {
		var orders []string
		for _, order := range spec.Orders {
			if pos := strings.LastIndex(order, " "); pos > -1 {
				field := order[:pos]
				if inArray(field, clientsSortableFields) {
					sort := order[pos+1:]
					switch strings.ToUpper(sort) {
					case "ASC", "DESC", "":
						orders = append(orders, field+" "+sort)
					}
				}
			}
		}
		if len(orders) > 0 {
			str = str + " ORDER BY " + strings.Join(orders, ",")
		}
	}

	str = fmt.Sprintf("%s LIMIT %d OFFSET %d", str, spec.Limit, (spec.Page-1)*spec.Limit)

	clients = []oauth.Client{}
	err = withDbQuery(func(db dber) error {
		return db.Select(&clients, str)
	})
	if err != nil {
		return
	}

	return
}

func (s *DbStorage) CountClients() (total uint) {
	qs := func(db dber) error {
		return db.QueryRow("SELECT COUND(id) FROM oauth_client").Scan(&total)
	}
	withDbQuery(qs)
	return
}

func (s *DbStorage) SaveClient(client *oauth.Client) error {
	if client.ID == "" || client.Secret == "" {
		return valueError
	}
	qs := func(tx dbTxer) error {
		str := `INSERT INTO oauth_client(id, secret, redirect_uri, meta)
		 VALUES($1, $2, $3, $4) ON CONFLICT (id) DO UPDATE
		 SET secret = $5, redirect_uri = $6, meta = $7 RETURNING created`
		err := tx.QueryRow(str,
			client.ID,
			client.Secret,
			client.RedirectURI,
			client.Meta,
			client.Secret,
			client.RedirectURI,
			client.Meta).
			Scan(&client.CreatedAt)

		if err != nil {
			logger().Warnw("save client failed ", "client", client, "err", err)
		}
		return err
	}
	return withTxQuery(qs)
}

func (s *DbStorage) RemoveClient(id string) error {
	return withTxQuery(func(tx dbTxer) error {
		_, err := tx.Exec("DELETE FROM oauth_client WHERE id = $1", id)
		if err != nil {
			logger().Warnw("remove client failed ", "id", id, "err", err)
		}
		return err
	})
}

func (s *DbStorage) LoadScopes() (scopes []oauth.Scope, err error) {
	scopes = make([]oauth.Scope, 0)

	if err = withDbQuery(func(db dber) error {
		return db.Select(&scopes, "SELECT name, label, description, is_default FROM oauth_scope")
	}); err != nil {
		return nil, err
	}

	return scopes, nil
}

func (s *DbStorage) IsAuthorized(clientID, username string) bool {
	var (
		created time.Time
	)
	if err := withDbQuery(func(db dber) error {
		return db.QueryRow("SELECT created FROM oauth_client_user_authorized WHERE client_id = $1 AND username = $2",
			clientID, username).Scan(&created)
	}); err != nil {
		logger().Infow("load isAuthorized fail", "clientID", clientID, "username", username, "err", err)
		return false
	}
	return true
}

func (s *DbStorage) SaveAuthorized(clientID, username string) error {
	return withDbQuery(func(db dber) error {
		_, err := db.Exec("INSERT INTO oauth_client_user_authorized(client_id, username) VALUES($1, $2) ",
			clientID, username)
		return err
	})
}
