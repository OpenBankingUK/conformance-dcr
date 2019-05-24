package ssa_test

import (
	"testing"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/ssa"
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

func TestValidateSSA(t *testing.T) {
	ssaJwt := `eyJ0eXAiOiJKV1QiLCJhbGciOiJQUzI1NiIsImtpZCI6IjEyMzQ1Njc4OSJ9.eyJpc3MiOiJPcGVuQmFua2luZyBMdGQiLCJpYXQiOjE0OTI3NTYzMzEsImV4cCI6MTU5NTc1NzU1MCwiYXVkIjoiT3BlbkJhbmtpbmcgVFBQIFVuaXF1ZSBJRCIsInN1YiI6Ik9wZW5CYW5raW5nIFRQUCBTb2Z0d2FyZSBVbmlxdWUgSUQiLCJjbGllbnRfbmFtZSI6IkFtYXpvbiBQcmltZSBMdGQiLCJjbGllbnRfdXJpIjoiaHR0cHM6Ly9wcmltZS5hbWF6b24uY29tIiwic29mdHdhcmVfaWQiOiJPcGVuQmFua2luZyBUUFAgU29mdHdhcmUgVW5pcXVlIElEIiwic29mdHdhcmVfcm9sZXMiOlsiUElTUCIsIkFJU1AiXSwib3JnX2lkIjoiT3BlbkJhbmtpbmcgVFBQIFVuaXF1ZSBJRCIsIm9yZ19uYW1lIjoiT3BlbkJhbmtpbmcgVFBQIFJlZ2lzdGVyZWQgTmFtZSIsInRwcF9qd2tzX2VuZHBvaW50IjoiaHR0cHM6Ly9qd2tzLm9wZW5iYW5raW5nLm9yZy51ay90cHBfaWQuamt3cyJ9.e4niy0SYnyr8Nhqoaz3npI1Jkl76soj1z1UYqx5bh8OHIfCpRnxcfi_Qp-udi2A8IV2md8RH2ohr7WkuCvt2NjeR473iika92hPF1Wt8F6HQAG_PmXqe0OFhMFWdNguuPnlRHCHi3HXKizsfWojBNdLOEsTZ5c9HLfo3d1D1-wC_74shgdzGWEQ1fW23Gp4SP4hzrfd6I1MKp3a24GVnP_PJS4vPYJakfdMiIdRIQ773P5cDJchxM9tPGuKHTxcp2G-fiskK_RI0kXvWl4UgJaCfSm8q9ls5YXpu6bEHmo4Usa6DC9GPG7wJ1TIZntkoDgoYPTQISSMz5xc8XkiYRA`
	ssaValidator := ssa.NewSSAValidator(ssa.PublicKeyLookupFromByteSlice([]byte(PubKeyTest)))
	ssaValue, err := ssaValidator.Validate(ssaJwt)
	if err != nil {
		t.Errorf("should not expect error. Got %s", err.Error())
	}
	if ssaValue.Issuer != "OpenBanking Ltd" {
		t.Errorf("Issuer should be Open Banking Ltd. Got %#v", ssaValue.Issuer)
	}
}
