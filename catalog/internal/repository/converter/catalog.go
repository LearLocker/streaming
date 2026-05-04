package converter

import (
	"github.com/LearLocker/streaming/catalog/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func MovieFromBson(doc bson.M) *model.Movie {
	movie := &model.Movie{
		ImdbID:     ToString(doc["imdb_id"]),
		Title:      ToString(doc["title"]),
		PosterPath: ToString(doc["poster_path"]),
		YouTubeID:  ToString(doc["youtube_id"]),
	}

	if genres, ok := doc["genre"].(bson.A); ok {
		for _, g := range genres {
			if gm, ok := g.(bson.M); ok {
				movie.Genre = append(movie.Genre, model.Genre{
					GenreID:   ToInt(gm["genre_id"]),
					GenreName: ToString(gm["genre_name"]),
				})
			}
		}
	}

	if ranking, ok := doc["ranking"].(bson.M); ok {
		movie.Ranking = model.Ranking{
			RankingValue: ToInt(ranking["ranking_value"]),
			RankingName:  ToString(ranking["ranking_name"]),
		}
	}

	return movie
}

func PlanFromBson(doc bson.M) *model.Plan {
	return &model.Plan{
		PlanID:       ToString(doc["plan_id"]),
		Name:         ToString(doc["name"]),
		Price:        ToInt64(doc["price"]),
		Currency:     ToString(doc["currency"]),
		DurationDays: int32(ToInt(doc["duration_days"])),
		IsActive:     ToBool(doc["is_active"]),
	}
}

func GenresToBson(genres []model.Genre) bson.A {
	result := make(bson.A, 0, len(genres))
	for _, g := range genres {
		result = append(result, bson.M{
			"genre_id":   g.GenreID,
			"genre_name": g.GenreName,
		})
	}
	return result
}

func ToString(v any) string {
	s, _ := v.(string)
	return s
}

func ToInt(v any) int {
	switch n := v.(type) {
	case int:
		return n
	case int32:
		return int(n)
	case int64:
		return int(n)
	case float64:
		return int(n)
	}
	return 0
}

func ToInt64(v any) int64 {
	switch n := v.(type) {
	case int64:
		return n
	case int32:
		return int64(n)
	case float64:
		return int64(n)
	}
	return 0
}

func ToBool(v any) bool {
	b, _ := v.(bool)
	return b
}
