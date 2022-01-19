package cms

import (
	"context"
	"fmt"

	"github.com/k-yomo/kagu-miru/backend/pkg/logging"
	sanity "github.com/sanity-io/client-go"
)

type Client interface {
	GetFeaturedPosts(ctx context.Context) (*GetFeaturedPostsResponse, error)
}

type cmsClient struct {
	sanityClient *sanity.Client
}

func NewCMSClient(sanityClient *sanity.Client) *cmsClient {
	return &cmsClient{
		sanityClient: sanityClient,
	}
}

type GetFeaturedPostsResponse struct {
	Title string  `json:"title"`
	Posts []*Post `json:"posts"`
}

func (c *cmsClient) GetFeaturedPosts(ctx context.Context) (*GetFeaturedPostsResponse, error) {
	query := `*[_type == "postsGroup" && id == "featuredPosts"][0]{
  title,
  posts[]->{
    "slug": slug.current,
    title,
    description,
    mainImage,
    publishedAt,
    categories,
  }
}`
	result, err := c.sanityClient.Query(query).Do(ctx)
	if err != nil {
		return nil, logging.Error(ctx, fmt.Errorf("failed to query to sanity: %w", err))
	}

	var resp GetFeaturedPostsResponse
	if err := result.Unmarshal(&resp); err != nil {
		return nil, logging.Error(ctx, fmt.Errorf("failed to query to sanity: %w", err))
	}

	return &resp, nil
}
