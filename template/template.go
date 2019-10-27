package template

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"text/template"
)

// TemplateData is the data available to a ContentTypeSpec.Example template string
type TemplateData struct {
	Params    map[string]string
	ItemCount int
	Index     int
}

func (e TemplateData) JSONArray(count int, s string) string {
	out := "["
	for i := 1; i <= count; i++ {
		itd := e
		itd.Index = i
		ib, err := ExampleOutput(s, itd)
		if err != nil {
			fmt.Println(err)
			continue
		}

		out += string(ib)

		if i != count {
			out += ","
		}
	}
	out += "]"
	return out
}

func (e TemplateData) ToLower(in string) string {
	return strings.ToLower(in)
}

func (e TemplateData) RandomFirstName() string {
	names := []string{
		"Aaron",
		"Barbara",
		"Charles",
		"Diane",
		"Edward",
		"Felicity",
		"Greg",
		"Harriet",
		"Idris",
		"Jacqueline",
		"Ken",
		"Lisa",
		"Mark",
		"Nina",
		"Orlando",
		"Pierre",
		"Quaid",
		"Ryan",
		"Stacy",
		"Timothy",
		"Ursula",
		"Victor",
		"Wanda",
		"Xavier",
		"Yolanda",
		"Zachary",
	}

	return names[rand.Intn(len(names))]
}

func (e TemplateData) RandomFrom(f ...string) string {
	return f[rand.Intn(len(f))]
}

func (e TemplateData) RandomPassword(n int) string {
	alphaString := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-="
	alphaSlice := strings.Split(alphaString, "")

	out := ""
	for i := 0; i < n; i++ {
		out += alphaSlice[rand.Intn(len(alphaSlice))]
	}

	return out
}

func ExampleOutput(input string, td TemplateData) ([]byte, error) {
	tpl, err := template.New("").Parse(input)
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer

	err = tpl.Execute(&out, td)
	if err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
