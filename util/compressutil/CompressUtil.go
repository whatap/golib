package compressutil

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
)

func DoZip(in []byte) (output []byte, err error) {
	if in == nil {
		err = fmt.Errorf("error input data is nil ")
		return
	}
	buf := new(bytes.Buffer)

	gz := gzip.NewWriter(buf)
	gz.Write(in)
	gz.Flush()

	// gz.Close 가 호출 되어야만 buf.Bytes 내용이 정상 출력 됨
	err = gz.Close()
	if err == nil {
		output = buf.Bytes()
	}
	return
}

func UnZip(in []byte) ([]byte, error) {
	r, err := gzip.NewReader(ioutil.NopCloser(bytes.NewBuffer(in)))
	if err != nil {
		return make([]byte, 0), err
	}
	defer r.Close()
	if b, err := ioutil.ReadAll(r); err != nil {
		return make([]byte, 0), err
	} else {
		return b, nil
	}
}
