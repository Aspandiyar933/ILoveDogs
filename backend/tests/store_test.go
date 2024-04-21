package tests

import "github.com/Aspandiyar933/Ilovedogs/typeslink"

type MockStore struct{}

func (s *MockStore) CreateUser() error {
	return nil
}

func (s *MockStore) CreatePost(p *typeslink.Post) (*typeslink.Post, error) {
	return &typeslink.Post{}, nil
}

func (s *MockStore) GetPost(id string) (*typeslink.Post, error) {
	return &typeslink.Post{}, nil
}