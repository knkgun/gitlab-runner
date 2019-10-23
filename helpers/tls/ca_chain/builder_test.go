package ca_chain

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testCACert = `-----BEGIN CERTIFICATE-----
MIIFjTCCA3WgAwIBAgIUdC7ewPrKJksR4FvSUhjdtolff6IwDQYJKoZIhvcNAQEL
BQAwVTELMAkGA1UEBhMCVVMxCjAIBgNVBAgMASAxCjAIBgNVBAoMASAxCjAIBgNV
BAsMASAxEDAOBgNVBAMMB1Rlc3QgQ0ExEDAOBgkqhkiG9w0BCQEWASAwIBcNMTkx
MDE4MDU1NzI5WhgPMjExOTA5MjQwNTU3MjlaMFUxCzAJBgNVBAYTAlVTMQowCAYD
VQQIDAEgMQowCAYDVQQKDAEgMQowCAYDVQQLDAEgMRAwDgYDVQQDDAdUZXN0IENB
MRAwDgYJKoZIhvcNAQkBFgEgMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKC
AgEArXISLnSKP2Az5LDx9PSBgnca8Rwu3wA6EoK5YEB01M21TS2PlOmF8pls1Ojl
d8OiSbiio8clhERikUsj6/schKXIv7JX0paqmSbMi++VRimXz8LakTBj58QAV53p
fnPc6InbSVXdq1jK8HIh1/8zFBbeMaZTTeV3cuX3Ue0kXWRUPtHKuJor6vksYgGS
GI4kLM5N7PMfgLQlCc4bVxXqst2HZvimPOpL5DZAYg8fEz3EIqXyIgQfxSLCcUWs
mELhPP1XD3hkPPlc1pCL/ANmNEw0bU0TLuh3h7i+cC0yVE9xKne3v1HkdmnsUiBC
gJzmqlAvb1PbVUmpubvCimuC8nvJbuQYZfglqIuRVtGOnPkpAOeyxTdbA2bvZA8L
8fj7mdnCJIOOKqdfW/Nh2TpSTcL++pHW1qW5M4I8v9y/NE3+t42ur4VMLXkFyFrS
Ygm1Jsi9+qht0q0YllaEmpXCthD+uxlulMBsrUZHZ9T8nPPVXHzEF4DHEnYWWeco
emuz+uksIn2Jlh7FZIjUHfIhtkK3Gxw9xgSrhirdfP5lSBb1qUe+d1jZWo+t9Ftj
gS4FDFmN5uZlNLNs6LutB2gHxaGcSgtZ73shgp6sOpCDU7OxyLzdNjWdQy0MM50M
cuaOfMKhJaWFqn9pQbQAWeUkouUKYvLIky2bjZalqg2M+A8CAwEAAaNTMFEwHQYD
VR0OBBYEFCtSc7nrSk/ugFmuO+/A8BvkYT95MB8GA1UdIwQYMBaAFCtSc7nrSk/u
gFmuO+/A8BvkYT95MA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggIB
AAl2Ohrfi6ZCF3kdAUG3j5ujQpMkPvVyxWRHf/Nyef9TBcWOQdVpT47ckW1QvyPO
U/+XsTy/3+paZuejWnG/t44ITz+Zilt4cpby1GcQOWLZzlTVciL8wPiUA+P8AD8s
yZ5Sk6rBQBooMWKOrzNA3OdMEe5NbMT0//TrzJHu5mMKZierYzhBPo22SH3Onwwq
icypW8DLKpJIp1r7JWquVWiux4349Y514tH5Hn3lq5C3k21ioYuXrg5zlUz5sTx2
9T09DmyNu1GF+UYF85gyc6rBTQFMBi/ZX8GGG709lAgdcDd46O1rI32DIpzn9XMo
O6vk58UIbedbdjPeURx1+qa39tR6jVURodTNLqbzhusNmSzJHxNtOtCa5ygFOUUJ
oMiMvSitZ+HbPPjsS8uXq+c0/08HYqODidw5DGj/KzhwCfIl2gKn4k4ikWWD9OED
54eTRpt6m0SCLXRfIWSLLJoU7AlqZ9jvenH/9vtuMPG1IXc3/YISacqxBZq/yfI9
nJu5mzOPRdKPVcI/I+0Bqnqg1x7cMf7kkippUg+GygL24hLw5xVrcyembk6ca9RH
Jrz2TngQylcfjMtWKTvn9TcRuCgYy5CRYSm9+ZphpsQdYpmQG5278q2lKH3AvIo1
pmNh6pRdOvIQX2i8UFDrD+tD7qSYciwRrEJbp1mc6zfw
-----END CERTIFICATE-----`
	testCert = `-----BEGIN CERTIFICATE-----
MIIEEjCCAfoCFBhRTszftYHtN+HOfbU/q3zvYBYOMA0GCSqGSIb3DQEBCwUAMFUx
CzAJBgNVBAYTAlVTMQowCAYDVQQIDAEgMQowCAYDVQQKDAEgMQowCAYDVQQLDAEg
MRAwDgYDVQQDDAdUZXN0IENBMRAwDgYJKoZIhvcNAQkBFgEgMCAXDTE5MTAxODA2
MDA1MloYDzIxMTkwOTI0MDYwMDUyWjA0MQswCQYDVQQGEwJVUzEKMAgGA1UECAwB
IDEKMAgGA1UECgwBIDENMAsGA1UEAwwEdGVzdDCCASIwDQYJKoZIhvcNAQEBBQAD
ggEPADCCAQoCggEBALc0+Xo61c0xCvebNg1OJl4iXC5blzGlbDfejWKn7266g+UU
Z3xscCDWMNruojd+7EbkQmAyUtdGifNw+xIHyNA/jiyIsB3KteN84X+toA4mjY1t
SpqlNMOUW0EZ9f0KZNn4GZnA/TyFWI3EC4gOcJyuuL7YfE7Qu1e3LeBwDcRYpJ3W
Zw1k3+aClC1N7iTPEP9scr64+KA0d5xIkrtl5t8qiSR8Tn+JLPygGre0G0hhIZeH
pfPQWX6iILbJMgPnbPmCivklkyUIE8WHh2qGbOGaO3LVKSS6/YfOshw4g/RQyusI
Ii65iXnFa/VvRY2dkn5w9EehZzbT8kQa7U39NwkCAwEAATANBgkqhkiG9w0BAQsF
AAOCAgEAMAfp7FRBHm9t4byRfWrUYblI7eQOlcixXHSPc16VX93HTsNWwZV1EBiO
GWcTRcts5FQr9HWGHVukQ+4iXLWtb/Og+hHrjyLmOGvx7sgPeHuyWB89npSABden
rpMHPePMzsO/YTw1QuYJOijNYpLCL83YWk62DCSGwQ2HO1KKLDw3suBHudV80cHV
nav7Q0VW+iA+3apdrgediCHCtc6PQDHPzdrXQSVA+OF2itX3Xhc6Mm3dn4D3Hhqo
WYJNeI0naNHTguoKFYdJHHjv07nX+1I+CAk6kjEv17VEKsU7SjhOizLYdtb9OrOS
gnQ6KTkPfCeIlK2PNguwxgeLBNYQyTnUxr1QxgVkKFsBfwFV4hq9podEbjrgUSu1
KZSdU7u7WMCjLYpyC5kbRmd/Qkdo/45wifomJNP3/16NSNZ0gatKVUJ6q6UjRsZl
3va4QcB3QuNtGiQZqEuc/+KM21MSvC8cC/bIOaKZlWbKtEV+tsbuIIhng0opJrEw
+5ZqVqrwIVjbsGaw/NPROth/XDJp5jzpwxnf5HDQhLV04sfdN9IRw005WC+l0f19
iG9V6qslKJvNR8A8A+RqvyfIJ0gjNzVLQHrZyTsEbC62w1IcxkBG7lR6W7ZCXal1
RSKf+3OIln1a6DKx+zEzL20uwW5L/5l3FsLwwvOLybX4mAhiyxY=
-----END CERTIFICATE-----`
)

