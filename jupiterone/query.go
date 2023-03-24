package jupiterone

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jupiterone/jupiterone-client-go/jupiterone/domain"
)

type QueryService service

var ErrNetworkMessage = errors.New("error at network level")

func NetworkError(errm string) error {
	return fmt.Errorf("NetworkError %w : %s", ErrNetworkMessage, errm)
}

// Finished is the status of a query when it has completed
// this is helpful for the consumer to know if their job is
// actually complete or if it just failed.
const (
	Finished   = "FINISHED"
	inProgress = "IN_PROGRESS"
	sleepTime  = 5
)

// Cannot move this to domain until we fix the import cycle issue.
type QueryInput struct {
	Query            string                 `json:"query"`
	Cursor           string                 `json:"cursor"`
	DeferredFormat   DeferredResponseFormat `json:"deferredFormat"`
	DeferredResponse DeferredResponseOption `json:"deferredResponse"`
	DryRun           bool                   `json:"dryRun"`
	Flags            *QueryV1Flags          `json:"flags"`
	IncludeDeleted   bool                   `json:"includeDeleted"`
	Remember         bool                   `json:"remember"`
	Variables        map[string]interface{} `json:"variables"`
}

func (q *QueryService) Query(qi QueryInput) (interface{}, error) {
	var queryResults interface{}

	if qi.Flags == nil {
		qi.Flags = &QueryV1Flags{
			AllPages:           true,
			ComputedProperties: false,
			RowMetadata:        false,
			VariableResultSize: false,
		}
	}

	if qi.DeferredFormat == "" {
		qi.DeferredFormat = DeferredResponseFormatJson
	}

	if qi.DeferredResponse == "" {
		qi.DeferredResponse = DeferredResponseOptionForce
	}

	graphQLResponse, err := QueryJupiterOne(
		context.Background(),
		q.client.gqlClient,
		qi.Query,
		qi.Cursor,
		qi.DeferredFormat,
		qi.DeferredResponse,
		qi.DryRun,
		*qi.Flags,
		qi.IncludeDeleted,
		qi.Remember,
		qi.Variables,
	)
	if err != nil {
		fmt.Println("in query: graphql failure: ", err)
		return queryResults, err
	}

	deferredResponse, err := q.pollDeferredURL(graphQLResponse.QueryV1.Url)
	if err != nil {
		fmt.Println("deferred request failed", err)
		return queryResults, err
	}

	queryResults, err = q.getQueryResults(deferredResponse)

	if err != nil {
		fmt.Println("in query: failure to retrieve results ", err)
		return queryResults, err
	}

	return queryResults, nil
}

func (q *QueryService) AsList(queryResults interface{}) (domain.QueryResult[[]domain.QueryDataVertex], error) {
	var queryResultsList domain.QueryResult[[]domain.QueryDataVertex]

	b, err := json.Marshal(queryResults)
	if err != nil {
		return queryResultsList, err
	}
	err = json.Unmarshal(b, &queryResultsList)
	if err != nil {
		return queryResultsList, err
	}

	return queryResultsList, nil
}

func (q *QueryService) AsTree(queryResults interface{}) (domain.QueryResult[domain.QueryDataTreeResultFormat], error) {
	var queryResultsTree domain.QueryResult[domain.QueryDataTreeResultFormat]

	b, err := json.Marshal(queryResults)
	if err != nil {
		return queryResultsTree, err
	}
	err = json.Unmarshal(b, &queryResultsTree)
	if err != nil {
		return queryResultsTree, err
	}

	return queryResultsTree, nil
}

func (q *QueryService) getQueryResults(d domain.DeferredQueryURLResponse) (interface{}, error) {
	var queryResults interface{}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, d.URL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := q.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, NetworkError(resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&queryResults)
	if err != nil {
		return nil, err
	}

	return queryResults, nil
}

func (q *QueryService) pollDeferredURL(url string) (domain.DeferredQueryURLResponse, error) {
	var deferredResults domain.DeferredQueryURLResponse

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return deferredResults, err
	}

	resp, err := q.client.httpClient.Do(req)
	if err != nil {
		return deferredResults, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return deferredResults, NetworkError(resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&deferredResults)
	if err != nil {
		return deferredResults, err
	}

	if deferredResults.Status == inProgress {
		fmt.Println("deferred results are in progress. sleeping...")
		time.Sleep(sleepTime * time.Second)
		return q.pollDeferredURL(url)
	}

	return deferredResults, nil
}
