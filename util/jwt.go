package util

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/jwtauth"
)

// GetClaimFromJWT get claim value using id
func GetClaimFromJWT(r *http.Request, ClaimID string) (interface{}, error) {
	var claim interface{}
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		return claim, err
	}
	if _, ok := claims[ClaimID]; !ok {
		return claim, errors.New(fmt.Sprintf(ErrUtilJWTMapping, ClaimID))
	}
	claim, _ = claims[ClaimID]

	return claim, nil
}

// GetClaimStringFromJWT get ID From claim
func GetClaimStringFromJWT(r *http.Request, claimID string) (string, error) {
	claim, err := GetClaimFromJWT(r, claimID)
	if err != nil {
		return "", err
	}
	return claim.(string), nil
}
