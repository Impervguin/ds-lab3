package usecase

import (
	"context"

	"github.com/Impervguin/ds-lab2/gateway/internal/domain"
)

type RatingUseCase struct {
	ratings RatingService
}

func NewRatingUseCase(ratings RatingService) *RatingUseCase {
	return &RatingUseCase{ratings: ratings}
}

func (uc *RatingUseCase) GetRating(ctx context.Context, username string) (*domain.Rating, error) {
	return uc.ratings.GetRating(ctx, username)
}
