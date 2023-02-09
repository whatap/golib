package urlutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURL(t *testing.T) {
	p := NewURL("http://www.naver.com/")
	assert.Equal(t, "http", p.Protocol)
	assert.Equal(t, "www.naver.com", p.Host)
	assert.Equal(t, 80, p.Port)
	assert.Equal(t, "/", p.Path)
	assert.Equal(t, "", p.Query)

	p = NewURL("http://www.naver.com/?http://sss.com/aaa")
	assert.Equal(t, "http", p.Protocol)
	assert.Equal(t, "www.naver.com", p.Host)
	assert.Equal(t, 80, p.Port)
	assert.Equal(t, "/", p.Path)
	assert.Equal(t, "http://sss.com/aaa", p.Query)

	p = NewURL("hammaa.test.com:8080/roundTripper/inputCallUrl?url=http://c7default.test.com/test/curl/curl.php")
	assert.Equal(t, "", p.Protocol)
	assert.Equal(t, "hammaa.test.com", p.Host)
	assert.Equal(t, 8080, p.Port)
	assert.Equal(t, "/roundTripper/inputCallUrl", p.Path)
	assert.Equal(t, "url=http://c7default.test.com/test/curl/curl.php", p.Query)

	p = NewURL("/roundTripper/inputCallUrl?url=http://c7default.test.com/test/curl/curl.php")
	assert.Equal(t, "", p.Protocol)
	assert.Equal(t, "", p.Host)
	assert.Equal(t, 80, p.Port)
	assert.Equal(t, "/roundTripper/inputCallUrl", p.Path)
	assert.Equal(t, "url=http://c7default.test.com/test/curl/curl.php", p.Query)

	p = NewURL("https://www.naver.com/a/b/c/d/index.php?aal=3&bbb=3&url=http://c7default.test.com/test/curl/curl.php")
	assert.Equal(t, "https", p.Protocol)
	assert.Equal(t, "www.naver.com", p.Host)
	assert.Equal(t, 443, p.Port)
	assert.Equal(t, "/a/b/c/d/index.php", p.Path)
	assert.Equal(t, "aal=3&bbb=3&url=http://c7default.test.com/test/curl/curl.php", p.Query)

	p = NewURL("http://www.naver.com/a/b/c/d/index.php?aal=3&bbb=3")
	//fmt.Println("url=", p.Url , "\r\n", p.Protocol, ", ", p.Host, ", ", p.Port , ", ", p.Path , ", ", p.File, ", ", p.Query)

	p = NewURL("https://www.naver.com:80")
	//fmt.Println("url=", p.Url , "\r\n", p.Protocol, ", ", p.Host, ", ", p.Port , ", ", p.Path , ", ", p.File, ", ", p.Query)

	p = NewURL("https://www.naver.com:80/")
	//fmt.Println("url=", p.Url , "\r\n", p.Protocol, ", ", p.Host, ", ", p.Port , ", ", p.Path , ", ", p.File, ", ", p.Query)

	p = NewURL("https://www.naver.com:80/a/b/c/")
	//fmt.Println("url=", p.Url , "\r\n", p.Protocol, ", ", p.Host, ", ", p.Port , ", ", p.Path , ", ", p.File, ", ", p.Query)

	p = NewURL("http://www.naver.com:80/a/b/c/index.php")
	//fmt.Println("url=", p.Url , "\r\n", p.Protocol, ", ", p.Host, ", ", p.Port , ", ", p.Path , ", ", p.File, ", ", p.Query)

	p = NewURL("http://www.naver.com:80/a/b/c/d/index.php?aal=3&bbb=3")
	//fmt.Println("url=", p.Url , "\r\n", p.Protocol, ", ", p.Host, ", ", p.Port , ", ", p.Path , ", ", p.File, ", ", p.Query)
}
