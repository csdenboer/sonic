package main

import (
	"fmt"

	"github.com/csdenboer/sonic"
)

func main() {
	ioc := sonic.MustIO()

	conn, err := sonic.Dial(ioc, "tcp", "localhost:8082")
	if err != nil {
		panic(err)
	}

	b := []byte("hello, sonic!")

	conn.AsyncWrite(b, func(err error, n int) {
		if err != nil {
			panic(err)
		}

		fmt.Println("wrote:", string(b))

		b = make([]byte, 32)
		conn.AsyncRead(b, func(err error, n int) {
			if err != nil {
				panic(err)
			}
			fmt.Println("read:", string(b))
		})
	})

	ioc.RunPending()
}