func TestDefaultBuilder_FetchCertificatesFromVerifiedChain(t *testing.T) {
	testError := errors.New("test-error")

	block, _ := pem.Decode([]byte(testCert))
	testCertificate, err := x509.ParseCertificate(block.Bytes)
	require.NoError(t, err)

	block, _ = pem.Decode([]byte(testCACert))
	testCACertificate, err := x509.ParseCertificate(block.Bytes)
	require.NoError(t, err)

	tests := map[string]struct {
		chains                    [][]*x509.Certificate
		fetchCertificateChainMock certificateChainFetcher
		addRootCAMock             rootCAAdder
		expectedError             string
		expectedChainLength       int
	}{
		"no chains": {
			chains:              [][]*x509.Certificate{},
			expectedChainLength: 0,
		},
		"empty chain": {
			chains:              [][]*x509.Certificate{{}},
			expectedChainLength: 0,
		},
		"error on chain fetching": {
			chains: [][]*x509.Certificate{{testCertificate}},
			fetchCertificateChainMock: func(cert *x509.Certificate) ([]*x509.Certificate, error) {
				return nil, testError
			},
			expectedError:       "error while fetching certificates into the CA Chain: couldn't fetch certificates chain: test-error",
			expectedChainLength: 0,
		},
		"error on root CA adding": {
			chains: [][]*x509.Certificate{{testCertificate}},
			fetchCertificateChainMock: func(cert *x509.Certificate) ([]*x509.Certificate, error) {
				return []*x509.Certificate{cert}, nil
			},
			addRootCAMock: func(certs []*x509.Certificate) ([]*x509.Certificate, error) {
				return nil, testError
			},
			expectedError:       "error while fetching certificates into the CA Chain: couldn't add root CA to the chain: test-error",
			expectedChainLength: 0,
		},
		"certificates chain prepared properly": {
			chains: [][]*x509.Certificate{{testCertificate}},
			fetchCertificateChainMock: func(cert *x509.Certificate) ([]*x509.Certificate, error) {
				return []*x509.Certificate{cert}, nil
			},
			addRootCAMock: func(certs []*x509.Certificate) ([]*x509.Certificate, error) {
				return append(certs, testCACertificate), nil
			},
			expectedChainLength: 2,
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			var err error

			builder := NewBuilder().(*defaultBuilder)
			builder.fetchCertificateChain = tc.fetchCertificateChainMock
			builder.addRootCA = tc.addRootCAMock

			TLS := new(tls.ConnectionState)
			TLS.VerifiedChains = tc.chains

			err = builder.FetchCertificatesFromTLSConnectionState(TLS)

			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, builder.certificates, tc.expectedChainLength)
		})
	}
}

