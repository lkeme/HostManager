package util

import "github.com/google/uuid"

// GenerateKey 生成Key
func GenerateKey() string {
	return uuid.New().String()[9:23]
}
