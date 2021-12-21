package queryclassifier

import (
	"context"
	"fmt"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	aiplatformpb "google.golang.org/genproto/googleapis/cloud/aiplatform/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

const scoreThreshold = 0.6

type QueryClassifier interface {
	CategorizeQuery(ctx context.Context, query string) ([]string, error)
}

type queryClassifierClient struct {
	predictionClient               *aiplatform.PredictionClient
	categoryClassificationEndpoint string
}

func NewQueryClassifierClient(
	predictionClient *aiplatform.PredictionClient,
	gcpProjectName string,
	categoryClassificationEndpointID string,
) *queryClassifierClient {
	return &queryClassifierClient{
		predictionClient:               predictionClient,
		categoryClassificationEndpoint: fmt.Sprintf("projects/%s/locations/us-central1/endpoints/%s", gcpProjectName, categoryClassificationEndpointID),
	}
}

// CategorizeQuery predicts the query's intended categories and returns the list of ids.
func (q *queryClassifierClient) CategorizeQuery(ctx context.Context, query string) ([]string, error) {
	ctx, span := otel.Tracer("").Start(ctx, "queryclassifier.queryClassifierClient_CategorizeQuery")
	defer span.End()

	predictionInstance, err := structpb.NewStruct(map[string]interface{}{
		"content":  query,
		"mimeType": "text/plain",
	})
	resp, err := q.predictionClient.Predict(ctx, &aiplatformpb.PredictRequest{
		Endpoint: q.categoryClassificationEndpoint,
		Instances: []*structpb.Value{
			structpb.NewStructValue(predictionInstance),
		},
	})
	if err != nil {
		return nil, err
	}

	predictions := resp.GetPredictions()
	if len(predictions) == 0 {
		return nil, nil
	}

	var categoryIDs []string
	fields := predictions[0].GetStructValue().GetFields()
	displayNameValues := fields["displayNames"].GetListValue().Values
	confidenceValues := fields["confidences"].GetListValue().Values
	for i, categoryIDValue := range displayNameValues {
		categoryID := categoryIDValue.GetStringValue()
		if confidenceValues[i].GetNumberValue() >= scoreThreshold {
			categoryIDs = append(categoryIDs, categoryID)
		}

	}

	span.SetAttributes(attribute.StringSlice("categoryIds", categoryIDs))
	return categoryIDs, nil
}
