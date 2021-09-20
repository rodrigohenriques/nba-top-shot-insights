package topshot

import (
	"context"
	"fmt"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
)

func GetSaleMomentFromOwnerAtBlock(flowClient *client.Client, blockHeight uint64, ownerAddress flow.Address, momentFlowID uint64) (*SaleMoment, error) {
	getSaleMomentScript := `
		import TopShotMarketV3 from 0xc1e4f4f4c4257510
		import TopShot from 0x0b2a3299cc857e29
		import Market from 0xc1e4f4f4c4257510

        pub struct SaleMoment {
          pub var id: UInt64
		  pub var playId: UInt32
		  pub var play: {String: String}
		  pub var serialNumber: UInt32
		  pub var setId: UInt32
          pub var setName: String
		  pub var series: UInt32
          
		  init(moment: &TopShot.NFT) {
            self.id = moment.id
            self.playId = moment.data.playID
            self.play = TopShot.getPlayMetaData(playID: self.playId) ?? panic("Play doesn't exist")
            self.serialNumber = moment.data.serialNumber
			self.setId = moment.data.setID
            self.setName = TopShot.getSetName(setID: self.setId) ?? panic("Could not find the specified set")
			self.series = TopShot.getSetSeries(setID: self.setId) ?? panic("Could not find the specified set series")
          }
        }

		pub fun main(sellerAddress:Address, momentID:UInt64): SaleMoment? {
			let acct = getAccount(sellerAddress)
		
			let momentCollectionRef = acct.getCapability(/public/MomentCollection).borrow<&{TopShot.MomentCollectionPublic}>()
				?? panic("Could not get public moment collection reference")
		
			let moment = momentCollectionRef.borrowMoment(id: momentID)
				?? panic("Could not borrow a reference to the specified moment")

			return SaleMoment(moment: moment!)
		}
`
	res, err := flowClient.ExecuteScriptAtBlockHeight(context.Background(), blockHeight, []byte(getSaleMomentScript), []cadence.Value{
		cadence.BytesToAddress(ownerAddress.Bytes()),
		cadence.UInt64(momentFlowID),
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching sale moment from flow: %w", err)
	}
	optional := res.(cadence.Optional)
	saleMoment := SaleMoment(optional.Value.(cadence.Struct))
	return &saleMoment, nil
}

type SaleMoment cadence.Struct

func (s SaleMoment) ID() uint64 {
	return uint64(s.Fields[0].(cadence.UInt64))
}

func (s SaleMoment) PlayID() uint32 {
	return uint32(s.Fields[1].(cadence.UInt32))
}

// func (s SaleMoment) SetName() string {
// 	return string(s.Fields[4].(cadence.String))
// }

// func (s SaleMoment) SetID() uint32 {
// 	return uint32(s.Fields[3].(cadence.UInt32))
// }

func (s SaleMoment) Play() map[string]string {
	dict := s.Fields[2].(cadence.Dictionary)
	res := map[string]string{}
	for _, kv := range dict.Pairs {
		res[string(kv.Key.(cadence.String))] = string(kv.Value.(cadence.String))
	}
	return res
}

func (s SaleMoment) SerialNumber() uint32 {
	return uint32(s.Fields[3].(cadence.UInt32))
}

func (s SaleMoment) SetID() uint32 {
	return uint32(s.Fields[4].(cadence.UInt32))
}

func (s SaleMoment) SetName() string {
	return string(s.Fields[5].(cadence.String))
}

func (s SaleMoment) Series() uint32 {
	return uint32(s.Fields[6].(cadence.UInt32))
}

func (s SaleMoment) FullName() string {
	return s.Play()["FullName"]
}

func (s SaleMoment) TeamAtMoment() string {
	return s.Play()["TeamAtMoment"]
}

func (s SaleMoment) JerseyNumber() string {
	return s.Play()["JerseyNumber"]
}

func (s SaleMoment) isRookie() bool {
	return s.Play()["TotalYearsExperience"] == "1"
}

func (s SaleMoment) String() string {
	return fmt.Sprintf("SaleMoment:\n\tplayer: %s, serialNumber: %d, team: %s, set: %s (Series %d)", s.FullName(), s.SerialNumber(), s.TeamAtMoment(), s.SetName(), s.Series())
}
