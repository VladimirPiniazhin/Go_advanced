package link

import (
	"go/order-api/internals/stat"
	"math/rand"

	"gorm.io/gorm"
)

type Link struct {
	gorm.Model
	Url   string      `json:"url"`
	Hash  string      `json:"hash" gorm:"uniqueIndex"`
	Stats []stat.Stat `gorm:"constraint:OnUpdate:CASCADE, OnDelete:SET NULL"`
}

func NewLink(url string, linkRepo *LinkRepository) *Link {
	hash := ""
	for {
		hash = RandsStringRunes(6)
		_, err := linkRepo.GetByHash(hash)
		if err != nil {
			break
		}
	}
	return &Link{
		Url:  url,
		Hash: hash,
	}
}

var letterRunes = []rune("qwertyuioplkjhgfdszxcbnmQWERTYUIOPLKJHGFDSAZXCVBNM")

func RandsStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
