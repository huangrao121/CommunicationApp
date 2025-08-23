package types

type LoginResp struct {
	AppJwt  string `json:"app_jwt"`
	MqttJwt string `json:"mqtt_jwt"`
}

type MqttClaims struct {
	ID       string        `json:"id"`
	Username string        `json:"username"`
	Email    string        `json:"email"`
	ACL      []interface{} `json:"acl"`
}
