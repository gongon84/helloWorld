package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/quickfixgo/quickfix"
)

type MyApplication struct{}

func (app MyApplication) OnCreate(sessionID quickfix.SessionID)                       {}
func (app MyApplication) OnLogon(sessionID quickfix.SessionID)                        {}
func (app MyApplication) OnLogout(sessionID quickfix.SessionID)                       {}
func (app MyApplication) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) {}
func (app MyApplication) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) error {
	return nil
}
func (app MyApplication) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	return nil
}
func (app MyApplication) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	log.Printf("メッセージ受信しました: %s\n", msg.String())
	return nil
}

func main() {
	app := MyApplication{}
	storeFactory := quickfix.NewMemoryStoreFactory()
	logFactory := quickfix.NewScreenLogFactory()

	cfgAcceptor := "./config/acceptor.cfg"
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println(dir)
	cfgAc, err := os.Open(cfgAcceptor)
	if err != nil {
		log.Fatalf("acceptor.cfgの読み込みに失敗: %v", err)
	}
	defer cfgAc.Close()

	settingsAc, err := quickfix.ParseSettings(cfgAc)
	if err != nil {
		log.Fatalf("設定の解析に失敗: %v", err)
	}

	// アクセプタの起動
	acceptor, err := quickfix.NewAcceptor(app, storeFactory, settingsAc, logFactory)
	if err != nil {
		log.Fatalf("アクセプタの作成に失敗: %v", err)
	}
	go func() {
		if err := acceptor.Start(); err != nil {
			log.Fatalf("アクセプタの起動に失敗: %v", err)
		}
		fmt.Printf("アクセプタの起動開始\n")
		defer acceptor.Stop()

		select {}
	}()

	cfgInitiator := "./config/initiator.cfg"
	cfgIn, err := os.Open(cfgInitiator)
	if err != nil {
		log.Fatalf("initiator.cfgの読み込みに失敗: %v", err)
	}
	defer cfgIn.Close()

	settingsIn, err := quickfix.ParseSettings(cfgIn)
	if err != nil {
		log.Fatalf("設定の解析に失敗: %v", err)
	}

	time.Sleep(5 * time.Second)
	// イニシエータの起動
	initiator, err := quickfix.NewInitiator(app, storeFactory, settingsIn, logFactory)
	if err != nil {
		log.Fatalf("イニシエータの作成に失敗: %v", err)
	}

	go func() {
		if err := initiator.Start(); err != nil {
			log.Fatalf("イニシエータの起動に失敗: %v", err)
		}
		fmt.Printf("イニシエータの起動開始\n")
		defer initiator.Stop()

		// テストメッセージを送信
		// time.Sleep(5 * time.Second) // 適当な遅延を設ける
		testMsg := quickfix.NewMessage()
		// メッセージの設定（例：Logonメッセージ）
		testMsg.Header.SetField(quickfix.Tag(35), quickfix.FIXString("A"))
		quickfix.SendToTarget(testMsg, quickfix.SessionID{
			BeginString:  "FIX.4.2",
			SenderCompID: "SENDER",
			TargetCompID: "TARGET",
		})

		select {}
	}()

	// サーバーを継続的に動作させる
	select {}
}
