package backends

import (
	"database/sql"
	"errors"
	"github.com/RangelReale/osin"
	_ "github.com/lib/pq"
	"log"
	"time"
	"tuluu.com/liut/staffio/models"
)

var (
	dbError = errors.New("db error")
)

type DbStorage struct {
	// clients   map[string]osin.Client
	// authorize map[string]*osin.AuthorizeData
	// access  map[string]*osin.AccessData
	refresh map[string]string
}

func NewStorage() *DbStorage {

	r := &DbStorage{
		// clients:   make(map[string]osin.Client),
		// authorize: make(map[string]*osin.AuthorizeData),
		// access:  make(map[string]*osin.AccessData),
		refresh: make(map[string]string),
	}

	return r
}

func (s *DbStorage) Clone() osin.Storage {
	return s
}

func (s *DbStorage) Close() {
}

func (s *DbStorage) GetClient(id string) (osin.Client, error) {
	log.Printf("GetClient: '%s'", id)
	c, err := GetClientWithCode(id)
	if err == nil {
		return c, nil
	}
	return nil, errors.New("Client not found")
}

// func (s *DbStorage) SetClient(id string, client osin.Client) error {
// 	log.Printf("SetClient: %s\n", id)
// 	// unused
// 	// s.clients[id] = client
// 	return nil
// }

func (s *DbStorage) SaveAuthorize(data *osin.AuthorizeData) error {
	log.Printf("SaveAuthorize: '%s'\n", data.Code)
	qs := func(tx *sql.Tx) error {
		sql := `INSERT INTO
		 oauth_authorization_code(code, client_id, username, redirect_uri, expires_in, scopes, created)
		 VALUES($1, $2, $3, $4, $5, $6, $7);`
		r, err := tx.Exec(sql, data.Code, data.Client.GetId(), data.UserData.(string),
			data.RedirectUri, data.ExpiresIn, data.Scope, data.CreatedAt)
		if err != nil {
			return err
		}

		log.Printf("save authorizeData code %s OK %v", data.Code, r)

		return nil
	}
	return withTxQuery(qs)
}

func (s *DbStorage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	log.Printf("LoadAuthorize: '%s'\n", code)
	var (
		client_id string
		username  string
		err       error
	)
	a := &osin.AuthorizeData{Code: code}
	err = getDb().QueryRow("SELECT client_id, username, redirect_uri, expires_in, scopes, created FROM oauth_authorization_code WHERE code = $1",
		code).Scan(&client_id, &username, &a.RedirectUri, &a.ExpiresIn, &a.Scope, &a.CreatedAt)
	if err == nil {
		a.UserData = username
		a.Client, err = GetClientWithCode(client_id)
		if err != nil {
			return nil, err
		}
		log.Printf("loaded authorization ok, createdAt %s", a.CreatedAt)
		return a, nil
	}

	log.Printf("load authorize error: %s", err)
	return nil, errors.New("Authorize not found")
}

func (s *DbStorage) RemoveAuthorize(code string) error {
	log.Printf("RemoveAuthorize: '%s'\n", code)
	if code == "" {
		log.Print("authorize code is empty")
		return nil
	}
	qs := func(tx *sql.Tx) error {
		sql := `DELETE FROM oauth_authorization_code WHERE code = $1;`
		r, err := tx.Exec(sql, code)
		if err != nil {
			return err
		}

		log.Printf("delete authorizeData code %s OK %v", code, r)

		return nil
	}
	return withTxQuery(qs)
}

func (s *DbStorage) SaveAccess(data *osin.AccessData) error {
	log.Printf("SaveAccess: '%s'\n", data.AccessToken)
	qs := func(tx *sql.Tx) error {
		sql := `INSERT INTO
		 oauth_access_token(client_id, username, access_token, refresh_token, expires_in, scopes, created)
		 VALUES($1, $2, $3, $4, $5, $6, $7);`
		r, err := tx.Exec(sql, data.Client.GetId(), data.UserData.(string),
			data.AccessToken, data.RefreshToken, data.ExpiresIn, data.Scope, data.CreatedAt)
		if err != nil {
			return err
		}

		log.Printf("save AccessData token %s OK %v", data.AccessToken, r)

		if data.RefreshToken != "" {
			s.refresh[data.RefreshToken] = data.AccessToken
		}
		return nil
	}
	return withTxQuery(qs)
}

func (s *DbStorage) LoadAccess(code string) (*osin.AccessData, error) {
	log.Printf("LoadAccess: '%s'", code)
	var (
		client_id string
		username  string
		err       error
		is_frozen bool
		id        int
	)
	a := &osin.AccessData{AccessToken: code}
	err = getDb().QueryRow(`SELECT id, client_id, username, refresh_token, expires_in, scopes, is_frozen, created
		 FROM oauth_access_token WHERE access_token = $1`,
		code).Scan(&id, &client_id, &username, &a.RefreshToken, &a.ExpiresIn, &a.Scope, &is_frozen, &a.CreatedAt)
	if err == nil {
		a.UserData = username
		a.Client, err = GetClientWithCode(client_id)
		if err != nil {
			return nil, err
		}
		log.Printf("access token '%d' expires: \n\t%s created \n\t%s expire_at \n\t%s now \n\tis_expired %v", id, a.CreatedAt, a.ExpireAt(), time.Now(), a.IsExpired())
		return a, nil
	}

	log.Printf("load access error: %s", err)
	return nil, errors.New("Access not found")
}

func (s *DbStorage) RemoveAccess(code string) error {
	log.Printf("RemoveAccess: %s\n", code)
	qs := func(tx *sql.Tx) error {
		sql := `DELETE FROM oauth_access_token WHERE access_token = $1;`
		r, err := tx.Exec(sql, code)
		if err != nil {
			return err
		}

		log.Printf("delete accessToken %s OK %v", code, r)

		return nil
	}
	return withTxQuery(qs)
}

func (s *DbStorage) LoadRefresh(code string) (*osin.AccessData, error) {
	log.Printf("LoadRefresh: %s\n", code)
	if d, ok := s.refresh[code]; ok {
		return s.LoadAccess(d)
	}
	return nil, errors.New("Refresh not found")
}

func (s *DbStorage) RemoveRefresh(code string) error {
	log.Printf("RemoveRefresh: %s\n", code)
	delete(s.refresh, code)
	return nil
}

func GetClientWithCode(code string) (*models.Client, error) {
	c := new(models.Client)
	err := getDb().QueryRow("SELECT id, name, code, secret, redirect_uri, created FROM oauth_client WHERE code = $1",
		code).Scan(&c.Id, &c.Name, &c.Code, &c.Secret, &c.RedirectUri, &c.Created)
	if err == nil {
		return c, nil
	}
	log.Printf("GetClientWithCode ERROR: %s", err)
	return nil, dbError
}

func LoadScopes() (scopes []*Scope, err error) {
	scopes = make([]*Scope, 0)
	rows, err := getDb().Query("SELECT name, label, description, is_default FROM oauth_scope")
	if err != nil {
		log.Fatalf("load scopes error: %s", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		s := new(Scope)
		err = rows.Scan(&s.Name, &s.Label, &s.Description, &s.IsDefault)
		if err != nil {
			log.Printf("rows scan error: %s", err)
		}
		scopes = append(scopes, s)
	}
	return scopes, rows.Err()
}
