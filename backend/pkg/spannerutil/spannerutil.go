package spannerutil

import "fmt"

func BuildSpannerDBPath(projectID, instanceID, databaseID string) string {
	return fmt.Sprintf("projects/%s/instances/%s/databases/%s", projectID, instanceID, databaseID)
}
