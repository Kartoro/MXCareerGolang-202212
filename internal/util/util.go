package util

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"path"
	"strconv"
	"strings"
	"time"
)

func createRandomNumber() string {
	return fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
}

func getCurTime() string {
	now := time.Now()
	dateStr := fmt.Sprintf("%d%d%d%d%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute())
	return dateStr
}

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func SHA256(str string) string {
	s := sha256.Sum256([]byte(str))
	return fmt.Sprintf("%x", s)
}

func GetToken(userName string) string {
	return MD5(userName + getCurTime() + createRandomNumber())
}

func CheckAndCreateFileName(oldName string) (newName string, isLegal bool) {
	ext := path.Ext(oldName)
	if strings.ToLower(ext) == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" {
		h := md5.New()
		h.Write([]byte(oldName + strconv.FormatInt(time.Now().Unix(), 10)))
		newName = hex.EncodeToString(h.Sum(nil)) + ext
		isLegal = true
	}
	return newName, isLegal
}
