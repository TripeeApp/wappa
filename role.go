package wappa

import (
	"context"
	"net/http"
)

const roleEndpoint endpoint = `cargo`

var roleFields = map[string]string{
	"id": "idCargo",
	"desc": "descricao",
}

// Role is the  struct representing the role
// entity in the API.
type Role struct {
	ID		int `json:"IdCargo,omtempty"`
	Description	int `json:"Descricao,omitempty"`
	Filter		string `json:"Filtro,omitempty"`
}

// RoleResponse is the API reading response payload for role entity.
type RoleResponse struct {
	DefaultResponse

	Response []*Role
}

// RoleService is responsible for handling
// the requests to the role resource.
type RoleService struct {
	client requester
}

// Read returns the RoleResponse for the passed filters.
func (rs *RoleService) Read(ctx context.Context, f Filter) (*RoleResponse, error) {
	r := &RoleResponse{}

	if err := rs.client.Request(ctx, http.MethodGet, roleEndpoint.Action(read).Query(f.Values(roleFields)), nil, r); err != nil {
		return nil, err
	}

	return r, nil
}

// Create creates a role resource in the API.
func (rs *RoleService) Create(ctx context.Context, r *Role) (*DefaultResponse, error) {
	res := &DefaultResponse{}

	if err := rs.client.Request(ctx, http.MethodPost, roleEndpoint.Action(create), r, res); err != nil {
		return nil, err
	}

	return res, nil
}
