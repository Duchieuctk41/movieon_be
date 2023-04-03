package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"movieon_be/pkg/model"
)

func HasRole(headerRole string, role uint32) bool {
	return ConvertStringToUint32(headerRole)&role > 0
}

func ParseIDFromUri(c *gin.Context) *uuid.UUID {
	tID := model.UriParse{}
	if err := c.ShouldBindUri(&tID); err != nil {
		_ = c.Error(err)
		return nil
	}
	if len(tID.ID) == 0 {
		_ = c.Error(fmt.Errorf("error: Empty when parse ID from URI"))
		return nil
	}
	if id, err := uuid.Parse(tID.ID[0]); err != nil {
		_ = c.Error(err)
		return nil
	} else {
		return &id
	}
}
func GetIdFromUri(c *gin.Context) *string {
	tID := model.UriParse{}
	if err := c.ShouldBindUri(&tID); err != nil {
		_ = c.Error(err)
		return nil
	}
	if len(tID.ID) == 0 {
		_ = c.Error(fmt.Errorf("error: Empty when parse ID from URI"))
		return nil
	}
	return &tID.ID[0]
}

func ParseUUID(in string) uuid.UUID {
	id, err := uuid.Parse(in)
	if err != nil {
		return uuid.Nil
	}
	return id
}

func ValidPhoneFormat(phone string) bool {
	internationalPhone := regexp.MustCompile("^\\+[1-9]\\d{1,14}$")
	if !internationalPhone.MatchString(phone) {
		return false
	}
	return true
}

func ValidOtpFormat(otp string) bool {
	re := regexp.MustCompile("([0-9]{6})$")
	if !re.MatchString(otp) {
		return false
	}
	return true
}

func GetRandNum(n int) string {
	const letters = "0123456789"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return ""
		}
		ret[i] = letters[num.Int64()]
	}
	return string(ret)
}

func GetTime(in int) time.Time {
	return time.Now().Add(time.Duration(in) * time.Second)
}

func CheckOTPIsExpired(in time.Time) bool {
	return in.Unix() < time.Now().Unix()
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func CurrentUser(c *http.Request) (uuid.UUID, error) {
	userIdStr := c.Header.Get("x-user-id")
	if strings.Contains(userIdStr, "|") {
		userIdStr = strings.Split(userIdStr, "|")[0]
	}
	res, err := uuid.Parse(userIdStr)
	if err != nil {
		return uuid.Nil, err
	}
	return res, nil
}

func Roles(r *http.Request) byte {
	s := r.Header.Get("x-user-roles")
	userRoles, err := strconv.ParseInt(s, 10, 8)
	if err != nil {
		return 0
	}
	return byte(userRoles)
}

func ConvertUnixMiliToTime(in int64) time.Time {
	t := time.UnixMilli(in)
	fmt.Println(t)
	return t
}

func DayTime(i time.Time) *time.Time {
	if i.IsZero() {
		return nil
	} else {
		return &i
	}
}

func IsErrNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func ConvertStringToUint32(req string) (rs uint32) {
	tmpRole, _ := strconv.Atoi(req)
	return uint32(tmpRole)
}

func ConvertFloatToString(req float64) (rs string) {
	s := fmt.Sprintf("%f", req)
	return s
}

func CurrentFunctionName(level int) string {
	pc, _, _, _ := runtime.Caller(1 + level)
	strArr := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	if len(strArr) > 0 {
		return strArr[len(strArr)-1]
	}
	return ""
}

func GetCurrentCaller(caller interface{}, level int) string {
	strArr := strings.Split(reflect.TypeOf(caller).String(), ".")
	if caller != nil && len(strArr) > 0 {
		return fmt.Sprintf("%s.%s", strArr[len(strArr)-1], CurrentFunctionName(1+level))
	}
	return CurrentFunctionName(1)
}

func GenerateRandomString(n int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return ""
		}
		ret[i] = letters[num.Int64()]
	}
	return string(ret)
}

func SliceContains(value, slice interface{}) error {
	valueInterface := reflect.ValueOf(slice)
	if valueInterface.Kind() == reflect.Slice {
		for i := 0; i < valueInterface.Len(); i++ {
			if value == valueInterface.Index(i).Interface() {
				return nil
			}
		}
	} else {
		return fmt.Errorf("invalid type for slice")
	}
	return fmt.Errorf("slice doesn't contain value %v", value)
}

func CheckExistPath(path string) (err error) {
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, 0755)
		if err != nil {
			return err
		}
	}

	return err
}

func GetCurrentTime() time.Time {
	return time.Now().UTC()
}

func GetToday() time.Time {
	return time.Now().UTC().Truncate(24 * time.Hour)
}
