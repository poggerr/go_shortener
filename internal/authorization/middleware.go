package authorization

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// AuthMiddleware мидлваря авторизации
func AuthMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		c, err := r.Cookie("session_token")
		if err != nil {

			uuidUserID := uuid.New()
			jwtString, err := BuildJWTString(&uuidUserID)
			if err != nil {
				fmt.Println(err.Error())
			}

			cook := &http.Cookie{
				Name:    "session_token",
				Value:   jwtString,
				Path:    "/",
				Domain:  "localhost",
				Expires: time.Now().Add(120 * time.Second),
			}

			http.SetCookie(w, cook)

			ur := r.WithContext(NewContext(r.Context(), &uuidUserID))
			h.ServeHTTP(w, ur)
			return
		}
		user := GetUserID(c.Value)

		ur := r.WithContext(NewContext(r.Context(), user))

		h.ServeHTTP(w, ur)
	}
	return http.HandlerFunc(fn)
}

type userID string

const ReqUserKey = userID("userKey")

func NewContext(ctx context.Context, user *uuid.UUID) context.Context {
	return context.WithValue(ctx, ReqUserKey, user)
}

func FromContext(ctx context.Context) *uuid.UUID {
	u := ctx.Value(ReqUserKey)
	if u != nil {
		return u.(*uuid.UUID)
	}
	id := uuid.New()
	return &id
}
