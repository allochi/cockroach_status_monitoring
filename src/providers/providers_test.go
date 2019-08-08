package providers

import (
	"io/ioutil"
	"testing"
)

func TestHTTPParseNodes(t *testing.T) {
	filename := "mocks/http_response.json"
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Error("couldn't find mock file: ", filename)
	}

	var sp HTTPStatusProvider
	nodes, err := sp.parseNodes(content)

	if err != nil {
		t.Error("failed to parse well formatted response")
	}

	if len(nodes) != 4 {
		t.Errorf("failed to parse well formatted response, expected 4 nodes for %d", len(nodes))
	}

	// test for default http address
	if nodes[0].HTTPAddress != "localhost:8080" {
		t.Errorf("failed to retrieve http address, expected 'localhost:8080' got '%s'", nodes[0].HTTPAddress)
	}

	// test low in memory
	if !nodes[3].IsLowInMemory {
		t.Errorf("failed to report low in memory for node 4")
	}

}

func TestCMDParseNodes(t *testing.T) {
	filename := "mocks/cmd_response.csv"
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Error("couldn't find mock file: ", filename)
	}

	var sp CmdStatusProvider
	nodes, err := sp.parseNodes(content)

	if err != nil {
		t.Error("failed to parse well formatted response")
	}

	if len(nodes) != 4 {
		t.Errorf("failed to parse well formatted response, expected 4 nodes for %d", len(nodes))
	}

	// test not live node 3
	if nodes[2].IsLive {
		t.Errorf("failed to report not live for node 3")
	}
}

// func TestHTTPHealthCalls(t *testing.T) {
// 	healthEndpoint := func(w http.ResponseWriter, r *http.Request) {
// 		if len(r.URL.Query()["ready"]) > 0 {
// 			http.Error(w, "", http.StatusServiceUnavailable)
// 		} else {
// 			http.Error(w, "", http.StatusServiceUnavailable)
// 		}

// 	}

// 	ts := httptest.NewServer(http.HandlerFunc(healthEndpoint))
// 	defer ts.Close()

// 	{
// 		client, _ := http.Get(ts.URL + "/health?ready=1")
// 		fmt.Println(client.Status)
// 	}

// }
