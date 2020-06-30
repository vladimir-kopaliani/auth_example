package serverhttp

import (
	"encoding/base64"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	jwtKey = []byte("secret_salt")

	accessTokenName           = "token_a"
	accessTokenExpirationTime = 30 * time.Minute

	refreshTokenName           = "token_r"
	refreshTokenExpirationTime = 30 * 24 * time.Hour
)

func (s *serverHTTP) auth(w http.ResponseWriter, r *http.Request) {
	// retrive guid from url query
	queries := r.URL.Query()
	guid := queries["guid"]
	if len(guid) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("guid is not provided"))
		return
	}

	// sign access token
	expirationTime := time.Now().Add(accessTokenExpirationTime)

	token := jwt.NewWithClaims(jwt.SigningMethodHS512,
		jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	)

	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// check credentails
	refreshToken, err := s.service.AuthorizeUser(r.Context(), guid[0], signedToken)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode refresh token
	refreshTokenDecoded := base64.StdEncoding.EncodeToString([]byte(refreshToken))

	// set refresh login in cookie
	http.SetCookie(w, &http.Cookie{
		Name:    refreshTokenName,
		Value:   refreshTokenDecoded,
		Expires: time.Now().Add(refreshTokenExpirationTime),
	})

	// set access login in cookie
	http.SetCookie(w, &http.Cookie{
		Name:    accessTokenName,
		Value:   signedToken,
		Expires: expirationTime,
	})
}

func (s *serverHTTP) refresh(w http.ResponseWriter, r *http.Request) {
	cookieAccessToken, err := r.Cookie(accessTokenName)
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cookieRefreshToken, err := r.Cookie(refreshTokenName)
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// retrive guid from url query
	queries := r.URL.Query()
	guid := queries["guid"]
	if len(guid) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("guid is not provided"))
		return
	}

	// set tokens in cookie
	expirationTime := time.Now().Add(accessTokenExpirationTime)

	token := jwt.NewWithClaims(jwt.SigningMethodHS512,
		jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	)

	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// decode refresh token
	refreshTokenDecoded, err := base64.StdEncoding.DecodeString(cookieRefreshToken.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	refreshToken, err := s.service.RefreshUserToken(r.Context(), guid[0], cookieAccessToken.Value, string(refreshTokenDecoded), signedToken)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// encode refresh token
	refreshTokenEncoded := base64.StdEncoding.EncodeToString([]byte(refreshToken))

	// set refresh login
	http.SetCookie(w, &http.Cookie{
		Name:    refreshTokenName,
		Value:   refreshTokenEncoded,
		Expires: time.Now().Add(refreshTokenExpirationTime),
	})

	// set access login
	http.SetCookie(w, &http.Cookie{
		Name:    accessTokenName,
		Value:   signedToken,
		Expires: expirationTime,
	})
}

func (s *serverHTTP) removeRefreshToken(w http.ResponseWriter, r *http.Request) {
	cookieAccessToken, err := r.Cookie(accessTokenName)
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cookieRefreshToken, err := r.Cookie(refreshTokenName)
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// retrive guid from url query
	queries := r.URL.Query()
	guid := queries["guid"]
	if len(guid) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("guid is not provided"))
		return
	}

	// decode refresh token
	refreshTokenDecoded, err := base64.StdEncoding.DecodeString(cookieRefreshToken.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.service.RemoveRefreshUserToken(r.Context(), guid[0], cookieAccessToken.Value, string(refreshTokenDecoded))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// remove refresh login
	http.SetCookie(w, &http.Cookie{
		Name:    refreshTokenName,
		Expires: time.Now().Add(-1 * time.Minute),
	})

	// remove access login
	http.SetCookie(w, &http.Cookie{
		Name:    accessTokenName,
		Expires: time.Now().Add(-1 * time.Minute),
	})
}

func (s *serverHTTP) removeAllResfreshToken(w http.ResponseWriter, r *http.Request) {
	// retrive guid from url query
	queries := r.URL.Query()
	guid := queries["guid"]
	if len(guid) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("guid is not provided"))
		return
	}

	err := s.service.RemoveAllTokens(r.Context(), guid[0])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// remove refresh login
	http.SetCookie(w, &http.Cookie{
		Name:    refreshTokenName,
		Expires: time.Now().Add(-1 * time.Minute),
	})

	// remove access login
	http.SetCookie(w, &http.Cookie{
		Name:    accessTokenName,
		Expires: time.Now().Add(-1 * time.Minute),
	})
}
