package dbrepo

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type PostgresDBRepo struct {
	DB *sql.DB
}

// 3 seconds counter to interact with database, if longer - too late
const dbTimeout = time.Second * 3

func (m *PostgresDBRepo) Connection() *sql.DB {
	return m.DB
}

// AllMovies - Hybrid function that fetch all movies whatever the genres, but also
// all movies by specified genre(s)
func (m *PostgresDBRepo) AllMovies(genre ...int) ([]*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	where := ""
	if len(genre) > 0 {
		where = fmt.Sprintf("WHERE id in (SELECT movie_id FROM movies_genres WHERE genre_id = %d)", genre[0])
	}

	var movies []*models.Movie
	query := fmt.Sprintf(`SELECT id, title, release_date, runtime, mpaa_rating,
		description, coalesce(image, ''), created_at, updated_at FROM movies %s
		ORDER BY title`, where)

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var movie models.Movie
		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.ReleaseDate,
			&movie.RunTime,
			&movie.MPAARating,
			&movie.Description,
			&movie.Image,
			&movie.CreatedAt,
			&movie.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		movies = append(movies, &movie)
	}
	return movies, nil
}

func (m *PostgresDBRepo) OneMovieForEdit(id int) (*models.Movie, []*models.Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var movie models.Movie
	query := `SELECT id, title, release_date, runtime, mpaa_rating,
		description, coalesce(image, ''), created_at, updated_at FROM movies WHERE id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.ReleaseDate,
		&movie.RunTime,
		&movie.MPAARating,
		&movie.Description,
		&movie.Image,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)
	if err != nil {
		return nil, nil, err
	}

	// get genres, if any
	query = `SELECT g.id, g.genre FROM movies_genres mg
	LEFT JOIN genres g ON (mg.genre_id = g.id)
	WHERE mg.movie_id = $1
	`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, nil, err
	}
	defer rows.Close()

	var genres []*models.Genre
	var genresArray []int

	for rows.Next() {
		var g models.Genre
		err := rows.Scan(
			&g.ID,
			&g.Genre,
		)
		if err != nil {
			return nil, nil, err
		}
		genres = append(genres, &g)
		genresArray = append(genresArray, g.ID)
	}
	// finally append the slice of pointers to genres to the movie
	movie.Genres = genres
	movie.GenresArray = genresArray

	var allGenres []*models.Genre

	query = `SELECT id, genre FROM genres`
	gRows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, nil, err
	}
	defer gRows.Close()

	for gRows.Next() {
		var g models.Genre
		err := gRows.Scan(
			&g.ID,
			&g.Genre,
		)
		if err != nil {
			return nil, nil, err
		}
		allGenres = append(allGenres, &g)
	}
	return &movie, allGenres, nil
}

func (m *PostgresDBRepo) OneMovie(id int) (*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var movie models.Movie
	query := `SELECT id, title, release_date, runtime, mpaa_rating,
		description, coalesce(image, ''), created_at, updated_at FROM movies WHERE id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.ReleaseDate,
		&movie.RunTime,
		&movie.MPAARating,
		&movie.Description,
		&movie.Image,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// get genres, if any
	query = `SELECT g.id, g.genre FROM movies_genres mg
	LEFT JOIN genres g ON (mg.genre_id = g.id)
	WHERE mg.movie_id = $1
	ORDER BY g.genre`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	var genres []*models.Genre
	for rows.Next() {
		var g models.Genre
		err := rows.Scan(
			&g.ID,
			&g.Genre,
		)
		if err != nil {
			return nil, err
		}
		genres = append(genres, &g)
	}
	// finally append the slice of pointers to genres to the movie
	movie.Genres = genres

	return &movie, nil
}

// GetUserByEmail - check in Database if user exists
func (m *PostgresDBRepo) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var user models.User
	query := `SELECT id, email, first_name, last_name, password, created_at, updated_at FROM users WHERE email = $1`
	row := m.DB.QueryRowContext(ctx, query, email)
	// Note the 'QueryRowContext' defers the potentials errors until `row.Scan()` is called ;)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *PostgresDBRepo) GetUserByID(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var user models.User
	query := `SELECT id, email, first_name, last_name, password, created_at, updated_at FROM users WHERE id = $1`
	row := m.DB.QueryRowContext(ctx, query, id)
	// Note the 'QueryRowContext' defers the potentials errors until `row.Scan()` is called ;)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *PostgresDBRepo) AllGenres() ([]*models.Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var genres []*models.Genre

	query := `SELECT id, genre FROM genres
		ORDER BY genre`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var g models.Genre
		rows.Scan(
			&g.ID,
			&g.Genre,
		)
		genres = append(genres, &g)
	}
	defer rows.Close()

	return genres, nil
}

func (m *PostgresDBRepo) InsertMovie(movie models.Movie) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var movieID int
	stmt := `INSERT INTO movies (title, description, release_date, runtime,
		mpaa_rating, created_at, updated_at, image)
		VALUES ($1, $2, $3, $4, $5,$6, $7, $8) RETURNING id`

	// `QueryRowContext` a query that is actually an Insert
	// ==>  Errors are deferred until -> [Row]'s Scan method is called.
	// That's why we call `err` the return value of Query here
	err := m.DB.QueryRowContext(ctx, stmt,
		movie.Title, movie.Description, movie.ReleaseDate, movie.RunTime,
		movie.MPAARating, movie.CreatedAt, movie.UpdatedAt, movie.Image).Scan(&movieID)

	if err != nil {
		return 0, err
	}

	return movieID, nil
}

func (m *PostgresDBRepo) UpdateMovie(movie models.Movie) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	stmt := `UPDATE movies SET title = $1, description = $2, release_date = $3, runtime = $4, mpaa_rating = $5,
	updated_at = $6, image = $7 WHERE id = $8`

	_, err := m.DB.ExecContext(ctx, stmt,
		movie.Title,
		movie.Description,
		movie.ReleaseDate,
		movie.RunTime,
		movie.MPAARating,
		movie.UpdatedAt,
		movie.Image,
		movie.ID)

	if err != nil {
		return err
	}
	return nil
}

func (m *PostgresDBRepo) UpdateMovieGenres(id int, genreIDs []int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// Easiest way -> Delete from table `genres` everything that has the movie ID that
	// we receive as a call to this function
	stmt := `DELETE FROM movies_genres where movie_id = $1`
	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	// At this point, no genres associated with the movie with id 'id'
	for _, n := range genreIDs {
		stmt := `INSERT INTO movies_genres (movie_id, genre_id) VALUES ($1, $2)`
		_, err := m.DB.ExecContext(ctx, stmt, id, n)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *PostgresDBRepo) DeleteMovie(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	stmt := `DELETE FROM movies where id = $1`
	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}
	return nil
}
