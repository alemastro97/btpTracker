package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"citation-graph/backend/database"
	"citation-graph/backend/request"
)

const port string = ":8080"

func assert(cond bool) {
	if !cond {
		panic("Assertion failed")
	}
}

// Number of CPUs to use to parallelize the work.
var numCPU = runtime.NumCPU() * 2

// The Semantic Scholar API to use to query for a specific DOI.
const AllInfoRoute = "https://api.semanticscholar.org/graph/v1/paper/%s?fields=title,authors,references,abstract,openAccessPdf,externalIds"
const ReferencesRoute = "https://api.semanticscholar.org/graph/v1/paper/%s?fields=references"
const QuerySearchRoute = "https://api.semanticscholar.org/graph/v1/paper/search?query=%s&fields=title,authors,references,abstract,openAccessPdf,externalIds&limit=%s"

func Extract_citation(cites []request.Reference, ancestors []string, subtree *ConcurrentMap, c chan []*request.Paper) {
	var papers_list []*request.Paper = make([]*request.Paper, len(cites))

	for i, cite := range cites {
		if cite.PaperId == "" {
			// Skip invalid paperId, cannot request them.
			continue
		}
		var got request.Paper
		subtree.RLock()
		paper, exists := subtree.m[cite.PaperId]
		subtree.RUnlock()
		if exists {
			got = paper
		} else {
			paperInDb, err := database.IsPresent("papers", "PaperId", cite.PaperId)
			// If the paper is not in the db, request it to the Semantic Scholar API.
			if errors.Is(err, database.ErrNoDocuments) {
				err := request.MakeRequest(fmt.Sprintf(AllInfoRoute, cite.PaperId), &got)
				if err != nil {
					log.Printf("Error while retrieving PaperID from Scholar. Error: %s", err)
					continue
				}
				// Make sure to update the ancestor array
				got.Ancestors = ancestors
				if got.PaperId == "" {
					got.PaperId = cite.PaperId
					log.Printf("Catch empty PaperId, replace with %s\n", got.PaperId)
				}
				err = database.Insert_element("papers", got)
				if err != nil {
					log.Printf("Error while inserting paper in db: ERROR %s", err)
					continue
				}
				// Save the new value in cache.
				subtree.Lock()
				subtree.m[got.PaperId] = got
				subtree.Unlock()

			} else {
				got = *paperInDb
			}
		}

		// Add paper to citation.
		papers_list[i] = &got
	}

	c <- papers_list // Send response to channel.
}

// Tree data structure to return to frontend.
type Node struct {
	PaperId    string             `json:"paperId" bson:"PaperId"`
	Title      string             `json:"title" bson:"Title"`
	Abstract   string             `json:"abstract" bson:"Abstract"`
	ExternalId request.ExternalId `json:"externalIds,omitempty" bson:"ExternalId"`
	Pdf        request.Pdf        `json:"openAccessPdf,omitempty" bson:"Pdf"`
	Authors    []request.Author   `json:"authors" bson:"Authors"`
	References []*Node            `json:"references" bson:"References"`
}

// Create a new Node starting from the input Paper.
func FromPaper(paper *request.Paper) Node {
	return Node{
		PaperId:    paper.PaperId,
		Title:      paper.Title,
		Abstract:   paper.Abstract,
		ExternalId: paper.ExternalId,
		Pdf:        paper.Pdf,
		Authors:    paper.Authors,
		References: []*Node{},
	}
}

type ConcurrentMap struct {
	sync.RWMutex
	m map[string]request.Paper
}

