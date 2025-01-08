// package api

// import (
// 	// "net/http"
// 	// "github.com/julienschmidt/httprouter"
// 	"github.com/LingFengJ/WASAText/service/api/reqcontext"
// 	"github.com/LingFengJ/WASAText/service/database"
// )

// // GetUserID helper function to extract user ID from context
// func (rt *_router) GetUserID(ctx reqcontext.RequestContext) (string, error) {
//     if ctx.UserID == "" {
//         return "", database.ErrUserNotFound
//     }

//     // TODO:  add additional validation here
//     // For example, checking if the ID matches a UUID format
//     // or if it exists in the database

//     return ctx.UserID, nil
// }

package api

import (
	"errors"
	"github.com/LingFengJ/WASAText/service/api/reqcontext"
)

var (
	ErrIdentifierNotFound = errors.New("identifier not found in context")
	ErrUserIDNotFound     = errors.New("user ID not found in context")
)

func GetIdentifierFromContext(ctx reqcontext.RequestContext) (string, error) {
	if ctx.Identifier == "" {
		return "", ErrIdentifierNotFound
	}
	return ctx.Identifier, nil
}

func GetUserIDFromContext(ctx reqcontext.RequestContext) (string, error) {
	if ctx.UserID == "" {
		return "", ErrUserIDNotFound
	}
	return ctx.UserID, nil
}
