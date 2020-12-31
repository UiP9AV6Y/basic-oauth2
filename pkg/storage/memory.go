package storage

import (
	"github.com/openshift/osin"
)

type ClientsMemory struct {
	clients   map[string]osin.Client
	authorize map[string]*osin.AuthorizeData
	access    map[string]*osin.AccessData
	refresh   map[string]string
}

func NewClientsMemory() *ClientsMemory {
	r := &ClientsMemory{
		clients:   make(map[string]osin.Client),
		authorize: make(map[string]*osin.AuthorizeData),
		access:    make(map[string]*osin.AccessData),
		refresh:   make(map[string]string),
	}

	return r
}

func (s *ClientsMemory) AddClient(c osin.Client) bool {
	_, ok := s.clients[c.GetId()]
	if ok {
		return false
	}

	s.clients[c.GetId()] = c
	return true
}

func (s *ClientsMemory) Clone() osin.Storage {
	return s
}

func (s *ClientsMemory) Close() {
}

func (s *ClientsMemory) GetClient(id string) (osin.Client, error) {
	if c, ok := s.clients[id]; ok {
		return c, nil
	}
	return nil, osin.ErrNotFound
}

func (s *ClientsMemory) SetClient(id string, client osin.Client) error {
	s.clients[id] = client
	return nil
}

func (s *ClientsMemory) SaveAuthorize(data *osin.AuthorizeData) error {
	s.authorize[data.Code] = data
	return nil
}

func (s *ClientsMemory) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	if d, ok := s.authorize[code]; ok {
		return d, nil
	}
	return nil, osin.ErrNotFound
}

func (s *ClientsMemory) RemoveAuthorize(code string) error {
	delete(s.authorize, code)
	return nil
}

func (s *ClientsMemory) SaveAccess(data *osin.AccessData) error {
	s.access[data.AccessToken] = data
	if data.RefreshToken != "" {
		s.refresh[data.RefreshToken] = data.AccessToken
	}
	return nil
}

func (s *ClientsMemory) LoadAccess(code string) (*osin.AccessData, error) {
	if d, ok := s.access[code]; ok {
		return d, nil
	}
	return nil, osin.ErrNotFound
}

func (s *ClientsMemory) RemoveAccess(code string) error {
	delete(s.access, code)
	return nil
}

func (s *ClientsMemory) LoadRefresh(code string) (*osin.AccessData, error) {
	if d, ok := s.refresh[code]; ok {
		return s.LoadAccess(d)
	}
	return nil, osin.ErrNotFound
}

func (s *ClientsMemory) RemoveRefresh(code string) error {
	delete(s.refresh, code)
	return nil
}
