package httpclient

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/Impervguin/ds-lab2/gateway/internal/domain"
	"github.com/Impervguin/ds-lab2/gateway/internal/usecase"
)

type RatingClient struct {
	baseClient
}

var _ usecase.RatingService = (*RatingClient)(nil)

func NewRatingClient(baseURL string) *RatingClient {
	return &RatingClient{baseClient: newBaseClient(baseURL)}
}

func (c *RatingClient) GetRating(ctx context.Context, username string) (*domain.Rating, error) {
	var response ratingResponse
	if err := c.do(ctx, call{
		method:   http.MethodGet,
		path:     "/api/v1/rating",
		username: username,
		out:      &response,
	}); err != nil {
		return nil, err
	}

	return &domain.Rating{Stars: response.Stars, MaxBooks: response.MaxBooks}, nil
}

func (c *RatingClient) CloseReservation(
	ctx context.Context,
	username string,
	closed domain.ClosedReservation,
) (*domain.RatingChange, error) {
	var response ratingChangeResponse
	if err := c.do(ctx, call{
		method:   http.MethodPost,
		path:     "/api/v1/rating/reservation-closed",
		username: username,
		body: reservationClosedRequest{
			ReservationUID:    closed.ReservationUID,
			ReservationStatus: string(closed.Status),
			ConditionAtRent:   string(closed.ConditionAtRent),
			ConditionOnReturn: string(closed.ConditionOnReturn),
		},
		out: &response,
	}); err != nil {
		return nil, err
	}

	return &domain.RatingChange{Delta: response.Delta, Stars: response.Stars}, nil
}

type ratingResponse struct {
	Stars    int `json:"stars"`
	MaxBooks int `json:"maxBooks"`
}

type ratingChangeResponse struct {
	Delta int `json:"delta"`
	Stars int `json:"stars"`
}

type reservationClosedRequest struct {
	ReservationUID    uuid.UUID `json:"reservationUid"`
	ReservationStatus string    `json:"reservationStatus"`
	ConditionAtRent   string    `json:"conditionAtRent"`
	ConditionOnReturn string    `json:"conditionOnReturn"`
}
