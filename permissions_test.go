package permissions

import (
	"log"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCheckRoleWithBackend(t *testing.T) {
	// Creation And Starting Endpoint
	router := gin.Default()
	router.POST("/has-permission", func(c *gin.Context) {
		var request RequestBody
		_ = c.ShouldBindJSON(&request)

		c.JSON(http.StatusOK, gin.H{
			"has_permission": request.Permissions[0].Name == "permissions_valid",
		})

	})
	go func() {
		log.Fatal(router.Run("localhost:8181"))
	}()

	// Checking
	testCases := map[string]bool{
		"permissions_valid":   true,
		"permissions_invalid": false,
	}

	for permission, v := range testCases {
		has := CheckRoleWithBackend("admin", permission, "http://localhost:8181/has-permission", "POST")
		if has != v {
			t.Fatal("Unexpected behavior", permission, has)
		}
	}
}
