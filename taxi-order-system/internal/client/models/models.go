package models

type ClientType int8

const (
	Mobile ClientType = iota
	Browser
)

type Client struct {
	Id   int64
	Imei string
	Type ClientType
}

type Trip struct {
	ClientId string
	From     string
	To       string
}
