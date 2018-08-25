package pb

import (
	"fmt"
	"github.com/golang/protobuf/proto"
)

func (t *TripDescriptor) GetDirection() string {
	tripInterface, err := proto.GetExtension(t, E_NyctTripDescriptor)
	if err != nil {
		fmt.Printf("Could not get NyctTripDescriptor extension: %v\n", err)
	}

	nyctTrip, found := tripInterface.(*NyctTripDescriptor)
	if !found {
		fmt.Println("NYCT trip not present")
		return ""
	}

	return nyctTrip.GetDirection().String()
}
