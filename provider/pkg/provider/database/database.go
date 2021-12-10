package database

import (
	"context"
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

func CreateDatabase(ctx context.Context, req *pulumirpc.CreateRequest) (*pulumirpc.CreateResponse, error) {
	inputs, err := plugin.UnmarshalProperties(req.GetProperties(), plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}

	if !inputs["name"].IsString() {
		return nil, fmt.Errorf("expected input property 'name' of type 'string' but got '%s", inputs["name"].TypeString())
	}

	if !inputs["region"].IsString() {
		return nil, fmt.Errorf("expected input property 'region' of type 'string' but got '%s", inputs["name"].TypeString())
	}

	databaseName := inputs["name"].StringValue()
	databaseRegion := inputs["region"].StringValue()

	// Actually "create" the random number
	result := databaseName + databaseRegion

	outputs := map[string]interface{}{}

	outputProperties, err := plugin.MarshalProperties(
		resource.NewPropertyMapFromMap(outputs),
		plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true},
	)
	if err != nil {
		return nil, err
	}

	return &pulumirpc.CreateResponse{
		Id:         result,
		Properties: outputProperties,
	}, nil
}
