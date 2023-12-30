package middleware

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type authMiddleware struct {
	scopes     []string
	authUrl    string
	httpClient http.Client
}

type AuthData struct {
	ID         string   `json:"id"`
	FirstName  string   `json:"first_name"`
	MiddleName string   `json:"middle_name"`
	LastName   string   `json:"last_name"`
	Email      string   `json:"email"`
	Scopes     []string `json:"scopes"`
}

func (a *authMiddleware) Auth(c *gin.Context) {
	header := c.GetHeader("Authorization")
	token, err := extractBearerToken(header)

	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{
			"error": gin.H{
				"code":    1,
				"message": "Unauthorized",
			},
		})
		return
	}

	data, err := a.fetchAuthData(token)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"error": gin.H{
				"code":    1,
				"message": "Failed to Fetch Token",
			},
		})
		return
	}

	if !checkScopes(data.Scopes, a.scopes) {
		c.AbortWithStatusJSON(403, gin.H{
			"error": gin.H{
				"code":    1,
				"message": "Unauthorized",
			},
		})
		return
	}
	c.Set("auth_data", data)
}

func (a *authMiddleware) fetchAuthData(token string) (AuthData, error) {
	authData := AuthData{}
	req, err := http.NewRequest("GET", a.authUrl, nil)
	if err != nil {
		return authData, err
	}

	response, err := a.httpClient.Do(req)
	if err != nil {
		return authData, err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return authData, err
	}
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return authData, err
	}
	err = json.Unmarshal(bytes, &authData)
	return authData, err
}

func extractBearerToken(bearer string) (string, error) {
	if bearer == "" {
		return "", errors.New("bad header value given")
	}

	jwtToken := strings.Split(bearer, " ")
	if len(jwtToken) != 2 {
		return "", errors.New("incorrectly formatted authorization header")
	}

	return jwtToken[1], nil
}

func checkScopes(scopes []string, requiredScopes []string) bool {
	requiredMap := make(map[string]bool)
	for _, requiredScope := range requiredScopes {
		requiredMap[requiredScope] = true
	}

	for _, scope := range scopes {
		if requiredMap[scope] {
			return true
		}
	}

	return false
}
