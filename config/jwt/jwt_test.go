package jwtconfig

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	models "linkaja.com/e-wallet/lib/base_models"
)

func Test_CreateToken(t *testing.T) {
	var result *models.Result

	result = CreateToken(1)
	require.NotNil(t, result.Data)

	os.Setenv("ACCES_TOKEN", "\\d")
	result = CreateToken(1)
	require.NotNil(t, result.Data)
}
