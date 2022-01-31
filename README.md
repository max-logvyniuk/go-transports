# Golang  transports

## 1.TCP

Inside there is TCP router for handling tcp connection from clients written in different languages(tested with node.js client)

### 1.1 handleSimpleConnection
Import module and use for routing tcp connections like this:


```bash

/** in package */
type TContext struct {
	Conn    net.Conn
	Options string
	Data    []byte
}

routes := map[string]func(tc transports.TContext){
		"calculate":         controllers.CalculationGo,
}

transports.TCPRouter(ln, routes)
```

routes - is  map[string]func(tc transports.TContext){} that describe controllers for "handleSimpleConnection" function

ln - is instance of net.Listener

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

This method start to execute when you specify additional params in beginning of tcp stream. It helps to process files to golang server

#### Example of string and buffer data that sends from client server to golang server:

`message_length={LENGTHOFBUFFER}@{FILEOPTIONS}#` - {METADATA} in string representation (you need to convert it to buffer)

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

