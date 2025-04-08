package cache

import (
	"math/rand"
	"regexp"
	"strings"
	"time"
)

// safeValue get not nil value.
func safeValue[T any](value *T, fallback T) T {
	if value == nil {
		return fallback
	}
	return *value
}

// cacheKey generate normalized cache key.
func cacheKey(prefix string, keys ...string) string {
	prefix = slugify(prefix)
	key := slugify(keys...)

	if prefix != "" {
		return prefix + ":" + key
	}

	return key
}

// slugify make slug-format-text from strings.
func slugify(keys ...string) string {
	rxChars := regexp.MustCompile("[^a-zA-Z0-9\\-]")
	rxSpaces := regexp.MustCompile(`\s+`)
	rxDashes := regexp.MustCompile(`\-+`)
	content := strings.Join(keys, "-")
	content = rxChars.ReplaceAllString(content, "")
	content = rxSpaces.ReplaceAllString(content, "-")
	content = rxDashes.ReplaceAllString(content, "-")
	return content
}

// randomString generate a random string.
func randomString(n uint, letters string) (string, error) {
	if len(letters) == 0 {
		return "", nil
	}

	randomer := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, n)

	for i := range bytes {
		bytes[i] = letters[randomer.Intn(len(letters))]
	}

	return string(bytes), nil
}
