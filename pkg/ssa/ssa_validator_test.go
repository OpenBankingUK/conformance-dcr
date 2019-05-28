package ssa_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/ssa"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"
	"github.com/dgrijalva/jwt-go"
)

const PubKeyTest = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzSkSfBX6BqUU2xOLKH6M
E5R3ZEPi7Ab7fXCRGUibKVpzVfC/+ldIuFkmbmpnrtCgEGMmeZK2DdcHt/PTk+6R
R5Jn7Huhxtn4tjLqe4Rle53QGIfdIwQyZAgKvNa9rjrQ5no2Ux2ViSu0csS0FaHg
2RK06JGBBwk1HZShDRxQpwuoW7fwPmhybpzZGypM4LsHGIzKrW2ygYI/1u5zDeFJ
Lj8qr7Qvkny4z/X3Na+mkgsQ8z/smDyBG6ZRRpMzFriFN+F+scq6nUGDYSuQ4RSm
R95fGAKyibC6uuAN3bcWbIEtuLTTPr9mPMJv0YZLPQ/ItC8gyMdK9MMSqRzbtP/V
4QIDAQAB
-----END PUBLIC KEY-----`

const PrivKeyTest = `-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAzSkSfBX6BqUU2xOLKH6ME5R3ZEPi7Ab7fXCRGUibKVpzVfC/
+ldIuFkmbmpnrtCgEGMmeZK2DdcHt/PTk+6RR5Jn7Huhxtn4tjLqe4Rle53QGIfd
IwQyZAgKvNa9rjrQ5no2Ux2ViSu0csS0FaHg2RK06JGBBwk1HZShDRxQpwuoW7fw
PmhybpzZGypM4LsHGIzKrW2ygYI/1u5zDeFJLj8qr7Qvkny4z/X3Na+mkgsQ8z/s
mDyBG6ZRRpMzFriFN+F+scq6nUGDYSuQ4RSmR95fGAKyibC6uuAN3bcWbIEtuLTT
Pr9mPMJv0YZLPQ/ItC8gyMdK9MMSqRzbtP/V4QIDAQABAoIBAGrUbkmq7+yx3XBO
dvH5I1u2vYy6RNl+sUoaaZPu2MzpypD/shfbh4Rm97obMi43eIHP/2Li4GXIiL5i
+tNXMNNCC1d68Tyol/fX/32d7XI3NPrxFFd9gffQsDaw40hsXAcHsG4LQ6iP14DD
qLpy9rBSK13HFpbhVoU2tb4r9ltkqat3oJfko0wg04z5TyJlASElckDEQRkhuRt4
Fg+kw/oY9pKD/7fC4fjJ1QZ4cDQWOv1yYZrGmoIs0dxbhokfVC197oC2O4M2n6pa
Gc14KL1JkBOlpZXHOuRv0+Tnv3qSZcAMEADfiumnH8eptjzjdZbJ0FM1T2ds/pJ8
tbF2bRkCgYEA8KDn9U7p5J2CmxdgB1r9bOPg46xjSA8J3Yc3Su4VKZ62otlBL5mR
yKwbhU4bXaOgHfj7cuTsrVmDEzVXvxXUX6hnvdXLNlEJXHz5XkUzuawZcZ0jpIv5
GQWVbc85ln2YdRk6riuCTt5DDrgv6JAE3UXbZBI0sWT9XQsgBLunE4MCgYEA2kQk
ZyRyTWrPQi/sU/zdfN+Rot9x1YPHHoC8rnM5XVZk3OjtzK2Y16hCGoYIRQ559ZNZ
nz6s2j+4PTu52T2JBbuCXYDRaQLSktE4ALYemdb/7NNkdD4LhlCVn7S+Xux6vsGF
VjvTZSq4nu+1ElBuFl5LCry94SqeGdeyCEX3n8sCgYA6yWTB9oyH0L9Wuog4Y89k
KewIU1ZSBXKIj38/rBi5eU/vSxp00ZTfLMTwdVuULeRxTiHIOQtlcmfmanLMeT3Q
POlTZIbn9zZNRS77C/cOFnCE5DoP+i5aIZYXJLhR/s8fVJGUeYa1U/GYCAGUVJML
qARoV1ZOPHj1oUEqRtoTlwKBgC9g7038WVt4vfiuEmzAzQtYNHLYcgtZCZYTd+Ge
XWtnX4mcflIZtL3LZl3/jjf/RnYKQEATCM5vWnzgRB1mACJga5IEbnCPDkqUY8Wz
wry+MNuln36kIThMsc3zHAfa6WIS+/CWF/Mz7NODURjioKL2YO+5vLXt3FfbvGeT
WIc1AoGAMzB1b1C9kU9eurmSArSPooT0m7y1ff+KIU1I7G3oaCx1yNAvOJ89Iw9g
kk244ZP5FygXJ2pRhOzjRlORUhVUwfzqcxDJwOsk9Jq+Z4fqeYWKI9vkGYSZu6K3
0hGp0Szn+zNkFdxlCkO5BVILCPJO6htP3HcFiLhuAxIMMFiTR8Q=
-----END RSA PRIVATE KEY-----`

func TestValidateSSA(t *testing.T) {
	require := test.NewRequire(t)
	claims := jwt.MapClaims{
		"iss":                         "OpenBanking Ltd",
		"iat":                         1492756331,
		"exp":                         1595757550,
		"jti":                         "id12345685439487678",
		"software_environment":        "production",
		"software_mode":               "live",
		"software_id":                 "65d1f27c-4aea-4549-9c21-60e495a7a86f",
		"software_client_id":          "OpenBanking TPP Client Unique ID",
		"software_client_name":        "Amazon Prime Movies",
		"software_client_description": "Amazon Prime Movies is a moving streaming service",
		"software_version":            "2.2",
		"software_client_uri":         "https://prime.amazon.com",
		"software_redirect_uris": []string{
			"https://prime.amazon.com/cb",
			"https://prime.amazon.co.uk/cb",
		},
		"software_roles": []string{
			"PISP",
			"AISP",
		},
		"software_logo_uri":              "https://prime.amazon.com/logo.png",
		"software_jwks_endpoint":         "https://jwks.openbanking.org.uk/org_id/software_id.jkws",
		"software_jwks_revoked_endpoint": "https://jwks.openbanking.org.uk/org_id/revoked/software_id.jkws",
		"software_policy_uri":            "https://tpp.com/policy.html",
		"software_tos_uri":               "https://tpp.com/tos.html",
		"software_on_behalf_of_org":      "https://api.openbanking.org.uk/scim2/OBTrustedPaymentParty/1234567789",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodPS256, claims)
	token.Header["kid"] = "veryUniqueJwtKey"
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(PrivKeyTest))
	require.Nil(err)
	ssaJwt, err := token.SignedString(privKey)
	require.Nil(err)
	ssaValidator := ssa.NewSSAValidator(ssa.PublicKeyLookupFromByteSlice([]byte(PubKeyTest)))
	ssaValue, err := ssaValidator.Validate(ssaJwt)
	require.Nil(err)
	require.Equal(claims["iss"], ssaValue.Issuer)
	require.Equal(claims["software_id"], ssaValue.SoftwareID)
	require.Equal(claims["software_roles"], ssaValue.SoftwareRoles)
}

func TestSSAJwtIsInvalid(t *testing.T) {
	require := test.NewRequire(t)
	token := jwt.NewWithClaims(jwt.SigningMethodPS256, jwt.MapClaims{
		"iss":                         "OpenBanking Ltd",
		"iat":                         1492756331,
		"exp":                         10,
		"jti":                         "id12345685439487678",
		"software_environment":        "production",
		"software_mode":               "live",
		"software_id":                 "65d1f27c-4aea-4549-9c21-60e495a7a86f",
		"software_client_id":          "OpenBanking TPP Client Unique ID",
		"software_client_name":        "Amazon Prime Movies",
		"software_client_description": "Amazon Prime Movies is a moving streaming service",
		"software_version":            "2.2",
		"software_client_uri":         "https://prime.amazon.com",
		"software_redirect_uris": []string{
			"https://prime.amazon.com/cb",
			"https://prime.amazon.co.uk/cb",
		},
		"software_roles": []string{
			"PISP",
			"AISP",
		},
		"software_logo_uri":              "https://prime.amazon.com/logo.png",
		"software_jwks_endpoint":         "https://jwks.openbanking.org.uk/org_id/software_id.jkws",
		"software_jwks_revoked_endpoint": "https://jwks.openbanking.org.uk/org_id/revoked/software_id.jkws",
		"software_policy_uri":            "https://tpp.com/policy.html",
		"software_tos_uri":               "https://tpp.com/tos.html",
		"software_on_behalf_of_org":      "https://api.openbanking.org.uk/scim2/OBTrustedPaymentParty/1234567789",
	})
	token.Header["kid"] = "veryUniqueJwtKey"
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(PrivKeyTest))
	require.Nil(err)
	ssaJwt, err := token.SignedString(privKey)
	require.Nil(err)
	ssaValidator := ssa.NewSSAValidator(ssa.PublicKeyLookupFromByteSlice([]byte(PubKeyTest)))
	_, err = ssaValidator.Validate(ssaJwt)
	require.NotNil(err)
	require.Contains(err.Error(), "token is expired by")
}

func TestRemotePubKeyJwt(t *testing.T) {
	require := test.NewRequire(t)
	certServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, PubKeyTest)
	}))
	defer certServer.Close()
	jwkServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"keys": [
			{
				"e": "AQAB",
				"kid": "GyVVcMPbU4QucpelwnDNiUJR4qQ",
				"kty": "RSA",
				"n": "vakAE3hb8opMX3zP6o929xh2ncsqAa9UtlbwZluVRFYZJb5s7-n4zqR2tqadaG57Fd6ZvhSqzq5qwd8ZvQeVM5N70ISwwXD5u9MFupjtmgLS3ioFucIbTNEmnobXppQC3eDTZI8x3DMkxy5H3za2e8ZFRrHwu6boNFQ-c7eibOQpmSAhD0G2CRm6sEK2uJuBEvUKQXZ5L6sli3Zd1TxsYxmO2x9fYkoml5Q_SK-OKi6x_MvDWxVOE1Ld1i4YhiPczDSgrWxPbMGh5iUdFT3Jikc3ppiE6E2h0HjQ0r1jQstlGScR5zul4-WQr9b8JEqYRK9uOE8dlW6zXu4mGtH36Q",
				"use": "tls",
				"x5c": [
				  "MIIFODCCBCCgAwIBAgIEWcV+HzANBgkqhkiG9w0BAQsFADBTMQswCQYDVQQGEwJHQjEUMBIGA1UEChMLT3BlbkJhbmtpbmcxLjAsBgNVBAMTJU9wZW5CYW5raW5nIFByZS1Qcm9kdWN0aW9uIElzc3VpbmcgQ0EwHhcNMTkwNDA5MTA0ODU2WhcNMjAwNTA5MTExODU2WjBhMQswCQYDVQQGEwJHQjEUMBIGA1UEChMLT3BlbkJhbmtpbmcxGzAZBgNVBAsTEjAwMTU4MDAwMDEwNDFSYkFBSTEfMB0GA1UEAxMWUkVmWktvN3pOMkllRTBYMlJGR1RiNDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAL2pABN4W/KKTF98z+qPdvcYdp3LKgGvVLZW8GZblURWGSW+bO/p+M6kdramnWhuexXemb4Uqs6uasHfGb0HlTOTe9CEsMFw+bvTBbqY7ZoC0t4qBbnCG0zRJp6G16aUAt3g02SPMdwzJMcuR982tnvGRUax8Lum6DRUPnO3omzkKZkgIQ9BtgkZurBCtribgRL1CkF2eS+rJYt2XdU8bGMZjtsfX2JKJpeUP0ivjiousfzLw1sVThNS3dYuGIYj3Mw0oK1sT2zBoeYlHRU9yYpHN6aYhOhNodB40NK9Y0LLZRknEec7pePlkK/W/CRKmESvbjhPHZVus17uJhrR9+kCAwEAAaOCAgQwggIAMA4GA1UdDwEB/wQEAwIHgDAgBgNVHSUBAf8EFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwgeAGA1UdIASB2DCB1TCB0gYLKwYBBAGodYEGAWQwgcIwKgYIKwYBBQUHAgEWHmh0dHA6Ly9vYi50cnVzdGlzLmNvbS9wb2xpY2llczCBkwYIKwYBBQUHAgIwgYYMgYNVc2Ugb2YgdGhpcyBDZXJ0aWZpY2F0ZSBjb25zdGl0dXRlcyBhY2NlcHRhbmNlIG9mIHRoZSBPcGVuQmFua2luZyBSb290IENBIENlcnRpZmljYXRpb24gUG9saWNpZXMgYW5kIENlcnRpZmljYXRlIFByYWN0aWNlIFN0YXRlbWVudDBtBggrBgEFBQcBAQRhMF8wJgYIKwYBBQUHMAGGGmh0dHA6Ly9vYi50cnVzdGlzLmNvbS9vY3NwMDUGCCsGAQUFBzAChilodHRwOi8vb2IudHJ1c3Rpcy5jb20vb2JfcHBfaXNzdWluZ2NhLmNydDA6BgNVHR8EMzAxMC+gLaArhilodHRwOi8vb2IudHJ1c3Rpcy5jb20vb2JfcHBfaXNzdWluZ2NhLmNybDAfBgNVHSMEGDAWgBRQc5HGIXLTd/T+ABIGgVx5eW4/UDAdBgNVHQ4EFgQUanhMVcNxUI03lzhtM0Ap9Uqe9MYwDQYJKoZIhvcNAQELBQADggEBAA+Pxffl5XELhA5X2k7eL4nqqnR82DWn5iG6sHfdJOUwUlsIewyTB7M6seYiSu8ezrWfyVASqYJUqQacNVc1Q0DncmqURBetAsGNWh1hBVB7mTci54CGnqc3WAZZ9Mkl326uceNVEcE5HQ/wbynDqaZzJb7kqJlfaSZgSptV22dYnSX8ZWG7AWFYWWXytCUw29KLUZv4QDtSpOUZOP98GWkDXgEo082GaJjr4IS7BlNUVtICQGVFZ9RvJr7yAiscQTSKII+viHa+8jtaGweHKr69oAaIzvMQ1hK9jFaNRaYSK6eNgEncQSddd9U04x65N+uyHUd1qG39gtEipxOVlMs="
				],
				"x5t": "47LacKAUQ_OcuAmsSomIywM9e4g=",
				"x5u": "`+certServer.URL+`",
				"x5t#S256": "5G7DWO0Omk1GxnM_PTnpq29fY3FT81EVEAIvkYii-BI="
			  }
		]}`)
	}))
	defer jwkServer.Close()
	token := jwt.NewWithClaims(jwt.SigningMethodPS256, jwt.MapClaims{
		"iss":                         "OpenBanking Ltd",
		"iat":                         1492756331,
		"exp":                         1692756331,
		"jti":                         "id12345685439487678",
		"software_environment":        "production",
		"software_mode":               "live",
		"software_id":                 "65d1f27c-4aea-4549-9c21-60e495a7a86f",
		"software_client_id":          "OpenBanking TPP Client Unique ID",
		"software_client_name":        "Amazon Prime Movies",
		"software_client_description": "Amazon Prime Movies is a moving streaming service",
		"software_version":            "2.2",
		"software_client_uri":         "https://prime.amazon.com",
		"software_redirect_uris": []string{
			"https://prime.amazon.com/cb",
			"https://prime.amazon.co.uk/cb",
		},
		"software_roles": []string{
			"PISP",
			"AISP",
		},
		"software_logo_uri":              "https://prime.amazon.com/logo.png",
		"software_jwks_endpoint":         jwkServer.URL,
		"software_jwks_revoked_endpoint": "https://jwks.openbanking.org.uk/org_id/revoked/software_id.jkws",
		"software_policy_uri":            "https://tpp.com/policy.html",
		"software_tos_uri":               "https://tpp.com/tos.html",
		"software_on_behalf_of_org":      "https://api.openbanking.org.uk/scim2/OBTrustedPaymentParty/1234567789",
	})
	token.Header["kid"] = "GyVVcMPbU4QucpelwnDNiUJR4qQ"
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(PrivKeyTest))
	require.Nil(err)
	ssaJwt, err := token.SignedString(privKey)
	require.Nil(err)
	ssaValidator := ssa.NewSSAValidator(ssa.PublicKeyLookupFromJWKSEndpoint(jwkServer.Client()))
	_, err = ssaValidator.Validate(ssaJwt)
	require.Nil(err)
}

