package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/Impervguin/ds-lab2/gateway/internal/domain"
)

type LibraryService interface {
	ListLibraries(ctx context.Context, city string, page, size int) (Page[domain.Library], error)
	GetLibrary(ctx context.Context, libraryUID uuid.UUID) (*domain.Library, error)
	ListBooks(ctx context.Context, libraryUID uuid.UUID, page, size int, showAll bool) (Page[domain.LibraryBook], error)
	GetBook(ctx context.Context, bookUID uuid.UUID) (*domain.Book, error)
	TakeBook(ctx context.Context, libraryUID, bookUID uuid.UUID) (*domain.BookCopy, error)
	ReturnBook(ctx context.Context, libraryUID, bookUID uuid.UUID, condition domain.BookCondition) (*domain.BookCopy, error)
}

type ReservationService interface {
	ListReservations(ctx context.Context, username string, status domain.ReservationStatus) ([]domain.Reservation, error)
	CountReservations(ctx context.Context, username string, status domain.ReservationStatus) (int, error)
	CreateReservation(ctx context.Context, username string, request domain.NewReservation) (*domain.Reservation, error)
	ReturnReservation(ctx context.Context, username string, reservationUID uuid.UUID, date domain.Date) (*domain.Reservation, error)
}

type RatingService interface {
	GetRating(ctx context.Context, username string) (*domain.Rating, error)
	CloseReservation(ctx context.Context, username string, closed domain.ClosedReservation) (*domain.RatingChange, error)
}
