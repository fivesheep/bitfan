//go:generate bitfanDoc -codec jsonlinesDecoder
package jsonlinescodec

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"

	"github.com/mitchellh/mapstructure"
)

type jsonlinesDecoder struct {
	more    bool
	r       *bufio.Scanner
	options options
}

type options struct {
	// Change the delimiter that separates lines
	// @Default "\n"
	Delimiter string
}

func New(r io.Reader, opt map[string]interface{}) *jsonlinesDecoder {

	d := &jsonlinesDecoder{
		r:    bufio.NewScanner(r),
		more: true,
		options: options{
			Delimiter: "\n",
		},
	}

	if err := mapstructure.Decode(opt, &d.options); err != nil {
		return nil
	}

	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		// Return nothing if at end of file and no data passed
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		// Find the index of the input of a newline followed by a
		// pound sign.
		if i := strings.Index(string(data), d.options.Delimiter); i >= 0 {
			return i + 1, data[0:i], nil
		}

		// If at end of file with data return the data
		if atEOF {
			return len(data), data, nil
		}

		// Request more data.
		return 0, nil, nil
	}

	d.r.Split(split)

	return d
}

func (c *jsonlinesDecoder) Decode() (map[string]interface{}, error) {
	data := map[string]interface{}{}

	if true == c.r.Scan() {
		c.more = true

		json.Unmarshal([]byte(c.r.Text()), &data)
	} else {
		c.more = false
		return data, io.EOF
	}

	return data, nil
}

func (c *jsonlinesDecoder) More() bool {
	return c.more
}