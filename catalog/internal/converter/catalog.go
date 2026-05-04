package converter

import (
	"github.com/LearLocker/streaming/catalog/internal/model"
	catalogV1 "github.com/LearLocker/streaming/shared/pkg/proto/catalog/v1"
)

func AddMovieInfoFromProto(req *catalogV1.AddMovieRequest) model.AddMovie {
	genres := make([]model.Genre, 0, len(req.GetGenre()))
	for _, g := range req.GetGenre() {
		genres = append(genres, model.Genre{
			GenreID:   int(g.GetGenreId()),
			GenreName: g.GetGenreName(),
		})
	}

	return model.AddMovie{
		ImdbID:     req.GetImdbId(),
		Title:      req.GetTitle(),
		PosterPath: req.GetPosterPath(),
		YouTubeID:  req.GetYoutubeId(),
		Genre:      genres,
		Ranking: model.Ranking{
			RankingValue: int(req.GetRanking().GetRankingValue()),
			RankingName:  req.GetRanking().GetRankingName(),
		},
	}
}

func MovieToProto(m *model.Movie) *catalogV1.Movie {
	genres := make([]*catalogV1.Genre, 0, len(m.Genre))
	for _, g := range m.Genre {
		genres = append(genres, &catalogV1.Genre{
			GenreId:   int32(g.GenreID),
			GenreName: g.GenreName,
		})
	}

	return &catalogV1.Movie{
		ImdbId:     m.ImdbID,
		Title:      m.Title,
		PosterPath: m.PosterPath,
		YoutubeId:  m.YouTubeID,
		Genre:      genres,
		Ranking: &catalogV1.Ranking{
			RankingValue: int32(m.Ranking.RankingValue),
			RankingName:  m.Ranking.RankingName,
		},
	}
}

func MoviesToProto(movies []*model.Movie) []*catalogV1.Movie {
	result := make([]*catalogV1.Movie, 0, len(movies))
	for _, m := range movies {
		result = append(result, MovieToProto(m))
	}
	return result
}

func PlanToProto(p *model.Plan) *catalogV1.Plan {
	return &catalogV1.Plan{
		PlanId:       p.PlanID,
		Name:         p.Name,
		Price:        p.Price,
		Currency:     p.Currency,
		DurationDays: p.DurationDays,
		IsActive:     p.IsActive,
	}
}

func PlansToProto(plans []*model.Plan) []*catalogV1.Plan {
	result := make([]*catalogV1.Plan, 0, len(plans))
	for _, p := range plans {
		result = append(result, PlanToProto(p))
	}
	return result
}

func GenresToProto(genres []*model.Genre) []*catalogV1.Genre {
	result := make([]*catalogV1.Genre, 0, len(genres))
	for _, g := range genres {
		result = append(result, &catalogV1.Genre{
			GenreId:   int32(g.GenreID),
			GenreName: g.GenreName,
		})
	}
	return result
}

func toString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func toInt(v any) int {
	switch n := v.(type) {
	case int:
		return n
	case int32:
		return int(n)
	case int64:
		return int(n)
	case float64:
		return int(n) // MongoDB числа часто приходят как float64
	}
	return 0
}

func toInt64(v any) int64 {
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

func toBool(v any) bool {
	if b, ok := v.(bool); ok {
		return b
	}
	return false
}