// Recursively retrieve all the nodes starting from `root` paper up to
// `maxDepth` height.  Also, check if the element is already in DB using the
// lookup table `subtree`.
func insert_citations(root *request.Paper, maxDepth int, subtree map[string]request.Paper) Node {
	node := FromPaper(root)

	var subtree_map = ConcurrentMap{m: subtree}

	// Helper function to recur over all the citations.
	var innerRecur func([]request.Reference, []string, int) []*Node

	innerRecur = func(cites []request.Reference, ancestor []string, deepLimit int) []*Node {
		// Max depth reached. Abort.
		if deepLimit > maxDepth {
			return []*Node{}
		}

		c := make(chan []*request.Paper, numCPU) // Buffering optional but sensible.

		for i := 0; i < numCPU; i++ {
			// Divide into chunk.
			go Extract_citation(cites[i*len(cites)/numCPU:(i+1)*len(cites)/numCPU], ancestor, &subtree_map, c)
		}

		var papers_list []*request.Paper
		// Retrive all the citations.
		for i := 0; i < numCPU; i++ {
			tmp := <-c // wait for one task to complete
			papers_list = append(papers_list, tmp...)
		}

		// Construct the tree node and recur to find all the subtree.
		children := []*Node{}
		for i := 0; i < len(papers_list); i++ {
			if papers_list[i] == nil {
				continue
			}
			child := FromPaper(papers_list[i])
			child.References = innerRecur(papers_list[i].References, append(ancestor, child.PaperId), deepLimit+1)
			children = append(children, &child)
		}

		return children
	}

	node.References = innerRecur(root.References, []string{root.PaperId}, 1)
	return node
}

// Generic request paper function.
//
// This route expects the following query parameters:
//   - 'id': the identifier of the paper. This can be a DOI, ArXiv link or a a query term. How the actual value is read depeneds on the 'type' parameter
//   - 'type': the type of the query. One of "DOI", "ARXIV" or "QUERY"
//   - 'maxDepth': the depth of the tree to retrieve.
//   - 'limit': how many papers to retrieve. Only valid if type == "QUERY"
func request_paper(w http.ResponseWriter, r *http.Request) {
	// Enable CORS.
	(w).Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "GET" {
		http.Error(w, "Only GET root is supported", http.StatusMethodNotAllowed)
		return
	}
	queryValues := r.URL.Query()

	// Get the id of the paper to retrieve.
	link, present := queryValues["id"]
	if !present {
		http.Error(w, "Missing query param 'id'", http.StatusBadRequest)
		return
	}
	assert(len(link) == 1)

	// Get the type of the link.
	link_type, present := queryValues["type"]
	if !present {
		http.Error(w, "Missing query param 'type'", http.StatusBadRequest)
		return
	}
	assert(len(link_type) == 1)

	maxDepthStr, present := queryValues["maxDepth"]
	maxDepth := 2
	if !present {
		log.Printf("Max depth was no supplied, falling back to default value (2)\n")
	} else {
		assert(len(maxDepthStr) == 1)
		parsed, err := strconv.Atoi(maxDepthStr[0])
		if err != nil {
			http.Error(w, "Cannot parse maxDepth", http.StatusBadRequest)
		}
		maxDepth = parsed
	}
	log.Printf("ID: %s, Type = %s!, Max Depth: %d\n", link[0], link_type[0], maxDepth)

	// Query the database for DOI.
	var paper *request.Paper

	// Init the query column in db.
	var matching_column string
	db_value := link[0] // How to read the value in db.
	switch link_type[0] {
	case "DOI":
		matching_column = "ExternalId.DOI"
	case "ARXIV":
		matching_column = "ExternalId.ArXiv"
		// @NOTE: This trick is required because Semantic Scholars returns only the
		// digits that compose the arxiv, without the 'arxiv' prefix. If the
		// input query has this value, we need to trim it, otherwise the db
		// won't see match anything and it'll insert a new value.
		parts := strings.Split(link[0], ":")
		if len(parts) == 2 {
			db_value = parts[1]
		}
	case "QUERY":
		// We want to perform a query term. Pass the handler to another method.
		request_by_query(w, r, link[0], maxDepth)
		return
	default:
		http.Error(w, fmt.Sprintf("Invalid 'type' %s", link_type[0]), http.StatusBadRequest)
		return
	}
	assert(len(matching_column) > 0)

	paper, err := database.IsPresent("papers", matching_column, db_value)

	if errors.Is(err, database.ErrNoDocuments) {
		// Request the DOI.
		err := request.MakeRequest(fmt.Sprintf(AllInfoRoute, link[0]), &paper)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error while retrieving link from Scholar. Error: %s", err), http.StatusInternalServerError)
		}

		// Insert the element in db.
		err = database.Insert_element("papers", paper)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error while inserting paper in db: ERROR %s", err), http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		http.Error(w, fmt.Sprintf("Error while retrieving link: ERROR %s", err), http.StatusInternalServerError)
		return
	}

	start := time.Now()
	// Get all the citations for the input paper.
	var subtree map[string]request.Paper = make(map[string]request.Paper)
	citation_tree := insert_citations(paper, maxDepth, subtree)
	elapsed := time.Since(start)
	log.Printf("All citations have been computed in %s\n", elapsed)

	fmt.Fprint(w, request.PrettyPrint(citation_tree))
}

