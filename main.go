package main

import (
	"fmt"
	"net"

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
		_, err := res.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		writer := resp.NewWriter(conn)
		writer.Write(resp.Value{
			Typ: resp.String,
			Str: "OK",
		})
	}
}
