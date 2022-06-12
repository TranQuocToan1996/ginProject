package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"image"
	"reflect"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	minCost = 10
)

//

// map data from interface to struct
func MapItoM(i interface{}, s interface{}) error {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errors.New("Map() expects a pointer to a struct")
	}
	v = v.Elem()
	t := v.Type()
	iValue := reflect.ValueOf(i)
	if iValue.Kind() != reflect.Map {
		return errors.New("Map() expects an interface{} of type map")
	}
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).CanSet() {
			f := iValue.MapIndex(reflect.ValueOf(t.Field(i).Name))
			if !f.IsValid() {
				continue
			}
			v.Field(i).Set(f)
		}
	}
	return nil
}

//resize image
// func Resize(img image.Image, width, height int) *image.NRGBA {
// 	srcBounds := img.Bounds()
// 	srcAspectRatio := float64(srcBounds.Dx()) / float64(srcBounds.Dy())
// 	dstBounds := image.Rect(0, 0, width, height)
// 	dstAspectRatio := float64(dstBounds.Dx()) / float64(dstBounds.Dy())
// 	var dst image.Image
// 	if srcAspectRatio > dstAspectRatio {
// 		dstBounds.Min.X = dstBounds.Min.Y * srcAspectRatio
// 		dstBounds.Max.X = dstBounds.Max.Y * srcAspectRatio
// 	} else {
// 		dstBounds.Min.Y = dstBounds.Min.X / srcAspectRatio
// 		dstBounds.Max.Y = dstBounds.Max.X / srcAspectRatio
// 	}
// 	dst = image.NewNRGBA(dstBounds)
// 	draw.Draw(dst, dstBounds, img, srcBounds.Min, draw.Src)
// 	return dst.(*image.NRGBA)
// }

// func Resize(img image.Image, width, height int) image.Image {
// 	return resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
// }

// decode base64 to image
func DecodeBase64Image(data []byte) (image.Image, error) {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(string(data)))
	img, _, err := image.Decode(reader)
	return img, err
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// loop through struct counting the number of fields
func CountFieldsStruct(s interface{}) int {
	count := 0
	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Interface() != nil {
			count++
		}
	}
	return count
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

// NewSHA256 ...
func NewSHA256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

// HashPassword using bcrypt
func HashPassword(password string, cost int) (string, error) {
	if cost < minCost {
		cost = minCost
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}

func ValidPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	// bcrypt.ErrMismatchedHashAndPassword
	return err == nil
}

/*
A bcrypt cost of 6 means 64 rounds (2^6 = 64).
 1/23/2014  Intel Core i7-2700K CPU @ 3.50 GHz

| Cost | Iterations        |    Duration |
|------|-------------------|-------------|
|  8   |    256 iterations |     38.2 ms | <-- minimum allowed by BCrypt
|  9   |    512 iterations |     74.8 ms |
| 10   |  1,024 iterations |    152.4 ms |
| 11   |  2,048 iterations |    296.6 ms |
| 12   |  4,096 iterations |    594.3 ms |
| 13   |  8,192 iterations |  1,169.5 ms |
| 14   | 16,384 iterations |  2,338.8 ms |<-- current default (BCRYPT_COST=10)
| 15   | 32,768 iterations |  4,656.0 ms |
| 16   | 65,536 iterations |  9,302.2 ms |
*/
