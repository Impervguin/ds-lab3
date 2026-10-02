package common

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Impervguin/ds-lab2/gateway/internal/domain"
	"github.com/Impervguin/ds-lab2/gateway/internal/logger"
)

type ErrorWriter struct {
	log *slog.Logger
}

func NewErrorWriter(name string) ErrorWriter {
	return ErrorWriter{log: logger.Named(name)}
}

func (e ErrorWriter) Write(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, domain.ErrLibraryNotFound):
		WriteError(w, http.StatusNotFound, "Library not found")
	case errors.Is(err, domain.ErrBookNotFound):
		WriteError(w, http.StatusNotFound, "Book not found")
	case errors.Is(err, domain.ErrReservationNotFound):
		WriteError(w, http.StatusNotFound, "Reservation not found")
	case errors.Is(err, domain.ErrNoAvailableCopies):
		WriteError(w, http.StatusConflict, "No available copies of the book left")
	case errors.Is(err, domain.ErrReservationClosed):
		WriteError(w, http.StatusConflict, "Reservation is already closed")
	case errors.Is(err, domain.ErrBookLimitReached):
		WriteValidationError(w, "The number of books taken has reached the limit", nil)
	default:
		e.log.ErrorContext(r.Context(), "request failed",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("error", err.Error()),
		)
		WriteError(w, http.StatusInternalServerError, "Internal server error")
	}
}
