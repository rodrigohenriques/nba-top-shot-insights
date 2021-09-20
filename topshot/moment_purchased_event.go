package topshot

import (
	"fmt"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
)

// pub event MomentPurchased(id: UInt64, price: UFix64, seller: Address?)
type MomentListedEvent cadence.Event

func (evt MomentListedEvent) Id() uint64 {
	return uint64(evt.Fields[0].(cadence.UInt64))
}

func (evt MomentListedEvent) Price() float64 {
	return float64(evt.Fields[1].(cadence.UFix64).ToGoValue().(uint64)) / 1e8 // ufixed 64 have 8 digits of precision
}

func (evt MomentListedEvent) Seller() *flow.Address {
	optionalAddress := (evt.Fields[2]).(cadence.Optional)
	if cadenceAddress, ok := optionalAddress.Value.(cadence.Address); ok {
		sellerAddress := flow.BytesToAddress(cadenceAddress.Bytes())
		return &sellerAddress
	}
	return nil
}

func (evt MomentListedEvent) String() string {
	return fmt.Sprintf("Moment Listed(id: %d, price: %2f, seller: %s)",
		evt.Id(), evt.Price(), evt.Seller())
}
