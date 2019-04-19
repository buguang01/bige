package util

import (
	"bytes"
)

/*
  StringBuilder struct.
	  Usage:
		builder := NewStringBuilder()
		builder.Append("a").Append("b").Append("c")
		s := builder.String()
		println(s)
	  print:
		abc
*/
type StringBuilder struct {
	buffer bytes.Buffer
}

func NewStringBuilder() *StringBuilder {
	var builder StringBuilder
	return &builder
}

func NewStringBuilderString(str *String) *StringBuilder {
	var builder StringBuilder
	builder.buffer.WriteString(str.ToString())
	return &builder
}

func (builder *StringBuilder) Append(s string) *StringBuilder {
	builder.buffer.WriteString(s)
	return builder
}

func (builder *StringBuilder) AppendStrings(ss ...string) *StringBuilder {
	for i := range ss {
		builder.buffer.WriteString(ss[i])
	}
	return builder
}

func (builder *StringBuilder) AppendInt(i int) *StringBuilder {
	builder.buffer.WriteString(NewStringInt(i).ToString())
	return builder
}

func (builder *StringBuilder) AppendInt64(i int64) *StringBuilder {
	builder.buffer.WriteString(NewStringInt64(i).ToString())
	return builder
}

func (builder *StringBuilder) AppendFloat64(f float64) *StringBuilder {
	builder.buffer.WriteString(NewStringFloat64(f).ToString())
	return builder
}

func (builder *StringBuilder) Replace(old, new string) *StringBuilder {
	str := NewString(builder.ToString()).Replace(old, new)
	builder.Clear()
	builder.buffer.WriteString(str.ToString())
	return builder
}

func (builder *StringBuilder) RemoveLast() *StringBuilder {
	str1 := NewString(builder.ToString())
	builder.Clear()
	str2 := str1.Substring(0, str1.Len()-1)
	builder.buffer.WriteString(str2.ToString())
	return builder
}

func (builder *StringBuilder) Clear() *StringBuilder {
	var buffer bytes.Buffer
	builder.buffer = buffer
	return builder
}

func (builder *StringBuilder) ToString() string {
	return builder.buffer.String()
}
