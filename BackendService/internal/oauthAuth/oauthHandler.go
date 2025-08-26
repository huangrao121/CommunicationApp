package oauthauth

type OauthHandler struct {
	oauthStore *OauthStore
}

func NewOauthHandler(oauthStore *OauthStore) *OauthHandler {
	return &OauthHandler{oauthStore: oauthStore}
}
