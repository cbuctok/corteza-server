package rbac

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/crusttech/crust/internal/rbac/types"
)

type (
	Users struct {
		*Client
	}

	UsersInterface interface {
		Create(username, password string) (*types.User, error)
		Get(userID string) (*types.User, error)
		Delete(userID string) error

		AddRole(userID string, roles ...string) error
		RemoveRole(userID string, roles ...string) error
	}
)

const (
	usersCreate = "/users/"
	usersGet    = "/users/%s"
	usersDelete = "/users/%s"
	// @todo: plural for users, but singular for sessions
	usersAddRole    = "/users/%s/assignRoles"
	usersRemoveRole = "/users/%s/deassignRoles"
)

func (u *Users) Create(username, password string) (*types.User, error) {
	body := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{username, password}

	resp, err := u.Client.Post(usersCreate, body)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
		user := &types.User{}
		return user, errors.Wrap(json.NewDecoder(resp.Body).Decode(user), "decoding json failed")
	default:
		return nil, toError(resp)
	}
}

func (u *Users) AddRole(userID string, roles ...string) error {
	body := struct {
		Roles []string `json:"roles"`
	}{roles}

	resp, err := u.Client.Patch(fmt.Sprintf(usersAddRole, userID), body)
	if err != nil {
		return errors.Wrap(err, "request failed")
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
		return nil
	default:
		return toError(resp)
	}
}

func (u *Users) RemoveRole(userID string, roles ...string) error {
	body := struct {
		Roles []string `json:"roles"`
	}{roles}

	resp, err := u.Client.Patch(fmt.Sprintf(usersRemoveRole, userID), body)
	if err != nil {
		return errors.Wrap(err, "request failed")
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
		return nil
	default:
		return toError(resp)
	}
}

func (u *Users) Get(userID string) (*types.User, error) {
	resp, err := u.Client.Get(fmt.Sprintf(usersGet, userID))
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
		user := &types.User{}
		return user, errors.Wrap(json.NewDecoder(resp.Body).Decode(user), "decoding json failed")
	default:
		return nil, toError(resp)
	}
}

func (u *Users) Delete(userID string) error {
	resp, err := u.Client.Delete(fmt.Sprintf(usersDelete, userID))
	if err != nil {
		return errors.Wrap(err, "request failed")
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
		return nil
	default:
		return toError(resp)
	}
}

var _ UsersInterface = &Users{}
