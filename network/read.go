package network

import (
	"encoding/binary"
	"io"
	"log"
	"net"

	"github.com/TeruMiyake/match/game"
)

func ReadFull(conn net.Conn, size int) ([]byte, error) {
	buffer := make([]byte, size)
	offset := 0
	for offset < size {
		n, err := conn.Read(buffer[offset:])
		if err != nil {
			if err != io.EOF {
				return nil, err // ネットワークエラー
			}
			break
		}
		offset += n
	}
	return buffer[:offset], nil
}

func ReadGame(conn net.Conn) (*game.Game, error) {
	log.Println("Reading game data...")

	// まず、ゲームデータのサイズを読み込む
	sizeBuf, err := ReadFull(conn, 4)
	if err != nil {
		return nil, err
	}
	size := binary.BigEndian.Uint32(sizeBuf)
	log.Println("Game data size:", size)

	// ゲームデータを読み込む
	gameData, err := ReadFull(conn, int(size))
	if err != nil {
		return nil, err
	}
	log.Println("Game data:", gameData)

	// game.Game.FromBytes()を使用してゲームデータを解析する
	game := game.FromBytes(gameData)
	return game, nil
}