func TestRemotePubKeyJwtFailsOnMissingJWKS(t *testing.T) {
	require := test.NewRequire(t)
	token := jwt.NewWithClaims(jwt.SigningMethodPS256, jwt.MapClaims{
		"iss":                         "OpenBanking Ltd",
		"iat":                         1492756331,
		"exp":                         1692756331,
		"jti":                         "id12345685439487678",
		"software_environment":        "production",
		"software_mode":               "live",
		"software_id":                 "65d1f27c-4aea-4549-9c21-60e495a7a86f",
		"software_client_id":          "OpenBanking TPP Client Unique ID",
		"software_client_name":        "Amazon Prime Movies",
		"software_client_description": "Amazon Prime Movies is a moving streaming service",
		"software_version":            "2.2",
		"software_client_uri":         "https://prime.amazon.com",
		"software_redirect_uris": []string{
			"https://prime.amazon.com/cb",
			"https://prime.amazon.co.uk/cb",
		},
		"software_roles": []string{
			"PISP",
			"AISP",
		},
		"software_logo_uri":              "https://prime.amazon.com/logo.png",
		"software_jwks_endpoint":         "invalid url",
		"software_jwks_revoked_endpoint": "https://jwks.openbanking.org.uk/org_id/revoked/software_id.jkws",
		"software_policy_uri":            "https://tpp.com/policy.html",
		"software_tos_uri":               "https://tpp.com/tos.html",
		"software_on_behalf_of_org":      "https://api.openbanking.org.uk/scim2/OBTrustedPaymentParty/1234567789",
	})
	token.Header["kid"] = "GyVVcMPbU4QucpelwnDNiUJR4qQ"
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(PrivKeyTest))
	require.Nil(err)
	ssaJwt, err := token.SignedString(privKey)
	require.Nil(err)
	ssaValidator := ssa.NewSSAValidator(ssa.PublicKeyLookupFromJWKSEndpoint(&http.Client{}))
	_, err = ssaValidator.Validate(ssaJwt)
	require.NotNil(err)
	require.Contains(err.Error(), "unable to retrieve data from jwks endpoint invalid url")
}

