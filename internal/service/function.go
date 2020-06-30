package service

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/vladimir-kopaliani/auth_example/internal/token"
)

func (s service) AuthorizeUser(ctx context.Context, guid, accessToken string) (string, error) {
	refreshToken := generateToken(30)

	// encrypt token
	encryptedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return "", err
	}

	tkn := token.Token{
		GUID:             guid,
		Access:           accessToken,
		Refresh:          string(encryptedToken),
		CreatedAt:        time.Now(),
		RefreshExpiredAt: time.Now().Add(30 * 24 * time.Hour),
	}

	err = s.repository.SaveToken(ctx, &tkn)
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func (s service) RefreshUserToken(ctx context.Context, guid, oldAccessToken, oldRefreshToken, newAccessToken string) (string, error) {
	refreshToken := generateToken(30)

	// encrypt token
	encryptedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return "", err
	}

	tkn, err := s.repository.GetToken(ctx, guid, oldAccessToken)
	if err != nil {
		log.Println(err)
		return "", err
	}
	if tkn == nil || tkn.RefreshExpiredAt.Before(time.Now()) {
		return "", errors.New("token is expired")
	}
	// check refresh token
	if err = bcrypt.CompareHashAndPassword([]byte(tkn.Refresh), []byte(oldRefreshToken)); err != nil {
		log.Println(err)
		return "", err
	}

	err = s.repository.ChangeToken(ctx, guid, oldAccessToken, string(encryptedToken), newAccessToken)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return refreshToken, nil
}

func (s service) RemoveRefreshUserToken(ctx context.Context, guid, accessToken, refreshToken string) error {
	tkn, err := s.repository.GetToken(ctx, guid, accessToken)
	if err != nil {
		log.Println(err)
		return err
	}
	if tkn == nil || tkn.RefreshExpiredAt.Before(time.Now()) {
		return errors.New("token is expired")
	}
	// check refresh token
	if err = bcrypt.CompareHashAndPassword([]byte(tkn.Refresh), []byte(refreshToken)); err != nil {
		log.Println(err)
		return err
	}

	err = s.repository.RemoveToken(ctx, guid, accessToken)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s service) RemoveAllTokens(ctx context.Context, guid string) error {
	err := s.repository.RemoveAllTokens(ctx, guid)
	if err != nil {
		return err
	}

	return nil
}

var charatersForToken = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func generateToken(length int) string {
	rand.Seed(time.Now().UnixNano())

	str := make([]rune, length)

	for i := range str {
		str[i] = charatersForToken[rand.Intn(len(charatersForToken))]
	}

	return string(str)
}
