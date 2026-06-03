package model

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Genre struct {
	GenreID   int    `bson:"genre_id" json:"genre_id" validate:"required"`
	GenreName string `bson:"genre_name" json:"genre_name" validate:"required,min=2,max=100"`
}

type Review struct {
	ReviewID int    `bson:"review_id" json:"review_id"`
	AuthorID int    `bson:"author_id" json:"author_id"`
	Text     string `bson:"text" json:"text"`
}

type Ranking struct {
	RankingValue int    `bson:"ranking_value" json:"ranking_value" validate:"required"`
	RankingName  string `bson:"ranking_name" json:"ranking_name" validate:"required"`
}

type Movie struct {
	ID         bson.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ImdbID     string        `bson:"imdb_id" json:"imdb_id" validate:"required"`
	Title      string        `bson:"title" json:"title" validate:"required,min=2,max=500"`
	PosterPath string        `bson:"poster_path" json:"poster_path" validate:"required,url"`
	YouTubeID  string        `bson:"youtube_id" json:"youtube_id" validate:"required"`
	Genre      []Genre       `bson:"genre" json:"genre" validate:"required,dive"`
	Review     []Review      `bson:"review" json:"review"`
	Ranking    Ranking       `bson:"ranking" json:"ranking" validate:"required"`
}

type AddMovie struct {
	ImdbID     string
	Title      string
	PosterPath string
	YouTubeID  string
	Genre      []Genre
	Ranking    Ranking
}

type Plan struct {
	PlanID       string
	Name         string
	Price        int64
	Currency     string
	DurationDays int32
	IsActive     bool
}

type UpdateReviewResult struct {
	RankingName  string
	RankingValue int32
	Text         string
}

type GetMoviesFilter struct {
	Page       int32
	PageSize   int32
	GenreNames []string
}
