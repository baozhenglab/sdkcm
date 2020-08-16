package sdkcm

import (
	"fmt"
	"regexp"
)

type (
	Base64Decode struct {
		TypeBase64 string
		Data       string
	}
	OptionsIsBase64 struct {
		MimeRequired    bool
		AllowMime       bool
		PaddingRequired bool
	}
	Base64 string
)

//IsBase64 check string is base64
func (base64 Base64) IsBase64(options OptionsIsBase64) bool {
	regex := `(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\/]{3}=)?`
	mimeRegex := `(data:\\w+\\/[a-zA-Z\\+\\-\\.]+;base64,)`
	if options.MimeRequired {
		regex = mimeRegex + regex
	} else if options.AllowMime {
		regex = mimeRegex + "?" + regex
	}
	if !options.PaddingRequired {
		regex = `(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}(==)?|[A-Za-z0-9+\\/]{3}=?)?`
	}
	r, _ := regexp.Compile(regex)
	return r.MatchString(fmt.Sprintf("%s", base64))
}

//DecodeBase64PDF return string base64 and type
func (base64 Base64) DecodeBase64PDF() Base64Decode {
	r, _ := regexp.Compile(`^data:application\/([\w+]+);base64,([\s\S]+)`)
	matches := r.FindAllStringSubmatch(fmt.Sprintf("%s", base64), -1)
	if len(matches[0]) != 3 {
		return Base64Decode{}
	}
	return Base64Decode{
		TypeBase64: matches[0][1],
		Data:       matches[0][2],
	}
}

//DecodeBase64Image return type and buffer base64
func (base64 Base64) DecodeBase64Image() Base64Decode {
	r, _ := regexp.Compile(`^data:image\/([\w+]+);base64,([\s\S]+)`)
	matches := r.FindAllStringSubmatch(fmt.Sprintf("%s", base64), -1)
	if len(matches[0]) != 3 {
		return Base64Decode{}
	}
	return Base64Decode{
		TypeBase64: matches[0][1],
		Data:       matches[0][2],
	}
}
