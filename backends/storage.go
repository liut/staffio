package backends

import (
	"database/sql"
	"errors"
	"github.com/RangelReale/osin"
	_ "github.com/lib/pq"
	"log"
	"tuluu.com/liut/staffio/models"
)

type DbStorage struct {
	// clients   map[string]osin.Client
	// authorize map[string]*osin.AuthorizeData
	access  map[string]*osin.AccessData
	refresh map[string]string
}

func NewStorage() *DbStorage {

	r := &DbStorage{
		// clients:   make(map[string]osin.Client),
		// authorize: make(map[string]*osin.AuthorizeData),
		access:  make(map[string]*osin.AccessData),
		refresh: make(map[string]string),
	}

	// testing
	// r.clients["1234"] = &osin.DefaultClient{
	// 	Id:          "1234",
	// 	Secret:      "aabbccdd",
	// 	RedirectUri: "http://localhost:3000/appauth",
	// }

	// log.Printf("clients: %v", r.clients)
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
	log.Printf("SaveAuthorize: %s\n", data.Code)
	// s.authorize[data.Code] = data
	qs := func(tx *sql.Tx) error {
		sql := `INSERT INTO
		 oauth_authorization_code(code, client_id, username, redirect_uri, expires_in, scopes)
		 VALUES($1, $2, $3, $4, $5, $6);`
		result, err := tx.Exec(sql, data.Code, data.Client.GetId(), data.UserData.(string),
			data.RedirectUri, data.ExpiresIn, data.Scope)
		if err != nil {
			log.Printf("save authorizedData error %s", err)
			return errors.New("storage failed")
		}

		n, err := result.RowsAffected()
		if err != nil {
			log.Printf("RowsAffected error %s", err)
		}
		log.Printf("save authorizeData code %s result %v", data.Code, n)

		return nil
	}
	return withTxQuery(qs)
}

func (s *DbStorage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	log.Printf("LoadAuthorize: %s\n", code)
	var (
		client_id string
		username  string
		err       error
	)
	a := new(osin.AuthorizeData)
	err = getDb().QueryRow("SELECT client_id, username, redirect_uri, expires_in, scopes, created FROM oauth_authorization_code WHERE code = $1",
		code).Scan(&client_id, &username, &a.RedirectUri, &a.ExpiresIn, &a.Scope, &a.CreatedAt)
	if err == nil {
		a.Client, err = GetClientWithCode(client_id)
		if err != nil {
			log.Printf("get client error: %s", err)
		}
		a.UserData = username
		log.Printf("loaded authorization ok, createdAt %s", a.CreatedAt)
		return a, nil
	}
	return nil, errors.New("Authorize not found")
}

func (s *DbStorage) RemoveAuthorize(code string) error {
	log.Printf("RemoveAuthorize: %s\n", code)
	qs := func(tx *sql.Tx) error {
		sql := `DELETE FROM oauth_authorization_code WHERE code = $1;`
		result, err := tx.Exec(sql, code)
		if err != nil {
			log.Printf("delete authorizedData error %s", err)
			return err
		}

		n, err := result.RowsAffected()
		if err != nil {
			log.Printf("RowsAffected error %s", err)
		}
		log.Printf("delete authorizeData code %s result %v", code, n)

		return nil
	}
	return withTxQuery(qs)
	// delete(s.authorize, code)
	return nil
}

func (s *DbStorage) SaveAccess(data *osin.AccessData) error {
	log.Printf("SaveAccess: %s\n", data.AccessToken)
	s.access[data.AccessToken] = data
	if data.RefreshToken != "" {
		s.refresh[data.RefreshToken] = data.AccessToken
	}
	return nil
}

func (s *DbStorage) LoadAccess(code string) (*osin.AccessData, error) {
	log.Printf("LoadAccess: %s\n", code)
	if d, ok := s.access[code]; ok {
		return d, nil
	}
	return nil, errors.New("Access not found")
}

func (s *DbStorage) RemoveAccess(code string) error {
	log.Printf("RemoveAccess: %s\n", code)
	delete(s.access, code)
	return nil
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
	return nil, err
}
