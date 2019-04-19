package util

/*
 string utils
 @author Tony Tian
 @date 2018-03-17
 @version 1.0.0
*/

import (
	"strconv"
	"strings"
)

/*
  String struct.
	  Usage:
		str1 := NewString("13990521234")
		str2 := NewString("14")
		str3 := str1.Substring(0, 2).Append(str2).AppendString("520")
		println(str1.ToString())
		println(str2.ToString())
		println(str3.ToString())
	  print:
		13990521234
		14
		1314520
*/
type String struct {
	value string
}

func NewString(s string) *String {
	var str String
	str.value = s
	return &str
}

func NewStringInt(i int) *String {
	var str String
	str.value = strconv.Itoa(i)
	return &str
}

func NewStringInt64(i int64) *String {
	var str String
	str.value = strconv.FormatInt(i, 10)
	return &str
}

func NewStringFloat64(f float64) *String {
	var str String
	str.value = strconv.FormatFloat(f, 'E', -1, 64)
	return &str
}

func NewStringAny(f interface{}) *String {
	var str *String
	switch f.(type) {
	case string:
		str = NewString(f.(string))
	case int:
		str = NewStringInt(f.(int))
	case int64:
		str = NewStringInt64(f.(int64))
	case float64:
		str = NewStringFloat64(f.(float64))

	}
	return str
}

func (str *String) ToString() string {
	return str.value
}

func (str *String) Clear() *String {
	var newStr string
	str.value = newStr
	return str
}

/*
 "123" -> 3
*/
func (str *String) Len() int {
	return len(str.value)
}

/*
	"123xxxbbb5990" -> "123x" = true
*/
func (str *String) StartsWith(s string) bool {
	return str.SubstringEnd(len(s)).ToString() == s
}

/*
	"123xxxbbb5990" -> "5990" = true
*/
func (str *String) EndsWith(s string) bool {
	return str.SubstringBegin(str.Len()-len(s)).ToString() == s
}

/*
  " 123 " -> "123"
  " 1 23 " -> "1 23"
*/
func (str *String) Trim() *String {
	return NewString(strings.Trim(str.value, SPACE))
}

/*
  "%111%abc%987%" -> ("%", "$") -> "$111$abc$987$"
*/
func (str *String) Replace(old, new string) *String {
	return NewString(strings.Replace(str.value, old, new, -1))
}

/*
	"abc" -> 1 -> "ac"
*/
func (str *String) Remove(index int) *String {
	strTmp := NewStringBuilder().Append(str.SubstringEnd(index).ToString()).Append(str.SubstringBegin(index + 1).ToString()).ToString()
	return NewString(strTmp)
}

/*
	"abc" -> "ab"
*/
func (str *String) RemoveLast() *String {
	return str.Substring(0, str.Len()-1)
}

/*
  If a string contains a string, return true, and ignore case.
  eg: "strings insert chars"
     chars = "insert" -> true
     chars = "Insert" -> true
     chars = "key" -> false
*/
func (str *String) ContainsIgnoreCase(chars string) bool {
	return str.ToLower().Contains(strings.ToLower(chars))
}

/*
  If a string contains a string, return true
  eg: "strings insert chars"
     chars = "insert" -> true
     chars = "Insert" -> false
     chars = "key" -> false
*/
func (str *String) Contains(chars string) bool {
	return strings.Contains(str.value, chars)
}

/*
  abcdef -> b = 1
*/
func (str *String) LastIndex(chars string) int {
	return strings.LastIndex(str.value, chars)
}

/*
  abcdef -> e = 4
*/
func (str *String) Index(chars string) int {
	return strings.Index(str.value, chars)
}

/*
   "12345" -> 12345
*/
func (str *String) ToInt() (int, error) {
	return strconv.Atoi(str.value)
}

func (str *String) ToInt64() (int64, error) {
	return strconv.ParseInt(str.value, 10, 64)
}

func (str *String) ToFloat() (float64, error) {
	return strconv.ParseFloat(str.value, 64)
}

/*
  str := NewString("abcde")
  str.Substring(0, 2)
  return: "ab"
*/
func (str *String) Substring(beginIndex, endIndex int) *String {
	return NewString(str.value[beginIndex:endIndex])
}

/*
  str := NewString("abcde")
  str.SubstringBegin(2)
  return: "cde"
*/
func (str *String) SubstringBegin(beginIndex int) *String {
	return str.Substring(beginIndex, str.Len())
}

/*
  str := NewString("abcde")
  str.SubstringEnd(3)
  return: "abc"
*/
func (str *String) SubstringEnd(endIndex int) *String {
	return str.Substring(0, endIndex)
}

func (str *String) Append(arg *String) *String {
	strTmp := NewStringBuilder().Append(str.value).Append(arg.ToString()).ToString()
	return NewString(strTmp)
}

func (str *String) AppendString(arg string) *String {
	strTmp := NewStringBuilder().Append(str.value).Append(arg).ToString()
	return NewString(strTmp)
}

func (str *String) AppendInt(i int) *String {
	strTmp := str.value + strconv.Itoa(i)
	return NewString(strTmp)
}

func (str *String) AppendInt64(i int64) *String {
	strTmp := str.value + strconv.FormatInt(i, 10)
	return NewString(strTmp)
}

func (str *String) AppendFloat64(f float64) *String {
	strTmp := str.value + strconv.FormatFloat(f, 'E', -1, 64)
	return NewString(strTmp)
}

/*
  "460364431014955c2483ec91230e5435" -> [4 6 0 3 6 4 4 3 1 0 1 4 9 5 5 c 2 4 8 3 e c 9 1 2 3 0 e 5 4 3 5]
*/
func (str *String) ToArray() []string {
	return strings.Split(str.value, "")
}

/*
  "aaa" -> "AAA"
*/
func (str *String) ToLower() *String {
	return NewString(strings.ToLower(str.value))
}

/*
  "BBB" -> "bbb"
*/
func (str *String) ToUpper() *String {
	return NewString(strings.ToUpper(str.value))
}

/*
  first = false: "aaa_bbb_ccc" -> "aaaBbbCcc"
  first = true: "aaa_bbb_ccc" -> "AaaBbbCcc"
*/
func FirstCaseToUpper(str string, first bool) string {
	temp := strings.Split(str, "_")
	var upperStr string
	for y := 0; y < len(temp); y++ {
		vv := []rune(temp[y])
		if y == 0 && !first {
			continue
		}
		for i := 0; i < len(vv); i++ {
			if i == 0 {
				vv[i] -= 32
				upperStr += string(vv[i])
			} else {
				upperStr += string(vv[i])
			}
		}
	}
	if first {
		return upperStr
	} else {
		return temp[0] + upperStr
	}
}

/*
 [9 9 8 4 2 9 1 7 - a 5 4 b - 3 3 1 6 - c d f 3 - 8 7 d 9 f b 5 7] -> "99842917-a54b-3316-cdf3-87d9fb57"
*/
func ArrayToString(arrays []string) string {
	return strings.Join(arrays, "")
}
