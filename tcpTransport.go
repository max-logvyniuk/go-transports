package transport

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"
)

type resS struct {
	Pattern string
	Data    string
}

type FileOptions struct {
	Fieldname    string
	Originalname string
	Encoding     string
	Mimetype     string
	Size         int
}

func TCPRouter(ln net.Listener, routes map[string]func(conn net.Conn, data string)) {
	for {
		conn, err := ln.Accept()

		// timeoutDuration := 500 * time.Second

		// conn.SetDeadline(time.Now().Add(timeoutDuration))

		// conn.SetReadDeadline(time.Now().Add(timeoutDuration))

		fmt.Println("TCPRouter conn - ", conn)

		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn, routes)
	}
}

func handleConnection(conn net.Conn, routes map[string]func(conn net.Conn, data string)) {
	defer conn.Close()

	// fmt.Println("handleConnection ---", conn)

	handleMessages(conn, routes)
}

func handleMessages(conn net.Conn, routes map[string]func(conn net.Conn, data string)) {
	defer conn.Close()

	bs, err := ioutil.ReadAll(conn)

	// buf := make([]byte, 100000000)
	// n, err := conn.Read(buf)

	if err != nil {
		fmt.Println(err)
	}

	//TODO: get message_length
	// incData := string(buf[:n])
	incData := string(bs)

	// fmt.Println("incData ---", incData)

	incDataSlice := strings.Split(incData, "#")

	messageAndFileMetaSlice := strings.Split(incDataSlice[0], "@")

	messageMeta := messageAndFileMetaSlice[0]
	// fmt.Println("messageMeta0 ---", messageMeta)

	messageMetaSlice := strings.Split(messageMeta, "=")

	metaLength := len([]byte(incDataSlice[0])) + 1

	// var fileOptions *FileOptions

	// fmt.Println("messageAndFileMetaSlice ---", messageAndFileMetaSlice, messageAndFileMetaSlice[1]);

	// err = json.Unmarshal([]byte(messageAndFileMetaSlice[1]), &fileOptions)

	// if err != nil {
	// 	fmt.Println("FileOptions unmarshal error ---", err)
	// }

	fileOptions := ""

	if len(messageAndFileMetaSlice) > 1 {
		fileOptions = messageAndFileMetaSlice[1]
	}

	fmt.Println("fileOptions ---", fileOptions)

	// fmt.Println("messageMeta1 ---", messageMetaSlice[0], messageMetaSlice, metaLength, len(bs))

	n := len(bs)

	switch messageMetaSlice[0] {
	case "message_length":
		handleStreamConnection(conn, bs[metaLength:], n, fileOptions, routes)
	default:
		handleSimpleConnection(conn, bs[metaLength:], routes)

	}
}

func handleSimpleConnection(conn net.Conn, data []byte, routes map[string]func(conn net.Conn, data string)) {

	jsonRequest := string(data)

	var parsedRequest resS

	err := json.Unmarshal([]byte(jsonRequest), &parsedRequest)

	if err != nil {
		fmt.Println("handleSimpleConnection: Request parsing error ---", err)

		conn.Write([]byte("Request parsing error"))
	}

	// fmt.Println("handleConnection parsedRequest ---", parsedRequest)

	router(conn, routes, parsedRequest)

}

func handleStreamConnection(conn net.Conn, data []byte, bl int, fileOptionsS string, routes map[string]func(conn net.Conn, data string)) {
	defer conn.Close()

	var fileOptions *FileOptions

	fmt.Println("messageAndFileMetaSlice ---", fileOptionsS)

	err := json.Unmarshal([]byte(fileOptionsS), &fileOptions)

	fmt.Println("handleStreamConnection len ---", len(data))

	// FOR STREAMING: need to think how implement streams handling
	err = ioutil.WriteFile("/tmp/"+fileOptions.Originalname, data, 0644)

	if err != nil {
		fmt.Println("Error ioutil.WriteFile ---", err)
		conn.Write([]byte("Error ioutil.WriteFile ---"))
	} else {
		conn.Write([]byte("Data saved"))
	}

}

func router(conn net.Conn, r map[string]func(conn net.Conn, data string), parsedRequest resS) {
	useDefaultPattern := true

	for k, v := range r {

		if k == parsedRequest.Pattern {
			v(conn, parsedRequest.Data)

			useDefaultPattern = false
			break
		}
	}

	if useDefaultPattern {
		conn.Write([]byte("Unknowen pattern"))
	}
}

// func handleConnection(conn net.Conn, routes map[string]func(conn net.Conn, data string)) {
// 	defer conn.Close()

// 	// scanner := bufio.NewScanner(conn)

// 	// fmt.Println("scanner ----", scanner, scanner.Text())

// 	// for scanner.Scan() {
// 	// 	message := scanner.Text()
// 	// 	fs := strings.Fields(message)
// 	// 	bsL := strings.Split(fs[0], "=")
// 	// 	fmt.Println("message ---", fs[0], bsL[1])

// 	// 	if len(fs) < 2 {
// 	// 		continue
// 	// 	}
// 	// }

// 	buf := make([]byte, 10000000)

// 	n, err := conn.Read(buf)

// 	if err != nil {
// 		fmt.Println("conn.Read error", err)
// 	}

// 	println("handleConnection data from client 0 -", n)

// 	bs, err := ioutil.ReadAll(conn)

// 	if err != nil {
// 		fmt.Println("Error ioutil.ReadAll ---", err)
// 	}

// 	fmt.Println("ioutil.ReadAll ---", bs, string(bs))

// 	// FOR STREAMING: need to think how implement streams handling
// 	err = ioutil.WriteFile("/tmp/dat2", bs, 0644)

// 	if err != nil {
// 		fmt.Println("Error ioutil.WriteFile ---", err)
// 	}

// 	jsonRequest := string(buf[:n])
// 	println("HandleConnection data from client 1 -", jsonRequest, buf)

// 	// i := strings.Index(jsonRequest, "#")

// 	var parsedRequest resS

// 	// println("Data from client 2---", i, jsonRequest[i+1:], string(buf[:n]))

// 	// If use nest.js ClientPoxy
// 	// err = json.Unmarshal([]byte(jsonRequest[i+1:]), &parsedRequest)
// 	err = json.Unmarshal([]byte(jsonRequest), &parsedRequest)

// 	if err != nil {
// 		fmt.Println("Request parsing error ---", err)
// 	}

// 	fmt.Println("handleConnection parsedRequest ---", parsedRequest)

// 	router(conn, routes, parsedRequest)
// }

// func router(conn net.Conn, r map[string]func(conn net.Conn, data string), parsedRequest resS) {
// 	useDefaultPattern := true

// 	for k, v := range r {

// 		if k == parsedRequest.Pattern {
// 			v(conn, parsedRequest.Data)

// 			useDefaultPattern = false
// 			break
// 		}
// 	}

// 	if useDefaultPattern {
// 		conn.Write([]byte("Unknowen pattern"))
// 	}
// }
