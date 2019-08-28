package wappa

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"
)

func TestQuote(t *testing.T) {
	testCases := []testTable{
		{
			"Estimate()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&QuoteService{req}).Estimate(ctx, Filter{"latOrigin": []string{"3.14"}})
				return
			},
			context.Background(),
			http.MethodGet,
			quoteEndpoint.Query(url.Values{"LatitudeOrigin": []string{"3.14"}}),
			nil,
			&QuoteResult{
				Categories: []*Category{
					&Category{ID: 1},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, test(tc))
	}
}

func TestQuoteError(t *testing.T) {
	testCases := []testTableError{
		{
			"Estimate()",
			func(req requester) error {
				_, err := (&QuoteService{req}).Estimate(context.Background(), nil)
				return err
			},
			errors.New("Error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, testError(tc))
	}
}
