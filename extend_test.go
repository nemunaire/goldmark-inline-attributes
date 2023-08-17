package attributes

import (
	"log"
	"os"
	"testing"

	"github.com/yuin/goldmark"
)

func TestAttributes(t *testing.T) {
	source := []byte(`
[Text underlined]{.underline}
`)

	var md = goldmark.New(Enable)
	err := md.Convert(source, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}

}
