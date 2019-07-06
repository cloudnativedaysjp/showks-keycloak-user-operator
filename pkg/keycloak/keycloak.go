package keycloak

import (
	"github.com/Nerzal/gocloak"
)

type KeyCloakClientInterface interface {
	GetUsers(realm string, param gocloak.GetUsersParams) (*[]gocloak.User, error)
	GetUserByID(realm string, id string) (*gocloak.User, error)
	CreateUser(realm string, user gocloak.User) (string, error)
	DeleteUser(realm string, id string) error
}

func NewClient(basePath string, username string, password string, realm string) (KeyCloakClientInterface, error) {
	client := gocloak.NewClient(basePath)
	token, err := client.LoginAdmin(username, password, realm)
	if err != nil {
		return nil, err
	}
	return &KeyCloak{
		client: client,
		token:  token,
	}, nil
}

type KeyCloak struct {
	client gocloak.GoCloak
	token  *gocloak.JWT
}

func (c *KeyCloak) GetUsers(realm string, param gocloak.GetUsersParams) (*[]gocloak.User, error) {
	users, err := c.client.GetUsers(c.token.AccessToken, realm, param)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (c *KeyCloak) GetUserByID(realm string, id string) (*gocloak.User, error) {
	user, err := c.client.GetUserByID(c.token.AccessToken, realm, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (c *KeyCloak) CreateUser(realm string, user gocloak.User) (string, error) {
	id, err := c.client.CreateUser(c.token.AccessToken, realm, user)
	if err != nil {
		return "", err
	}
	return *id, nil
}

func (c *KeyCloak) DeleteUser(realm string, id string) error {
	err := c.client.DeleteUser(c.token.AccessToken, realm, id)
	if err != nil {
		return err
	}
	return nil
}
