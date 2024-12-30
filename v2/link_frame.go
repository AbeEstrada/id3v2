// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import "io"

type LinkFrame struct {
	Encoding Encoding
	Url      string
}

func (tf LinkFrame) Size() int {
	return 1 + encodedSize(tf.Url, tf.Encoding) + len(tf.Encoding.TerminationBytes)
}

func (tf LinkFrame) UniqueIdentifier() string {
	return "ID"
}

func (tf LinkFrame) WriteTo(w io.Writer) (int64, error) {
	return useBufWriter(w, func(bw *bufWriter) {
		bw.WriteByte(tf.Encoding.Key)
		bw.EncodeAndWriteText(tf.Url, tf.Encoding)

		// https://github.com/bogem/id3v2/pull/52
		// https://github.com/bogem/id3v2/pull/33
		bw.Write(tf.Encoding.TerminationBytes)
	})
}

func parseLinkFrame(br *bufReader) (Framer, error) {
	encoding := getEncoding(br.ReadByte())

	if br.Err() != nil {
		return nil, br.Err()
	}

	buf := getBytesBuffer()
	defer putBytesBuffer(buf)
	if _, err := buf.ReadFrom(br); err != nil {
		return nil, err
	}

	lf := LinkFrame{
		Encoding: encoding,
		Url:      decodeText(buf.Bytes(), encoding),
	}

	return lf, nil
}
