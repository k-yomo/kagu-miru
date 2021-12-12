package rakutenichiba

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/k-yomo/kagu-miru/backend/pkg/httputil"
	"github.com/k-yomo/kagu-miru/backend/pkg/urlutil"
)

type Genre struct {
	ID        int    `json:"genreId"`
	Name      string `json:"genreName"`
	Level     int    `json:"genreLevel"`
	Parent    *Genre
	Children  []*Genre
	TagGroups []*struct {
		TagGroup *TagGroup `json:"tagGroup"`
	} `json:"tagGroups"`
}

func (g *Genre) GenreIDs() []string {
	genreIDs := []string{strconv.Itoa(g.ID)}
	parentGenre := g.Parent
	for parentGenre != nil {
		genreIDs = append([]string{strconv.Itoa(parentGenre.ID)}, genreIDs...)
		parentGenre = parentGenre.Parent
	}
	return genreIDs
}

func (g *Genre) GenreNames() []string {
	genreNames := []string{g.Name}
	parentGenre := g.Parent
	for parentGenre != nil {
		genreNames = append([]string{parentGenre.Name}, genreNames...)
		parentGenre = parentGenre.Parent
	}
	return genreNames
}

const GenreFurnitureID = "100804"

type SearchGenreResponse struct {
	Parents  []*Genre `json:"parents"`
	Current  *Genre   `json:"current"`
	Children []struct {
		Child *struct {
			ID    int    `json:"genreId"`
			Name  string `json:"genreName"`
			Level int    `json:"genreLevel"`
		} `json:"child"`
	} `json:"children"`
	TagGroups []*struct {
		TagGroup *TagGroup `json:"tagGroup"`
	} `json:"tagGroups"`
}

// SearchGenre searches parent, current and children genre of given ID
// https://webservice.rakuten.co.jp/api/ichibagenresearch/
func (c *Client) SearchGenre(ctx context.Context, genreID string) (*SearchGenreResponse, error) {
	u := urlutil.CopyWithQueries(c.genreSearchAPIURL, c.buildParams(map[string]string{"genreId": genreID}))
	var resp SearchGenreResponse
	if err := httputil.GetAndUnmarshal(ctx, c.httpClient, u, &resp); err != nil {
		return nil, fmt.Errorf("httputil.GetAndUnmarshal: %w", err)
	}
	resp.Current.TagGroups = resp.TagGroups
	return &resp, nil
}

// GetGenreWithAllChildren gets genre with all lower hierarchy genre
func (c *Client) GetGenreWithAllChildren(ctx context.Context, genreID string) (*Genre, error) {
	resp, err := c.SearchGenre(ctx, genreID)
	if err != nil {
		return nil, fmt.Errorf("httputil.GetAndUnmarshal: %w", err)
	}

	genre := &Genre{
		ID:        resp.Current.ID,
		Name:      resp.Current.Name,
		Level:     resp.Current.Level,
		TagGroups: resp.TagGroups,
	}
	if err := c.setChildGenres(ctx, genre); err != nil {
		return nil, fmt.Errorf("c.setChildGenres: %w", err)
	}
	return genre, nil
}

func (c *Client) setChildGenres(ctx context.Context, genre *Genre) error {
	time.Sleep((time.Duration(1000.0 / float64(c.ApplicationIDNum()))) * time.Millisecond)

	rakutenGenre, err := c.SearchGenre(ctx, strconv.Itoa(genre.ID))
	if err != nil {
		return fmt.Errorf("rakutenIchibaAPIClient.SearchGenre: %w", err)
	}
	genres := make([]*Genre, 0, len(genre.Children))
	for _, child := range rakutenGenre.Children {
		childGenre := child.Child
		g := &Genre{
			ID:     childGenre.ID,
			Name:   childGenre.Name,
			Level:  childGenre.Level,
			Parent: genre,
		}
		if err := c.setChildGenres(ctx, g); err != nil {
			return err
		}
		genres = append(genres, g)
	}

	genre.TagGroups = rakutenGenre.TagGroups
	genre.Children = genres

	return nil
}
