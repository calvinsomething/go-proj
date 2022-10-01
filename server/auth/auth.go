package auth

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"log"
	"math/rand"
	"reflect"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/calvinsomething/go-proj/db"
)

const (
	saltLen      = 12
	checksumSize = 36
	// SessionMaxAge is the max age of the session and should be used as the MaxAge property of the cookie.
	SessionMaxAge = time.Hour * 48
)

var (
	hmacKey []byte

	// ErrBadLogin ...
	ErrBadLogin = errors.New("Invalid email/password")
	// ErrSessionExpired ...
	ErrSessionExpired = errors.New("Session expired")
	// ErrBadMAC ...
	ErrBadMAC = errors.New("MAC does not match")
	// ErrUserExists ...
	ErrUserExists = errors.New("A User with that email already exists")
)

func init() {
	rand.Seed(time.Now().UnixNano())
	hmacKey = randBytes(32)
	gob.Register(User{})
}

type (
	// User is the main type for user data.
	User struct {
		Email string
	}
)

// PasswordValidator checks the password string for at least one lower, upper, digit and symbol character.
func PasswordValidator(fl validator.FieldLevel) bool {
	if fl.Field().Kind() != reflect.String {
		return false
	}

	var hasLower, hasUpper, hasSym, hasDigit bool

	am1, zp1, AM1, ZP1, d0m1, d9p1 :=
		int('a')-1, int('z')+1, int('A')-1, int('Z')+1, int('0')-1, int('9')+1

	for _, r := range fl.Field().String() {
		ir := int(r)
		if ir > 127 {
			// string is not ascii
			return false
		} else if ir > am1 && ir < zp1 {
			hasLower = true
		} else if ir > AM1 && ir < ZP1 {
			hasUpper = true
		} else if ir > d0m1 && ir < d9p1 {
			hasDigit = true
		} else {
			hasSym = true
		}
		if hasLower && hasUpper && hasSym && hasDigit {
			return true
		}
	}

	return false
}

func randBytes(len int) []byte {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := rand.Intn(62) // 10 digits + 26 letters * 2
		if b > 36 {
			// lowercase
			bytes[i] = byte(b + 60) // 'a' - 36
		} else if b > 10 {
			// uppercase
			bytes[i] = byte(b + 55) // 'A' - 10
		} else {
			// digits
			bytes[i] = byte(b + 48)
		}
	}
	return bytes
}

func hashPassword(b []byte) ([]byte, []byte) {
	salt := randBytes(saltLen)
	return salt, sha256.New().Sum(append(b, salt...))
}

func deleteSession(ctx context.Context, sid string) error {
	return db.Pool.MustAffect(ctx, `
		DELETE FROM sessions
		WHERE [id] = ?;
	`, sid)
}

func getSession(ctx context.Context, sidB64 []byte) (data []byte, err error) {
	var sid []byte
	var updatedAt time.Time
	base64.StdEncoding.Decode(sid, sidB64)
	err = db.Pool.QueryRowContext(ctx, `
		SELECT [data], [updated_at]
		FROM sessions
		WHERE [id] = ?;
	`, sid).Scan(&data, &updatedAt)
	if err != nil {
		return
	}

	if time.Now().After(updatedAt.Add(SessionMaxAge)) {
		err = deleteSession(ctx, string(sid))
		if err != nil {
			log.Println("UNHANDLED:", err)
		}
		return nil, ErrSessionExpired
	}

	macStart := len(data) - checksumSize
	contents := data[:macStart]
	mac := data[macStart:]

	hashed := getHMAC(append(contents, []byte(updatedAt.String())...))
	if !hmac.Equal(hashed, mac) {
		log.Println("deleting corrupted session:", sid)
		err = deleteSession(ctx, string(sid))
		if err != nil {
			log.Println("UNHANDLED:", err)
		}
		return nil, ErrBadMAC
	}
	return
}

func getHMAC(message []byte) []byte {
	h := hmac.New(sha256.New, hmacKey)
	return h.Sum(message)
}

// CreateUser creates a new user in the database, hashing the password.
func CreateUser(ctx context.Context, email, password string) error {
	res, err := db.Pool.ExecContext(ctx, `
		SELECT NULL
		FROM users
		WHERE email = ?;
	`, email)
	if err != nil {
		return err
	} else if ra, err := res.RowsAffected(); err != nil {
		return err
	} else if ra != 0 {
		return ErrUserExists
	}

	salt, hashedPass := hashPassword([]byte(password))
	return db.Pool.MustAffect(ctx, `
		INSERT INTO users (email, password)
		VALUES (?, ?);
	`, email, append(salt, hashedPass...))
}

func setLoginAttempts(ctx context.Context, email string, attempts int) error {
	return db.Pool.MustAffect(ctx, `
		UPDATE users
		SET failed_attempts = ?;
	`, attempts)
}

// LogIn logs the User in by checking their password, recording failed attempts, and creating a session.
func LogIn(ctx context.Context, email, password string) (string, error) {
	var hashedPass []byte
	var attempts int
	err := db.Pool.QueryRowContext(ctx, `
		SELECT password, failed_attempts
		FROM users
		WHERE email = ?;
	`, email).Scan(&hashedPass, &attempts)
	if err == sql.ErrNoRows {
		return "", ErrBadLogin
	} else if err != nil {
		return "", err
	}

	pwAndSalt := append([]byte(password), hashedPass[:saltLen]...)

	if subtle.ConstantTimeCompare(sha256.New().Sum(pwAndSalt), hashedPass[saltLen:]) != 1 {
		attempts++
		if err = setLoginAttempts(ctx, email, attempts); err != nil {
			log.Println("UNHANDLED:", err)
		}
		log.Printf("Failed login attempt %d for user: %s\n", attempts, email)
		return "", ErrBadLogin
	}

	if err = setLoginAttempts(ctx, email, 0); err != nil {
		log.Println("UNHANDLED:", err)
	}

	u := &User{Email: email}

	sid, err := createSession(ctx, u)
	if err != nil {
		return "", err
	}

	var cookie []byte
	base64.StdEncoding.Encode(cookie, []byte(sid))

	return string(cookie), nil
}

func createSession(ctx context.Context, u *User) (string, error) {
	sid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	timestamp := time.Now()

	data, err := encodeSession(u, timestamp)
	if err != nil {
		return "", err
	}

	err = db.Pool.MustAffect(ctx, `
		INSERT INTO sessions (id, data, created_at, updated_at)
		VAlUES (?, ?, ?, ?);
	`, sid.String(), data, timestamp, timestamp)
	if err != nil {
		return "", err
	}

	return sid.String(), nil
}

func encodeSession(u *User, timestamp time.Time) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(u); err != nil {
		return nil, err
	}

	gobData := buf.Bytes()
	mac := getHMAC(append(gobData, []byte(timestamp.String())...))

	return append(gobData, mac...), nil
}

// GetUser decodes the session data and returns the User struct pointer.
func GetUser(ctx context.Context, sid string) (u *User, err error) {
	data, err := getSession(ctx, []byte(sid))
	if err != nil {
		return
	}

	var buf bytes.Buffer
	buf.Write(data)
	decoder := gob.NewDecoder(&buf)

	decoder.Decode(&u)
	return
}
