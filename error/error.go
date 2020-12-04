package error

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

type ApplicationError interface {
	Dump() string
}

type AppError struct {
	Original error
	ErrCode  int
	Remark   string
	File     string
	Line     int
	Context  map[string]string
}

//Error method will instantiate new struct of AppError and populates Original (error) and ErrCode (int) values
//The resulting struct can then be used to enrich information, e.g. adding remarks and contexts (key-value)
func Error(err error, errCode int) *AppError {
	e := &AppError{}
	e.Original = err
	e.ErrCode = errCode
	if _, file, line, ok := runtime.Caller(1); ok {
		f := strings.Split(file, "/")
		e.File = f[len(f)-1]
		e.Line = line
	}
	return e
}

//Errorc method will instantiate new struct of AppError and populates ErrCode (int) value
//The resulting struct can then be used to enrich information, e.g. adding remarks and contexts (key-value)
func Errorc(errCode int) *AppError {
	e := &AppError{}
	e.ErrCode = errCode
	if _, file, line, ok := runtime.Caller(1); ok {
		f := strings.Split(file, "/")
		e.File = f[len(f)-1]
		e.Line = line
	}
	return e
}

func (e *AppError) Rem(msg string, a ...interface{}) *AppError {
	e.Remark = fmt.Sprintf(msg, a...)
	return e
}

func (e *AppError) SetString(key string, val string) *AppError {
	e.Context[key] = val
	return e
}

func (e *AppError) Dump() string {
	var buff bytes.Buffer

	if e.File != "" {
		buff.WriteString(e.File)
		buff.WriteString(":")
		buff.WriteString(strconv.Itoa(e.Line))
		buff.WriteString("/")
	}

	if e.ErrCode != -255 {
		buff.WriteString("Err-")
		buff.WriteString(strconv.Itoa(e.ErrCode))
		buff.WriteString("/")
	}

	if e.Remark != "" {
		buff.WriteString(e.Remark)
		buff.WriteString("/")
	}

	if e.Original != nil && e.Original.Error() != e.Remark {
		buff.WriteString(e.Original.Error())
		buff.WriteString("/")
	}

	if e.Context != nil {
		for key, val := range e.Context {
			buff.WriteString(fmt.Sprintf("%s=%s/", key, val))
		}
	}
	return buff.String()
}
