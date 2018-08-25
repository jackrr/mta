package api

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jackrr/mta/pb"
	"io/ioutil"
	"net/http"
)

type FeedGetter struct {
	ApiKey string
}

func NewFeedGetter(key string) FeedGetter {
	return FeedGetter{ApiKey: key}
}

func (f FeedGetter) GetFeed() pb.FeedMessage {
	url := fmt.Sprintf("http://datamine.mta.info/mta_esi.php?key=%s&feed_id=1", f.ApiKey)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("HTTP Request failed: %s", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Body to string failed: %s", err)
	}

	feed := &pb.FeedMessage{}
	err = proto.Unmarshal(body, feed)
	if err != nil {
		fmt.Printf("Unmarshall failed: %s", err)
	}

	return *feed
}
