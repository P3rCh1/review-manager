package models

import (
	"net/http"
)

type ErrorResponce struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *ErrorResponce) Error() string {
	return e.Message
}

var (
	ErrTeamExists = &ErrorResponce{http.StatusBadRequest, "TEAM_EXISTS", "team_name already exists"}

	ErrPRExists    = &ErrorResponce{http.StatusConflict, "PR_EXISTS", "PR id already exists"}
	ErrPRMerged    = &ErrorResponce{http.StatusConflict, "PR_MERGED", "cannot reassign on merged PR"}
	ErrNotAssigned = &ErrorResponce{http.StatusConflict, "NOT_ASSIGNED", "reviewer is not assigned to this PR"}
	ErrNoCandidate = &ErrorResponce{http.StatusConflict, "NO_CANDIDATE", "no active replacement candidate in team"}

	ErrPRNotFound   = &ErrorResponce{http.StatusNotFound, "PR_NOT_FOUND", "pull request not found"}
	ErrTeamNotFound = &ErrorResponce{http.StatusNotFound, "TEAM_NOT_FOUND", "team not found"}
	ErrUserNotFound = &ErrorResponce{http.StatusNotFound, "USER_NOT_FOUND", "user not found"}

	ErrInvalidInput  = &ErrorResponce{http.StatusBadRequest, "INVALID_INPUT", "invalid request body"}
	ErrRepeatableIDs = &ErrorResponce{http.StatusBadRequest, "REPEATABLE_IDS", "repeatable IDs"}

	ErrInternal = &ErrorResponce{http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error"}
)
