package ssa_test

import (
	"testing"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/ssa"
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
	token := jwt.NewWithClaims(jwt.SigningMethodPS256, jwt.MapClaims{
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
	})
	token.Header["kid"] = "veryUniqueJwtKey"
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(PrivKeyTest))
	if err != nil {
		t.Errorf("unable to parse private key: %v", err)
	}
	ssaJwt, err := token.SignedString(privKey)
	if err != nil {
		t.Errorf("unable to sign jwt: %v", err)
	}
	ssaValidator := ssa.NewSSAValidator(ssa.PublicKeyLookupFromByteSlice([]byte(PubKeyTest)))
	ssaValue, err := ssaValidator.Validate(ssaJwt)
	if err != nil {
		t.Errorf("should not expect error. Got %s", err.Error())
	}
	if ssaValue.Issuer != "OpenBanking Ltd" {
		t.Errorf("Issuer should be Open Banking Ltd. Got %#v", ssaValue.Issuer)
	}
}
