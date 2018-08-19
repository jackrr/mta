package api

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jackrr/mta/pb"
	"io/ioutil"
	"net/http"
)

type MTA struct {
	ApiKey string
}

func unmarshall(data []byte) (string, error) {
	res := &pb.NyctFeedHeader{}
	err := proto.Unmarshal(data, res)
	if err != nil {
		fmt.Printf("Unmarshall failed: %s", err)
	}
	return res.String(), nil
}

func (m MTA) GetData() {
	url := fmt.Sprintf("http://datamine.mta.info/mta_esi.php?key=%s&feed_id=1", m.ApiKey)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("HTTP Request failed: %s", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Body to string failed: %s", err)
	}

	fmt.Println(unmarshall(body))
}
