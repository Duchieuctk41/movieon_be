package valid

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func String(in *string) string {
	if in == nil {
		return ""
	}
	return *in
}

func Int(in *int) int {
	if in == nil {
		return 0
	}
	return *in
}

func Float64(in *float64) float64 {
	if in == nil {
		return 0
	}
	return *in
}

func Float32(in *float32) float32 {
	if in == nil {
		return 0
	}
	return *in
}

func Bool(in *bool) bool {
	if in == nil {
		return false
	}
	return *in
}

func Int32(in *int32) int32 {
	if in == nil {
		return 0
	}
	return *in
}

func StringArray(in *pq.StringArray) pq.StringArray {
	var strArr []string
	if in == nil {
		return strArr
	}
	return *in
}

func Int64(in *int64) int64 {
	if in == nil {
		return 0
	}
	return *in
}

func Byte(in *byte) byte {
	if in == nil {
		return 0
	}
	return *in
}

func UUID(req *uuid.UUID) uuid.UUID {
	if req == nil {
		return uuid.Nil
	}
	return *req
}

func DayTime(in *time.Time) time.Time {
	if in == nil {
		return time.Time{}
	}
	return *in
}

//------------------------------------POINTER------------------------------------------------------------------------------------

// String returns pointer to s.
func StringPointer(s string) *string {
	return &s
}

// Bool returns a pointer to b.
func BoolPointer(b bool) *bool {
	return &b
}

// Int returns a pointer to i.
func IntPointer(i int) *int {
	return &i
}

func Float64Pointer(i float64) *float64 {
	return &i
}

func DayTimePointer(i time.Time) *time.Time {
	if i.IsZero() {
		return nil
	} else {
		return &i
	}
}

func UUIDPointer(i uuid.UUID) *uuid.UUID {
	return &i
}

func StringToPointerUUID(s string) *uuid.UUID {
	sMarshal, _ := json.Marshal(s)
	sUUID, _ := uuid.FromBytes(sMarshal)
	return &sUUID
}

func NullString(in string) string {
	if in == "null" {
		return ""
	}
	return in
}
