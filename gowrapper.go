// Copyright 2017 Google LLC All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	sdk "agones.dev/agones/sdks/go"
	"github.com/kidoman/go-steam"
)

type interceptor struct {
	forward   io.Writer
	intercept func(p []byte)
}

// Write will intercept the incoming stream, and forward
// the contents to its `forward` Writer.
func (i *interceptor) Write(p []byte) (n int, err error) {
	if i.intercept != nil {
		i.intercept(p)
	}

	return i.forward.Write(p)
}

func getIP() (addr string) {

	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println(err)
	}

	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	split := strings.Split(localAddr.String(), ":")
	return split[0] + ":27015"
}

// main intercepts the stdout of the csgo gameserver and uses it
// to determine if the game server is ready or not.
func main() {

	fmt.Println(">>> Connecting to Agones with the SDK")
	s, err := sdk.NewSDK()
	if err != nil {
		log.Fatalf(">>> Could not connect to sdk: %v", err)
	}

	fmt.Println(">>> Starting health checking")
	go doHealth(s)
	go doPing(s)
	fmt.Println(">>> Starting wrapper for csgo!")
	cmd := exec.Command("/home/csgo/hlserver/csgo.sh") // #nosec
	cmd.Stderr = &interceptor{forward: os.Stderr}
	cmd.Stdout = &interceptor{
		forward: os.Stdout,
		intercept: func(p []byte) {

			str := strings.TrimSpace(string(p))
			// csgo will say "Server listening" 4 times before being ready,
			// once for ipv4 and once for ipv6.
			// but it does it each twice because it loads the maps between
			// each one, and resets state as it does so
			fmt.Println(str)
		}}
	err = cmd.Start()
	if err != nil {
		log.Fatalf(">>> Error Starting Cmd %v", err)
	}

	err = cmd.Wait()
	log.Fatal(">>> csgo shutdown unexpectantly", err)
}

// doHealth sends the regular Health Pings
func doHealth(sdk *sdk.SDK) {
	tick := time.Tick(2 * time.Second)
	for {
		err := sdk.Health()
		if err != nil {
			log.Fatalf("[wrapper] Could not send health ping, %v", err)
		}
		<-tick
	}
}

func doPing(sdk *sdk.SDK) {
	tick := time.Tick(2 * time.Second)
	for {
		fmt.Printf(">>> doPing() got ip: %s \n", getIP())
		addr := getIP()
		server, err := steam.Connect(addr)
		if err != nil {
			fmt.Printf("err!= nil => %v \n", err)
		}
		defer server.Close()

		info, err := server.Info()
		if err != nil {
			fmt.Printf("steam: could not get server info from %v: %v\n", addr, err)

		}
		if info != nil {
			fmt.Printf("steam: info of %v: %v\n", addr, info)
			sdk.Ready()
			break
		}
		<-tick

	}

}
