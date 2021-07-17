package socket

import (
	"bytes"
	"encoding/binary"
	"github.com/pkg/errors"
	"net"
)

const (
	_packLenSize = 4
)

type Worker struct {
	ip    string
	index int64

	conn *net.TCPConn
	step uint64

	notifyStop     chan error
	stopSign       chan struct{}
	stopped        chan struct{}
	replyCompleted chan struct{}

	reader Reader

	replyChan chan []byte
}


func (w *Worker) read() (buf []byte, err error) {
	lb, err := w.reader.ReadFull(_packLenSize)
	if err != nil {
		err = errors.Wrap(err, "tcp: read packSize error.")
		return
	}

	var packLen int32
	err = binary.Read(bytes.NewReader(lb), binary.LittleEndian, &packLen)
	if err != nil {
		return
	}

	buf, err = w.reader.ReadFull(int(packLen))
	if err != nil {
		err = errors.Wrap(err, "bufreader read pack failed")
	}
	return
}

func (w *Worker) write(pack []byte) (err error) {
	var buf bytes.Buffer

	err = binary.Write(&buf, binary.LittleEndian, uint32(len(pack)))
	if err != nil {
		err = errors.Wrap(err, "write pack len to buffer failed")
		return
	}

	if _, err = buf.Write(pack); err != nil {
		err = errors.Wrap(err, "write pack to buffer failed")
		return
	}

	if _, err = w.conn.Write(buf.Bytes()); err != nil {
		err = errors.Wrap(err, "write pack to conn failed")
		return
	}
	return
}

type Reader interface {
	ReadFull(size int) ([]byte, error)
}
