package tencentcloud

import "errors"

var ErrRecordNotFound = errors.New("record not found")
var ErrNotValid = errors.New("returned value is not valid")

type Provider struct {
	SecretId  string
	SecretKey string
}

type CreateModifyRecordRequest struct {
	Domain     string `json:"Domain"`
	SubDomain  string `json:"SubDomain,omitempty"`
	RecordType string `json:"RecordType,omitempty"`
	RecordLine string `json:"RecordLine,omitempty"`
	Value      string `json:"Value,omitempty"`
	TTL        int64  `json:"TTL,omitempty"`
	RecordId   uint64 `json:"RecordId,omitempty"`
}

type FindRecordRequest struct {
	Domain     string `json:"Domain"`
	RecordType string `json:"RecordType,omitempty"`
	RecordLine string `json:"RecordLine,omitempty"`
	Subdomain  string `json:"Subdomain,omitempty"`
	Limit      int64  `json:"Limit,omitempty"`
}

type DeleteRecordRequest struct {
	Domain   string `json:"Domain"`
	RecordId string `json:"RecordId"`
}

type Response struct {
	Response ResponseData `json:"Response"`
}

type ResponseData struct {
	RecordList []RecordInfo `json:"RecordList,omitempty"`
	RecordId   uint64       `json:"RecordId,omitempty"`
	Error      *ErrorInfo   `json:"Error,omitempty"`
}

type RecordInfo struct {
	RecordId int64  `json:"RecordId"`
	Type     string `json:"Type"`
	Name     string `json:"Name"`
	Value    string `json:"Value"`
	TTL      int64  `json:"TTL"`
}

type ErrorInfo struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
}
