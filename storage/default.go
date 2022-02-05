package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/txsvc/observer/pkg/provider"
)

type (
	// defaultStorageImpl provides a simple implementation in the absence of any configuration
	defaultStorageImpl struct {
	}

	bucketImpl struct {
		name string
	}

	objImpl struct {
		bucket string
		name   string
		file   *os.File
	}
)

var (
	// Interface guards.

	// This enforces a compile-time check of the provider implmentation,
	// making sure all the methods defined in the interfaces are implemented.

	_ provider.GenericProvider = (*defaultStorageImpl)(nil)

	_ StorageProvider = (*defaultStorageImpl)(nil)

	// the instance, a singleton
	theDefaultProvider *defaultStorageImpl
)

func init() {
	Init()
}

func Init() {
	// force a reset
	theDefaultProvider = nil

	// initialize the observer with a NULL provider that prevents NPEs in case someone forgets to initialize the platform with a real provider
	storageConfig := provider.WithProvider("storage.default.storage", TypeStorage, NewDefaultProvider)

	NewConfig(storageConfig)
}

// a default provider that does nothing but prevents NPEs in case someone forgets to actually initializa the 'real' provider
func NewDefaultProvider() interface{} {
	if theDefaultProvider == nil {
		theDefaultProvider = &defaultStorageImpl{}
	}
	return theDefaultProvider
}

func (np *defaultStorageImpl) Close() error {
	return nil
}

func (np *defaultStorageImpl) Bucket(name string) BucketHandle {
	bkt := &bucketImpl{
		name: name,
	}
	return bkt
}

func (bkt *bucketImpl) Object(name string) ObjectHandle {
	obj := &objImpl{
		bucket: bkt.name,
		name:   name,
	}
	return obj
}

func (obj *objImpl) Close() error {
	if obj.file != nil {
		return obj.file.Close()
	}
	return nil
}

func (obj *objImpl) NewReader(context.Context) (io.Reader, error) {
	path := filepath.Join(obj.bucket, obj.name)
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}

	if obj.file != nil {
		obj.file.Close()
		obj.file = nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	obj.file = f
	return f, nil
}

func (obj *objImpl) NewWriter(context.Context) (io.Writer, error) {
	path := filepath.Join(obj.bucket, obj.name)
	if err := os.MkdirAll(filepath.Dir(path), 0770); err != nil {
		return nil, err
	}

	if obj.file != nil {
		obj.file.Close()
		obj.file = nil
	}

	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	obj.file = f
	return f, nil
}
