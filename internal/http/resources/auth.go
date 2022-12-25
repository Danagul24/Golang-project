package resources

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	lru "github.com/hashicorp/golang-lru"
	"net/http"
	"project/internal/models"
	"project/internal/pkg/auth"
	"project/internal/store"
	"time"
)

const (
	accessTokenTTL  = 2 * time.Hour
	refreshTokenTTL = 168 * time.Hour
)

type AuthResource struct {
	store        store.Store
	cache        *lru.TwoQueueCache
	tokenManager auth.TokenManager
}

func NewAuthResource(store store.Store, cache *lru.TwoQueueCache, tokenManager auth.TokenManager) *AuthResource {
	return &AuthResource{
		store:        store,
		cache:        cache,
		tokenManager: tokenManager,
	}
}

func (a *AuthResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/login", a.LoginUser)
	return r
}

func (a *AuthResource) LoginUser(w http.ResponseWriter, r *http.Request) {
	user := new(models.LogInDTO)

	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown error: %v", err)
		return
	}

	u, err := a.store.Users().ByEmail(r.Context(), user.Email)
	if err != nil || !u.ComparePassword(user.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Incorrect email or password")
		return
	}

	tokens, err := a.CreateSession(&models.AuthorizedInfo{
		Id:   u.ID,
		Role: *u.Role,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unknown error: %v", err)
		return
	}
	render.JSON(w, r, tokens)

}

func (a *AuthResource) CreateSession(userInfo *models.AuthorizedInfo) (*models.Tokens, error) {
	var token models.Tokens
	var err error

	if token.AccessToken, err = a.tokenManager.NewJWT(userInfo, accessTokenTTL); err != nil {
		return nil, err
	}

	token.RefreshToken, err = a.tokenManager.NewRefreshToken()

	if err != nil {
		return nil, err
	}

	session := models.Session{
		RefreshToken: token.RefreshToken,
		ExpiresAt:    time.Now().Add(refreshTokenTTL),
	}

	a.cache.Add(userInfo.Id, session)
	return &token, err
}

func (a *AuthResource) RefreshTokens(userInfo *models.AuthorizedInfo, refreshToken string) (*models.Tokens, error) {
	tokenRaw, ok := a.cache.Get(userInfo.Id)
	if !ok {
		return nil, errors.New("Cache error : user not registered")
	}

	token, ok := tokenRaw.(models.Session)

	if !ok {
		return nil, errors.New("refresh token err")
	}

	if token.RefreshToken != refreshToken && time.Now().Unix() > token.ExpiresAt.Unix() {
		return nil, errors.New("refresh token err expired")
	}

	return a.CreateSession(userInfo)
}
