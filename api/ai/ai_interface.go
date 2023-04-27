package ai

import (
	"context"
	"movieon_be/pkg/model"
)

type AiApiInterface interface {
	GetListSuggest(ctx context.Context, idOld string) (*model.ListSuggestResponse, error)
}
