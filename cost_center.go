package wappa

import (
	"context"
	"net/http"
	"net/url"
)

const costCenterEndpoint endpoint = `centrocusto`

// CostCenter is the  struct representing the cost center
// entity in the API.
type CostCenter struct {
	ID			int `json:"IdCentroCusto,omtempty"`
	Name			string `json:"Nome,omitempty"`
	Code			string `json:"Codigo,omitempty"`
	HigherCostCenterID	string `json:"IdCentroCustoSuperior,omitempty"`
	HigherCostCenterCode	string `json:"CodigoCCSuperior,omitempty"`
	HigherCostCenterName	string `json:"NomeCCSuperior,omitempty"`
	CNPJ			string `json:"CNJPEmpresaGrupo,omitempty"`
}

// CostCenterResponse is the API response payload.
type CostCenterResponse struct {
	DefaultResponse

	Response []*CostCenter
}

// CostCenterService is responsible for handling
// the requests to the cost center resource.
type CostCenterService struct {
	client requester
}

// Read returns the CostCenterResponse for the passed filters.
func (cs *CostCenterService) Read(ctx context.Context, id, desc string) (*CostCenterResponse, error) {
	cr := &CostCenterResponse{}
	vals := url.Values{}
	vals.Set("idCentroCusto", id)
	vals.Set("codDescricao", desc)

	if err := cs.client.Request(ctx, http.MethodGet, costCenterEndpoint.Action(read).Query(vals), nil, cr); err != nil {
		return nil, err
	}

	return cr, nil
}

// Create creates a unity resource in the API.
func (cs *CostCenterService) Create(ctx context.Context, cc *CostCenter) (*DefaultResponse, error) {
	res := &DefaultResponse{}

	if err := cs.client.Request(ctx, http.MethodPost, costCenterEndpoint.Action(create), cc, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Updated edits the cost center information.
func (cs *CostCenterService) Update(ctx context.Context, cc *CostCenter) (*DefaultResponse, error) {
	res := &DefaultResponse{}

	if err := cs.client.Request(ctx, http.MethodPost, costCenterEndpoint.Action(update), cc, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Inactivate inactivates the cost center in the API.
func (cs *CostCenterService) Inactivate(ctx context.Context, id int) (*DefaultResponse, error) {
	res := &DefaultResponse{}
	cc := &CostCenter{ID: id}

	if err := cs.client.Request(ctx, http.MethodPost, costCenterEndpoint.Action(inactivate), cc, res); err != nil {
		return nil, err
	}

	return res, nil
}
