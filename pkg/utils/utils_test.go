package utils_test

import (
	"testing"
	"time"

	"ETM/pkg/utils"
	"ETM/pkg/vars"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func init() {
	vars.SecretKey = "test-secret-key-for-unit-tests-only"
}

func makeToken(claims jwt.MapClaims, method jwt.SigningMethod) string {
	var key interface{}
	if method == jwt.SigningMethodNone {
		key = jwt.UnsafeAllowNoneSignatureType
	} else {
		key = []byte(vars.SecretKey)
	}
	token := jwt.NewWithClaims(method, claims)
	signed, _ := token.SignedString(key)
	return signed
}

func validClaims(userID uint, userUUID uuid.UUID) jwt.MapClaims {
	return jwt.MapClaims{
		"sub":  float64(userID),
		"uuid": userUUID.String(),
		"exp":  time.Now().Add(30 * time.Minute).Unix(),
		"iss":  "etm",
	}
}

// ── Password hashing ──────────────────────────────────────────────────────────

func TestGenerateHashPassword_NonEmpty(t *testing.T) {
	hash, err := utils.GenerateHashPassword("hunter2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hash == "" || hash == "hunter2" {
		t.Fatalf("expected bcrypt hash, got %q", hash)
	}
}

func TestCompareHashPassword_Correct(t *testing.T) {
	hash, _ := utils.GenerateHashPassword("hunter2")
	if !utils.CompareHashPassword("hunter2", hash) {
		t.Error("expected correct password to match")
	}
}

func TestCompareHashPassword_Wrong(t *testing.T) {
	hash, _ := utils.GenerateHashPassword("hunter2")
	if utils.CompareHashPassword("wrongpass", hash) {
		t.Error("expected wrong password to not match")
	}
}

func TestCompareHashPassword_EmptyInput(t *testing.T) {
	hash, _ := utils.GenerateHashPassword("hunter2")
	if utils.CompareHashPassword("", hash) {
		t.Error("expected empty input to not match")
	}
}

// ── Token parsing ─────────────────────────────────────────────────────────────

func TestParseToken_Valid(t *testing.T) {
	uid := uint(42)
	id := uuid.New()
	tokenStr := makeToken(validClaims(uid, id), jwt.SigningMethodHS256)

	claims, err := utils.ParseToken(tokenStr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if uint(claims["sub"].(float64)) != uid {
		t.Errorf("sub: got %v, want %d", claims["sub"], uid)
	}
	if claims["uuid"].(string) != id.String() {
		t.Errorf("uuid: got %v, want %v", claims["uuid"], id)
	}
}

func TestParseToken_TamperedSignature(t *testing.T) {
	tokenStr := makeToken(validClaims(1, uuid.New()), jwt.SigningMethodHS256)
	_, err := utils.ParseToken(tokenStr + "x")
	if err == nil {
		t.Error("expected error for tampered token signature")
	}
}

func TestParseToken_Expired(t *testing.T) {
	claims := jwt.MapClaims{
		"sub":  float64(1),
		"uuid": uuid.New().String(),
		"exp":  time.Now().Add(-1 * time.Hour).Unix(),
	}
	tokenStr := makeToken(claims, jwt.SigningMethodHS256)
	_, err := utils.ParseToken(tokenStr)
	if err == nil {
		t.Error("expected error for expired token")
	}
}

func TestParseToken_WrongSecret(t *testing.T) {
	claims := validClaims(1, uuid.New())
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte("different-secret"))
	_, err := utils.ParseToken(tokenStr)
	if err == nil {
		t.Error("expected error for token signed with wrong secret")
	}
}

func TestParseToken_AlgNone(t *testing.T) {
	tokenStr := makeToken(validClaims(1, uuid.New()), jwt.SigningMethodNone)
	_, err := utils.ParseToken(tokenStr)
	if err == nil {
		t.Error("alg:none token must be rejected — algorithm confusion attack not blocked")
	}
}

func TestParseTokenAllowExpired_AcceptsExpired(t *testing.T) {
	claims := jwt.MapClaims{
		"sub":  float64(1),
		"uuid": uuid.New().String(),
		"exp":  time.Now().Add(-1 * time.Hour).Unix(),
		"iss":  "etm",
	}
	tokenStr := makeToken(claims, jwt.SigningMethodHS256)
	got, err := utils.ParseTokenAllowExpired(tokenStr)
	if err != nil {
		t.Fatalf("expected expired token to be accepted: %v", err)
	}
	if uint(got["sub"].(float64)) != 1 {
		t.Errorf("sub: got %v, want 1", got["sub"])
	}
}

func TestParseTokenAllowExpired_RejectsBadSignature(t *testing.T) {
	claims := jwt.MapClaims{
		"sub": float64(1),
		"exp": time.Now().Add(-1 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte("wrong-secret"))
	if _, err := utils.ParseTokenAllowExpired(tokenStr); err == nil {
		t.Error("expected error for bad signature even if expired")
	}
}

func TestParseTokenAllowExpired_RejectsAlgNone(t *testing.T) {
	tokenStr := makeToken(validClaims(1, uuid.New()), jwt.SigningMethodNone)
	if _, err := utils.ParseTokenAllowExpired(tokenStr); err == nil {
		t.Error("alg:none must be rejected even by the lenient parser")
	}
}

// ── Claim extraction helpers ──────────────────────────────────────────────────

func TestGetUserID(t *testing.T) {
	uid := uint(7)
	tokenStr := makeToken(validClaims(uid, uuid.New()), jwt.SigningMethodHS256)
	got, err := utils.GetUserID("Bearer " + tokenStr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != uid {
		t.Errorf("got %d, want %d", got, uid)
	}
}

func TestGetUserUUID(t *testing.T) {
	id := uuid.New()
	tokenStr := makeToken(validClaims(1, id), jwt.SigningMethodHS256)
	got, err := utils.GetUserUUID("Bearer " + tokenStr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != id {
		t.Errorf("got %v, want %v", got, id)
	}
}
