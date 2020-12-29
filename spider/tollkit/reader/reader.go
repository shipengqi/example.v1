package reader

import (
	"bytes"
	"io"
	"io/ioutil"
)

// MultipleReader 代表多重读取器的接口
type MultipleReader interface {
	// Reader 用于获取一个可关闭读取器的实例
	// 后者会持有本多重读取器中的数据
	Reader() io.ReadCloser
}

type MultipleReaderImpl struct {
	data []byte
}

func NewMultipleReaderImpl(reader io.Reader) (*MultipleReaderImpl, error) {
	var data []byte
	var err error

	if reader != nil {
		data, err = ioutil.ReadAll(reader)
		if err != nil {
			return nil, err
		}
	} else {
		data = []byte{}
	}
	return &MultipleReaderImpl{data: data}, nil
}

func (m *MultipleReaderImpl) Reader() io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader(m.data))
}
