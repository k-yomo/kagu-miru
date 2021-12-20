package queryclassifier

import (
	"context"
	"fmt"

	automl "cloud.google.com/go/automl/apiv1"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	automlpb "google.golang.org/genproto/googleapis/cloud/automl/v1"
)

const scoreThreshold = 0.6

type QueryClassifier interface {
	CategorizeQuery(ctx context.Context, query string) ([]string, error)
}

type queryClassifierClient struct {
	predictionClient                *automl.PredictionClient
	categoryClassificationModelPath string
}

func NewQueryClassifierClient(
	predictionClient *automl.PredictionClient,
	gcpProjectName string,
	categoryClassificationModelName string,
) *queryClassifierClient {
	return &queryClassifierClient{
		predictionClient:                predictionClient,
		categoryClassificationModelPath: fmt.Sprintf("projects/%s/locations/us-central1/models/%s", gcpProjectName, categoryClassificationModelName),
	}
}

// CategorizeQuery predicts the query's intended categories and returns the list of ids.
func (q *queryClassifierClient) CategorizeQuery(ctx context.Context, query string) ([]string, error) {
	ctx, span := otel.Tracer("").Start(ctx, "queryclassifier.queryClassifierClient_CategorizeQuery")
	defer span.End()

	resp, err := q.predictionClient.Predict(ctx, &automlpb.PredictRequest{
		Name: q.categoryClassificationModelPath,
		Payload: &automlpb.ExamplePayload{
			Payload: &automlpb.ExamplePayload_TextSnippet{
				TextSnippet: &automlpb.TextSnippet{
					Content:  query,
					MimeType: "text/plain",
				}},
		},
	})
	if err != nil {
		return nil, err
	}

	var categoryIDs []string
	for _, payload := range resp.Payload {
		detail, ok := payload.Detail.(*automlpb.AnnotationPayload_Classification)
		if ok && detail.Classification.Score >= scoreThreshold {
			categoryIDs = append(categoryIDs, payload.DisplayName)
		}
	}

	span.SetAttributes(attribute.StringSlice("categoryIds", categoryIDs))
	return categoryIDs, nil
}
