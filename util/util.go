package util

import (
	"bytes"
	srand "crypto/rand"
	"encoding/binary"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"
	"errors"
)


const (
	StatusUnKnown     = "NA"
)

//CreateRandomString create random string
func CreateRandomString() string {

	rand.Seed(time.Now().UnixNano())
	randInt := rand.Int63()
	randStr := strconv.FormatInt(randInt, 36)

	return randStr
}

// Encode encodes the values into ``URL encoded'' form
// ("acl&bar=baz&foo=quux") sorted by key.
func Encode(v url.Values) string {
	if v == nil {
		return ""
	}
	var buf bytes.Buffer
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		prefix := url.QueryEscape(k)
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(prefix)
			if v != "" {
				buf.WriteString("=")
				buf.WriteString(url.QueryEscape(v))
			}
		}
	}
	return buf.String()
}

func GetGMTime() string {
	return time.Now().UTC().Format(http.TimeFormat)
}

//

func randUint32() uint32 {
	return randUint32Slice(1)[0]
}

func randUint32Slice(c int) []uint32 {
	b := make([]byte, c*4)

	_, err := srand.Read(b)

	if err != nil {
		// fail back to insecure rand
		rand.Seed(time.Now().UnixNano())
		for i := range b {
			b[i] = byte(rand.Int())
		}
	}

	n := make([]uint32, c)

	for i := range n {
		n[i] = binary.BigEndian.Uint32(b[i*4 : i*4+4])
	}

	return n
}

func toByte(n uint32, st, ed byte) byte {
	return byte(n%uint32(ed-st+1) + uint32(st))
}

func toDigit(n uint32) byte {
	return toByte(n, '0', '9')
}

func toLowerLetter(n uint32) byte {
	return toByte(n, 'a', 'z')
}

func toUpperLetter(n uint32) byte {
	return toByte(n, 'A', 'Z')
}

type convFunc func(uint32) byte

var convFuncs = []convFunc{toDigit, toLowerLetter, toUpperLetter}

// tools for generating a random ECS instance password
// from 8 to 30 char MUST contain digit upper, case letter and upper case letter
// http://docs.aliyun.com/#/pub/ecs/open-api/instance&createinstance
func GenerateRandomECSPassword() string {

	// [8, 30]
	l := int(randUint32()%23 + 8)

	n := randUint32Slice(l)

	b := make([]byte, l)

	b[0] = toDigit(n[0])
	b[1] = toLowerLetter(n[1])
	b[2] = toUpperLetter(n[2])

	for i := 3; i < l; i++ {
		b[i] = convFuncs[n[i]%3](n[i])
	}

	s := make([]byte, l)
	perm := rand.Perm(l)
	for i, v := range perm {
		s[v] = b[i]
	}

	return string(s)

}


func SimpleLoopCall(attempts AttemptStrategy, f func() bool) error {
	for attempt := attempts.Start(); attempt.Next(); {

		stop := f()

		if stop {
			return nil
		}

		if attempt.HasNext() {
			continue;
		}

		return errors.New("timeout execution ")
	}
	panic("unreachable")
}

func LoopCall(attempts AttemptStrategy,f func() (bool,interface{},error))(interface{}, error){

	for attempt := attempts.Start(); attempt.Next(); {
		needStop,status,err := f()

		if(err != nil) {
			return nil, err
		}
		if(needStop){
			return status,nil;
		}

		if attempt.HasNext() {
			continue;
		}
		return nil,errors.New("timeout execution ")
	}
	panic("unreachable")
}
