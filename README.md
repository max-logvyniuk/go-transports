# Golang  transports

## 1.TCP

Inside there is TCP router for handling tcp connection from clients written in different languages(tested with node.js client)

### 1.1 handleSimpleConnection
Import module and use for routing tcp connections like this:

```bash
routes := transport.RouterMap{
		"calculate":         controllers.CalculationGo,
}

transport.TCPRouter(ln, routes)
```

routes - of type:

````bash
type RouterMap map[string]func(conn net.Conn, options string, data []byte)
````
ln - is instence of net.Listener

```bash
ln, err := net.Listen("tcp", ":4444")
```

### Example of incomming message that can process handleSimpleConnection:

```bash

messageData: 
{
  pattern: string,
  data: any,
}

in stringify and converted to []byte;  

 "{LENGTH_OF_BUFFER}#messageData"
```

In you only specify buffer length and "#" separator {LENGTH_OF_BUFFER}

### 1.2 handleStreamConnection

This method start to execute when you spesify additional params in begining of tcp stream. It helps to process files to golang server 

#### Example of string and buffer data that sends from client server to golang server:

`message_length={LENGTH_OF_BUFFER}@pattern={PATTERN}@{OPTIONS}#` - {METADATA} in string representation (you need to convert it to buffer)

`{METADATA}{BUFFER}` - you need to concat two buffers into one

{LENGTH_OF_BUFFER} - its length of buffer

{PATTERN} - name of appropriate handler func

{OPTIONS} - options = {
           [string]: any
          }

 @ - separator between {LENGTH_OF_BUFFER}, {PATTERN} and {OPTIONS}

 `#` - separator between message meta and {BUFFER}       

