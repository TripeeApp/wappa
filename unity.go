package wappa

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

const unityEndpoint endpoint = `unidade`

var unityFields = map[string]string{
	"id": "idUnidade",
	"code": "codDescricao",
}

// Unity is the  struct representing the unity
// entity in the API.
type Unity struct {
	ID		int `json:"IdUniddade,omtempty"`
	Customer	int `json:"IdEmpresaCliente,omitempty"`
	Code		string `json:"Codigo,omitempty"`
	Name		string `json:"Nome,omitempty"`
	Address		string `json:"Endereco,omitempty"`
	Number		string `json:"Numero,omitempty"`
	Complement	string `json:"Complemento,omitempty"`
	Neighborhood	string `json:"Bairro,omitempty"`
	StateID		string `json:"SiglaUf,omitempty"`
	ZipCode		string `json:"Cep,omitempty"`
	AreaCode	string `json:"Ddd,omitempty"`
	Phone		string `json:"Telefone,omitempt"`
	AreaCode2	string `json:"Ddd2,omitempty"`
	Phone2		string `json:"Telefone2,omitempty"`
	CityID		int `json:"IdCidade,omitempty"`
	City		string `json:"Cidade,omitempty"`
}

// UnitResponse is the API response payload.
type UnityResponse struct {
	DefaultResponse

	Response []*Unity
}

// UnityService is responsible for handling
// the requests to the unity resource.
type UnityService struct {
	client requester
}

// Read returns the UnityResponse for the passed filters.
func (us *UnityService) Read(ctx context.Context, f Filter) (*UnityResponse, error) {
	ur := &UnityResponse{}

	if err := us.client.Request(ctx, http.MethodGet, unityEndpoint.Action(read).Query(f.Values(unityFields)), nil, ur); err != nil {
		return nil, err
	}

	return ur, nil
}

// Create creates a unity resource in the API.
func (us *UnityService) Create(ctx context.Context, u *Unity) (*DefaultResponse, error) {
	res := &DefaultResponse{}

	if err := us.client.Request(ctx, http.MethodPost, unityEndpoint.Action(create), u, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Updated edits the unity information.
func (us *UnityService) Update(ctx context.Context, u *Unity) (*DefaultResponse, error) {
	res := &DefaultResponse{}

	if err := us.client.Request(ctx, http.MethodPost, unityEndpoint.Action(update), u, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Inactivate inactivates the unity in the API.
func (us *UnityService) Inactivate(ctx context.Context, id int) (*DefaultResponse, error) {
	res := &DefaultResponse{}
	vals := url.Values{}
	// Converts to string in order to keep the pattern of receiving integer in the Inactivate method.
	vals.Set("idUnidade", strconv.Itoa(id))

	if err := us.client.Request(ctx, http.MethodPost, unityEndpoint.Action(inactivate).Query(vals), nil, res); err != nil {
		return nil, err
	}

	return res, nil
}
