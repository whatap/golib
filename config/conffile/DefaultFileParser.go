package conffile

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/magiconair/properties"
)

type DefaultFileParser struct {
}

func NewDefaultFileParser() *DefaultFileParser {
	p := new(DefaultFileParser)
	return p
}

func (this *DefaultFileParser) Read(filePath string) (map[string]string, error) {
	prop := properties.MustLoadFile(filePath, properties.UTF8)
	m := make(map[string]string)
	for _, k := range prop.Keys() {
		v := prop.GetString(k, "")
		if v != "" {
			m[k] = v
		}
	}
	return m, nil
}

func (this *DefaultFileParser) Write(filePath string, m *map[string]string) error {
	props := properties.MustLoadFile(filePath, properties.UTF8)

	for key, value := range *m {
		props.Set(key, value)
	}

	line := ""
	if f, err := os.OpenFile(filePath, os.O_RDWR, 0644); err != nil {
		return err
	} else {
		defer f.Close()

		r := bufio.NewReader(f)
		new_keys := props.Keys()
		old_keys := map[string]bool{}
		for {
			data, _, err := r.ReadLine()
			if err != nil { // new key
				for _, key := range new_keys {
					if old_keys[key] {
						continue
					}
					match, _ := regexp.MatchString("^\\w", key)
					if match {
						value, _ := props.Get(key)
						if strings.TrimSpace(value) != "" {
							tmp := strings.Replace(value, "\\\\", "\\", -1)
							tmp = strings.Replace(tmp, "\\", "\\\\", -1)
							line += fmt.Sprintf("%s=%s\n", key, tmp)
						}
					}
				}
				break
			}
			if strings.Index(string(data), "=") == -1 {
				line += fmt.Sprintf("%s\n", string(data))
				//io.WriteString(f, line)
			} else {
				datas := strings.Split(string(data), "=")
				key := strings.Trim(datas[0], " ")
				value := strings.Trim(datas[1], " ")
				old_keys[key] = true

				match, _ := regexp.MatchString("^\\w", key)
				if match {
					value, _ = props.Get(key)
				}
				// value 가 없는 경우 항목 추가 안함(삭제)
				if strings.TrimSpace(value) != "" {
					tmp := strings.Replace(value, "\\\\", "\\", -1)
					tmp = strings.Replace(tmp, "\\", "\\\\", -1)

					line += fmt.Sprintf("%s=%s\n", key, tmp)
				}
				//io.WriteString(f, line)
			}
		}
	}

	if f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644); err != nil {
		return err
	} else {
		defer f.Close()
		io.WriteString(f, line)

		// flush
		f.Sync()
	}
	return nil
}
