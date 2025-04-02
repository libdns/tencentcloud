package tencentcloud

import "errors"

var ErrRecordNotFound = errors.New("record not found")
var ErrNotValid = errors.New("returned value is not valid")

type Provider struct {
	SecretId  string
	SecretKey string

	id uint64
}
