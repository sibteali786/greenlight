package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"greenlight.sibteali.net/internal/validator"
)

type MovieModel struct {
	DB *sql.DB
}

func (m MovieModel) Insert(movie *Movie) error {
	insertMovieQuery := `
		INSERT INTO movies (title, year, runtime, genres)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version`
	queryArgumentsSlice := []interface{}{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}
	return m.DB.QueryRow(insertMovieQuery, queryArgumentsSlice...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (m MovieModel) Get(id int64) (*Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	getMovieQuery := `
	SELECT id, created_at, title, year, runtime, genres, version
	FROM movies
	WHERE id = $1`
	var movie Movie

	err := m.DB.QueryRow(getMovieQuery, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &movie, nil
}

func (m MovieModel) Update(movie *Movie) error {
	updateMovieQuery := `
		UPDATE movies 
		SET title=$1, year=$2, runtime=$3, genres=$4, version=version+1
		WHERE id=$5
		RETURNING version`
	updateMovieQueryArguments := []interface{}{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
	}

	return m.DB.QueryRow(updateMovieQuery, updateMovieQueryArguments...).Scan(&movie.Version)
}

func (m MovieModel) Delete(id int64) error {
	return nil
}

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`    // force it to be string
	Runtime   Runtime   `json:"runtime,omitempty"` // chose to parse field as a string
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes")

	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 500, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")

}
