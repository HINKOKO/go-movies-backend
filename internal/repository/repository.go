package repository

import (
	"backend/internal/models"
	"database/sql"
)

/*
*
The two packages 'repository' & 'dbrepo' are linked through the "DatabaseRepo interface" defined in the **repository package**
and implemented by the "PostgresDBRepo type" in the **dbrepo package.**
==> This separation allows for modularity and flexibility in the design, enabling different implementations of
the repository interface for different database backends, without directly exposing the database details to higher-level application logic.
*/

type DatabaseRepo interface {
	Connection() *sql.DB
	AllMovies(genre ...int) ([]*models.Movie, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)

	OneMovieForEdit(id int) (*models.Movie, []*models.Genre, error)
	OneMovie(id int) (*models.Movie, error)

	AllGenres() ([]*models.Genre, error)

	InsertMovie(movie models.Movie) (int, error)
	UpdateMovie(movie models.Movie) error
	UpdateMovieGenres(id int, genreIDs []int) error

	DeleteMovie(id int) error
}
