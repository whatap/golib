package zip

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

type DefaultZipMod struct {
	id byte
}

func NewDefaultZipMod() *DefaultZipMod {
	p := new(DefaultZipMod)
	return p
}
func (this *DefaultZipMod) ID() byte {
	return ZIP_MOD_DEFULAT_GZIP
}

func (this *DefaultZipMod) Compress(in []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	func() {
		gz := gzip.NewWriter(buf)
		defer gz.Close()

		_, err := gz.Write(in)
		if err != nil {
			panic(err)
		}
	}()

	return buf.Bytes(), nil
}

func (this *DefaultZipMod) Decompress(in []byte) ([]byte, error) {
	r, err := gzip.NewReader(ioutil.NopCloser(bytes.NewBuffer(in)))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	if b, err := ioutil.ReadAll(r); err != nil {
		return nil, err
	} else {
		return b, nil
	}
}
