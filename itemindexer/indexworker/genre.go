package indexworker

import (
	"strconv"

	"github.com/k-yomo/kagu-miru/pkg/rakutenichiba"
)

type Genre struct {
	*rakutenichiba.Genre
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
	if len(genreIDs) > 1 {
		// remove root genre id since it's same for all items
		genreIDs = genreIDs[1:]
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
	if len(genreNames) > 1 {
		// remove root genre name since it's same for all items
		genreNames = genreNames[1:]
	}
	return genreNames
}
