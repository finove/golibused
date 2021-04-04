package randkey

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"runtime"
	"strings"
	"time"
)

var numberUpperEncode = base32.NewEncoding("0123456789ABCDEFGHJKLMNPQRSTUVWX").WithPadding(base32.NoPadding)
var numberLowerEncode = base32.NewEncoding("0123456789abcdefghijkmnpqrstuvwx").WithPadding(base32.NoPadding)
var numberPassEncode = base64.NewEncoding("0123456789abcdefghijkmnpqrstuvwx!@#$%^&*()ABCDEFGHJKLMNPQRSTUVWX").WithPadding(base64.NoPadding)

func randNum() (num uint16) {
	var err error
	var b = make([]byte, 2)
	if _, err = rand.Read(b); err == nil {
		num = binary.BigEndian.Uint16(b)
	} else {
		num = uint16(time.Now().UnixNano() % 100000)
	}
	return
}

func randPuts(ctx context.Context, codes []byte, num int) (out chan byte) {
	out = make(chan byte, 1)
	if len(codes) == 0 {
		num = 0
	}
	go func(num int) {
		for i := 0; i < num; i++ {
			out <- codes[int(randNum())%len(codes)]
			runtime.Gosched()
		}
		<-ctx.Done()
	}(num)
	return
}

func GeneratePassword(counts ...int) (pass string) {
	var ssrc0 = "abcdefghijkmnpqrstuvwxABCDEFGHJKLMNPQRSTUVWX"
	var ssrc1 = "0123456789"
	var ssrc2 = "!@#$%^&*()"
	if len(counts) < 3 {
		counts = append(counts, 0, 0, 0)
	}
	pass = GenPassword([]string{ssrc0, ssrc1, ssrc2}, counts)
	return
}

func GenPassword(src []string, counts []int) (pass string) {
	var out []byte
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if len(src) < 3 {
		src = append(src, "", "", "")
	}
	if len(counts) < 3 {
		counts = append(counts, 0, 0, 0)
	}
	src0 := randPuts(ctx, []byte(src[0]), counts[0])
	src1 := randPuts(ctx, []byte(src[1]), counts[1])
	src2 := randPuts(ctx, []byte(src[2]), counts[2])
	for {
		select {
		case c0 := <-src0:
			out = append(out, c0)
		case c1 := <-src1:
			out = append(out, c1)
		case c2 := <-src2:
			out = append(out, c2)
		}
		if len(out) >= (counts[0] + counts[1] + counts[2]) {
			break
		}
	}
	pass = string(out)
	return
}

// NumbersOnly 只返回数字
func NumbersOnly(count int) (code string) {
	var nums []string
	if count <= 0 {
		return
	}
	for i := 0; i < count; i += 5 {
		nums = append(nums, fmt.Sprintf("%05d", randNum()))
	}
	code = strings.Join(nums, "")
	if len(code) > count {
		code = code[:count]
	}
	return
}

// NumberUpper 数字和大写字母
func NumberUpper(count int) (code string) {
	var err error
	var b = make([]byte, count)
	if _, err = rand.Read(b); err != nil {
		return
	}
	code = numberUpperEncode.EncodeToString(b)
	if len(code) > count {
		code = code[:count]
	}
	return
}

// NumberLower 数字和小写字母
func NumberLower(count int) (code string) {
	var err error
	var b = make([]byte, count)
	if _, err = rand.Read(b); err != nil {
		return
	}
	code = numberLowerEncode.EncodeToString(b)
	if len(code) > count {
		code = code[:count]
	}
	return
}

// NumberPass 数字和大写字母、特殊字符
func NumberPass(count int) (code string) {
	var err error
	var b = make([]byte, count)
	if _, err = rand.Read(b); err != nil {
		return
	}
	code = numberPassEncode.EncodeToString(b)
	if len(code) > count {
		code = code[:count]
	}
	return
}
