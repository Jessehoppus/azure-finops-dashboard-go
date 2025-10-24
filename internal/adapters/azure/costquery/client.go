package costquery

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/costmanagement/armcostmanagement"
)

// Client implementa o serviço de custos usando a API Cost Management Query.
type Client struct {
	query *armcostmanagement.QueryClient
}

func NewClient(ctx context.Context) (*Client, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}
	factory, err := armcostmanagement.NewClientFactory("", cred, nil)
	if err != nil {
		return nil, err
	}
	return &Client{query: factory.NewQueryClient()}, nil
}

// CostByDimension executa uma consulta agregada por dimensão.
func (c *Client) CostByDimension(ctx context.Context, scope string, from, to time.Time, dimension, granularity string) ([][]any, []string, error) {
	var gran armcostmanagement.GranularityType
	switch strings.ToLower(granularity) {
	case "daily":
		gran = armcostmanagement.GranularityTypeDaily
	case "monthly":
		gran = armcostmanagement.GranularityTypeMonthly
	default:
		gran = armcostmanagement.GranularityTypeNone
	}

	grpName := dimension
	grpType := armcostmanagement.QueryColumnTypeDimension
	if strings.HasPrefix(strings.ToLower(dimension), "tagkey:") {
		key := dimension[len("TagKey:"):]
		if key == "" {
			return nil, nil, fmt.Errorf("TagKey vazio")
		}
		grpName = fmt.Sprintf("TagKey:%s", key)
		grpType = armcostmanagement.QueryColumnTypeTag
	}

	q := armcostmanagement.QueryDefinition{
		Type:      toPtr(armcostmanagement.ExportTypeUsage),
		Timeframe: toPtr(armcostmanagement.TimeframeTypeCustom),
		TimePeriod: &armcostmanagement.QueryTimePeriod{
			From: toPtr(from),
			To:   toPtr(to),
		},
		Dataset: &armcostmanagement.QueryDataset{
			Granularity: toPtr(gran),
			Aggregation: map[string]*armcostmanagement.QueryAggregation{
				"totalCost": {Name: toPtr("PreTaxCost"), Function: toPtr(armcostmanagement.FunctionTypeSum)},
			},
			Grouping: []*armcostmanagement.QueryGrouping{
				{Type: toPtr(grpType), Name: toPtr(grpName)},
			},
		},
	}

	res, err := c.query.Usage(ctx, scope, q, nil)
	if err != nil {
		return nil, nil, err
	}

	cols := res.Properties.Columns
	rows := res.Properties.Rows
	headers := make([]string, len(cols))
	for i, c := range cols {
		headers[i] = *c.Name
	}
	return rows, headers, nil
}

func toPtr[T any](v T) *T { return &v }
