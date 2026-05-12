package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Signer interface {
	Claims() (string, error)
}

type jwtSigner struct {
	signingAlgorithm        jwt.SigningMethod
	ssa                     string
	issuer                  string
	audience                string
	kID                     string
	tokenEndpointAuthMethod string
	requestObjectSignAlg    string
	redirectURIs            []string
	responseTypes           []string
	privateKey              *rsa.PrivateKey
	jwtExpiration           time.Duration
	transportCert           *x509.Certificate
	transportSubjectDn      string
	useOID                  bool
}

func NewJwtSigner(
	signingAlgorithm jwt.SigningMethod,
	ssa,
	issuer,
	audience,
	kID,
	tokenEndpointAuthMethod string,
	requestObjectSignAlg string,
	redirectURIs []string,
	responseTypes []string,
	privateKey *rsa.PrivateKey,
	jwtExpiration time.Duration,
	transportCert *x509.Certificate,
	transportSubjectDn string,
	useOID bool,
) Signer {
	return jwtSigner{
		signingAlgorithm:        signingAlgorithm,
		ssa:                     ssa,
		issuer:                  issuer,
		audience:                audience,
		kID:                     kID,
		tokenEndpointAuthMethod: tokenEndpointAuthMethod,
		requestObjectSignAlg:    requestObjectSignAlg,
		redirectURIs:            redirectURIs,
		responseTypes:           responseTypes,
		privateKey:              privateKey,
		jwtExpiration:           jwtExpiration,
		transportCert:           transportCert,
		transportSubjectDn:      transportSubjectDn,
		useOID:                  useOID,
	}
}

func (s jwtSigner) Claims() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", errors.Wrap(err, "generating claims")
	}

	iat := time.Now().UTC()
	exp := iat.Add(s.jwtExpiration)
	claims := jwt.MapClaims{
		// This should be the unique identifier for the ASPSP
		// issued by the issuer of the software statement.
		// An ASPSP processing the software statement may validate the
		// value of the claim and reject software statements for which the ASPSP is not the audience.
		// The value must be a Base62 encoded GUID.
		"aud": s.audience,

		"exp": exp.Unix(),
		"jti": id.String(),
		"iat": iat.Unix(),

		// Identifier for the TPP.
		// This value must be unique for each TPP registered by the issuer of the SSA.
		// The value must be a Base62 encoded GUID.
		// For SSAs issued by the OB Directory, this must be the software_id
		"iss": s.issuer,

		// metadata

		"grant_types": []string{
			"authorization_code",
			"client_credentials",
		},

		"application_type":             "web",
		"redirect_uris":                s.redirectURIs,
		"token_endpoint_auth_method":   s.tokenEndpointAuthMethod,
		"software_statement":           s.ssa,
		"scope":                        "accounts openid",
		"request_object_signing_alg":   s.requestObjectSignAlg,
		"id_token_signed_response_alg": s.signingAlgorithm.Alg(),
	}

	if s.responseTypes != nil {
		claims["response_types"] = s.responseTypes
	}

	// Instead of potentially custom ASN/OID parsing to get exact, expected value of Subject DN
	// we use a config entry
	if err = s.addTlsClientAuthClaims(claims); err != nil {
		return "", err
	}

	s.addSigningAlgClaims(claims)

	token := jwt.NewWithClaims(s.signingAlgorithm, claims)
	token.Header["kid"] = s.kID

	signedJwt, err := token.SignedString(s.privateKey)
	if err != nil {
		return "", errors.Wrap(err, "signing claims")
	}

	return signedJwt, nil
}

// oidNames maps OID strings to their attribute type names for Subject DN formatting.
var oidNames = map[string]string{
	"2.5.4.3":  "CN",
	"2.5.4.6":  "C",
	"2.5.4.7":  "L",
	"2.5.4.8":  "ST",
	"2.5.4.10": "O",
	"2.5.4.11": "OU",
	"2.5.4.97": "organizationIdentifier",
}

// escapeRFC2253 escapes special characters in a DN attribute value per RFC 2253.
func escapeRFC2253(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == ',' || c == '+' || c == '"' || c == '\\' || c == '<' || c == '>' || c == ';':
			b.WriteByte('\\')
			b.WriteByte(c)
		case c == ' ' && (i == 0 || i == len(s)-1):
			b.WriteByte('\\')
			b.WriteByte(c)
		case c == '#' && i == 0:
			b.WriteByte('\\')
			b.WriteByte(c)
		default:
			b.WriteByte(c)
		}
	}
	return b.String()
}

// subjectDN formats a raw ASN.1 subject into an RFC 2253 string preserving wire order.
// Multi-valued RDNs are joined with '+'; RDNs are joined with ','.
// When useOID is true, all attribute type labels use the numeric OID string;
// otherwise, known OIDs are replaced with their friendly names.
func subjectDN(rawSubject []byte, useOID bool) (string, error) {
	var rdnSeq pkix.RDNSequence
	rest, err := asn1.Unmarshal(rawSubject, &rdnSeq)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse transport cert raw subject")
	}
	if len(rest) > 0 {
		return "", errors.Errorf("failed to parse transport cert raw subject: trailing data (%d bytes)", len(rest))
	}

	// RDNSequence is in wire order (most general first).
	// Iterate in reverse to produce CN first, C last (RFC 2253 display order).
	rdnParts := make([]string, 0, len(rdnSeq))
	for i := len(rdnSeq) - 1; i >= 0; i-- {
		atvParts := make([]string, 0, len(rdnSeq[i]))
		for _, atv := range rdnSeq[i] {
			oidStr := atv.Type.String()
			var label string
			if useOID {
				label = oidStr
			} else if name, ok := oidNames[oidStr]; ok {
				label = name
			} else {
				label = oidStr
			}
			atvParts = append(atvParts, label+"="+escapeRFC2253(fmt.Sprintf("%v", atv.Value)))
		}
		rdnParts = append(rdnParts, strings.Join(atvParts, "+"))
	}
	return strings.Join(rdnParts, ","), nil
}

func (s jwtSigner) addSigningAlgClaims(claims jwt.MapClaims) {
	// We should only provide signing alg when it makes sense
	if s.tokenEndpointAuthMethod == "private_key_jwt" || s.tokenEndpointAuthMethod == "client_secret_jwt" {
		claims["token_endpoint_auth_signing_alg"] = s.signingAlgorithm.Alg()
	}
}

func (s jwtSigner) addTlsClientAuthClaims(claims jwt.MapClaims) error {
	if s.tokenEndpointAuthMethod != "tls_client_auth" {
		return nil
	}

	if s.transportCert == nil {
		return errors.New("transport cert not available")
	}

	if s.transportSubjectDn != "" {
		claims["tls_client_auth_subject_dn"] = s.transportSubjectDn
	} else if len(s.transportCert.RawSubject) == 0 {
		// Fallback for manually constructed certs (e.g. in tests) where RawSubject
		// is not populated from parsed ASN.1 bytes.
		claims["tls_client_auth_subject_dn"] = s.transportCert.Subject.ToRDNSequence().String()
	} else {
		dn, err := subjectDN(s.transportCert.RawSubject, s.useOID)
		if err != nil {
			return err
		}
		claims["tls_client_auth_subject_dn"] = dn
	}

	return nil
}
