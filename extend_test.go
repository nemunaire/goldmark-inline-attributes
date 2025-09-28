package attributes

import (
	"log"
	"os"
	"testing"

	"github.com/yuin/goldmark"
)

func TestAttributes(t *testing.T) {
	source := []byte(`
[This is *some text*]{.class key="val"} outside text
`)

	var md = goldmark.New(Enable)
	err := md.Convert(source, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}

}

func TestWithNonAttributes(t *testing.T) {
	source := []byte(`
[This is *some text*] outside text [This is *some text*]{.class key="val"} outside text
[This is *some text* outside text [This is *some text*]{.class key="val"} outside text
`)

	var md = goldmark.New(Enable)
	err := md.Convert(source, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}

}
