package service

import (
	argon2 "github.com/alexedwards/argon2id"
	"go.uber.org/zap"
)

func Encryption(rawPwd string) (string, error) {
	hash, err := argon2.CreateHash(rawPwd, argon2.DefaultParams)
	if err != nil {
		zap.L().Error("Error while creating hash", zap.Error(err))
		return "", err
	}
	return hash, err
}

func CheckPwd(rawPwd string, hash string) bool {
	match, err := argon2.ComparePasswordAndHash(rawPwd, hash)

	if err != nil {
		zap.L().Error("check password error", zap.Error(err))
		return false
	}
	return match
}
