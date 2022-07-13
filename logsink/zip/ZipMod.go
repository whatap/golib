package zip

import ()

const (
	ZIP_MOD_DEFULAT_GZIP = 1
)

type ZipMod interface {
	ID() byte
	Compress(b []byte) ([]byte, error)
}
