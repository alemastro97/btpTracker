package request

import (
	"encoding/json"
	// "fmt"
	"io"
	"net/http"
	// "time"
)

// Semantic scholar response

type Response struct {
	Data []Paper `json:"data"`
}

type Paper struct {
	PaperId    string      `json:"paperId" bson:"PaperId"`
	Title      string      `json:"title" bson:"Title"`
	Abstract   string      `json:"abstract" bson:"Abstract"`
	ExternalId ExternalId  `json:"externalIds,omitempty" bson:"ExternalId"`
	Pdf        Pdf         `json:"openAccessPdf,omitempty" bson:"Pdf"`
	Authors    []Author    `json:"authors" bson:"Authors"`
	References []Reference `json:"references" bson:"References"`
	Ancestors  []string    `json:"ancestors" bson:"Ancestors"`
}

type ExternalId struct {
	DOI   string `json:",omitempty" bson:"DOI"`
	ArXiv string `json:",omitempty" bson:"ArXiv"`
	MAG   string `json:",omitempty" bson:"MAG"`
	ACL   string `json:",omitempty" bson:"ACL"`
}

type Pdf struct {
	Url string `json:"url" bson:"Url"`
}

type Author struct {
	AuthorId string `json:"authorId" bson:"AuthorId"`
	Name     string `json:"name" bson:"Name"`
}

type Reference struct {
	PaperId string `json:"paperId" bson:"PaperId"`
	Title   string `json:"title" bson:"Title"`
}

// Pretty print any json-serializable interface.
func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

// Perform a request
func MakeRequest(url string, dest any) error {
	// start := time.Now()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Read as byte
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Parse the response
	err = json.Unmarshal(body, dest)
	if err != nil {
		return err
	}

	// elapsed := time.Since(start)
	// fmt.Printf("Elapsed time: %s\n", elapsed)

	return nil
}
