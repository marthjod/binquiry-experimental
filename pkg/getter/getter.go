package getter

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/marthjod/binquiry-experimental/pkg/reader"
	"golang.org/x/net/html"
	"gopkg.in/xmlpath.v2"
)

// getter builds query URLs and HTTP requests against a data source.
type getter struct {
	Client         *http.Client
	URLPrefix      string
	responseBodies [][]byte
	logger         *log.Entry
}

type result struct {
	err  error
	body []byte
}

func NewGetter(urlPrefix string, client *http.Client, correlationID string) *getter {
	return &getter{
		Client:    client,
		URLPrefix: urlPrefix,
		logger: log.WithFields(log.Fields{
			"cid":  correlationID,
			"task": "getter",
		}),
	}
}

// WordQuery builds a URL for querying a word.
func (g *getter) WordQuery(word string) (query string) {
	v := url.Values{}
	v.Set("q", word)
	return g.URLPrefix + "?" + v.Encode()
}

// IDQuery builds a URL for querying a search ID.
func (g *getter) IDQuery(id int) (query string) {
	return fmt.Sprintf("%s?id=%d", g.URLPrefix, id)
}

// GetWord makes an HTTP request for a word against the data source.
func (g *getter) GetWord(word string) (responses [][]byte, err error) {
	var emptyResult = [][]byte{}

	query := g.WordQuery(word)
	g.logger.WithFields(log.Fields{
		"query": query,
	}).Debug()

	req, err := http.NewRequest(http.MethodGet, query, nil)
	if err != nil {
		log.Error(err)
		return emptyResult, err
	}

	resp, err := g.Client.Do(req)
	if err != nil {
		log.Error(err)
		return emptyResult, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return emptyResult, err
	}
	defer resp.Body.Close()

	return g.dispatch(reader.Sanitize(b))
}

// GetID makes an HTTP request for a search ID against the data source.
func (g *getter) GetID(id int) (*http.Response, error) {
	query := g.IDQuery(id)
	g.logger.WithFields(log.Fields{
		"query": query,
	}).Debug()
	req, err := http.NewRequest(http.MethodGet, query, nil)
	if err != nil {
		return nil, err
	}

	return g.Client.Do(req)
}

func (g *getter) fetchLink(link string, resultChan chan<- result) {
	id, err := getSearchID(link)
	if err != nil {
		resultChan <- result{
			err:  err,
			body: []byte{},
		}
		return
	}

	r, err := g.GetID(id)
	if err != nil {
		resultChan <- result{
			err:  err,
			body: []byte{},
		}
		return
	}

	body, err := readSanitized(r.Body)
	if err != nil {
		resultChan <- result{
			err:  err,
			body: []byte{},
		}
		return
	}

	resultChan <- result{
		err:  err,
		body: body,
	}
}

func (g *getter) dispatch(r []byte) (responses [][]byte, err error) {
	root, err := xmlpath.Parse(bytes.NewReader(r))
	if err != nil {
		return
	}

	// did we land on a multiple-choice page?
	qLinks := xmlpath.MustCompile("/ul/li/strong/a")
	if qLinks.Exists(root) {
		doc, err := html.Parse(bytes.NewReader(r))
		if err != nil {
			return [][]byte{}, err
		}

		links := getLinkNodes(doc)
		resultChan := make(chan result, len(links))

		for _, link := range links {
			go g.fetchLink(link, resultChan)
			select {
			case result := <-resultChan:
				if result.err != nil {
					g.logger.WithFields(log.Fields{
						"error": err,
					}).Error()
					continue
				}
				g.responseBodies = append(g.responseBodies, result.body)
			}
		}

		return g.responseBodies, nil
	}

	// add original response body if we did not land on a multiple-choice page
	g.responseBodies = append(g.responseBodies, r)

	return g.responseBodies, nil
}

func getSearchID(val string) (int, error) {
	var searchID = regexp.MustCompile(`leit_id\('(\d+)'\)`)
	groups := searchID.FindStringSubmatch(val)
	if len(groups) > 1 {
		return strconv.Atoi(groups[1])
	}
	return 0, errors.New("unable to convert search ID")
}

func getLinkNodes(doc *html.Node) []string {

	var (
		f     func(*html.Node)
		links []string
	)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "onclick" {
					links = append(links, a.Val)
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	return links
}

func readSanitized(r io.Reader) ([]byte, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return []byte{}, err
	}
	return reader.Sanitize(b), nil
}
