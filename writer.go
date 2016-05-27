package quicktemplate

import (
	"fmt"
	"io"
	"strconv"
	"sync"
)

// Writer implements auxiliary writer used by quicktemplate functions.
//
// Use AcquireWriter for creating new writers.
type Writer struct {
	e QWriter
	n QWriter
}

// W returns the underlying writer passed to AcquireWriter.
func (qw *Writer) W() io.Writer {
	return qw.n.w
}

// E returns QWriter with enabled html escaping.
func (qw *Writer) E() *QWriter {
	return &qw.e
}

// N returns QWriter without html escaping.
func (qw *Writer) N() *QWriter {
	return &qw.n
}

// AcquireWriter returns new writer from the pool.
//
// Return unneeded writer to the pool by calling ReleaseWriter
// in order to reduce memory allocations.
func AcquireWriter(w io.Writer) *Writer {
	v := writerPool.Get()
	if v == nil {
		v = &Writer{}
	}
	qw := v.(*Writer)
	qw.e.w = acquireHTMLEscapeWriter(w)
	qw.n.w = w
	return qw
}

// ReleaseWriter returns the writer to the pool.
//
// Do not access released writer, otherwise data races may occur.
func ReleaseWriter(qw *Writer) {
	releaseHTMLEscapeWriter(qw.e.w)
	qw.e.w = nil
	qw.n.w = nil
	writerPool.Put(qw)
}

var writerPool sync.Pool

// QWriter is auxiliary writer used by Writer.
type QWriter struct {
	w io.Writer
}

// S writes s to w.
func (w *QWriter) S(s string) {
	w.w.Write(unsafeStrToBytes(s))
}

// Z writes z to w.
func (w *QWriter) Z(z []byte) {
	w.w.Write(z)
}

// SZ is a synonym to Z.
func (w *QWriter) SZ(z []byte) {
	w.w.Write(z)
}

// D writes n to w.
func (w *QWriter) D(n int) {
	writeQuick(w.w, func(dst []byte) []byte {
		return strconv.AppendInt(dst, int64(n), 10)
	})
}

// F writes f to w.
func (w *QWriter) F(f float64) {
	w.FPrec(f, -1)
}

// FPrec writes f to w using the given floating point precision.
func (w *QWriter) FPrec(f float64, prec int) {
	writeQuick(w.w, func(dst []byte) []byte {
		return strconv.AppendFloat(dst, f, 'f', prec, 64)
	})
}

// Q writes quoted json-safe s to w.
func (w *QWriter) Q(s string) {
	ww := w.w
	ww.Write(strQuote)
	writeJSONString(w.w, s)
	ww.Write(strQuote)
}

var strQuote = []byte(`"`)

// QZ writes quoted json-safe z to w.
func (w *QWriter) QZ(z []byte) {
	w.Q(unsafeBytesToStr(z))
}

// J writes json-safe s to w.
//
// Unlike Q it doesn't qoute resulting s.
func (w *QWriter) J(s string) {
	writeJSONString(w.w, s)
}

// JZ writes json-safe z to w.
//
// Unlike Q it doesn't qoute resulting z.
func (w *QWriter) JZ(z []byte) {
	w.J(unsafeBytesToStr(z))
}

// V writes v to w.
func (w *QWriter) V(v interface{}) {
	fmt.Fprintf(w.w, "%v", v)
}

// U writes url-encoded s to w.
func (w *QWriter) U(s string) {
	writeQuick(w.w, func(dst []byte) []byte {
		return appendURLEncode(dst, s)
	})
}

// UZ writes url-encoded z to w.
func (w *QWriter) UZ(z []byte) {
	w.U(unsafeBytesToStr(z))
}

func writeQuick(w io.Writer, f func(dst []byte) []byte) {
	bb, ok := w.(*ByteBuffer)
	if !ok {
		bb = AcquireByteBuffer()
	}
	bb.B = f(bb.B)
	if !ok {
		w.Write(bb.B)
		ReleaseByteBuffer(bb)
	}
}
