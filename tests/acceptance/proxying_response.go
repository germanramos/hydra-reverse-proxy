package acceptance

// import (
// 	"net/http"
// 	"testing"
// )

// func TestProxyingApplication(t *testing.T) {
// 	ts := httptest.NewServer(http.HandlerFunc("/app/app1", func(w http.ResponseWriter, r *http.Request) {
// 		sortedInstanceUris := []string{"http://www.q1.com", "http://www.q2.com", "http://www.q3.com"}
// 		jsonOutput, _ = json.Marshal(sortedInstanceUris)
// 		w.WriteHeader(http.StatusOK)
// 		w.Header().Set("Content-Type", "application/json")
// 		w.Write(jsonOutput)
// 	}))
// 	defer ts.Close()

// 	res, err := http.Get(ts.URL)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	instances, err := ioutil.ReadAll(res.Body)
// 	res.Body.Close()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if instances == []byte{"http://www.q1.com", "http://www.q2.com", "http://www.q3.com"} {
// 		t.Error("ERROR")
// 	}
// }
