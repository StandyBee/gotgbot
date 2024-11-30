package server

import (
	"github.com/StandyBee/gotgbot/database/credis"
	pocket "github.com/StandyBee/pocketSDK"
	"log"
	"net/http"
	"strconv"
)

type AuthorizationServer struct {
	server          *http.Server
	pocketClient    *pocket.Client
	tokenRepository *credis.TokenRepository
	redirectURL     string
}

func NewAuthorizationServer(pocketClient *pocket.Client, redirectURL string, tokenRepository *credis.TokenRepository) *AuthorizationServer {
	return &AuthorizationServer{
		pocketClient:    pocketClient,
		tokenRepository: tokenRepository,
		redirectURL:     redirectURL,
	}
}

func (s *AuthorizationServer) Start() error {
	s.server = &http.Server{
		Addr:    ":80",
		Handler: s,
	}

	return s.server.ListenAndServe()
}

func (s *AuthorizationServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	chatIdParam := r.URL.Query().Get("chat_id")

	if chatIdParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chatId, err := strconv.ParseInt(chatIdParam, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestToken, err := s.tokenRepository.GetRequestToken(chatId)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	log.Println(requestToken)

	authResponse, err := s.pocketClient.Authorize(r.Context(), requestToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println(authResponse.AccessToken)
	err = s.tokenRepository.SaveAccessToken(chatId, authResponse.AccessToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("CHAT ID: %d, ACCESS TOKEN: %s, REQUEST TOKEN: %s", chatId, authResponse.AccessToken, requestToken)

	w.Header().Add("Location", s.redirectURL)
	w.WriteHeader(http.StatusMovedPermanently)
}
