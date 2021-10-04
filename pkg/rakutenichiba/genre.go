package rakutenichiba

import "strconv"

type Genre struct {
	ID       int    `json:"genreId"`
	Name     string `json:"genreName"`
	Level    int    `json:"genreLevel"`
	Parent   *Genre
	Children []*Genre
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
