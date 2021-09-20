package main

import (
	"context"
	"fmt"
	"nba-top-shot-insights/topshot"

	"github.com/onflow/flow-go-sdk/client"
	"google.golang.org/grpc"
)

const version = "0.1.0"

func handleErr(err error, message string) {
	if err != nil {
		fmt.Printf(message+"\n\nError: %v.\n", err)
		panic(err)
	}
}

const blocksToInspect = 500
const maxBlockRange = 249
const TopShotAccountAddress = "c1e4f4f4c4257510"
const TopShotMarketContract = "Market"

func main() {
	fmt.Printf("NBA Top Shot Insights v%v.\n", version)

	// var MomentPurchased = fmt.Sprintf("A.%s.%s.MomentPurchased", TopShotAccountAddress, TopShotMarketContract)
	var MomentListed = "A.c1e4f4f4c4257510.TopShotMarketV3.MomentListed"

	flowClient, err := client.New("access.mainnet.nodes.onflow.org:9000", grpc.WithInsecure())

	handleErr(err, "Connection failed")

	// fetch latest block
	latestBlock, err := flowClient.GetLatestBlock(context.Background(), false)

	handleErr(err, "Can't get last block")

	fmt.Println("current height: ", latestBlock.Height)

	startHeight := latestBlock.Height - blocksToInspect
	endHeight := startHeight + maxBlockRange

	for endHeight <= latestBlock.Height {
		fmt.Printf("Searching from %v to %v\n", startHeight, endHeight)

		purchaseEvents, err := flowClient.GetEventsForHeightRange(context.Background(), client.EventRangeQuery{
			Type:        MomentListed,
			StartHeight: startHeight,
			EndHeight:   endHeight,
		})

		handleErr(err, "Can't get listed events from blockchain")

		for _, blockEvent := range purchaseEvents {
			if len(blockEvent.Events) != 0 {
				for _, purchaseEvent := range blockEvent.Events {
					e := topshot.MomentListedEvent(purchaseEvent.Value)

					saleMoment, _ := topshot.GetSaleMomentFromOwnerAtBlock(flowClient, blockEvent.Height-1, *e.Seller(), e.Id())

					// if err != nil {
					// 	fmt.Println(err)
					// }

					if saleMoment != nil && saleMoment.Series() == 1 {
						// serialNumber := strconv.FormatUint(uint64(saleMoment.SerialNumber()), 10)

						// if isPalindrome(serialNumber) {
						fmt.Println(e)
						fmt.Println(saleMoment)
						fmt.Println()
						// }
					}
				}
			}
		}

		startHeight = endHeight + 1
		endHeight = startHeight + maxBlockRange
	}

}

func isPalindrome(number string) bool {
	for i := 0; i < len(number); i++ {
		j := len(number) - 1 - i
		if number[i] != number[j] {
			return false
		}
	}
	return true
}
