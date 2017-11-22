package api

import (
	"io/ioutil"
	"log"
	"os"
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

	os.Setenv("workdir", "../data")
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

	buf, readerr := ioutil.ReadFile("../test/test1.zip")
	if readerr != nil {
		log.Fatal(readerr)
		t.FailNow()
	}

	client := NewFaasCliServiceClient(conn)
	res, err := client.Build(context.Background(), &BuildRequest{
		Archive: buf,
		Yaml:    "test1.yml",
		// Handler: "./test1"
		// Image:      "localhost:5000/test1",
		// Language:       "node",
		// Name:       "test1",
		NoCache:    true,
		Parallel:   5,
		Shrinkwrap: true,
		Squash:     true,
	})

	if err != nil {
		log.Fatal(err)
		t.FailNow()
	}

	log.Printf("Result: %d %s", res.Code, res.Message)

	Stop()
}
