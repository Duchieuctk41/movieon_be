package api

import "movieon_be/api/ai"

type ApiClient struct {
	AiClient ai.AiApiInterface
}

func GetApiClient() *ApiClient {
	return &ApiClient{
		AiClient: ai.GetAiRestApiClient(),
	}
}
