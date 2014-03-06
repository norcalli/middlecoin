package middlecoin

import (
	"encoding/json"
	"net/http"
	"testing"
)

const url = "http://www.middlecoin.com/json"

func TestParsing(t *testing.T) {
	addresses := []string{
		"1DgxRTofdbau7kpf3pQeRydcoTPG2L5NUX",
		"17Nt7rWiRZKDgcNp421zZ1FHGPWSnnT1bk",
	}
	var decoder *json.Decoder
	if true {
		response, err := http.Get(url)
		if err != nil {
			t.SkipNow()
			// log.Fatal(err)
		}
		defer response.Body.Close()
		decoder = json.NewDecoder(response.Body)
	}

	r := new(OverviewReport)
	err := decoder.Decode(&r)
	if err != nil {
		t.Fail()
		// log.Fatal(err)
	}
	total := new(AddressReport)
	for _, address := range addresses {
		report, ok := r.Report[address]
		if !ok {
			t.Fail()
		}
		total.Add(report)
	}
}
