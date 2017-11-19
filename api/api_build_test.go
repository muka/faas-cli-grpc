package api

import (
	"io/ioutil"
	"log"
	"testing"
	"time"

	context "golang.org/x/net/context"
	"google.golang.org/grpc"
)

//TODO use as test main()
// func TestStartStop(t *testing.T) {
// 	go func() {
// 		if err := Start(); err != nil {
// 			t.FailNow()
// 		}
// 	}()
// 	time.Sleep(time.Millisecond * 500)
// 	Stop()
// }

func TestBuild(t *testing.T) {

	const uri = "localhost:51001"
	go func() {
		if err := Start(uri); err != nil {
			t.FailNow()
		}
	}()
	time.Sleep(time.Millisecond * 500)

	conn, connerr := grpc.Dial(uri, grpc.WithInsecure())
	if connerr != nil {
		log.Fatal(connerr)
		t.FailNow()
	}

	const tmpArchive = "../tmp/tmp-test1.zip"
	//Create a tmp archive
	if ziperr := Zip("../test", tmpArchive); ziperr != nil {
		log.Fatal(ziperr)
		t.FailNow()
	}

	// read buffer
	buf, readerr := ioutil.ReadFile(tmpArchive)
	if readerr != nil {
		log.Fatal(readerr)
		t.FailNow()
	}

	client := NewFaasCliServiceClient(conn)
	client.Build(context.Background(), &BuildRequest{
		Archive: buf,
		// Image:      "localhost:5000/test1",
		// Lang:       "node",
		// Name:       "test1",
		NoCache:    true,
		Parallel:   5,
		Shrinkwrap: true,
		Squash:     true,
	})

	Stop()
}
