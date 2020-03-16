package env

import (
	"os"
)

// Env は環境変数をまとめた構造体
type Env struct {
	Region  string
	Salt    string
	BaseURL string
	Env     string
}

// Config は環境変数をまとめた構造体を返す
func Config() Env {
	return Env{
		Region:  os.Getenv("REGION"),
		Salt:    os.Getenv("SALT"),
		BaseURL: os.Getenv("BASE_URL"),
		Env:     os.Getenv("ENV"),
	}
}
