package middlewares

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

var errRequestTooLarge = errors.New("HTTP request too large")

type maxBytesReader struct {
	ctx        *gin.Context
	rdr        io.ReadCloser
	remaining  int64
	wasAborted bool
	sawEOF     bool
}

func (mbr *maxBytesReader) tooLarge() (int, error) {
	if !mbr.wasAborted {
		mbr.wasAborted = true
		_ = mbr.ctx.Error(errRequestTooLarge)
		mbr.ctx.Header("Connection", "close")
		mbr.ctx.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
			"error": "request too large",
		})
	}
	return 0, errRequestTooLarge
}

func (mbr *maxBytesReader) Read(p []byte) (int, error) {
	if mbr.wasAborted {
		return 0, errRequestTooLarge
	}
	toRead := mbr.remaining
	if mbr.remaining == 0 {
		if mbr.sawEOF {
			return mbr.tooLarge()
		}
		toRead = 1
	}
	if int64(len(p)) > toRead {
		p = p[:toRead]
	}
	n, err := mbr.rdr.Read(p)
	if err == io.EOF {
		mbr.sawEOF = true
	}
	if mbr.remaining == 0 {
		if n > 0 {
			return mbr.tooLarge()
		}
		return 0, err
	}
	mbr.remaining -= int64(n)
	if mbr.remaining < 0 {
		mbr.remaining = 0
	}
	return n, err
}

func (mbr *maxBytesReader) Close() error {
	return mbr.rdr.Close()
}

func RequestSizeLimiter(limit int64) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Request.Body = &maxBytesReader{
			ctx:       ctx,
			rdr:       ctx.Request.Body,
			remaining: limit,
		}
		ctx.Next()
	}
}
