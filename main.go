package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/davecgh/go-spew/spew"
	"golang.org/x/text/language"
	"google.golang.org/api/option"

	"cloud.google.com/go/translate"
)

var fprint = func(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, msg, args...)
}

var exiterror = func() {
	os.Exit(-1)
}

var exitsuccess = func() {
	os.Exit(0)
}

func cleanLang(in string) (language.Tag, error) {
	fprint("cleaning language %s\n", in)
	switch in {
	case "en", "english":
		fprint("using English\n")
		return language.English, nil
	case "fr", "french":
		fprint("using French\n")
		return language.French, nil
	default:
		return language.Tag{}, errors.New("Invalid language" + in)
	}
}

func main() {
	fromtxt := flag.String("from", "", "Text to translate from")
	fromlang := flag.String("fromlang", "french", "Source language [fr|french]")
	tolang := flag.String("tolang", "english", "Language to translate to [en|english]")
	flag.Parse()

	cleanFromLang, err := cleanLang(*fromlang)
	if err != nil {
		fprint(errors.Wrap(err, "Invalid from language").Error())
		exiterror()
	}

	cleanToLang, err := cleanLang(*tolang)
	if err != nil {
		fprint(errors.Wrap(err, "Invalid to language").Error())
		exiterror()
	}

	ctx := context.Background()
	c, err := translate.NewClient(ctx, option.WithAPIKey(os.Getenv("API_KEY")))
	if err != nil {
		fprint(err.Error())
		exiterror()
	}

	r, err := c.Translate(ctx,
		[]string{*fromtxt},
		cleanToLang,
		&translate.Options{
			Source: cleanFromLang,
			Format: translate.Text,
		},
	)

	if err != nil {
		fprint(err.Error())
		exiterror()
	}

	// fmt.Fprintf(os.Stdout, "%+v", r)
	spew.Dump(r)
	exitsuccess()
}
