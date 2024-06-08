package chat

import (
	"crypto/rand"
	"math/big"
)

var promoStrings = [10]string{
	"\n\n------\n\n💥📢 The Absolute Basstards - Drum&Bass label.\nhttps://t.me/+bS_eIEhkLuZkMDYy",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
}

var maxIndex = big.NewInt(int64(len(promoStrings)))

var lastRandomIndex int64 = -1

func getRandomIndex() int64 {
	randomIndex, _ := rand.Int(rand.Reader, maxIndex)
	return randomIndex.Int64()
}

// GetPromoString - Возвращает случайную строку из массива promoStrings без повторений
func GetPromoString() string {
	currentIndex := getRandomIndex()
	for currentIndex == lastRandomIndex {
		currentIndex = getRandomIndex()
	}
	lastRandomIndex = currentIndex
	return promoStrings[lastRandomIndex]
}
