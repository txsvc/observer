package storage

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	bucket = "../hack"
	file   = "test.txt"
)

func cleanup() {
	path := filepath.Join(bucket, file)
	if _, err := os.Stat(path); err == nil {
		os.Remove(path)
	}
}

func TestInitObserver(t *testing.T) {
	assert.NotNil(t, p)

	p, found := p.Find(TypeStorage)
	assert.True(t, found)
	assert.NotNil(t, p)
}

func TestBucketAndObject(t *testing.T) {
	cleanup()

	bkt := Bucket(bucket)
	assert.NotNil(t, bkt)

	obj := bkt.Object(file)
	assert.NotNil(t, obj)

	path := filepath.Join(bucket, file)
	_, err := os.Stat(path)
	assert.Error(t, err)
}

func TestWriter(t *testing.T) {
	cleanup()

	bkt := Bucket(bucket)
	assert.NotNil(t, bkt)

	obj := bkt.Object(file)
	assert.NotNil(t, obj)

	defer obj.Close()

	writer, err := obj.NewWriter(context.TODO())
	assert.NoError(t, err)
	assert.NotNil(t, writer)

	path := filepath.Join(bucket, file)
	_, err = os.Stat(path)
	assert.NoError(t, err)
}

func TestReader(t *testing.T) {

	bkt := Bucket(bucket)
	assert.NotNil(t, bkt)

	obj := bkt.Object(file)
	assert.NotNil(t, obj)

	defer obj.Close()

	reader, err := obj.NewReader(context.TODO())
	assert.NoError(t, err)
	assert.NotNil(t, reader)
}

func TestReaderFail(t *testing.T) {
	cleanup()

	bkt := Bucket(bucket)
	assert.NotNil(t, bkt)

	obj := bkt.Object(file)
	assert.NotNil(t, obj)

	defer obj.Close()

	reader, err := obj.NewReader(context.TODO())
	assert.Error(t, err)
	assert.Nil(t, reader)
}
