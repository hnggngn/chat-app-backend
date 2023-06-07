package utils

import (
	"math/rand"
	"time"
)

// GetRandomImage from https://gopherize.me
func GetRandomImage() string {
	avatars := []string{
		"https://storage.googleapis.com/gopherizeme.appspot.com/gophers/30c621a657fb4a0bf4234e1f20f7ce91333fd712.png",
		"https://storage.googleapis.com/gopherizeme.appspot.com/gophers/2046af9c8e11b2cbb4b2645ade710820a25fdf5a.png",
		"https://storage.googleapis.com/gopherizeme.appspot.com/gophers/39d232350da7b7a14d6c2f77ca29e07e01621376.png",
		"https://storage.googleapis.com/gopherizeme.appspot.com/gophers/3c7ee5835c4757164348ea6f1632d98905eb8bf1.png",
		"https://storage.googleapis.com/gopherizeme.appspot.com/gophers/96382a6c0ebab94b2c25e825b243f5936a90ef0f.png",
		"https://storage.googleapis.com/gopherizeme.appspot.com/gophers/db248fdd04d02e221b1f72ea10c85ccba1797b18.png",
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))
	randomIndex := rand.Intn(len(avatars))

	return avatars[randomIndex]
}
