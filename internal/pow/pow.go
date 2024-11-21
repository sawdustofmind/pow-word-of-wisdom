package pow

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strconv"
)

const zeroByte = '0'

func sha1Hash(data string) string {
	h := sha1.New()
	h.Write([]byte(data))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func IsHashCorrect(hash string, zerosCount int) bool {
	if zerosCount > len(hash) {
		return false
	}

	for _, ch := range hash[:zerosCount] {
		if ch != zeroByte {
			return false
		}
	}
	return true
}

func ComputeHashcash(challenge string, maxIterations int, zerosCount int) (string, int, error) {
	counter := 0

	for counter <= maxIterations {
		hash := sha1Hash(challenge + strconv.Itoa(counter))

		if IsHashCorrect(hash, zerosCount) {
			return hash, counter, nil
		}

		counter++
	}

	return "", 0, fmt.Errorf("max iterations exceeded")
}

func GenerateChallenge(length int) string {
	bts := make([]byte, length)
	_, _ = rand.Read(bts)
	challenge := hex.EncodeToString(bts)
	return challenge
}

func IsChallengeSolved(challenge string, answer int, zerosCount int) bool {
	hash := sha1Hash(challenge + strconv.Itoa(answer))
	return IsHashCorrect(hash, zerosCount)
}