func request_by_query(w http.ResponseWriter, r *http.Request, searchTerm string, maxDepth int) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")

	queryValues := r.URL.Query()
	// Get the id of the paper to retrieve.
	limit, present := queryValues["limit"]
	if !present {
		limit = []string{"5"}
		log.Printf("Defult limit to %s\n", limit[0])
	}
	assert(len(limit) == 1)

	s := url.QueryEscape(searchTerm)

	_url := fmt.Sprintf(QuerySearchRoute, s, limit[0])
	var response request.Response

	err := request.MakeRequest(_url, &response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error while creating the request. Error: %s", err), http.StatusInternalServerError)
		return
	}

	// For each paper we got, check if have it in db. If so, replace it.
	// In this way, we can avoid the destruction of the graph.
	for idx, paper := range response.Data {
		paperInDb, err := database.IsPresent("papers", "PaperId", paper.PaperId)

		if err != nil && !errors.Is(err, database.ErrNoDocuments) {
			http.Error(w, fmt.Sprintf("Error while retrieving document from database: ERROR %s", err), http.StatusInternalServerError)
			return
		} else if !errors.Is(err, database.ErrNoDocuments) {
			log.Printf("Match in db for paper %s\n", paper.PaperId)
			// Paper is in db. Substitute.
			response.Data[idx] = *paperInDb
		}
	}

	trees := make([]Node, len(response.Data))
	var subtree map[string]request.Paper = make(map[string]request.Paper)
	for idx, paper := range response.Data {
		log.Printf("Processing paper number %d\n", idx+1)
		trees[idx] = insert_citations(&paper, maxDepth, subtree)
	}
	// To avoid unparsable characters
	json_string, err := json.Marshal(trees)
	w.Write(json_string)
}

func request_pdf(w http.ResponseWriter, r *http.Request) {
	// Enable CORS.
	(w).Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "GET" {
		http.Error(w, "Only GET root is supported", http.StatusMethodNotAllowed)
	} else {
		queryValues := r.URL.Query()
		pdf, present := queryValues["PDF"]

		if !present {
			http.Error(w, "Missing query param PDF url", http.StatusBadRequest)
			return
		}
		assert(len(pdf) == 1)
		log.Printf("PDF = %s!\n", pdf[0])

		// Generate the request to the pdf url
		req, err := http.NewRequest("GET", pdf[0], nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error while creating the request. Error: %s", err), http.StatusInternalServerError)
		}

		// Add some headers to avoid cors in sites
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.140 Safari/537.36 Edge/17.17134")

		// Generate the client to make the request
		client := &http.Client{}

		// Making an HTTP request to retrieve a PDF file from the URL.
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error while retrieving PDF. Error: %s", err), http.StatusInternalServerError)
		}
		defer resp.Body.Close()

		// Reading the response body from the HTTP request and storing it in the `pdfBytes` variable.
		pdfBytes, err := io.ReadAll(resp.Body)

		if err != nil {
			http.Error(w, fmt.Sprintf("Error while reading the PDF. Error: %s", err), http.StatusInternalServerError)
		}
		w.Header().Set("Content-type", "application/pdf")
		w.Write(pdfBytes)
	}
}

func main() {
	log.Printf("Using %d CPUs\n", numCPU)

	// DB Connection
	var err error
	database.Client, err = database.Established_connection()
	if err != nil {
		panic(err)
	}

	log.Println("Connection established")
	defer func() {
		if err := database.Client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	database.Database = database.CreateDatabase("connected-papers")
	if err != nil {
		panic(err)
	}
	log.Println("Database created!")

	var collection = database.Database.Collection("papers")

	err = database.Create_index(collection, "PaperId")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/paper", request_paper)
	http.HandleFunc("/pdf", request_pdf)
	log.Printf("Starting the server on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
