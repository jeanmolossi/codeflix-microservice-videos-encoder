package utils_test

import (
	"github.com/jeanmolossi/codeflix-microservice-videos-encoder/framework/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsJson(t *testing.T) {
	json := `{
		"id": "41aabea8-ea76-459b-9726-6b38fe6a41c8",
		"file_path": "ForBiggerFun.mp4",
		"status": "pending"
	}`

	err := utils.IsJson(json)
	require.Nil(t, err)

	json = `john`
	err = utils.IsJson(json)
	require.Error(t, err)
}
