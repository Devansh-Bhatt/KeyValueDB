package test

import (
	"fmt"
	"net"
	"time"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

func Test() {
	serverAddr := "127.0.0.1:4000"

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	// request := "*2\r\n$4\r\necho\r\n$3\r\nhey\r\n"
	request2 := "+PING\r\n"
	reqset := "*5\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n$2\r\nPX\r\n:100\r\n"
	reqget := "*2\r\n$3\r\nGET\r\n$3\r\nfoo\r\n"
	_, err = conn.Write([]byte(reqset))
	_, err = conn.Write([]byte(request2))
	time.Sleep(2 * time.Second)
	_, err = conn.Write([]byte(reqget))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	respHandler := resp.NewRespHandler(conn)
	for {
		value, err := respHandler.ParseAny()
		if err != nil {
			fmt.Println(err)
		}
		switch value.Typ {
		case resp.StringType:
			fmt.Println("Response : ", value.Str)
		case resp.BulkStringType:
			fmt.Println("Response : ", value.Bulk)
		case resp.ErrorType:
			fmt.Println("Response : ", value.Err)
		}
	}

}
