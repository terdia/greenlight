package mock

import (
	"github.com/terdia/greenlight/infrastructures/dto"
	"github.com/terdia/greenlight/internal/data"
	"github.com/terdia/greenlight/src/movies/entities"
	"github.com/terdia/greenlight/src/movies/repositories"
)

type movieRepositoryMock struct {
	totalRecords int
}

func NewMovieRepoitoryMock(totalRecords int) repositories.MovieRepository {
	return &movieRepositoryMock{totalRecords: totalRecords}
}

func (repo *movieRepositoryMock) Insert(movie *entities.Movie) error {
	return nil
}

func (repo *movieRepositoryMock) Get(id int64) (*entities.Movie, error) {

	if id < 1 {
		return nil, data.ErrRecordNotFound
	}

	var movie entities.Movie

	return &movie, nil
}

func (repo *movieRepositoryMock) Update(movie *entities.Movie) error {
	return nil
}

func (repo *movieRepositoryMock) Delete(id int64) error {
	if id < 1 {
		return data.ErrRecordNotFound
	}

	return nil
}

func (repo *movieRepositoryMock) GetAll(r dto.ListMovieRequest) ([]*entities.Movie, data.Metadata, error) {

	filters := r.Filters

	movies := []*entities.Movie{}

	metadata := data.CalculateMetadata(repo.totalRecords, filters.Page, filters.PageSize)

	return movies, metadata, nil
}
