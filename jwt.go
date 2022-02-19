package permissions

import (
	"gopkg.in/square/go-jose.v2/jwt"
)

func GetRoleFromPayload(token string) (role string, err error) {
	// TODO: Flexible extraction of token key and validate token better

	var claims map[string]interface{}
	payload, err := jwt.ParseSigned(token)
	if err != nil || payload == nil {
		return
	}

	// This line retrieves the map of the payload
	err = payload.UnsafeClaimsWithoutVerification(&claims)
	if err != nil {
		return
	}

	// Role is inside of an Array
	role = claims["role"].([]interface{})[0].(string)
	return
}
