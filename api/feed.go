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

func AllFeeds() []int {
	return []int{1, 26, 16, 21, 2, 11, 31, 36, 51}
}

// GetFeed Returns the feed with id specified
// Possible ids:
// 1	--	123456S
// 26	--	ACES
// 16	--	NQRW
// 21	--	BDFM
// 2	--	L
// 11	--	Staten IS
// 31	--	G
// 36	--	JZ
// 51	--	7
func (f FeedGetter) GetFeed(id int) pb.FeedMessage {
	url := fmt.Sprintf("http://datamine.mta.info/mta_esi.php?key=%s&feed_id=%d", f.ApiKey, id)
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