func TestRemotePubKeyJwtFailsOnInvalidJWKSResponse(t *testing.T) {
	require := test.NewRequire(t)
	jwkServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `foobar`)
	}))
	defer jwkServer.Close()
	token := jwt.NewWithClaims(jwt.SigningMethodPS256, jwt.MapClaims{
		"iss":                         "OpenBanking Ltd",
		"iat":                         1492756331,
		"exp":                         1692756331,
		"jti":                         "id12345685439487678",
		"software_environment":        "production",
		"software_mode":               "live",
		"software_id":                 "65d1f27c-4aea-4549-9c21-60e495a7a86f",
		"software_client_id":          "OpenBanking TPP Client Unique ID",
		"software_client_name":        "Amazon Prime Movies",
		"software_client_description": "Amazon Prime Movies is a moving streaming service",
		"software_version":            "2.2",
		"software_client_uri":         "https://prime.amazon.com",
		"software_redirect_uris": []string{
			"https://prime.amazon.com/cb",
			"https://prime.amazon.co.uk/cb",
		},
		"software_roles": []string{
			"PISP",
			"AISP",
		},
		"software_logo_uri":              "https://prime.amazon.com/logo.png",
		"software_jwks_endpoint":         jwkServer.URL,
		"software_jwks_revoked_endpoint": "https://jwks.openbanking.org.uk/org_id/revoked/software_id.jkws",
		"software_policy_uri":            "https://tpp.com/policy.html",
		"software_tos_uri":               "https://tpp.com/tos.html",
		"software_on_behalf_of_org":      "https://api.openbanking.org.uk/scim2/OBTrustedPaymentParty/1234567789",
	})
	token.Header["kid"] = "GyVVcMPbU4QucpelwnDNiUJR4qQ"
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(PrivKeyTest))
	require.Nil(err)
	ssaJwt, err := token.SignedString(privKey)
	require.Nil(err)
	ssaValidator := ssa.NewSSAValidator(ssa.PublicKeyLookupFromJWKSEndpoint(&http.Client{}))
	_, err = ssaValidator.Validate(ssaJwt)
	require.NotNil(err)
	require.Contains(err.Error(), "unable to parse json from jwk endpoint")
}

