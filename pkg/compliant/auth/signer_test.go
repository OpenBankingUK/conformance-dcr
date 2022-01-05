package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"testing"
	"time"

	"github.com/OpenBankingUK/conformance-dcr/pkg/certs"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewJwtSigner(t *testing.T) {
	privateKey, err := certs.ParseRsaPrivateKeyFromPemFile("testdata/private-sign.key")
	require.NoError(t, err)
	signer := NewJwtSigner(
		jwt.SigningMethodRS256,
		"ssa",
		"issuer",
		"aud",
		"kid",
		"private_key_jwt",
		"none",
		[]string{"/redirect"},
		[]string{"code", "code id_token"},
		privateKey,
		time.Hour,
		&x509.Certificate{},
		"",
	)

	signedClaims, err := signer.Claims()
	require.NoError(t, err)

	token, err := jwt.Parse(signedClaims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return privateKey.Public(), nil
	})
	require.NoError(t, err)

	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, "web", claims["application_type"])
	assert.Equal(t, "aud", claims["aud"])
	assert.Equal(t, []interface{}{"authorization_code", "client_credentials"}, claims["grant_types"])
	assert.Equal(t, "RS256", claims["id_token_signed_response_alg"])
	assert.Equal(t, "issuer", claims["iss"])
	assert.Equal(t, "kid", token.Header["kid"])
	assert.Equal(t, []interface{}{"/redirect"}, claims["redirect_uris"])
	assert.Equal(t, []interface{}{"code", "code id_token"}, claims["response_types"])
	assert.Equal(t, "none", claims["request_object_signing_alg"])
	assert.Equal(t, []interface{}{"code", "code id_token"}, claims["response_types"])
	assert.Equal(t, "accounts openid", claims["scope"])
	assert.Equal(t, "ssa", claims["software_statement"])
	assert.Equal(t, "private_key_jwt", claims["token_endpoint_auth_method"])
	assert.Equal(t, nil, claims["tls_client_auth_subject_dn"])
}

func TestNewJwtSigner_TlsClientAuthAddSubjectToClaims(t *testing.T) {
	privateKey, err := certs.ParseRsaPrivateKeyFromPemFile("testdata/private-sign.key")
	require.NoError(t, err)

	signer := NewJwtSigner(
		jwt.SigningMethodRS256,
		"ssa",
		"issuer",
		"aud",
		"kid",
		"tls_client_auth",
		"none",
		[]string{"/redirect"},
		[]string{"code", "code id_token"},
		privateKey,
		time.Hour,
		&x509.Certificate{Subject: pkix.Name{Organization: []string{"OB"}}},
		"",
	)

	token, claims := getJwtClaims(t, signer, privateKey)

	assert.Equal(t, "web", claims["application_type"])
	assert.Equal(t, "aud", claims["aud"])
	assert.Equal(t, []interface{}{"authorization_code", "client_credentials"}, claims["grant_types"])
	assert.Equal(t, "RS256", claims["id_token_signed_response_alg"])
	assert.Equal(t, "issuer", claims["iss"])
	assert.Equal(t, "kid", token.Header["kid"])
	assert.Equal(t, []interface{}{"/redirect"}, claims["redirect_uris"])
	assert.Equal(t, []interface{}{"code", "code id_token"}, claims["response_types"])
	assert.Equal(t, "none", claims["request_object_signing_alg"])
	assert.Equal(t, []interface{}{"code", "code id_token"}, claims["response_types"])
	assert.Equal(t, "accounts openid", claims["scope"])
	assert.Equal(t, "ssa", claims["software_statement"])
	assert.Equal(t, "tls_client_auth", claims["token_endpoint_auth_method"])
	assert.Equal(t, "O=OB", claims["tls_client_auth_subject_dn"])
}

func TestNewJwtSigner_TlsClientAuthAddConfigurableSubjectToClaims(t *testing.T) {
	privateKey, err := certs.ParseRsaPrivateKeyFromPemFile("testdata/private-sign.key")
	require.NoError(t, err)
	signer := NewJwtSigner(
		jwt.SigningMethodRS256,
		"ssa",
		"issuer",
		"aud",
		"kid",
		"tls_client_auth",
		"none",
		[]string{"/redirect"},
		[]string{"code", "code id_token"},
		privateKey,
		time.Hour,
		&x509.Certificate{Subject: pkix.Name{Organization: []string{"OB"}}},
		"CN=Configured Subject DN",
	)

	_, claims := getJwtClaims(t, signer, privateKey)

	assert.Equal(t, "CN=Configured Subject DN", claims["tls_client_auth_subject_dn"])
}

func getJwtClaims(t *testing.T, signer Signer, privateKey *rsa.PrivateKey) (*jwt.Token, jwt.MapClaims) {
	signedClaims, err := signer.Claims()
	require.NoError(t, err)

	token, err := jwt.Parse(signedClaims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return privateKey.Public(), nil
	})
	require.NoError(t, err)

	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	return token, claims
}

func TestNewJwtSigner_TlsClientAuthDoesNotPanicOnMissingCert(t *testing.T) {
	privateKey, err := certs.ParseRsaPrivateKeyFromPemFile("testdata/private-sign.key")
	require.NoError(t, err)
	signer := NewJwtSigner(
		jwt.SigningMethodRS256,
		"ssa",
		"issuer",
		"aud",
		"kid",
		"tls_client_auth",
		"none",
		[]string{"/redirect"},
		[]string{"code", "code id_token"},
		privateKey,
		time.Hour,
		nil,
		"",
	)

	_, err = signer.Claims()

	assert.EqualError(t, err, "transport cert not available")
}

func TestNewJwtSigner_OmitsEmptyResponseTypes(t *testing.T) {
	privateKey, err := certs.ParseRsaPrivateKeyFromPemFile("testdata/private-sign.key")
	require.NoError(t, err)
	signer := NewJwtSigner(
		jwt.SigningMethodRS256,
		"ssa",
		"issuer",
		"aud",
		"kid",
		"tls_client_auth",
		"none",
		[]string{"/redirect"},

		// testing empty/nil
		nil,

		privateKey,
		time.Hour,
		&x509.Certificate{Subject: pkix.Name{Organization: []string{"OB"}}},
		"",
	)
	_, claims := getJwtClaims(t, signer, privateKey)

	_, exists := claims["response_types"]
	assert.False(t, exists)
}
