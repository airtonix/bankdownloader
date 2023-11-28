package clients

import "time"

type SecretsResolver interface {
	GetPassword(path string) (string, error)
	GetUsername(path string) (string, error)
	GetOtp(path string, timestampe time.Time) (string, error)
}
