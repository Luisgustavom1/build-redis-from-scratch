package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/Luisgustavom1/build-redis-from-scratch/resp"
)

func main() {
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	for {
		res := resp.NewResp(conn)
		value, err := res.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		if value.Typ != resp.Array {
			fmt.Println("Expected Array")
			return
		}

		if len(value.Array) == 0 {
			fmt.Println("Expected at least one element")
			continue
		}

		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		writer := resp.NewWriter(conn)

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(resp.Value{
				Typ: resp.Error,
				Str: "Invalid command",
			})
			continue
		}

		result := handler(args)
		writer.Write(result)
	}
}
