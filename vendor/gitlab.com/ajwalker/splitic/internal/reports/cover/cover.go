package cover

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
)

func Merge(inputs []string, w io.Writer) error {
	var coverMode string

	for _, input := range inputs {
		buf, err := ioutil.ReadFile(input)
		if err != nil {
			return err
		}

		idx := bytes.IndexByte(buf, '\n')
		hdr := bytes.TrimPrefix(buf[:idx], []byte("mode:"))
		hdr = bytes.TrimSpace(hdr)

		if coverMode == "" {
			coverMode = string(hdr)
		} else {
			if string(hdr) != coverMode {
				return fmt.Errorf("cover profile uses different coverMode")
			}
			buf = buf[idx+1:]
		}

		if _, err = w.Write(buf); err != nil {
			return err
		}
	}

	return nil
}
