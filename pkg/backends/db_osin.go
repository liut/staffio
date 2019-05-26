package backends

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/openshift/osin"

	"github.com/liut/staffio/pkg/models/oauth"
)

var (
	clientsSortableFields = []string{"id", "created"}

	_ OSINStore = (*DbStorage)(nil)
)

type OSINStore = oauth.OSINStore

type DbStorage struct {
	refresh *sync.Map
	isDebug bool
}

func NewStorage() *DbStorage {
	s := &DbStorage{
		refresh: new(sync.Map),
	}

	return s
}

func (s *DbStorage) Clone() osin.Storage {
	return s
}

func (s *DbStorage) Close() {
}

func (s *DbStorage) GetClient(id string) (osin.Client, error) {
	c, err := s.GetClientWithCode(id)
	if err == nil {
		return c, nil
	}
	logger().Infow("Client not found", "id", id, "err", err)
	return nil, fmt.Errorf("Client %q not found", id)
}

func (s *DbStorage) SaveAuthorize(data *osin.AuthorizeData) error {
	qs := func(tx dbTxer) error {
		sql := `INSERT INTO
		 oauth_authorization_code(code, client_id, username, redirect_uri, expires_in, scopes, created)
		 VALUES($1, $2, $3, $4, $5, $6, $7);`
		r, err := tx.Exec(sql, data.Code, data.Client.GetId(), data.UserData.(string),
			data.RedirectUri, data.ExpiresIn, data.Scope, data.CreatedAt)
		if err != nil {
			logger().Infow("save authorizeData fail", "code", data.Code, "err", err)
			return err
		}

		logger().Debugw("save authorizeData OK", "code", data.Code, "result", r)

		return nil
	}
	return withTxQuery(qs)
}

func (s *DbStorage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	var (
		client_id string
		username  string
		err       error
	)
	a := &osin.AuthorizeData{Code: code}
	qs := func(db dber) error {
		return db.QueryRow(`SELECT client_id, username, redirect_uri, expires_in, scopes, created
		 FROM oauth_authorization_code WHERE code = $1`,
			code).Scan(&client_id, &username, &a.RedirectUri, &a.ExpiresIn, &a.Scope, &a.CreatedAt)
	}
	err = withDbQuery(qs)
	if err == nil {
		a.UserData = username
		a.Client, err = s.GetClientWithCode(client_id)
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

func (s *DbStorage) SaveAccess(data *osin.AccessData) error {
	qs := func(tx dbTxer) error {
		str := `INSERT INTO
		 oauth_access_token(client_id, username, access_token, refresh_token, expires_in, scopes, created)
		 VALUES($1, $2, $3, $4, $5, $6, $7);`
		_, err := tx.Exec(str, data.Client.GetId(), data.UserData.(string),
			data.AccessToken, data.RefreshToken, data.ExpiresIn, data.Scope, data.CreatedAt)
		if err != nil {
			logger().Infow("save accessData fail", "data", data, "err", err)
			return err
		}

		if data.RefreshToken != "" {
			s.refresh.Store(data.RefreshToken, data.AccessToken)
		}
		return nil
	}
	return withTxQuery(qs)
}

func (s *DbStorage) LoadAccess(code string) (*osin.AccessData, error) {
	var (
		client_id string
		username  string
		err       error
		is_frozen bool
		id        int
	)
	a := &osin.AccessData{AccessToken: code}
	qs := func(db dber) error {
		return db.QueryRow(`SELECT id, client_id, username, refresh_token, expires_in, scopes, is_frozen, created
		 FROM oauth_access_token WHERE access_token = $1`,
			code).Scan(&id, &client_id, &username, &a.RefreshToken, &a.ExpiresIn, &a.Scope, &is_frozen, &a.CreatedAt)
	}
	err = withDbQuery(qs)
	if err == nil {
		a.UserData = username
		a.Client, err = s.GetClientWithCode(client_id)
	}
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
	return nil, fmt.Errorf("RefreshToken %q not found", code)
}

func (s *DbStorage) RemoveRefresh(code string) error {
	s.refresh.Delete(code)
	return nil
}

func (s *DbStorage) GetClientWithCode(code string) (c *oauth.Client, err error) {
	c = new(oauth.Client)
	err = withDbQuery(func(db dber) error {
		return db.Get(c, "SELECT * FROM oauth_client WHERE code = $1", code)
	})
	return
}

func (s *DbStorage) GetClientWithID(id int) (c *oauth.Client, err error) {
	c = new(oauth.Client)
	err = withDbQuery(func(db dber) error {
		return db.Get(c, "SELECT * FROM oauth_client WHERE id = $1", id)
	})
	return
}

func (s *DbStorage) LoadClients(spec *oauth.ClientSpec) (clients []oauth.Client, err error) {
	if spec.Limit < 1 {
		spec.Limit = 1
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

	clients = make([]oauth.Client, 0)
	err = withDbQuery(func(db dber) error {
		return db.Select(&clients, str)
	})
	if err != nil {
		return nil, err
	}

	return clients, nil
}

func (s *DbStorage) CountClients() (total uint) {
	qs := func(db dber) error {
		return db.QueryRow("SELECT COUND(id) FROM oauth_client").Scan(&total)
	}
	withDbQuery(qs)
	return
}

func (s *DbStorage) SaveClient(client *oauth.Client) error {
	if client.Name == "" || client.Code == "" || client.Secret == "" || client.RedirectURI == "" {
		return valueError
	}
	qs := func(tx dbTxer) error {
		var err error
		if client.ID > 0 {
			str := `UPDATE oauth_client SET name = $1, code = $2, secret = $3, redirect_uri = $4
			 WHERE id = $5`
			_, err = tx.Exec(str, client.Name, client.Code, client.Secret, client.RedirectURI, client.ID)
			logger().Infow("UPDATE client result", "err", err)
		} else {
			str := `INSERT INTO
		 oauth_client(name, code, secret, redirect_uri, grant_types, scopes, created)
		 VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id;`
			err = tx.QueryRow(str,
				client.Name,
				client.Code,
				client.Secret,
				client.RedirectURI,
				client.AllowedGrantTypes,
				client.AllowedScopes,
				client.CreatedAt).Scan(&client.ID)
		}
		if err != nil {
			logger().Warnw("save client failed ", "client", client, "err", err)
		}
		return err
	}
	return withTxQuery(qs)
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

func (s *DbStorage) IsAuthorized(clientId, username string) bool {
	var (
		created time.Time
	)
	if err := withDbQuery(func(db dber) error {
		return db.QueryRow("SELECT created FROM oauth_client_user_authorized WHERE client_id = $1 AND username = $2",
			clientId, username).Scan(&created)
	}); err != nil {
		logger().Warnw("load isAuthorized fail", "clientId", clientId, "err", err)
		return false
	}
	return true
}

func (s *DbStorage) SaveAuthorized(clientId, username string) error {
	return withDbQuery(func(db dber) error {
		_, err := db.Exec("INSERT INTO oauth_client_user_authorized(client_id, username) VALUES($1, $2) ",
			clientId, username)
		return err
	})
}