func TestDefaultBuilder_AddCertificate(t *testing.T) {
	block, _ := pem.Decode([]byte(testCert))
	testCertificate, err := x509.ParseCertificate(block.Bytes)
	require.NoError(t, err)

	b := NewBuilder().(*defaultBuilder)
	b.addCertificate(testCertificate)
	b.addCertificate(testCertificate)

	require.Len(t, b.certificates, 1)
	assert.Equal(t, testCertificate, b.certificates[0])
}

func TestDefaultBuilder_String(t *testing.T) {
	testError := errors.New("test-error")

	block, _ := pem.Decode([]byte(testCert))
	testCertificate, err := x509.ParseCertificate(block.Bytes)
	require.NoError(t, err)

	tests := map[string]struct {
		encodePEMMock        pemEncoder
		expectedOutput       string
		expectedLogToContain []string
	}{
		"encoding error": {
			encodePEMMock: func(out io.Writer, b *pem.Block) error {
				return testError
			},
			expectedOutput: "",
			expectedLogToContain: []string{
				"error=test-error",
				`msg="Failed to encode certificate from chain"`,
			},
		},
		"encoding succeeded": {
			encodePEMMock: func(out io.Writer, b *pem.Block) error {
				assert.Equal(t, pemTypeCertificate, b.Type)
				assert.Equal(t, testCertificate.Raw, b.Bytes)

				buf := bytes.NewBufferString(testCert)

				_, err := io.Copy(out, buf)
				require.NoError(t, err)

				return nil
			},
			expectedOutput: testCert,
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			out := new(bytes.Buffer)

			logger := logrus.New()
			logger.Out = out

			b := NewBuilder().(*defaultBuilder)
			b.encodePEM = tc.encodePEMMock
			b.logger = logger

			b.addCertificate(testCertificate)
			assert.Equal(t, tc.expectedOutput, b.String())

			output := out.String()

			if len(tc.expectedLogToContain) < 1 {
				assert.Empty(t, output)
				return
			}

			for _, part := range tc.expectedLogToContain {
				assert.Contains(t, output, part)
			}
		})
	}
}