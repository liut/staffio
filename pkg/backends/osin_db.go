package backends

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/RangelReale/osin"

	"lcgc/platform/staffio/pkg/models"
	"lcgc/platform/staffio/pkg/settings"
)

var (
	clients_sortable_fields           = []string{"id", "created"}
	_                       OSINStore = (*DbStorage)(nil)
)

type OSINStore interface {
	osin.Storage
	LoadClients(limit, offset int, sort map[string]int) ([]*models.Client, error)
	CountClients() uint
	GetClientWithCode(code string) (*models.Client, error)
	SaveClient(client *models.Client) error
	LoadScopes() (scopes []*models.Scope, err error)
	IsAuthorized(client_id, username string) bool
	SaveAuthorized(client_id, username string) error
}

type DbStorage struct {
	refresh map[string]string
	isDebug bool
}

func NewStorage() *DbStorage {

	s := &DbStorage{
		refresh: make(map[string]string),
		isDebug: settings.Debug,
	}

	return s
}

func (s *DbStorage) Clone() osin.Storage {
	return s
}

func (s *DbStorage) Close() {
}

func (s *DbStorage) logf(format string, args ...interface{}) {
	if s.isDebug {
		log.Printf(format, args...)
	}
}

func (s *DbStorage) GetClient(id string) (osin.Client, error) {
	s.logf("GetClient: '%s'", id)
	c, err := s.GetClientWithCode(id)
	if err == nil {
		return c, nil
	}
	return nil, fmt.Errorf("Client %q not found", id)
}

func (s *DbStorage) SaveAuthorize(data *osin.AuthorizeData) error {
	s.logf("SaveAuthorize: '%s'\n", data.Code)
	qs := func(tx dbTxer) error {
		sql := `INSERT INTO
		 oauth_authorization_code(code, client_id, username, redirect_uri, expires_in, scopes, created)
		 VALUES($1, $2, $3, $4, $5, $6, $7);`
		r, err := tx.Exec(sql, data.Code, data.Client.GetId(), data.UserData.(string),
			data.RedirectUri, data.ExpiresIn, data.Scope, data.CreatedAt)
		if err != nil {
			return err
		}

		s.logf("save authorizeData code %s OK %v", data.Code, r)

		return nil
	}
	return withTxQuery(qs)
}

func (s *DbStorage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	s.logf("LoadAuthorize: '%s'\n", code)
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
		if err != nil {
			return nil, err
		}
		s.logf("loaded authorization ok, createdAt %s", a.CreatedAt)
		return a, nil
	}

	s.logf("load authorize error: %s", err)
	return nil, fmt.Errorf("Authorize %q not found", code)
}

func (s *DbStorage) RemoveAuthorize(code string) error {
	s.logf("RemoveAuthorize: '%s'\n", code)
	if code == "" {
		log.Print("authorize code is empty")
		return nil
	}
	qs := func(tx dbTxer) error {
		sql := `DELETE FROM oauth_authorization_code WHERE code = $1;`
		r, err := tx.Exec(sql, code)
		if err != nil {
			return err
		}

		s.logf("delete authorizeData code %s OK %v", code, r)

		return nil
	}
	return withTxQuery(qs)
}

func (s *DbStorage) SaveAccess(data *osin.AccessData) error {
	s.logf("SaveAccess: '%s'\n", data.AccessToken)
	qs := func(tx dbTxer) error {
		str := `INSERT INTO
		 oauth_access_token(client_id, username, access_token, refresh_token, expires_in, scopes, created)
		 VALUES($1, $2, $3, $4, $5, $6, $7);`
		r, err := tx.Exec(str, data.Client.GetId(), data.UserData.(string),
			data.AccessToken, data.RefreshToken, data.ExpiresIn, data.Scope, data.CreatedAt)
		if err != nil {
			return err
		}

		s.logf("save AccessData token %s OK %v", data.AccessToken, r)

		if data.RefreshToken != "" {
			s.refresh[data.RefreshToken] = data.AccessToken
		}
		return nil
	}
	return withTxQuery(qs)
}

func (s *DbStorage) LoadAccess(code string) (*osin.AccessData, error) {
	s.logf("LoadAccess: '%s'", code)
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
		if err != nil {
			return nil, err
		}
		s.logf("access token '%d' expires: \n\t%s created \n\t%s expire_at \n\t%s now \n\tis_expired %v", id, a.CreatedAt, a.ExpireAt(), time.Now(), a.IsExpired())
		return a, nil
	}

	log.Printf("load access error: %s", err)
	return nil, fmt.Errorf("AccessToken %q not found", code)
}

func (s *DbStorage) RemoveAccess(code string) error {
	s.logf("RemoveAccess: %s\n", code)
	qs := func(tx dbTxer) error {
		str := `DELETE FROM oauth_access_token WHERE access_token = $1;`
		r, err := tx.Exec(str, code)
		if err != nil {
			return err
		}

		s.logf("delete accessToken %s OK %v", code, r)

		return nil
	}
	return withTxQuery(qs)
}

func (s *DbStorage) LoadRefresh(code string) (*osin.AccessData, error) {
	s.logf("LoadRefresh: %s\n", code)
	if d, ok := s.refresh[code]; ok {
		return s.LoadAccess(d)
	}
	return nil, fmt.Errorf("RefreshToken %q not found", code)
}

