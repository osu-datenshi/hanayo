// Package doc handles documentation.
package doc

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

func init() {
	// When we start the program, we should load the documentation files.
	err := loadDocFiles()
	if err != nil {
		fmt.Println("error while loading documentation:", err)
	}
}

// rawFile represents the data that may be provided at the top of files.
type rawFile struct {
	Title            string `yaml:"title"`
	OldID            int    `yaml:"old_id"`
	ReferenceVersion string `yaml:"reference_version"`
}

func loadDocFiles() error {
	langs, err := loadLanguagesAvailable()
	if err != nil {
		fmt.Println(err)
		return err
	}

	files, err := ioutil.ReadDir("website-docs/" + referenceLanguage)
	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, file := range files {
		data, err := ioutil.ReadFile("website-docs/" + referenceLanguage + "/" + file.Name())
		if err != nil {
			fmt.Println(err)
			return err
		}

		header := loadHeader(data)
		md5sum := fmt.Sprintf("%x", md5.Sum(data))

		doc := Document{
			OldID: header.OldID,
			Slug:  strings.TrimSuffix(file.Name(), ".md"),
		}

		doc.Languages, err = loadLanguages(langs, file.Name(), md5sum)
		if err != nil {
			fmt.Println(err)
			return err
		}

		docFiles = append(docFiles, doc)
		fmt.Println(docFiles)
	}

	return nil
}

func loadHeader(b []byte) rawFile {
	s := bufio.NewScanner(bytes.NewReader(b))
	var (
		isConf bool
		conf   string
	)

	for s.Scan() {
		line := s.Text()
		if !isConf {
			if line == "---" {
				isConf = true
			}
			continue
		}
		if line == "---" {
			break
		}
		conf += line + "\n"
	}

	var f rawFile
	err := yaml.Unmarshal([]byte(conf), &f)
	if err != nil {
		fmt.Println("Error unmarshaling yaml:", err)
		return rawFile{}
	}

	return f
}

func loadLanguagesAvailable() ([]string, error) {
	files, err := ioutil.ReadDir("website-docs")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	langs := make([]string, 0, len(files))
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		langs = append(langs, f.Name())
	}

	return langs, nil
}

func loadLanguages(langs []string, fname string, referenceMD5 string) (map[string]File, error) {
	m := make(map[string]File, len(langs))

	for _, lang := range langs {
		data, err := ioutil.ReadFile("website-docs/" + lang + "/" + fname)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			fmt.Println(err)
			return nil, err
		}

		header := loadHeader(data)

		m[lang] = File{
			IsUpdated:      lang == referenceLanguage || header.ReferenceVersion == referenceMD5,
			Title:          header.Title,
			referencesFile: "website-docs/" + lang + "/" + fname,
		}
	}

	resp(c, 200, "doc_content.html", &baseTemplateData{
		TitleBar:  header.Title,
		KyutGrill: "documentation.jpg",
	})

	return m, nil
}
