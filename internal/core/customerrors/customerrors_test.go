package customerrors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrStatusCode(t *testing.T) {
	err := ErrStatusCode{StatusCode: 404}
	require.Equal(t, 404, err.StatusCode)
	require.Equal(t, "status code: 404", err.Error())
}

func TestErrStatusCodeDifferentCodes(t *testing.T) {
	testCases := []struct {
		statusCode int
		expected   string
	}{
		{200, "status code: 200"},
		{400, "status code: 400"},
		{401, "status code: 401"},
		{403, "status code: 403"},
		{404, "status code: 404"},
		{500, "status code: 500"},
	}

	for _, tc := range testCases {
		err := ErrStatusCode{StatusCode: tc.statusCode}
		require.Equal(t, tc.expected, err.Error())
	}
}

func TestErrStatusCodeUnwrap(t *testing.T) {
	err := ErrStatusCode{StatusCode: 500}
	require.True(t, errors.Is(err, err))
}
