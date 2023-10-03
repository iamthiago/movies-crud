package service

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/iamthiago/movies-crud/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) GetMovies() ([]models.Movie, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]models.Movie), nil
}

func (m *mockRepo) GetMovieById(id int64) (*models.Movie, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Movie), nil
}

func (m *mockRepo) CreateMovie(movie *models.Movie) (*models.Movie, error) {
	args := m.Called(movie)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.Movie), nil
}

func (m *mockRepo) UpdateMovie(id int64, movie *models.Movie) (*models.Movie, error) {
	args := m.Called(id, movie)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.Movie), nil
}

func (m *mockRepo) DeleteMovie(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestGetMovies(t *testing.T) {
	mockRepository := new(mockRepo)
	service := Service{Repository: mockRepository}
	movies := []models.Movie{
		{ID: 1, Isbn: "9788401490040", Title: "Jaws", Director: "Steven Spielberg"},
	}

	testCases := []struct {
		name      string
		mockSetup func() []*mock.Call
		wantErr   bool
		err       string
	}{
		{
			name: "Should return a list of movies",
			mockSetup: func() []*mock.Call {
				return []*mock.Call{
					mockRepository.On("GetMovies").Return(movies, nil),
				}
			},
			wantErr: false,
			err:     "",
		},
		{
			name: "Should return error when calling get movies",
			mockSetup: func() []*mock.Call {
				return []*mock.Call{
					mockRepository.On("GetMovies").Return(nil, errors.New("failed")),
				}
			},
			wantErr: true,
			err:     "failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			calls := tc.mockSetup()

			resp, err := service.GetMovies()
			if tc.wantErr {
				assert.Nil(t, resp)
				assert.EqualError(t, err, tc.err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}

			for _, call := range calls {
				call.Unset()
			}
		})
	}
}

func TestGetMovieById(t *testing.T) {
	mockRepository := new(mockRepo)
	service := Service{Repository: mockRepository}
	var validId int64 = 1
	movie := models.Movie{ID: validId, Isbn: "9788401490040", Title: "Jaws", Director: "Steven Spielberg"}

	testCases := []struct {
		name      string
		mockSetup func(id int64) []*mock.Call
		wantErr   bool
		err       string
		args      struct {
			id int64
		}
	}{
		{
			name: "Should return a movie by id",
			mockSetup: func(id int64) []*mock.Call {
				return []*mock.Call{
					mockRepository.On("GetMovieById", mock.Anything).Return(&movie, nil),
				}
			},
			wantErr: false,
			err:     "",
			args:    struct{ id int64 }{id: validId},
		},
		{
			name: "Should return not found when trying to get a movie by id",
			mockSetup: func(id int64) []*mock.Call {
				return []*mock.Call{
					mockRepository.On("GetMovieById", mock.Anything).Return(nil, sql.ErrNoRows),
				}
			},
			wantErr: true,
			err:     "sql: no rows in result set",
			args:    struct{ id int64 }{id: 10},
		},
		{
			name: "Should return an error when getting movie by id",
			mockSetup: func(id int64) []*mock.Call {
				return []*mock.Call{
					mockRepository.On("GetMovieById", mock.Anything).Return(nil, errors.New("failed")),
				}
			},
			wantErr: true,
			err:     "failed",
			args:    struct{ id int64 }{id: 123},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			calls := tc.mockSetup(tc.args.id)

			resp, err := service.GetMovieById(tc.args.id)
			if tc.wantErr {
				assert.Nil(t, resp)
				assert.EqualError(t, err, tc.err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}

			for _, call := range calls {
				call.Unset()
			}
		})
	}
}

func TestCreateMovie(t *testing.T) {
	mockRepository := new(mockRepo)
	service := Service{Repository: mockRepository}

	testCases := []struct {
		name      string
		mockSetup func(movie *models.Movie) []*mock.Call
		req       models.Movie
		wantErr   bool
		err       string
	}{
		{
			name: "Should create a new movie",
			mockSetup: func(movie *models.Movie) []*mock.Call {
				return []*mock.Call{
					mockRepository.On("CreateMovie", mock.Anything).Return(movie, nil),
				}
			},
			req: models.Movie{
				ID:       123,
				Isbn:     "9788401490040",
				Title:    "Jaws",
				Director: "Steven Spielberg",
			},
			wantErr: false,
			err:     "",
		},
		{
			name: "Should return an error when creating a new movie",
			mockSetup: func(movie *models.Movie) []*mock.Call {
				return []*mock.Call{
					mockRepository.On("CreateMovie", mock.Anything).Return(nil, errors.New("failed")),
				}
			},
			req: models.Movie{
				ID:       123,
				Isbn:     "9788401490040",
				Title:    "Jaws",
				Director: "Steven Spielberg",
			},
			wantErr: true,
			err:     "failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			calls := tc.mockSetup(&tc.req)

			resp, err := service.CreateMovie(&tc.req)
			if tc.wantErr {
				assert.Nil(t, resp)
				assert.EqualError(t, err, tc.err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}

			for _, call := range calls {
				call.Unset()
			}
		})
	}
}

func TestUpdateMovie(t *testing.T) {
	mockRepository := new(mockRepo)
	service := Service{Repository: mockRepository}

	testCases := []struct {
		name      string
		mockSetup func(id int64, movie *models.Movie) []*mock.Call
		args      struct {
			id    int64
			movie models.Movie
		}
		wantErr bool
		err     string
	}{
		{
			name: "Should update movie",
			mockSetup: func(id int64, movie *models.Movie) []*mock.Call {
				return []*mock.Call{
					mockRepository.On("UpdateMovie", mock.Anything, mock.Anything).Return(movie, nil),
				}
			},
			args: struct {
				id    int64
				movie models.Movie
			}{
				id: 456,
				movie: models.Movie{
					ID:       456,
					Isbn:     "9788401490040",
					Title:    "Jaws",
					Director: "Steven Spielberg",
				},
			},
			wantErr: false,
			err:     "",
		},
		{
			name: "Should return an error when trying to update movie",
			mockSetup: func(id int64, movie *models.Movie) []*mock.Call {
				return []*mock.Call{
					mockRepository.On("UpdateMovie", mock.Anything, mock.Anything).Return(nil, errors.New("failed")),
				}
			},
			args: struct {
				id    int64
				movie models.Movie
			}{
				id: 890,
				movie: models.Movie{
					ID:       890,
					Isbn:     "9788401490040",
					Title:    "Jaws",
					Director: "Steven Spielberg",
				},
			},
			wantErr: true,
			err:     "failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			calls := tc.mockSetup(tc.args.id, &tc.args.movie)

			resp, err := service.UpdateMovie(tc.args.id, &tc.args.movie)
			if tc.wantErr {
				assert.Nil(t, resp)
				assert.EqualError(t, err, tc.err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}

			for _, call := range calls {
				call.Unset()
			}
		})
	}
}

func TestDeleteMovie(t *testing.T) {
	mockRepository := new(mockRepo)
	service := Service{Repository: mockRepository}

	testCases := []struct {
		name      string
		mockSetup func(id int64) []*mock.Call
		args      struct {
			id int64
		}
		wantErr bool
		err     string
	}{
		{
			name: "Should delete movie by id",
			mockSetup: func(id int64) []*mock.Call {
				return []*mock.Call{
					mockRepository.On("DeleteMovie", mock.Anything).Return(nil),
				}
			},
			args:    struct{ id int64 }{id: 567},
			wantErr: false,
			err:     "",
		},
		{
			name: "Should return an error when deleting by id",
			mockSetup: func(id int64) []*mock.Call {
				return []*mock.Call{
					mockRepository.On("DeleteMovie", mock.Anything).Return(errors.New("failed")),
				}
			},
			args:    struct{ id int64 }{id: 784},
			wantErr: true,
			err:     "failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			calls := tc.mockSetup(tc.args.id)

			err := service.DeleteMovie(tc.args.id)
			if tc.wantErr {
				assert.EqualError(t, err, tc.err)
			} else {
				assert.NoError(t, err)
			}

			for _, call := range calls {
				call.Unset()
			}
		})
	}
}
