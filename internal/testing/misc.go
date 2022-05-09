package testing

import "io"

type EmptyReader struct{}

var _ io.Reader = (*EmptyReader)(nil)

func (*EmptyReader) Read([]byte) (int, error) { return 0, io.EOF }
