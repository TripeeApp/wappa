package wappa

import (
	"context"
	"net/http"
)

const costCenterEndpoint endpoint = `centrocusto`

var costCenterFields = map[string]string{
	"id": "idCentroCusto",
	"code": "codDescricao",
}

// CostCenter is the  struct representing the cost center
// entity in the API.
type CostCenter struct {
	ID			int `json:"IdCentroCusto,omtempty"`
	Name			string `json:"Nome,omitempty"`
	Code			string `json:"Codigo,omitempty"`
	ParentCostCenterID	int `json:"IdCentroCustoSuperior,omitempty"`
	ParentCostCenterCode	string `json:"CodigoCCSuperior,omitempty"`
	ParentCostCenterName	string `json:"NomeCCSuperior,omitempty"`
	CNPJ			string `json:"CNPJEmpresaGrupo,omitempty"`
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
func (cs *CostCenterService) Read(ctx context.Context, f Filter) (*CostCenterResponse, error) {
	cr := &CostCenterResponse{}

	if err := cs.client.Request(ctx, http.MethodGet, costCenterEndpoint.Action(read).Query(f.Values(costCenterFields)), nil, cr); err != nil {
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
func (cs *CostCenterService) Update(ctx context.Context, cc *CostCenter) (*OperationDefaultResponse, error) {
	res := &OperationDefaultResponse{}

	if err := cs.client.Request(ctx, http.MethodPost, costCenterEndpoint.Action(update), cc, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Inactivate inactivates the cost center in the API.
func (cs *CostCenterService) Inactivate(ctx context.Context, id int) (*OperationDefaultResponse, error) {
	res := &OperationDefaultResponse{}
	cc := &CostCenter{ID: id}

	if err := cs.client.Request(ctx, http.MethodPost, costCenterEndpoint.Action(inactivate), cc, res); err != nil {
		return nil, err
	}

	return res, nil
}
