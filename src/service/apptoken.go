package service

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/devrajsinghrawat/invite_app/src/logging"
	"github.com/devrajsinghrawat/invite_app/src/middleware"
	"github.com/devrajsinghrawat/invite_app/src/model"
	random "github.com/devrajsinghrawat/invite_app/src/rand"
	"github.com/devrajsinghrawat/invite_app/src/repo"
	httputils "github.com/devrajsinghrawat/invite_app/src/utils"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// AppTokenService contains logger and repo
type AppTokenService struct {
	repo   repo.IAppToken
	logger *zap.Logger
}

type validationTokenRequest struct {
	AppToken string
}

// NewAppTokenService returns the app token service instance
func NewAppTokenService(repository repo.IAppToken) *AppTokenService {
	log := logging.GetLogger().Named("apptoken_service")
	return &AppTokenService{
		repo:   repository,
		logger: log,
	}
}

// GenToken generate the new app token
func (s AppTokenService) GenToken(w http.ResponseWriter, req *http.Request) {
	s.logger.Info("Getting the token")
	token := random.String(rand.Intn(6) + 6) // Random a length between 6 to 12
	ctx := req.Context()
	authToken, err := middleware.GetAuthorizationToken(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.logger.Debug("Found token in request", zap.Any("token", authToken))
	if authToken.Role != "ADMIN" && authToken.Username == "" {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	days, err := strconv.Atoi(os.Getenv("EXPIRE_IN_DAYS"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.repo.SaveAppToken(&model.AppToken{
		Token:    token,
		IsActive: true,
		ExpDate:  time.Now().AddDate(0, 0, days),
		Username: authToken.Username,
	})
	s.logger.Debug("App Token save into database", zap.String("token", token))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_ = httputils.SendSuccessHeader(w, token)
}

// ValidateToken check if given token is valid
func (s AppTokenService) ValidateToken(w http.ResponseWriter, req *http.Request) {
	s.logger.Info("Validated token")
	appToken := mux.Vars(req)["appToken"]
	if appToken == "" {
		http.Error(w, "invalid app token", http.StatusBadRequest)
		return
	}
	s.logger.Debug("Found app token in the request", zap.String("token", appToken))
	tokenlen := len(appToken)
	if tokenlen < 6 || tokenlen > 12 {
		http.Error(w, "invalid app token", http.StatusBadRequest)
		return
	}
	data, err := s.repo.GetAppToken(&model.AppToken{Token: appToken})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.logger.Info("Found app token in database")
	if !data.IsActive {
		http.Error(w, "token is inactive", http.StatusBadRequest)
		return
	}
	s.logger.Info("Token is valid")
	if time.Now().After(data.ExpDate) {
		data.IsActive = false
		_, _ = s.repo.UpdateAppToken(&data)
		http.Error(w, "app token has expired", http.StatusBadRequest)
		return
	}
	_ = httputils.SendSuccessHeader(w, data.IsActive)
}

// InvalidateToken deactivate the token
func (s AppTokenService) InvalidateToken(w http.ResponseWriter, req *http.Request) {
	s.logger.Info("Invalidate token")
	var vtr validationTokenRequest
	var err error
	err = json.NewDecoder(req.Body).Decode(&vtr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.logger.Debug("Found token in the request", zap.Any("token", vtr))
	data, err := s.repo.GetAppToken(&model.AppToken{Token: vtr.AppToken})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.logger.Info("Found token from database")
	if data.IsActive {
		data.IsActive = false
		_, err = s.repo.UpdateAppToken(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	_ = httputils.SendSuccessHeader(w, data)
}

// GetAllAppToken returns all the active and inactive tokens
func (s AppTokenService) GetAllAppToken(w http.ResponseWriter, req *http.Request) {
	s.logger.Info("Getting all token from database")
	data, err := s.repo.GetAllAppToken()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_ = httputils.SendSuccessHeader(w, data)
}
