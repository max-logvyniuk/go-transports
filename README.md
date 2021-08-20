# Golang  transports

## 1.TCP

Inside there is TCP router for handling tcp connection from clients written in different languages(tested with node.js client)

### 1.1 handleSimpleConnection
Import module and use for routing tcp connections like this:

```bash
routes := map[string]func(conn net.Conn, data string){
		"calculate":         controllers.CalculationGo,
}

transports.TCPRouter(ln, routes)
```

routes - is  map[string]func(conn net.Conn, data string){} that describe controllers for "handleSimpleConnection" function

ln - is instence of net.Listener

```bash
ln, err := net.Listen("tcp", ":4444")
```

### Example of incomming message that can process handleSimpleConnection:

```bash
 "{LENGTHOFBUFFER}#"{
           pattern: string,
           data: any,
       }""
```

In you only specify buffer length and "#" separator {LENGTHOFBUFFER}

### 1.2 handleStreamConnection

This method start to execute when you spesify additional params in begining of tcp stream. It helps to process files to golang server 

#### Example of string and buffer data that sends from client server to golang server:

`message_length={LENGTHOFBUFFER}@{FILEOPTIONS}#` - {METADATA} in string retpresentation (you need to convert it to buffer)

`{METADATA}{BUFFER}` - you need to concat two buffers into one

{LENGTHOFBUFFER} - its length of buffer

{FILEOPTIONS} - file options = {
           fieldname string,
           originalname string,
           encoding string,
           mimetype string,
           size int,
        }

 @ - separator between {LENGTHOFBUFFER} and {FILEOPTIONS}

 `#` - separator between message meta and {BUFFER}       

