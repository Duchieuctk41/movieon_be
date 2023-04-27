package ai

import (
	"context"
	"encoding/json"
	"github.com/praslar/cloud0/ginext"
	"github.com/praslar/cloud0/logger"
	"github.com/praslar/lib/common"
	"github.com/sendgrid/rest"
	"movieon_be/pkg/model"
	"movieon_be/pkg/utils"
	"net/http"
)

type AiRestApiClient struct{}

func GetAiRestApiClient() AiApiInterface {
	return &AiRestApiClient{}
}

func (c *AiRestApiClient) GetListSuggest(ctx context.Context, idOld string) (*model.ListSuggestResponse, error) {
	log := logger.WithCtx(context.Background(), utils.GetCurrentCaller(c, 0))

	// call api
	data, _, err := common.SendRestAPI("http://localhost:8000/api/v1/movie/get-list-suggest/"+idOld, rest.Get, nil, nil, nil)
	if err != nil {
		log.WithError(err).Errorf("error_400: error call api")
		return nil, ginext.NewError(http.StatusBadRequest, err.Error())
	}

	// unmarshal
	res := &model.ListSuggestResponse{}
	err = json.Unmarshal([]byte(data), &res)
	if err != nil {
		log.WithError(err).Errorf("error_400: error unmarshal api")
		return nil, ginext.NewError(http.StatusBadRequest, err.Error())
	}

	return res, nil
}
