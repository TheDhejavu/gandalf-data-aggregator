package webapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gandalf-data-aggregator/config"
	"gandalf-data-aggregator/pkg/crypto"

	"github.com/machinebox/graphql"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type GandalfClient struct {
	cfg    config.Config
	client *graphql.Client
}

func NewGandalfClient(cfg config.Config) *GandalfClient {
	client := graphql.NewClient(cfg.Gandalf.SauronURL)
	return &GandalfClient{cfg, client}
}

type GetActivityResponse[T NetflixActivityMetadata] struct {
	GetActivity ActivityResponse[T] `json:"getActivity"`
}

type ActivityResponse[T NetflixActivityMetadata] struct {
	Data  []*Activity[T] `json:"data"`
	Limit int64          `json:"limit"`
	Total int64          `json:"total"`
	Page  int64          `json:"page"`
}

type Activity[T NetflixActivityMetadata] struct {
	ID       string `json:"id"`
	Metadata T      `json:"metadata"`
}

type NetflixActivityMetadata struct {
	Title   string        `json:"title"`
	Subject []*Identifier `json:"subject,omitempty"`
	Date    string        `json:"date,omitempty"`
}

type Identifier struct {
	Value          string `json:"value"`
	IdentifierType string `json:"identifierType"`
}

func (g GandalfClient) QueryActivities(ctx context.Context, dataKey string, limit int, page int) (*ActivityResponse[NetflixActivityMetadata], error) {
	privateKey, err := crypto.HexToECDSAPrivateKey(g.cfg.Gandalf.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	req := graphql.NewRequest(`
	query getActivity($dataKey: String!, $source: Source!, $limit: Int64!, $page: Int64!) {
		getActivity(dataKey: $dataKey, source: $source, limit: $limit, page: $page) {
		  data {
			id
			metadata {
			  ...NetflixActivityMetadata
			}
		  }
		  limit
		  page
		  total
		}
	  }
	  fragment NetflixActivityMetadata on NetflixActivityMetadata {
		title
		subject {
		  value
		  identifierType
		}
		date
	  }
    `)
	req.Var("dataKey", dataKey)
	req.Var("source", "NETFLIX")
	req.Var("limit", limit)
	req.Var("page", page)

	requestBodyObj := struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables"`
	}{
		Query:     req.Query(),
		Variables: req.Vars(),
	}

	log.Info().Msgf("G_DataKey", dataKey)

	var requestBody bytes.Buffer
	if err := json.NewEncoder(&requestBody).Encode(requestBodyObj); err != nil {
		return nil, errors.Wrap(err, "encode body")
	}

	signatureB64 := crypto.SignMessageAsBase64(privateKey, requestBody.Bytes())
	req.Header.Set("X-Gandalf-Signature", signatureB64)

	var respData GetActivityResponse[NetflixActivityMetadata]
	if err := g.client.Run(ctx, req, &respData); err != nil {
		return nil, fmt.Errorf("failed to execute request: %v", err)
	}

	return &respData.GetActivity, nil
}
