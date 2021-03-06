package helpers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/steinfletcher/apitest"

	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/system/types"
)

var (
	jwtHandler auth.TokenHandler
)

func InitAuth() {
	if jwtHandler != nil {
		return
	}

	var err error
	jwtHandler, err = auth.JWT(string(rand.Bytes(32)), 10)
	if err != nil {
		panic(err)
	}
}

func BindAuthMiddleware(r chi.Router) {
	r.Use(
		jwtHandler.HttpVerifier(),
		jwtHandler.HttpAuthenticator(),
	)
}

func ReqHeaderAuthBearer(user *types.User) apitest.Intercept {
	return func(req *http.Request) {
		req.Header.Set("Authorization", "Bearer "+jwtHandler.Encode(user))
	}
}
