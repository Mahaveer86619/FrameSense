package utils

import (
	"errors"
	"sync"

	"github.com/Mahaveer86619/FrameSense/pkg/config"
	"github.com/speps/go-hashids/v2"
)

var (
	hasher *hashids.HashID
	once   sync.Once
)

func getHasher() *hashids.HashID {
	once.Do(func() {
		hd := hashids.NewData()
		hd.Salt = config.AppConfig.ID_SALT
		hd.MinLength = 10

		// Initialize the hasher
		var err error
		hasher, err = hashids.NewWithData(hd)
		if err != nil {
			panic("failed to initialize ID masker: " + err.Error())
		}
	})
	return hasher
}

func MaskID(id uint) string {
	h, _ := getHasher().Encode([]int{int(id)})
	return h
}

func UnmaskID(hash string) (uint, error) {
	ids, err := getHasher().DecodeWithError(hash)
	if err != nil || len(ids) == 0 {
		return 0, errors.New("invalid ID format")
	}
	return uint(ids[0]), nil
}
