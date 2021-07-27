package s3

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckFile(t *testing.T) {
	url, err := New(&Config{
		AccessKey: os.Getenv("AWS_ACCESS_KEY"),
		SecretKey: os.Getenv("AWS_SECRET_KEY"),
		Region:    os.Getenv("AWS_REGION"),
	}).CheckFile("ezielog-staging", "labels/cj_label_OR-BSDJ-913632")

	assert.NoError(t, err)
	fmt.Println(url)
}
