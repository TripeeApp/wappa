package wappa

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

const collaboratorEndpoint endpoint = `colaborador`

var collaboratorFields = map[string]string{
	"name": "nome",
	"enrollment": "matricula",
	"status": "status",
	"unity": "idUnidade",
	"CostCenterID": "idCentroCusto",
	"AdmID": "idAdministrador",
	"email": "email",
	"cpf": "cpf",
}


// Collaborator is the  struct representing the collaborator
// entity in the API.
type Collaborator struct {
	ID			 int     `json:"IdColaborador,omitempty"`
	Company			 int     `json:"Idempresa,omitempty"`
	Association              int     `json:"IdAssociacao,omitempty"`
	CostCenterID             int     `json:"IdCentroCusto,omitempty"`
	UnityID                  int     `json:"IdUnidade,omitempty"`
	RoleID                   int     `json:"IdCargo,omitempty"`
	Name                     string  `json:"Nome,omitempty"`
	CPF                      string  `json:"Cpf,omitempty"`
	Email                    string  `json:"Email,omitempty"`
	Role                     string  `json:"Cargo,omitempty"`
	Enrollment               string  `json:"Matricula,omitempty"`
	Active			 bool    `json:"StatusAtivo,omitempty"`
	Blocked			 bool    `json:"StatusBloqueado,omitempty"`
	MonthlyTaxiFare          float64 `json:"ValorMensalTaxi,omitempty"`
	MonthlyLimit             string  `json:"LimiteMensal,omitempty"`
	UnlimitedTaxi             bool    `json:"FlgTaxiIlimitado,omitempty"`
	ExtraMonthlyTaxiFare     float64 `json:"ValorAcrescimoMensalTaxi,omitempty"`
	AccumulatedTaxiFare      float64 `json:"ValorAcumuladoTaxi,omitempty"`
	AuthorizedFleet          bool    `json:"FlgFrotaAutorizado,omitempty"`
	SendUnblockSMS           bool    `json:"FlgEnviarSmsDesbloqueio,omitempty"`
	CostCenterCode           string  `json:"CodigoCentroCusto,omitempty"`
	CostCenterName           string  `json:"NomeCentroCusto,omitempty"`
	UnityCode                string  `json:"CodigoUnidade,omitempty"`
	UnityName                string  `json:"NomeUnidade,omitempty"`
	ChangedStatusAt          Time  `json:"DataStatus,omitempty"`
	CreatedAt                Time  `json:"DataCadastro,omitempty"`
	ActivatedAt              Time  `json:"DataAtivacao,omitempty"`
	ReactivatedAt            Time  `json:"DataReativacao,omitempty"`
	InactivatedAt            Time  `json:"DataDesativacao,omitempty"`
	BlockedAt                Time  `json:"DataBloqueio,omitempty"`
	PasswordResentAt         Time  `json:"DataReenvioSenha,omitempty"`
	Type                     int     `json:"TipoColaborador,omitempty"`
	Foreign                  bool    `json:"Estrangeiro,omitempty"`
	CardExpiration           string  `json:"ValidadeCartao,omitempty"`
	CardFinalDigits          string  `json:"FinalCartao,omitempty"`
	CountryCode              string  `json:"CodPaisNF,omitempty"`
	CountryName              string  `json:"NomePaisNF,omitempty"`
	UnlimitedRides           bool    `json:"QtdCorridasIlimitado,omitempty"`
	AnswerToID               int     `json:"IdRespondePara,omitempty"`
	AnswerToName             string  `json:"NomeRespondePara,omitempty"`
	AnwserToEnrollment       string  `json:"MatriculaRespondePara,omitempty"`
	AppVersion               string  `json:"VersaoApp,omitempty"`
	UserID                   int     `json:"IdUsuario,omitempty"`
	DDD                      string  `json:"Ddd,omitempty"`
	Login                    string  `json:"Login,omitempty"`
}

// UnitResponse is the API response payload.
type CollaboratorResponse struct {
	DefaultResponse

	Response []*Collaborator
}

// CollaboratorService is responsible for handling
// the requests to the collaborator resource.
type CollaboratorService struct {
	client requester
}

// Read returns the CollaboratorResponse for the passed filters.
func (cs *CollaboratorService) Read(ctx context.Context, f Filter) (*CollaboratorResponse, error) {
	cr := &CollaboratorResponse{}

	if err := cs.client.Request(ctx, http.MethodGet, collaboratorEndpoint.Action(read).Query(f.Values(collaboratorFields)), nil, cr); err != nil {
		return nil, err
	}

	return cr, nil
}

// Create creates a collaborator resource in the API.
func (cs *CollaboratorService) Create(ctx context.Context, c *Collaborator) (*DefaultResponse, error) {
	res := &DefaultResponse{}

	if err := cs.client.Request(ctx, http.MethodPost, collaboratorEndpoint.Action(create), c, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Updated edits the collaborator information.
func (cs *CollaboratorService) Update(ctx context.Context, c *Collaborator) (*OperationDefaultResponse, error) {
	res := &OperationDefaultResponse{}

	if err := cs.client.Request(ctx, http.MethodPost, collaboratorEndpoint.Action(update), c, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Inactivate inactivates the collaborator in the API.
func (cs *CollaboratorService) Inactivate(ctx context.Context, id int) (*OperationDefaultResponse, error) {
	res := &OperationDefaultResponse{}
	// Converts to string in order to keep the pattern of receiving integer in the Inactivate method.
	vals := url.Values{}
	vals.Set("idColaborador", strconv.Itoa(id))

	c := &Collaborator{ID: id}
	if err := cs.client.Request(ctx, http.MethodPost, collaboratorEndpoint.Action(inactivate).Query(vals), c, res); err != nil {
		return nil, err
	}

	return res, nil
}
