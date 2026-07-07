package atwtest

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/awootton/knotfreeiot/iot"
	"github.com/awootton/knotfreeiot/packets"
	"github.com/awootton/knotfreeiot/tokens"
	"github.com/coredns/coredns/plugin/pkg/log"
)

func TestContact(t *testing.T) {

	host := os.Getenv("TARGET_CLUSTER") // "localhost:8384"
	if host == "" {
		host = "knotfree.io:8384"
		fmt.Println("No TARGET_CLUSTER environment variable set, using default host", host)
	}
	token := os.Getenv("KNOTFREE_TOKEN")
	if token == "" {
		fmt.Println("No KNOTFREE_TOKEN environment variable set, using default token")
		tokenTmp, payload := tokens.GetImpromptuGiantTokenLocal("", "")
		fmt.Println("Using token", tokenTmp, "with payload", payload)
		token = tokenTmp
	}

	serviceController, err := iot.StartNewServiceContactTcp(host, token)
	if err != nil {
		log.Error("knotfree setup failed to start service contact", err)
		return // exit the test
	}

	for i := 0; i < 9999999999; i++ {
		// log.Info("knotfree setup started service contact, iteration ", i)

		command := "get option " + strings.ToUpper("TXT") + " " + "meta_group_id" // eg get option A

		subscriptionName := "testmain-0n0u0e16p-0_vr"

		cmd := packets.Lookup{}
		cmd.Address.FromString(subscriptionName)
		cmd.SetOption("cmd", []byte(command))
		// send it
		replyPacket, err := serviceController.Get(&cmd)
		if err != nil {
			log.Error("knotfree failed to get from service contact", err)
			time.Sleep(time.Duration(10) * time.Second)
			continue // to next iteration
		}
		//log.Info("knotfree returned from service contact", replyPacket.Sig())
		sendPacket, ok := replyPacket.(*packets.Send)
		if !ok {
			log.Error("knotfree failed to get 'send' from service contact", err, replyPacket.Sig())
			time.Sleep(time.Duration(10) * time.Second)
			continue // to next iteration
		}
		log.Info("knotfree returned message ", string(sendPacket.Payload))
		sleepTime := 30
		// log.Info("Sleeping for ", sleepTime, " seconds before next iteration")
		time.Sleep(time.Duration(sleepTime) * time.Second)
	}

}
