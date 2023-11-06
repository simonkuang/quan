package codec

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/kenshaw/baseconv"
	"github.com/simonkuang/quan/src/config"
)

type ShortUrlEntity struct {
	Url          string     `json:"url"`
	ShortUrl     string     `json:"short"`
	SecretPrefix string     `json:"prefix"`
	Time         *time.Time `json:"time"`
}

func (s *ShortUrlEntity) String() string {
	jsonStr, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return ""
	}
	return string(jsonStr)
}

func Encode(url string, startPosition int, length int) (string, error) {
	hash := sha256.New()
	hash.Write([]byte(config.SecretPrefix + url))
	hashSum := hash.Sum(nil)
	hexStr := fmt.Sprintf("%x", hashSum)
	outputBase := baseconv.Digits62
	if config.CharRange == 36 {
		outputBase = baseconv.Digits36
	}
	hashStr, err := baseconv.Convert(hexStr, baseconv.DigitsHex, outputBase)
	if err != nil {
		return "", err
	}
	if (startPosition+1)*length >= len(hashStr) {
		return "", errors.New("FATAL: Highly Conflict")
	}
	// fmt.Printf("%d\t%d\t%s\n", startPosition, config.Length, hashStr)
	return hashStr[startPosition*length : (startPosition+1)*length], nil
}
