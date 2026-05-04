package catalog

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/LearLocker/streaming/catalog/internal/model"
	"github.com/LearLocker/streaming/catalog/internal/repository"
	def "github.com/LearLocker/streaming/catalog/internal/service"
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms/openai"
)

var _ def.CatalogService = (*service)(nil)

type service struct {
	catalogRepository repository.CatalogRepository
}

func NewService(catalogRepository repository.CatalogRepository) *service {
	return &service{
		catalogRepository: catalogRepository,
	}
}

func (s *service) GetMovies(ctx context.Context, genreNames []string, page int32, pageSize int32) ([]*model.Movie, int32, error) {
	movies, total, err := s.catalogRepository.GetMovies(ctx, genreNames, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("get movies: %w", err)
	}

	return movies, total, nil
}

func (s *service) GetMovieById(ctx context.Context, imdbId string) (*model.Movie, error) {
	if imdbId == "" {
		return nil, fmt.Errorf("imdb_id is required")
	}

	movie, err := s.catalogRepository.GetMovieById(ctx, imdbId)
	if err != nil {
		return nil, fmt.Errorf("get movie: %w", err)
	}

	return movie, nil
}

func (s *service) AddMovie(ctx context.Context, info model.AddMovie) (string, error) {
	if info.ImdbID == "" || info.Title == "" {
		return "", fmt.Errorf("imdb_id and title are required")
	}

	insertedID, err := s.catalogRepository.AddMovie(ctx, info)
	if err != nil {
		return "", fmt.Errorf("add movie: %w", err)
	}

	return insertedID, nil
}

func (s *service) UpdateReview(ctx context.Context, imdbId string, authorId string, text string) (*model.UpdateReviewResult, error) {
	if imdbId == "" {
		return nil, fmt.Errorf("imdb_id is required")
	}
	if text == "" {
		return nil, fmt.Errorf("review text is required")
	}

	ranking, err := s.getReviewRanking(ctx, text)
	if err != nil {
		return nil, fmt.Errorf("get review ranking: %w", err)
	}

	if err = s.catalogRepository.UpdateReview(ctx, imdbId, authorId, text, *ranking); err != nil {
		return nil, fmt.Errorf("update review: %w", err)
	}

	return &model.UpdateReviewResult{
		RankingName:  ranking.RankingName,
		RankingValue: int32(ranking.RankingValue),
		Text:         text,
	}, nil
}

func (s *service) GetRecommendedMovies(ctx context.Context, userId string, limit int32) ([]*model.Movie, error) {
	if userId == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	movies, err := s.catalogRepository.GetRecommendedMovies(ctx, userId, limit)
	if err != nil {
		return nil, fmt.Errorf("get recommended movies: %w", err)
	}

	return movies, nil
}

func (s *service) GetGenres(ctx context.Context) ([]*model.Genre, error) {
	genres, err := s.catalogRepository.GetGenres(ctx)
	if err != nil {
		return nil, fmt.Errorf("get genres: %w", err)
	}

	return genres, nil
}

func (s *service) GetPlan(ctx context.Context, planId string) (*model.Plan, error) {
	if planId == "" {
		return nil, fmt.Errorf("plan_id is required")
	}

	plan, err := s.catalogRepository.GetPlan(ctx, planId)
	if err != nil {
		return nil, fmt.Errorf("get plan: %w", err)
	}

	return plan, nil
}

func (s *service) ListPlans(ctx context.Context, onlyActive bool) ([]*model.Plan, error) {
	plans, err := s.catalogRepository.ListPlans(ctx, onlyActive)
	if err != nil {
		return nil, fmt.Errorf("list plans: %w", err)
	}

	return plans, nil
}

func (s *service) getReviewRanking(ctx context.Context, review string) (*model.Ranking, error) {
	rankings, err := s.catalogRepository.GetRankings(ctx)
	if err != nil {
		return nil, err
	}

	var parts []string
	for _, r := range rankings {
		if r.RankingValue != 999 {
			parts = append(parts, r.RankingName)
		}
	}
	sentimentDelimited := strings.Join(parts, ",")

	_ = godotenv.Load("internal/config/.env")

	openAiApiKey := os.Getenv("OPEN_AI_API_KEY")
	if openAiApiKey == "" {
		return nil, errors.New("OPEN_AI_API_KEY is not set")
	}

	basePromptTemp := os.Getenv("BASE_PROMPT_TEMP")
	if basePromptTemp == "" {
		return nil, errors.New("BASE_PROMPT_TEMP is not set")
	}

	llm, err := openai.New(openai.WithToken(openAiApiKey))
	if err != nil {
		return nil, err
	}

	basePrompt := strings.Replace(basePromptTemp, "{rankings}", sentimentDelimited, 1)
	response, err := llm.Call(ctx, basePrompt+review)
	if err != nil {
		return nil, err
	}

	// найти ranking_value по имени из ответа LLM
	rankVal := 0
	for _, r := range rankings {
		if r.RankingName == response {
			rankVal = r.RankingValue
			break
		}
	}

	return &model.Ranking{RankingName: response, RankingValue: rankVal}, nil
}