func TestRemotePubKeyJwtFailsOnNonMatchingKid(t *testing.T) {
	require := test.NewRequire(t)
	certServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, PubKeyTest)
	}))
	defer certServer.Close()
	jwkServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"keys": [
			{
				"e": "AQAB",
				"kid": "foobar",
				"kty": "RSA",
				"n": "vakAE3hb8opMX3zP6o929xh2ncsqAa9UtlbwZluVRFYZJb5s7-n4zqR2tqadaG57Fd6ZvhSqzq5qwd8ZvQeVM5N70ISwwXD5u9MFupjtmgLS3ioFucIbTNEmnobXppQC3eDTZI8x3DMkxy5H3za2e8ZFRrHwu6boNFQ-c7eibOQpmSAhD0G2CRm6sEK2uJuBEvUKQXZ5L6sli3Zd1TxsYxmO2x9fYkoml5Q_SK-OKi6x_MvDWxVOE1Ld1i4YhiPczDSgrWxPbMGh5iUdFT3Jikc3ppiE6E2h0HjQ0r1jQstlGScR5zul4-WQr9b8JEqYRK9uOE8dlW6zXu4mGtH36Q",
				"use": "tls",
				"x5c": [
				  "MIIFODCCBCCgAwIBAgIEWcV+HzANBgkqhkiG9w0BAQsFADBTMQswCQYDVQQGEwJHQjEUMBIGA1UEChMLT3BlbkJhbmtpbmcxLjAsBgNVBAMTJU9wZW5CYW5raW5nIFByZS1Qcm9kdWN0aW9uIElzc3VpbmcgQ0EwHhcNMTkwNDA5MTA0ODU2WhcNMjAwNTA5MTExODU2WjBhMQswCQYDVQQGEwJHQjEUMBIGA1UEChMLT3BlbkJhbmtpbmcxGzAZBgNVBAsTEjAwMTU4MDAwMDEwNDFSYkFBSTEfMB0GA1UEAxMWUkVmWktvN3pOMkllRTBYMlJGR1RiNDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAL2pABN4W/KKTF98z+qPdvcYdp3LKgGvVLZW8GZblURWGSW+bO/p+M6kdramnWhuexXemb4Uqs6uasHfGb0HlTOTe9CEsMFw+bvTBbqY7ZoC0t4qBbnCG0zRJp6G16aUAt3g02SPMdwzJMcuR982tnvGRUax8Lum6DRUPnO3omzkKZkgIQ9BtgkZurBCtribgRL1CkF2eS+rJYt2XdU8bGMZjtsfX2JKJpeUP0ivjiousfzLw1sVThNS3dYuGIYj3Mw0oK1sT2zBoeYlHRU9yYpHN6aYhOhNodB40NK9Y0LLZRknEec7pePlkK/W/CRKmESvbjhPHZVus17uJhrR9+kCAwEAAaOCAgQwggIAMA4GA1UdDwEB/wQEAwIHgDAgBgNVHSUBAf8EFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwgeAGA1UdIASB2DCB1TCB0gYLKwYBBAGodYEGAWQwgcIwKgYIKwYBBQUHAgEWHmh0dHA6Ly9vYi50cnVzdGlzLmNvbS9wb2xpY2llczCBkwYIKwYBBQUHAgIwgYYMgYNVc2Ugb2YgdGhpcyBDZXJ0aWZpY2F0ZSBjb25zdGl0dXRlcyBhY2NlcHRhbmNlIG9mIHRoZSBPcGVuQmFua2luZyBSb290IENBIENlcnRpZmljYXRpb24gUG9saWNpZXMgYW5kIENlcnRpZmljYXRlIFByYWN0aWNlIFN0YXRlbWVudDBtBggrBgEFBQcBAQRhMF8wJgYIKwYBBQUHMAGGGmh0dHA6Ly9vYi50cnVzdGlzLmNvbS9vY3NwMDUGCCsGAQUFBzAChilodHRwOi8vb2IudHJ1c3Rpcy5jb20vb2JfcHBfaXNzdWluZ2NhLmNydDA6BgNVHR8EMzAxMC+gLaArhilodHRwOi8vb2IudHJ1c3Rpcy5jb20vb2JfcHBfaXNzdWluZ2NhLmNybDAfBgNVHSMEGDAWgBRQc5HGIXLTd/T+ABIGgVx5eW4/UDAdBgNVHQ4EFgQUanhMVcNxUI03lzhtM0Ap9Uqe9MYwDQYJKoZIhvcNAQELBQADggEBAA+Pxffl5XELhA5X2k7eL4nqqnR82DWn5iG6sHfdJOUwUlsIewyTB7M6seYiSu8ezrWfyVASqYJUqQacNVc1Q0DncmqURBetAsGNWh1hBVB7mTci54CGnqc3WAZZ9Mkl326uceNVEcE5HQ/wbynDqaZzJb7kqJlfaSZgSptV22dYnSX8ZWG7AWFYWWXytCUw29KLUZv4QDtSpOUZOP98GWkDXgEo082GaJjr4IS7BlNUVtICQGVFZ9RvJr7yAiscQTSKII+viHa+8jtaGweHKr69oAaIzvMQ1hK9jFaNRaYSK6eNgEncQSddd9U04x65N+uyHUd1qG39gtEipxOVlMs="
				],
				"x5t": "47LacKAUQ_OcuAmsSomIywM9e4g=",
				"x5u": "`+certServer.URL+`",
				"x5t#S256": "5G7DWO0Omk1GxnM_PTnpq29fY3FT81EVEAIvkYii-BI="
			  }
		]}`)
	}))
	defer jwkServer.Close()
	token := jwt.NewWithClaims(jwt.SigningMethodPS256, jwt.MapClaims{
		"iss":                         "OpenBanking Ltd",
		"iat":                         1492756331,
		"exp":                         1692756331,
		"jti":                         "id12345685439487678",
		"software_environment":        "production",
		"software_mode":               "live",
		"software_id":                 "65d1f27c-4aea-4549-9c21-60e495a7a86f",
		"software_client_id":          "OpenBanking TPP Client Unique ID",
		"software_client_name":        "Amazon Prime Movies",
		"software_client_description": "Amazon Prime Movies is a moving streaming service",
		"software_version":            "2.2",
		"software_client_uri":         "https://prime.amazon.com",
		"software_redirect_uris": []string{
			"https://prime.amazon.com/cb",
			"https://prime.amazon.co.uk/cb",
		},
		"software_roles": []string{
			"PISP",
			"AISP",
		},
		"software_logo_uri":              "https://prime.amazon.com/logo.png",
		"software_jwks_endpoint":         jwkServer.URL,
		"software_jwks_revoked_endpoint": "https://jwks.openbanking.org.uk/org_id/revoked/software_id.jkws",
		"software_policy_uri":            "https://tpp.com/policy.html",
		"software_tos_uri":               "https://tpp.com/tos.html",
		"software_on_behalf_of_org":      "https://api.openbanking.org.uk/scim2/OBTrustedPaymentParty/1234567789",
	})
	token.Header["kid"] = "GyVVcMPbU4QucpelwnDNiUJR4qQ"
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(PrivKeyTest))
	require.Nil(err)
	ssaJwt, err := token.SignedString(privKey)
	require.Nil(err)
	ssaValidator := ssa.NewSSAValidator(ssa.PublicKeyLookupFromJWKSEndpoint(jwkServer.Client()))
	_, err = ssaValidator.Validate(ssaJwt)
	require.NotNil(err)
	require.Contains(err.Error(), "unable to find key with kid")
}

func TestRemotePubKeyJwtFailsOnInvalidCertURL(t *testing.T) {
	require := test.NewRequire(t)
	jwkServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"keys": [
			{
				"e": "AQAB",
				"kid": "GyVVcMPbU4QucpelwnDNiUJR4qQ",
				"kty": "RSA",
				"n": "vakAE3hb8opMX3zP6o929xh2ncsqAa9UtlbwZluVRFYZJb5s7-n4zqR2tqadaG57Fd6ZvhSqzq5qwd8ZvQeVM5N70ISwwXD5u9MFupjtmgLS3ioFucIbTNEmnobXppQC3eDTZI8x3DMkxy5H3za2e8ZFRrHwu6boNFQ-c7eibOQpmSAhD0G2CRm6sEK2uJuBEvUKQXZ5L6sli3Zd1TxsYxmO2x9fYkoml5Q_SK-OKi6x_MvDWxVOE1Ld1i4YhiPczDSgrWxPbMGh5iUdFT3Jikc3ppiE6E2h0HjQ0r1jQstlGScR5zul4-WQr9b8JEqYRK9uOE8dlW6zXu4mGtH36Q",
				"use": "tls",
				"x5c": [
				  "MIIFODCCBCCgAwIBAgIEWcV+HzANBgkqhkiG9w0BAQsFADBTMQswCQYDVQQGEwJHQjEUMBIGA1UEChMLT3BlbkJhbmtpbmcxLjAsBgNVBAMTJU9wZW5CYW5raW5nIFByZS1Qcm9kdWN0aW9uIElzc3VpbmcgQ0EwHhcNMTkwNDA5MTA0ODU2WhcNMjAwNTA5MTExODU2WjBhMQswCQYDVQQGEwJHQjEUMBIGA1UEChMLT3BlbkJhbmtpbmcxGzAZBgNVBAsTEjAwMTU4MDAwMDEwNDFSYkFBSTEfMB0GA1UEAxMWUkVmWktvN3pOMkllRTBYMlJGR1RiNDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAL2pABN4W/KKTF98z+qPdvcYdp3LKgGvVLZW8GZblURWGSW+bO/p+M6kdramnWhuexXemb4Uqs6uasHfGb0HlTOTe9CEsMFw+bvTBbqY7ZoC0t4qBbnCG0zRJp6G16aUAt3g02SPMdwzJMcuR982tnvGRUax8Lum6DRUPnO3omzkKZkgIQ9BtgkZurBCtribgRL1CkF2eS+rJYt2XdU8bGMZjtsfX2JKJpeUP0ivjiousfzLw1sVThNS3dYuGIYj3Mw0oK1sT2zBoeYlHRU9yYpHN6aYhOhNodB40NK9Y0LLZRknEec7pePlkK/W/CRKmESvbjhPHZVus17uJhrR9+kCAwEAAaOCAgQwggIAMA4GA1UdDwEB/wQEAwIHgDAgBgNVHSUBAf8EFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwgeAGA1UdIASB2DCB1TCB0gYLKwYBBAGodYEGAWQwgcIwKgYIKwYBBQUHAgEWHmh0dHA6Ly9vYi50cnVzdGlzLmNvbS9wb2xpY2llczCBkwYIKwYBBQUHAgIwgYYMgYNVc2Ugb2YgdGhpcyBDZXJ0aWZpY2F0ZSBjb25zdGl0dXRlcyBhY2NlcHRhbmNlIG9mIHRoZSBPcGVuQmFua2luZyBSb290IENBIENlcnRpZmljYXRpb24gUG9saWNpZXMgYW5kIENlcnRpZmljYXRlIFByYWN0aWNlIFN0YXRlbWVudDBtBggrBgEFBQcBAQRhMF8wJgYIKwYBBQUHMAGGGmh0dHA6Ly9vYi50cnVzdGlzLmNvbS9vY3NwMDUGCCsGAQUFBzAChilodHRwOi8vb2IudHJ1c3Rpcy5jb20vb2JfcHBfaXNzdWluZ2NhLmNydDA6BgNVHR8EMzAxMC+gLaArhilodHRwOi8vb2IudHJ1c3Rpcy5jb20vb2JfcHBfaXNzdWluZ2NhLmNybDAfBgNVHSMEGDAWgBRQc5HGIXLTd/T+ABIGgVx5eW4/UDAdBgNVHQ4EFgQUanhMVcNxUI03lzhtM0Ap9Uqe9MYwDQYJKoZIhvcNAQELBQADggEBAA+Pxffl5XELhA5X2k7eL4nqqnR82DWn5iG6sHfdJOUwUlsIewyTB7M6seYiSu8ezrWfyVASqYJUqQacNVc1Q0DncmqURBetAsGNWh1hBVB7mTci54CGnqc3WAZZ9Mkl326uceNVEcE5HQ/wbynDqaZzJb7kqJlfaSZgSptV22dYnSX8ZWG7AWFYWWXytCUw29KLUZv4QDtSpOUZOP98GWkDXgEo082GaJjr4IS7BlNUVtICQGVFZ9RvJr7yAiscQTSKII+viHa+8jtaGweHKr69oAaIzvMQ1hK9jFaNRaYSK6eNgEncQSddd9U04x65N+uyHUd1qG39gtEipxOVlMs="
				],
				"x5t": "47LacKAUQ_OcuAmsSomIywM9e4g=",
				"x5u": "lol",
				"x5t#S256": "5G7DWO0Omk1GxnM_PTnpq29fY3FT81EVEAIvkYii-BI="
			  }
		]}`)
	}))
	defer jwkServer.Close()
	token := jwt.NewWithClaims(jwt.SigningMethodPS256, jwt.MapClaims{
		"iss":                         "OpenBanking Ltd",
		"iat":                         1492756331,
		"exp":                         1692756331,
		"jti":                         "id12345685439487678",
		"software_environment":        "production",
		"software_mode":               "live",
		"software_id":                 "65d1f27c-4aea-4549-9c21-60e495a7a86f",
		"software_client_id":          "OpenBanking TPP Client Unique ID",
		"software_client_name":        "Amazon Prime Movies",
		"software_client_description": "Amazon Prime Movies is a moving streaming service",
		"software_version":            "2.2",
		"software_client_uri":         "https://prime.amazon.com",
		"software_redirect_uris": []string{
			"https://prime.amazon.com/cb",
			"https://prime.amazon.co.uk/cb",
		},
		"software_roles": []string{
			"PISP",
			"AISP",
		},
		"software_logo_uri":              "https://prime.amazon.com/logo.png",
		"software_jwks_endpoint":         jwkServer.URL,
		"software_jwks_revoked_endpoint": "https://jwks.openbanking.org.uk/org_id/revoked/software_id.jkws",
		"software_policy_uri":            "https://tpp.com/policy.html",
		"software_tos_uri":               "https://tpp.com/tos.html",
		"software_on_behalf_of_org":      "https://api.openbanking.org.uk/scim2/OBTrustedPaymentParty/1234567789",
	})
	token.Header["kid"] = "GyVVcMPbU4QucpelwnDNiUJR4qQ"
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(PrivKeyTest))
	require.Nil(err)
	ssaJwt, err := token.SignedString(privKey)
	require.Nil(err)
	ssaValidator := ssa.NewSSAValidator(ssa.PublicKeyLookupFromJWKSEndpoint(jwkServer.Client()))
	_, err = ssaValidator.Validate(ssaJwt)
	require.NotNil(err)
	require.Contains(err.Error(), "unable to download certificate from URI")
}
