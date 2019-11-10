package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"time"
	"syscall"
)

func runServer(port int) {
	var str1 string = strconv.Itoa(port)

	// cmd := exec.Command("echo", str1)
	// stdoutStderr, err := cmd.CombinedOutput()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("%s\n", stdoutStderr)
	s, i := fmt.Printf("sudo docker run -d -p %s:%s -p %s:%s/udp csgo_test -console -usercon +game_type 0 game_mode 1 +mapgropup mg_active +port %s +map de_cache +sv_setsteamaccount %s\n", str1, str1, str1, str1, str1, str1)
	fmt.Print(s, i)

	fmt.Printf("%s \n", i)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func getLastupdatefromCSnet() (unixtime int64) {

	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://blog.counter-strike.net/index.php/category/updates/feed")
	t, err := time.Parse("2006-01-02 15:04:05 +0000 UTC", feed.UpdatedParsed.String())
	if err != nil {
		fmt.Println(err)
	}
	unixtime = t.Unix()
	return
}


func we_need_update() (result bool){
	if _, err := os.Stat("myfile"); err == nil {
		readfile, err := ioutil.ReadFile("myfile")
		check(err)
		String_read_from_file := string(readfile)
		Converted_str_rss := strconv.FormatInt(getLastupdatefromCSnet(), 10)

// we assume that file exists and that it has been run before

		if String_read_from_file == Converted_str_rss {
			return false
		} else {
			return true
		}

// we assume file does not exist, we need to make it.
	}else if os.IsNotExist(err) {
		f, err := os.Create("myfile")
		check(err)
	
		defer f.Close()
		f.WriteString(fmt.Sprintf("%d", getLastupdatefromCSnet()))
// shady errorhandling
	} else {

		fmt.Print("err.. something failed", err)
		// Schrodinger: file may or may not exist. See err for details.

		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence

	}
return true
}

func runUpdate() {
	fmt.Print("[*] runningrebuilding of new image [*]")
	binary, lookErr := exec.LookPath("docker")
    if lookErr != nil {
        panic(lookErr)
    }

    args := []string{"docker", "build", "-f", "Dockerfile_addons", "-t", "csgo_test", ".", "--no-cache"}


    env := os.Environ()

    execErr := syscall.Exec(binary, args, env)
    if execErr != nil {
        panic(execErr)
    }
		}

func main() {
	for {

	if we_need_update(){
		t := time.Now()
		fmt.Printf("[*] %d-%02d-%02d - %02d:%02d:%02d Updating [*]\n",
        t.Day(), t.Month(),  t.Year(),
		t.Hour(), t.Minute(), t.Second())
		
		runUpdate()		
	}else
	{
		t := time.Now()
		fmt.Printf("[*] %d-%02d-%02d - %02d:%02d:%02d nothing new. doing nothing [*]\n",
        t.Day(), t.Month(),  t.Year(),
        t.Hour(), t.Minute(), t.Second())
	}
	var sleeptime time.Duration = 5
	fmt.Printf("[*] checking again in %d minutes[*]\n", sleeptime)
	time.Sleep(sleeptime * time.Minute)
}
}
