package catalog

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/LearLocker/streaming/catalog/internal/model"
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms/openai"
)

func (s *Service) UpdateReview(ctx context.Context, imdbId string, authorId string, text string) (*model.UpdateReviewResult, error) {
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

func (s *Service) getReviewRanking(ctx context.Context, review string) (*model.Ranking, error) {
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
