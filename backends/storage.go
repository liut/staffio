package backends

import (
	"errors"
	"github.com/RangelReale/osin"
	"log"
)

type DbStorage struct {
	clients   map[string]osin.Client
	authorize map[string]*osin.AuthorizeData
	access    map[string]*osin.AccessData
	refresh   map[string]string
}

func NewStorage() *DbStorage {
	r := &DbStorage{
		clients:   make(map[string]osin.Client),
		authorize: make(map[string]*osin.AuthorizeData),
		access:    make(map[string]*osin.AccessData),
		refresh:   make(map[string]string),
	}

	// testing
	r.clients["1234"] = &osin.DefaultClient{
		Id:          "1234",
		Secret:      "aabbccdd",
		RedirectUri: "http://localhost:3000/appauth",
	}

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
	if c, ok := s.clients[id]; ok {
		return c, nil
	}
	return nil, errors.New("Client not found")
}

func (s *DbStorage) SetClient(id string, client osin.Client) error {
	log.Printf("SetClient: %s\n", id)
	s.clients[id] = client
	return nil
}

func (s *DbStorage) SaveAuthorize(data *osin.AuthorizeData) error {
	log.Printf("SaveAuthorize: %s\n", data.Code)
	s.authorize[data.Code] = data
	return nil
}

func (s *DbStorage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	log.Printf("LoadAuthorize: %s\n", code)
	if d, ok := s.authorize[code]; ok {
		return d, nil
	}
	return nil, errors.New("Authorize not found")
}

func (s *DbStorage) RemoveAuthorize(code string) error {
	log.Printf("RemoveAuthorize: %s\n", code)
	delete(s.authorize, code)
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
