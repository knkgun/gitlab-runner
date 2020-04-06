package kubernetes

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	k8sversion "k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/rest/fake"
)

func TestKubeClientFeatureChecker(t *testing.T) {
	kubeClientErr := errors.New("clientErr")

	type expectedCheckResult struct {
		supported bool
		err       error
	}

	version, _ := testVersionAndCodec()
	tests := map[string]struct {
		version                    k8sversion.Info
		clientErr                  error
		expectedHostAliasesResult  expectedCheckResult
		expectedRuntimeClassResult expectedCheckResult
	}{
		"version 1.7": {
			version: k8sversion.Info{
				Major: "1",
				Minor: "7",
			},
			clientErr: nil,
			expectedHostAliasesResult: expectedCheckResult{
				supported: true,
				err:       nil,
			},
			expectedRuntimeClassResult: expectedCheckResult{
				supported: false,
				err:       nil,
			},
		},
		"version 1.11": {
			version: k8sversion.Info{
				Major: "1",
				Minor: "11",
			},
			clientErr: nil,
			expectedHostAliasesResult: expectedCheckResult{
				supported: true,
				err:       nil,
			},
			expectedRuntimeClassResult: expectedCheckResult{
				supported: false,
				err:       nil,
			},
		},
		"version 1.6": {
			version: k8sversion.Info{
				Major: "1",
				Minor: "6",
			},
			clientErr: nil,
			expectedHostAliasesResult: expectedCheckResult{
				supported: false,
				err:       nil,
			},
			expectedRuntimeClassResult: expectedCheckResult{
				supported: false,
				err:       nil,
			},
		},
		"version 1.6 cleanup": {
			version: k8sversion.Info{
				Major: "1+535111",
				Minor: "6.^&5151111",
			},
			clientErr: nil,
			expectedHostAliasesResult: expectedCheckResult{
				supported: false,
				err:       nil,
			},
		},
		"version 1.14": {
			version: k8sversion.Info{
				Major: "1",
				Minor: "14",
			},
			clientErr: nil,
			expectedHostAliasesResult: expectedCheckResult{
				supported: true,
				err:       nil,
			},
			expectedRuntimeClassResult: expectedCheckResult{
				supported: true,
				err:       nil,
			},
		},
		"version 1.14 cleanup": {
			version: k8sversion.Info{
				Major: "1*)(535111",
				Minor: "14^^%&5151111",
			},
			clientErr: nil,
			expectedHostAliasesResult: expectedCheckResult{
				supported: true,
				err:       nil,
			},
			expectedRuntimeClassResult: expectedCheckResult{
				supported: true,
				err:       nil,
			},
		},
		"invalid version with leading characters": {
			version: k8sversion.Info{
				Major: "+1",
				Minor: "-14",
			},
			clientErr: nil,
			expectedHostAliasesResult: expectedCheckResult{
				supported: false,
				err:       new(badVersionError),
			},
			expectedRuntimeClassResult: expectedCheckResult{
				supported: false,
				err:       new(badVersionError),
			},
		},
		"invalid version": {
			version: k8sversion.Info{
				Major: "aaa",
				Minor: "bbb",
			},
			clientErr: nil,
			expectedHostAliasesResult: expectedCheckResult{
				supported: false,
				err:       new(badVersionError),
			},
			expectedRuntimeClassResult: expectedCheckResult{
				supported: false,
				err:       new(badVersionError),
			},
		},
		"empty version": {
			version: k8sversion.Info{
				Major: "",
				Minor: "",
			},
			clientErr: nil,
			expectedHostAliasesResult: expectedCheckResult{
				supported: false,
				err:       new(badVersionError),
			},
			expectedRuntimeClassResult: expectedCheckResult{
				supported: false,
				err:       new(badVersionError),
			},
		},
		"kube client error": {
			version: k8sversion.Info{
				Major: "",
				Minor: "",
			},
			clientErr: kubeClientErr,
			expectedHostAliasesResult: expectedCheckResult{
				supported: false,
				err:       kubeClientErr,
			},
			expectedRuntimeClassResult: expectedCheckResult{
				supported: false,
				err:       kubeClientErr,
			},
		},
	}

	for tn, tt := range tests {
		rt := func(request *http.Request) (response *http.Response, err error) {
			if tt.clientErr != nil {
				return nil, tt.clientErr
			}

			ver, _ := json.Marshal(tt.version)
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body: FakeReadCloser{
					Reader: bytes.NewReader(ver),
				},
			}
			resp.Header = make(http.Header)
			resp.Header.Add("Content-Type", "application/json")

			return resp, nil
		}
		fc := kubeClientFeatureChecker{
			kubeClient: testKubernetesClient(version, fake.CreateHTTPClient(rt)),
		}

		t.Run("host aliases "+tn, func(t *testing.T) {
			supported, err := fc.IsHostAliasSupported()
			assert.Equal(t, tt.expectedHostAliasesResult.supported, supported)
			assert.True(t, errors.Is(err, tt.expectedHostAliasesResult.err))
		})

		t.Run("runtime class "+tn, func(t *testing.T) {
			supported, err := fc.IsRuntimeClassSupported()
			assert.Equal(t, tt.expectedRuntimeClassResult.supported, supported)
			assert.True(t, errors.Is(err, tt.expectedRuntimeClassResult.err))
		})
	}
}
