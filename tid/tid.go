package tid

import (
	"context"
	"crypto/rand"
	"math/big"
)

type contextKey string

const TidContextKey contextKey = "logger.tid"

const tidAlphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const tidLength = 5

func TIDFromContext(ctx context.Context) string {
	if tid, ok := ctx.Value(TidContextKey).(string); ok {
		return tid
	}
	return ""
}

func WithContextTID(ctx context.Context) (context.Context, string) {
	tid := TIDFromContext(ctx)
	if tid == "" {
		tid = NewTID()
	}
	return context.WithValue(ctx, TidContextKey, tid), tid
}

func NewTID() string {
	return newRandomToken(tidLength)
}

func NewDerivedTID(parent string) string {
	if parent == "" {
		return NewTID()
	}

	return parent + "-" + newRandomToken(tidLength)
}

func newRandomToken(length int) string {
	buffer := make([]byte, length)

	for i := 0; i < length; i++ {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(tidAlphabet))))
		if err != nil {
			return "00000"
		}
		buffer[i] = tidAlphabet[index.Int64()]
	}

	return string(buffer)
}
