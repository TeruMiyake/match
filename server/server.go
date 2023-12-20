package server

import (
	"fmt"
	"log"
	"net"

	"github.com/TeruMiyake/match/game"
	"github.com/TeruMiyake/match/network"
)

const port int = 40800

const (
	width  uint8 = 5
	height uint8 = 3
)

var g *game.Game

type JoinResult uint8

const (
	JoinSuccess JoinResult = iota
	CapacityOver
	InvalidPlayerInfo
	GameAlreadyFinished
)

func RunServer() {
	fmt.Println("Running server...")

	g = game.NewGame(width, height)
	g.Illustrate()

	startTCPServer()
}

func startTCPServer() {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("Accepted %v\n", conn.RemoteAddr())

	// クライアントから playerInfo を受信する
	var result JoinResult
	var playerId uint8
	hasIdByte, err := network.ReadFull(conn, 1)
	log.Println("received hasIdByte:", hasIdByte)
	if err != nil {
		log.Println(err)
		return
	}
	switch hasIdByte[0] {
	case 0:
		// クライアントは ID を持っていない
		// クライアントに ID を割り当てる
		playerId = g.AddPlayer()
		// クライアントに ID を送信する
		result = JoinSuccess
	default:
		// TODO
		// クライアントが ID を持っていると主張している
		// が、まだこの処理は実装されていない
		log.Println("hasIdByte is not 0 (yet to be implemented)")
		return
	}

	// クライアントに結果を送信する
	log.Println("sending result... :", result)
	conn.Write([]byte{byte(result)})
	if result == JoinSuccess {
		// クライアントに playerId を送信する
		conn.Write([]byte{byte(playerId)})
	}
	log.Println("Result sent.")

	// game.Cells をバイト列で送信する
	bytes := g.Bytes()
	conn.Write(bytes)
	log.Println("Game data sent.: ", bytes)

	msgReceiver := NewMessageReceiver()

	// クライアントによりコネクションが切断されるまで Read し続ける
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				log.Println("Connection closed by client.")
			} else {
				log.Println(err)
			}
			return
		}
		msgReceiver.Receive(buf[:n])
	}
}
