package cache

import "time"

// VerificationCode defines the interface for managing verification codes in a cache.
type VerificationCode interface {
	// Set stores the verification code in the cache.
	Set(code string) error

	// Generate creates a random numeric code of the specified length and stores it in the cache.
	Generate(count uint) (string, error)

	// Clear removes the verification code from the cache.
	Clear() error

	// Get retrieves the verification code from the cache.
	Get() (string, error)

	// Validate checks if the provided code matches the stored verification code.
	Validate(code string) (bool, error)

	// Exists checks if the verification code exists in the cache.
	Exists() (bool, error)

	// TTL retrieves the time-to-live (TTL) of the verification code in the cache.
	TTL() (time.Duration, error)
}

// verification is the concrete implementation of the VerificationCode interface.
type verification struct {
	name  string
	ttl   time.Duration
	cache Cache
}

// NewVerification creates a new instance of the verification code.
func NewVerification(name string, ttl time.Duration, cache Cache) VerificationCode {
	return &verification{
		name:  "verify " + name,
		ttl:   ttl,
		cache: cache,
	}
}

func (v *verification) Set(code string) error {
	exists, err := v.cache.Update(v.name, code)
	if err != nil {
		return err
	}

	if !exists {
		return v.cache.Put(v.name, code, &v.ttl)
	}

	return nil
}

func (v *verification) Generate(count uint) (string, error) {
	code, err := randomString(count, "0123456789")
	if err != nil {
		return "", err
	}

	if err := v.Set(code); err != nil {
		return "", err
	}

	return code, nil
}

func (v *verification) Clear() error {
	return v.cache.Forget(v.name)
}

func (v *verification) Get() (string, error) {
	caster, err := v.cache.Cast(v.name)
	if err != nil {
		return "", err
	}

	if caster.IsNil() {
		return "", nil
	}

	return caster.String()
}

func (v *verification) Validate(code string) (bool, error) {
	c, err := v.Get()
	if err != nil || c == "" || code == "" {
		return false, err
	}

	return c == code, nil
}

func (v *verification) Exists() (bool, error) {
	return v.cache.Exists(v.name)
}

func (v *verification) TTL() (time.Duration, error) {
	return v.cache.TTL(v.name)
}
