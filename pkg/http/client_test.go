package http

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func certificatesFromPEMBlock(certBlock, keyBlock []byte) []tls.Certificate {
	crt, _ := tls.X509KeyPair(certBlock, keyBlock)
	return []tls.Certificate{crt}
}

func rootCAsFromPEMBlock(certsPEM []byte) *x509.CertPool {
	crtPool := x509.NewCertPool()
	crtPool.AppendCertsFromPEM(certsPEM)
	return crtPool
}

func TestNewMATLSClient(t *testing.T) {
	certBlock := `-----BEGIN CERTIFICATE-----
MIIFODCCBCCgAwIBAgIEWcWVhDANBgkqhkiG9w0BAQsFADBTMQswCQYDVQQGEwJH
QjEUMBIGA1UEChMLT3BlbkJhbmtpbmcxLjAsBgNVBAMTJU9wZW5CYW5raW5nIFBy
ZS1Qcm9kdWN0aW9uIElzc3VpbmcgQ0EwHhcNMTkwNzE4MTAzNzE5WhcNMjAwODE4
MTEwNzE5WjBhMQswCQYDVQQGEwJHQjEUMBIGA1UEChMLT3BlbkJhbmtpbmcxGzAZ
BgNVBAsTEjAwMTU4MDAwMDEwNDFSYkFBSTEfMB0GA1UEAxMWUWV5YjlUQzBJekx5
bXBBOW1Lb1NRMDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMtHBTP+
6Xw6SAY1hsj2tkm5qbVlXPMxGynvRrn00LQezGVmYwtca14uEq+SbfQoW2I7t4gA
KmFiSQIE4ZRAODQfKHfb6lkfCFa7EJxGAJL8e5/K5dS6hgIAc0pp5YB15/kSBSaO
DNkE//3EwN6r1iX95sKQwrmMTLw/1ggluXrTbC/bf17agnHEltjA7YVI/08ZGNj6
TDxidAUBBOvWCxly+eXFdfSkTHFA0facgPlnyY3Xwkb8zQmuPRcgBnnTwODHtzhU
aJO0UR4EU93kTq2Qj/PxrXxbptwRB2Od8uo2SnvyQ4QvFrxSNfKCS6BJytr1hYVA
FJ/P9y3MX++Ow/kCAwEAAaOCAgQwggIAMA4GA1UdDwEB/wQEAwIHgDAgBgNVHSUB
Af8EFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwgeAGA1UdIASB2DCB1TCB0gYLKwYB
BAGodYEGAWQwgcIwKgYIKwYBBQUHAgEWHmh0dHA6Ly9vYi50cnVzdGlzLmNvbS9w
b2xpY2llczCBkwYIKwYBBQUHAgIwgYYMgYNVc2Ugb2YgdGhpcyBDZXJ0aWZpY2F0
ZSBjb25zdGl0dXRlcyBhY2NlcHRhbmNlIG9mIHRoZSBPcGVuQmFua2luZyBSb290
IENBIENlcnRpZmljYXRpb24gUG9saWNpZXMgYW5kIENlcnRpZmljYXRlIFByYWN0
aWNlIFN0YXRlbWVudDBtBggrBgEFBQcBAQRhMF8wJgYIKwYBBQUHMAGGGmh0dHA6
Ly9vYi50cnVzdGlzLmNvbS9vY3NwMDUGCCsGAQUFBzAChilodHRwOi8vb2IudHJ1
c3Rpcy5jb20vb2JfcHBfaXNzdWluZ2NhLmNydDA6BgNVHR8EMzAxMC+gLaArhilo
dHRwOi8vb2IudHJ1c3Rpcy5jb20vb2JfcHBfaXNzdWluZ2NhLmNybDAfBgNVHSME
GDAWgBRQc5HGIXLTd/T+ABIGgVx5eW4/UDAdBgNVHQ4EFgQUkUBtOvDDZc8k+iHy
BzyB2oW4IoEwDQYJKoZIhvcNAQELBQADggEBAF8kleb9YFqu5OJdposfnAEJj/YX
JMi8rbYHq6V4+ETjeFQIXgmfTHss7sO6g0kNc97DslXUHrnDl2+XWB3e7uEWmxZN
t62J5WrWwHaa49VdMoKyHbg7bZICIVpS9+VX4J7F01ksqc523t3JMYMYfMDfsjue
b4X9tCaBwpmO100ChlcZMjkZAWSGejEHSSXgzE+BOAgEEpB9rDWd1WjgIl8Og/+P
Y68x95Z+TBiYPT6wOd9+xrhsHvf8VZnFHuXpWqyrjOb79yJLZ7trJOll/8HFlHdu
VaQXGsIVwrZ3N6M5qGIPYXhS6FtcGX4J0t+s6bFhH9GjBb3E+gkhWxrlGPM=
-----END CERTIFICATE-----
`
	keyBlock := `-----BEGIN PRIVATE KEY-----
MIIEwAIBADANBgkqhkiG9w0BAQEFAASCBKowggSmAgEAAoIBAQDLRwUz/ul8OkgG
NYbI9rZJuam1ZVzzMRsp70a59NC0HsxlZmMLXGteLhKvkm30KFtiO7eIACphYkkC
BOGUQDg0Hyh32+pZHwhWuxCcRgCS/HufyuXUuoYCAHNKaeWAdef5EgUmjgzZBP/9
xMDeq9Yl/ebCkMK5jEy8P9YIJbl602wv239e2oJxxJbYwO2FSP9PGRjY+kw8YnQF
AQTr1gsZcvnlxXX0pExxQNH2nID5Z8mN18JG/M0Jrj0XIAZ508Dgx7c4VGiTtFEe
BFPd5E6tkI/z8a18W6bcEQdjnfLqNkp78kOELxa8UjXygkugScra9YWFQBSfz/ct
zF/vjsP5AgMBAAECggEBAMElr0Tzfmjie9MigvpqqVQuDJnmQUZ2L2hOCbbhbnjn
FrF2cr+1xfML9gdLLBylSAdz1HRuK9aya8p2VtzllEG6JVqV4/tgaCz4u9SxExSW
wORZBr51qKU4RlO7gSzpW0wnGivDJ2QBwzcd/2DUh7s7oErY/50MOVFZmoXNSorv
rFfTp9ToD4GRcelENGOuECzuoLaLXo/G7oLJaQRcTs70NPcmUvkFiDrgde111bUD
kkHU4nthlUD59mqiV4U+AzmXEyN3la5/I9cpd8XATPt+jGbKdIyc3dnZW84cIhmk
4Ko3y03kYSeL56HqPyrm9vFxhFy/sKs2zA7qsC2Y5AECgYEA85RXxzlWQQoTwjww
JP0+ek+iwamoORgJPn4peGuVyrmiSJtYI7lK9yzCL5vpZ/OrAWWnY+WDgdpk1AI7
00NVsc8+nUR+bDKzxceeiFhYYbTPUkqLvN/35Kz6hAezXSlXq8/lC1Tm5g1NcoE3
zxgX+qT9fwgIerQUQiN7Kc4J1IECgYEA1aSUVWJZNPRDvt+buWpepAQRRYPJQhb4
E2kBFGBtCFbj2KgkQkeMa5OdBAiR2+7BZLNTyMDf9wImh9WCDladM81vP34TW/s5
q5ROC7C6mHYfpt/vp60euDQU/2AkjNM4sz1E2SV8k8bripVKwvVAl931EwweFhZw
7x5TqlVD03kCgYEAx8g/KsdraJMUW7a0IlKYAQf6TW+S66k8Q8aEyyEqzgjuAzFu
3HYo94z9hMETctCXzOCMp9HiyAnRs1ZVrVTIH7wE9kbsjmATtT+iVuBnNVRwy2Ub
MgJdN3FtVAdg5SN4phIxIdc0PzJf+G/lz3VKjajvxlZXZhT3nLuvVD2LMIECgYEA
mxUWDAkRQmxRxPuimeyJ+LtvIivw60WrHMPrYbRBUX1pdbtQXsB7QRftMaFa5/Cf
eA7oseC4cyCfgZjOCMR85r6ok8lcGjf6e/9yy++k88lDXqpN2ETF+ObtmxdaUNN+
5DWEhbA9hzQthPKsS2smUVdwcDwqltQBdMJp70pnqtkCgYEA8kOBV2eKedF5oTNG
N/ReqhSWneuWsd7JMuIP4y+tVy+/5WcEe82PIh9QcyYKOl3T65YTyJ5ndASM1Tz4
1rIsiOfeA2XlQrDznv+Bp/0Ixm2O+qsXJY+DuchiERdumdh5Q2v3l18rlCmRjp7i
1MZ+vNVBoKIUqWabrTYooOpTvBM=
-----END PRIVATE KEY-----
`
	rootCABlock := `-----BEGIN CERTIFICATE-----
MIIFQjCCAyqgAwIBAgIEWgGaAzANBgkqhkiG9w0BAQsFADBBMQswCQYDVQQGEwJH
QjEUMBIGA1UEChMLT3BlbkJhbmtpbmcxHDAaBgNVBAMTE09wZW5CYW5raW5nIFJv
b3QgQ0EwHhcNMTcxMTA3MTEwOTM2WhcNMzcxMTA3MTEzOTM2WjBBMQswCQYDVQQG
EwJHQjEUMBIGA1UEChMLT3BlbkJhbmtpbmcxHDAaBgNVBAMTE09wZW5CYW5raW5n
IFJvb3QgQ0EwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQDvJkaqdaIU
NgTnXcJ3lKLyjhTJSsNtYzmN7fvpn8oseBQXQDKzJAvLXhfUEVeuUu3Zv/TG+ab/
pSFdtiibh5PLIbB8nQDORl/fAA68wIjImsa2feUcq91Y+dKdKN8iW6zop8aDL8qw
EggAV/u3TRfOhF8LSKHOEZ/7/YRTuqinAxDkeHYh7G+uSReyP4NvehhDkSuhK44z
byEddOvvcAOrkYr9TtBj6iZ5OMVZGO9tY9gRkbiQOt1FozyuYB7XT0QzokIfBWE0
CZ1ypdu2bttDC7CuVhw9QSnyFHIG6HtQi2zKZH9OceMPJiG9RAdBUDZ3qqLFEVSv
w1Dgfu/iatPEgYTbRDA85EHeGCcTMCTGra0eoITekrq//CRW1e73lK40SFzmMK/l
KD3B2qWz/TxMvEH186s5REKPC6ptiQ4TxIp8Ls4gn2UHGwbS7i9ihryr0/ww9ILz
y3gkuahf1t6PaNwmU02dovfLG5LJrMnvn8P6SdPwgbt3TtMKPBTxawQK+4N7wcY3
slvh6bj9XLdyYKkqAk5QDiGoyZypZ6iH6P40gxJgJquF3kgYTSWunWkylDC6QgUU
5U+x43SorH3qBB/fN5+daI8PQo80gbvonnWDAelxMkNUTkt/469CBpOd0Ok5uhl6
g1cb9Tl1i3IR1c3Daa1hHK2eoKfsOMjwVwIDAQABo0IwQDAOBgNVHQ8BAf8EBAMC
AQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUKnp9q41DYWT2XgHagTSVGFM8
ny4wDQYJKoZIhvcNAQELBQADggIBAGX85+GSIXLJhZ6FwXZgrm+jKvdzxWP3qkwE
hNmxfA3Cl4oVzINkn8fQfz3LN9zwTqRusxXfdpSdxfMesB480sDUDy88VAIdNi5A
1DFFL02qZJxOH5cBRN+VVRPfRLSXK56LlbItM38GdhRVhd0FVnpG9+tqkmseF63r
DCP30BOidUEH1Ong+0Bt8vZOs/OcPyGswsQJS3/7I1QFPxm/0F7wwBxdZwODcz4T
Amw9EpePgNvI7ayhM7V/krMJeyG1bQ1sXu7LWdQIEEavrnV0fGgWPbG9L1QzhIxO
5PzUKsA09W3wweRVQJxcYRWw3L1orwrvKZktvsKq1K7PEsIzHd3N/L+gGNDdYCZg
eL+uv4aIoArPvJa06bVBSiunmkN4LuSRv0pVQPXkNzNkeTgJuCqE8DQavkjDY6Ov
hTjL54LGT8cv8wrgL9ZZWiol+LYABiF3ffdS7uXNAMEmHTAniBsw6t4VmoT6sjDD
7Y4QLG7mJ53MIFbBb/+Y3IJQj474Yl9bOk3lbEJ8fSj1DtuRrygxDjUFZ2Iqbuli
LN86nN9SMIr+WZBAIG3bT3I8EkAvVPPHiWXjZZV/oBQq3C4fZT7ELu1Y2Z4h3Z/O
W3/8OHbqKHnXS9MsOvJ1cVHHb/dRAeg2iKLbVikYKQM5mShYIJ0zIxKS7I/UKU5f
YtfkskMi
-----END CERTIFICATE-----
`

	type args struct {
		config MATLSConfig
	}
	tests := []struct {
		name       string
		args       args
		wantClient *http.Client
		wantErr    bool
	}{
		{
			name: "Correct setup",
			args: args{
				config: MATLSConfig{
					insecureSkipVerify: true,
					keyPEMBlock:        []byte(keyBlock),
					certPEMBlock:       []byte(certBlock),
					caCerts:            []byte(rootCABlock),
				},
			},
			wantClient: &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
						Certificates:       certificatesFromPEMBlock([]byte(certBlock), []byte(keyBlock)),
						MinVersion:         tls.VersionTLS12,
						Renegotiation:      tls.RenegotiateFreelyAsClient,
						RootCAs:            rootCAsFromPEMBlock([]byte(rootCABlock)),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid root ca certs",
			args: args{
				config: MATLSConfig{
					insecureSkipVerify: true,
					keyPEMBlock:        []byte(keyBlock),
					certPEMBlock:       []byte(certBlock),
					caCerts:            []byte{},
				},
			},
			wantClient: &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
						Certificates:       certificatesFromPEMBlock([]byte(certBlock), []byte(keyBlock)),
						MinVersion:         tls.VersionTLS12,
						Renegotiation:      tls.RenegotiateFreelyAsClient,
						RootCAs:            rootCAsFromPEMBlock([]byte("-- bad root ca certs --")),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid cert block",
			args: args{
				config: MATLSConfig{
					insecureSkipVerify: true,
					keyPEMBlock:        []byte(keyBlock),
					certPEMBlock:       []byte{},
					caCerts:            []byte(rootCABlock),
				},
			},
			wantClient: &http.Client{},
			wantErr: true,
		},
		{
			name: "Invalid key block",
			args: args{
				config: MATLSConfig{
					insecureSkipVerify: true,
					keyPEMBlock:        []byte{},
					certPEMBlock:       []byte(certBlock),
					caCerts:            []byte(rootCABlock),
				},
			},
			wantClient: &http.Client{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMATLSClient(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMATLSClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			trsActual := got.Transport.(*http.Transport)
			trsExpected := tt.wantClient.Transport.(*http.Transport)

			assert.Equal(t, trsExpected.TLSClientConfig.MinVersion, trsActual.TLSClientConfig.MinVersion)
			assert.Equal(t, trsExpected.TLSClientConfig.InsecureSkipVerify, trsActual.TLSClientConfig.InsecureSkipVerify)
			assert.Equal(t, trsExpected.TLSClientConfig.Renegotiation, trsActual.TLSClientConfig.Renegotiation)
			assert.Equal(t, trsExpected.TLSClientConfig.Certificates, trsActual.TLSClientConfig.Certificates)

			expectedRootCAs := trsExpected.TLSClientConfig.RootCAs.Subjects()
			actualRootCAs := trsActual.TLSClientConfig.RootCAs.Subjects()
			assert.Equal(t, expectedRootCAs, actualRootCAs)
		})
	}
}
