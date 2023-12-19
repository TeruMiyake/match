package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/TeruMiyake/match/network"
)

func RunClient() {
	fmt.Println("Running client...")

	// サーバーに接続する
	for connectToServer("127.0.0.1", 40800) != nil {
		fmt.Println("Failed to connect to the server. Retrying...")
		time.Sleep(5 * time.Second)
	}
}

func connectToServer(raddr string, port int) error {
	laddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:0")
	if err != nil {
		return err
	}
	dialer := net.Dialer{
		// 各 Read(), Write() ではなく接続の確立に使用されるタイムアウト
		Timeout:   time.Duration(5) * time.Second,
		LocalAddr: laddr,
	}

	conn, err := dialer.Dial("tcp4", fmt.Sprintf("%s:%d", raddr, port))
	if err != nil {
		return err
	}
	defer conn.Close()

	// サーバーへ playerInfo を送信する
	// hasId := false であると仮定して 0 を送信する
	// TODO: プロンプトで ID を所持しているか確認して場合分けする
	conn.Write([]byte{0})
	log.Println("Sent player info: 0")

	// サーバーから接続結果 joinResult を受信する
	joinResult, err := network.ReadFull(conn, 1)
	log.Println("Received join result: ", joinResult)
	if err != nil {
		return err
	}
	if joinResult[0] != 0 {
		log.Println("Failed to join the game")
		// TODO: 何らかの形でリトライする
		return err
	}

	// プレイヤー ID を受信する
	log.Println("Successfully joined the game")
	b, err := network.ReadFull(conn, 1)
	if err != nil {
		return err
	}
	me.Id = b[0]
	log.Println("Received player ID: ", me.Id)

	// サーバーからゲーム情報を受信する
	g, err := network.ReadGame(conn)
	if err != nil {
		return err
	}

	// ゲーム情報を表示する
	g.Illustrate()

	// プレイヤーのコンソール入力を受け付ける
	sc := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("Enter message: ")
		sc.Scan()

		msg := sc.Text()
		if strings.ToLower(msg) == "exit" {
			break
		}

		fmt.Println("Sending message: ", msg)
		conn.SetDeadline(time.Now().Add(5 * time.Second))
		_, err := conn.Write([]byte(msg))
		if err != nil {
			log.Println("Failed to send message: ", err)
			return err
		}
	}
	return nil
}
