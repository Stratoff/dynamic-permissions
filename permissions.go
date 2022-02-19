// This package manages permissions based on key found
// in the payload of a JWT
package permissions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RequestBody struct {
	Name        string `json:"name"`
	Permissions []struct {
		Name string `json:"name"`
	} `json:"permissions"`
}

type ResponseBody struct {
	HasPermission bool `json:"has_permission"`
}

func CheckRoleWithBackend(role, permission_code, backendURL, httpMethod string) (hasPermission bool) {
	// TODO: Flexible Interaction with Backend

	body, _ := json.Marshal(
		RequestBody{
			Name: role,
			Permissions: []struct {
				Name string `json:"name"`
			}{
				{Name: permission_code},
			},
		},
	)

	request, err := http.NewRequest(httpMethod, backendURL, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err)
		return
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return
	}

	responseBody, _ := ioutil.ReadAll(response.Body)

	var result ResponseBody
	if err = json.Unmarshal(responseBody, &result); err != nil {
		fmt.Println(err)
	}

	hasPermission = result.HasPermission
	return
}
