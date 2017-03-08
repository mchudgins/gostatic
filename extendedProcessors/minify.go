package extendedProcessors

import (
	"compress/gzip"
	"fmt"
	"path/filepath"

	gostatic "github.com/mchudgins/gostatic/lib"
	mini "github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
	"github.com/tdewolff/minify/svg"

	"strings"
	"os"
)

type MinifyProcessor struct {
}

func NewMinifyProcessor() *MinifyProcessor {
	fmt.Printf("MinifyProcessor constructed.\n")
	return &MinifyProcessor{}
}

func (p *MinifyProcessor) Process(page *gostatic.Page, args []string) error {
	fmt.Printf("Page:\n")
	fmt.Printf("\tOutputPath:  %s\n", page.OutputPath())

	// find the file's extension
	fmt.Printf("Base: %s\n", filepath.Base(page.OutputPath()))
	fmt.Printf("Ext : %s\n", filepath.Ext(page.OutputPath()))
	fmt.Printf("    : %s\n", strings.TrimSuffix(page.OutputPath(),filepath.Ext(page.OutputPath())))

	path := page.OutputPath()
	ext := filepath.Ext(path)
	raw := strings.TrimSuffix(path,ext)

	// change to -min.<extension>.gz

	minified := raw + "-min" + ext + ".gz"
	fmt.Printf("Minified Target:  %s\n", minified)

	f, err := os.OpenFile(minified,os.O_WRONLY | os.O_CREATE | os.O_TRUNC,0660)
	if err != nil {
		return err
	}
	defer f.Close()

	m := mini.New()
	m.AddFunc("text/html", html.Minify)
	m.AddFunc( "text/css", css.Minify)
	m.AddFunc( "text/javascript", js.Minify)
	m.AddFunc( "image/svg+xml", svg.Minify)

	s, err := m.String("text/html", page.Content())
	if err != nil {
		return err
	}
	fmt.Printf("minified: %s\n", s)

	// gzip it

	gzipper, err := gzip.NewWriterLevel(f,gzip.BestCompression)
	if err != nil {
		return err
	}
	defer gzipper.Close()

	// write it out along side the existing target

	_, err = gzipper.Write([]byte(s))

	return err
}

func (p *MinifyProcessor) Description() string {
	return "minify content"
}

func (p *MinifyProcessor) Mode() int {
	return 0
}


