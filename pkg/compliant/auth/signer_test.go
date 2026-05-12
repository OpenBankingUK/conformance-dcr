package auth

import (
	crypto_rand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/OpenBankingUK/conformance-dcr/pkg/certs"
	"github.com/golang-jwt/jwt/v4"
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
		false,
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
		false,
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
		false,
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
		false,
	)

	_, err = signer.Claims()

	assert.EqualError(t, err, "transport cert not available")
}

// newCertWithOrgIdentifier creates a real parsed certificate whose RawSubject
// includes an organizationIdentifier attribute (OID 2.5.4.97), so subjectDN
// can be exercised against genuine ASN.1-encoded bytes.
func newCertWithOrgIdentifier(t *testing.T) *x509.Certificate {
	t.Helper()
	key, err := rsa.GenerateKey(crypto_rand.Reader, 2048)
	require.NoError(t, err)

	orgIDOID := asn1.ObjectIdentifier{2, 5, 4, 97}
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Country:      []string{"GB"},
			Organization: []string{"Open Banking Limited"},
			CommonName:   "0015800001041RbAAI",
			ExtraNames: []pkix.AttributeTypeAndValue{
				{Type: orgIDOID, Value: "PSDGB-OB-Unknown001"},
			},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour),
	}
	certDER, err := x509.CreateCertificate(crypto_rand.Reader, template, template, &key.PublicKey, key)
	require.NoError(t, err)
	cert, err := x509.ParseCertificate(certDER)
	require.NoError(t, err)
	return cert
}

func TestNewJwtSigner_TlsClientAuthDerivedSubjectDN_UseOID(t *testing.T) {
	privateKey, err := certs.ParseRsaPrivateKeyFromPemFile("testdata/private-sign.key")
	require.NoError(t, err)

	cert := newCertWithOrgIdentifier(t)

	commonArgs := func(useOID bool) Signer {
		return NewJwtSigner(
			jwt.SigningMethodRS256,
			"ssa", "issuer", "aud", "kid",
			"tls_client_auth",
			"none",
			[]string{"/redirect"},
			[]string{"code"},
			privateKey,
			time.Hour,
			cert,
			"",
			useOID,
		)
	}

	t.Run("useOID=false uses friendly name", func(t *testing.T) {
		_, claims := getJwtClaims(t, commonArgs(false), privateKey)
		dn, _ := claims["tls_client_auth_subject_dn"].(string)
		assert.Contains(t, dn, "organizationIdentifier=PSDGB-OB-Unknown001")
		assert.NotContains(t, dn, "2.5.4.97=")
	})

	t.Run("useOID=true uses numeric OID", func(t *testing.T) {
		_, claims := getJwtClaims(t, commonArgs(true), privateKey)
		dn, _ := claims["tls_client_auth_subject_dn"].(string)
		assert.Contains(t, dn, "2.5.4.97=PSDGB-OB-Unknown001")
		assert.NotContains(t, dn, "organizationIdentifier=")
	})
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
		false,
	)
	_, claims := getJwtClaims(t, signer, privateKey)

	_, exists := claims["response_types"]
	assert.False(t, exists)
}
