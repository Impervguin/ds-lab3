package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/Impervguin/ds-lab2/gateway/internal/domain"
)

type LibraryUseCase struct {
	libraries LibraryService
}

func NewLibraryUseCase(libraries LibraryService) *LibraryUseCase {
	return &LibraryUseCase{libraries: libraries}
}

func (uc *LibraryUseCase) ListLibraries(ctx context.Context, city string, page, size int) (Page[domain.Library], error) {
	return uc.libraries.ListLibraries(ctx, city, page, size)
}

func (uc *LibraryUseCase) ListBooks(
	ctx context.Context,
	libraryUID uuid.UUID,
	page, size int,
	showAll bool,
) (Page[domain.LibraryBook], error) {
	return uc.libraries.ListBooks(ctx, libraryUID, page, size, showAll)
}
