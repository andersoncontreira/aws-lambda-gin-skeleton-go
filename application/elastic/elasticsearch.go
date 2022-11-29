package elastic

import (
	"context"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/olivere/elastic/v7"
	awsElk "github.com/olivere/elastic/v7/aws/v4"
	"log"
	"strings"
	"sync"
	"time"
)

type ElasticSearch struct {
	client       *elastic.Client
	bulkData     []ElasticDocs
	indexDefault string
}

type ElasticDocs struct {
	Channel   string    `json:"channel"`
	Extra     string    `json:"extra"`
	Level     int       `json:"level"`
	LevelName string    `json:"level_name"`
	Message   string    `json:"message"`
	DateTime  time.Time `json:"datetime"`
	Context   Context   `json:"context"`
	//Summary   string
	//Body      string
	//Timestamp int64
}

type Context struct {
	Command     string `json:"command"`
	Method      string `json:"method"`
	Type        string `json:"type"`
	Id          string `json:"id"`
	PayloadJson string `json:"payload_json"`
	Response    string `json:"response"`
}

type ElasticConfig struct {
	// Hosts struct of URLs passed initially to the client
	Hosts []string
	// Timeout sets the timeout for the sniffer that finds the nodes in a cluster default 2.
	// Timeout sets the timeout for the initial health check
	Timeout   int `default:"2"`
	AwsRegion string
	// IndexDefault set default index to elastic
	IndexDefault string
	AwsKey       string
	AwsSecret    string
}

var (
	sessionAws *session.Session
	once       sync.Once
)

// loadSession returns a new Session created from SDK aws defaults, config files.
func loadSession() {
	once.Do(func() {
		sessionAws, _ = session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		})
	})
}

// NewElasticSearch creates a new client to work with Elasticsearch.
// Uses the aws session to authenticate
// Example:
//
//	  esConfig := elastic.ElasticVars{Hosts: []string{http://127.0.0.1:9200, http://127.0.0.1:9201}, Timeout: 3, AwsRegion: "us-east-2"}
//		 elasticConnection, err := elastic.NewElasticSearch(esConfig)
func NewElasticSearch(elasticVars ElasticConfig) (ElasticSearch, error) {
	loadSession()

	signingClient := awsElk.NewV4SigningClient(sessionAws.Config.Credentials, elasticVars.AwsRegion)
	secondsTimeout := time.Second * time.Duration(elasticVars.Timeout)

	client, err := elastic.NewSimpleClient(
		elastic.SetURL(strings.Join(elasticVars.Hosts, ",")),
		elastic.SetSniff(false),
		elastic.SetHttpClient(signingClient),
		elastic.SetSnifferTimeout(secondsTimeout),
		elastic.SetSnifferTimeoutStartup(secondsTimeout),
		elastic.SetHealthcheckTimeoutStartup(secondsTimeout),
	)

	es := ElasticSearch{}
	es.indexDefault = elasticVars.IndexDefault + time.Now().Format("-2006-01")

	if err != nil {
		return es, err
	}

	es.client = client

	return es, nil
}

// Append ElasticDocs to send bulk data in Elastic
func (es *ElasticSearch) AddLog(doc ElasticDocs) {
	doc.DateTime = time.Now()
	es.bulkData = append(es.bulkData, doc)
}

// SendLogs sends the batched requests to Elasticsearch
// Index specifies the Elasticsearch index to be used for this index request.
// opType specifies if this request should follow create-only or upsert
// docType specifies the Elasticsearch type to be used for this index request.
func (es *ElasticSearch) SendLogs(ctx context.Context, index, opType, docType string) {
	bulk := es.client.Bulk()

	if index == "" {
		index = es.indexDefault
	}

	for _, doc := range es.bulkData {
		req := elastic.NewBulkIndexRequest()
		req.OpType(opType)
		req.Index(index)
		req.Type(docType)
		req.Doc(doc)
		bulk.Add(req)
	}

	_, err := bulk.Do(ctx)
	if err != nil {
		log.Println("Error sending bulk elastic", err)
		return
	}

}
