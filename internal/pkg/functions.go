package pkg

import (
	"context"
	"fmt"
	"net/http"
	"project/internal/models"
)

func IsUserAdmin(ctx context.Context, w http.ResponseWriter) bool {
	if userInfo := ctx.Value(CtxKeyUser).(*models.AuthorizedInfo); userInfo.Role != models.Admin {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Println(w, "err: Insufficient right to access data")
		return false
	}
	return true
}
