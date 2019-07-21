package keycloak

import (
	"github.com/Nerzal/gocloak"
)

type KeyCloakClientInterface interface {
	GetUsers(realm string, param gocloak.GetUsersParams) (*[]gocloak.User, error)
	GetUserByID(realm string, id string) (*gocloak.User, error)
	CreateUser(realm string, user gocloak.User) (string, error)
	DeleteUser(realm string, id string) error
	SetPassword(realm string, id string, password string) error
}

func NewClient(basePath string, username string, password string, realm string) (KeyCloakClientInterface, error) {
	client := gocloak.NewClient(basePath)
	token, err := client.LoginAdmin(username, password, realm)
	if err != nil {
		return nil, err
	}
	return &KeyCloak{
		client:   client,
		username: username,
		password: password,
		token:    token,
	}, nil
}

type KeyCloak struct {
	client   gocloak.GoCloak
	username string
	password string
	token    *gocloak.JWT
}

func (c *KeyCloak) GetUsers(realm string, param gocloak.GetUsersParams) (*[]gocloak.User, error) {
	token, err := c.client.LoginAdmin(c.username, c.password, realm)
	if err != nil {
		return nil, err
	}
	users, err := c.client.GetUsers(token.AccessToken, realm, param)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (c *KeyCloak) GetUserByID(realm string, id string) (*gocloak.User, error) {
	token, err := c.client.LoginAdmin(c.username, c.password, realm)
	if err != nil {
		return nil, err
	}
	user, err := c.client.GetUserByID(token.AccessToken, realm, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (c *KeyCloak) CreateUser(realm string, user gocloak.User) (string, error) {
	token, err := c.client.LoginAdmin(c.username, c.password, realm)
	if err != nil {
		return "", err
	}
	id, err := c.client.CreateUser(token.AccessToken, realm, user)
	if err != nil {
		return "", err
	}
	return *id, nil
}

func (c *KeyCloak) DeleteUser(realm string, id string) error {
	token, err := c.client.LoginAdmin(c.username, c.password, realm)
	if err != nil {
		return err
	}
	err = c.client.DeleteUser(token.AccessToken, realm, id)
	if err != nil {
		return err
	}
	return nil
}

func (c *KeyCloak) SetPassword(realm string, id string, password string) error {
	token, err := c.client.LoginAdmin(c.username, c.password, realm)
	if err != nil {
		return err
	}
	err = c.client.SetPassword(token.AccessToken, id, realm, password, false)
	if err != nil {
		return err
	}

	return nil
}