func (s *DbStorage) RemoveRefresh(code string) error {
	log.Printf("RemoveRefresh: %s\n", code)
	delete(s.refresh, code)
	return nil
}

func (s *DbStorage) GetClientWithCode(code string) (*models.Client, error) {
	c := new(models.Client)
	qs := func(db dber) error {
		return db.QueryRow("SELECT id, name, code, secret, redirect_uri, created FROM oauth_client WHERE code = $1",
			code).Scan(&c.Id, &c.Name, &c.Code, &c.Secret, &c.RedirectUri, &c.CreatedAt)
	}
	if err := withDbQuery(qs); err != nil {
		log.Printf("GetClientWithCode ERROR: %s", err)
		return nil, err
	}
	return c, nil
}

func (s *DbStorage) LoadClients(limit, offset int, sort map[string]int) (clients []*models.Client, err error) {
	if limit < 1 {
		limit = 1
	}
	if offset < 0 {
		offset = 0
	}

	var orders []string
	for k, v := range sort {
		if inArray(k, clients_sortable_fields) {
			var o string
			if v == ASCENDING {
				o = "ASC"
			} else {
				o = "DESC"
			}
			orders = append(orders, k+" "+o)
		}
	}

	str := `SELECT id, name, code, secret, redirect_uri, created
	  , allowed_grant_types, allowed_response_types, allowed_scopes
	   FROM oauth_client `

	if len(orders) > 0 {
		str = str + " ORDER BY " + strings.Join(orders, ",")
	}

	str = fmt.Sprintf("%s LIMIT %d OFFSET %d", str, limit, offset)

	clients = make([]*models.Client, 0)
	qs := func(db dber) error {
		rows, err := db.Query(str)
		if err != nil {
			log.Printf("db query error: %s for sql %s", err, str)
			return err
		}
		defer rows.Close()
		for rows.Next() {
			c := new(models.Client)
			var (
				grandTypes, responseTypes, scopes string
			)
			err = rows.Scan(&c.Id, &c.Name, &c.Code, &c.Secret, &c.RedirectUri, &c.CreatedAt,
				&grandTypes, &responseTypes, &scopes)
			if err != nil {
				log.Printf("rows scan error: %s", err)
				continue
			}
			c.AllowedGrantTypes = strings.Split(grandTypes, ",")
			c.AllowedResponseTypes = strings.Split(responseTypes, ",")
			c.AllowedScopes = strings.Split(scopes, ",")
			clients = append(clients, c)
		}
		return rows.Err()
	}

	if err := withDbQuery(qs); err != nil {
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

func (s *DbStorage) SaveClient(client *models.Client) error {
	log.Printf("SaveClient: id %d code %s", client.Id, client.Code)
	if client.Name == "" || client.Code == "" || client.Secret == "" || client.RedirectUri == "" {
		return valueError
	}
	qs := func(tx dbTxer) error {
		var err error
		if client.Id > 0 {
			str := `UPDATE oauth_client SET name = $1, code = $2, secret = $3, redirect_uri = $4
			 WHERE id = $5`
			var r sql.Result
			r, err = tx.Exec(str, client.Name, client.Code, client.Secret, client.RedirectUri, client.Id)
			log.Printf("UPDATE client result: %v", r)
		} else {
			str := `INSERT INTO
		 oauth_client(name, code, secret, redirect_uri, allowed_grant_types, allowed_scopes, created)
		 VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id;`
			err = tx.QueryRow(str,
				client.Name,
				client.Code,
				client.Secret,
				client.RedirectUri,
				strings.Join(client.AllowedGrantTypes, ","),
				strings.Join(client.AllowedScopes, ","),
				client.CreatedAt).Scan(&client.Id)
		}
		return err
	}
	return withTxQuery(qs)
}

func (s *DbStorage) LoadScopes() (scopes []*models.Scope, err error) {
	scopes = make([]*models.Scope, 0)
	qs := func(db dber) error {
		rows, err := db.Query("SELECT name, label, description, is_default FROM oauth_scope")
		if err != nil {
			log.Printf("load scopes error: %s", err)
			return err
		}
		defer rows.Close()
		for rows.Next() {
			s := new(models.Scope)
			err = rows.Scan(&s.Name, &s.Label, &s.Description, &s.IsDefault)
			if err != nil {
				log.Printf("rows scan error: %s", err)
			}
			scopes = append(scopes, s)
		}
		return rows.Err()
	}

	if err := withDbQuery(qs); err != nil {
		return nil, err
	}

	return scopes, nil
}

func (s *DbStorage) IsAuthorized(client_id, username string) bool {
	var (
		created time.Time
	)
	if err := withDbQuery(func(db dber) error {
		return db.QueryRow("SELECT created FROM oauth_client_user_authorized WHERE client_id = $1 AND username = $2",
			client_id, username).Scan(&created)
	}); err != nil {
		log.Printf("load IsAuthorized ERROR: %s", err)
		return false
	}
	return true
}

func (s *DbStorage) SaveAuthorized(client_id, username string) error {
	return withDbQuery(func(db dber) error {
		_, err := db.Exec("INSERT INTO oauth_client_user_authorized(client_id, username) VALUES($1, $2) ON CONFLICT DO NOTHING",
			client_id, username)
		return err
	})
}
