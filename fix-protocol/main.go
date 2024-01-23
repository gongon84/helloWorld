package main

import (
	"bytes"
	"fmt"

	"github.com/quickfixgo/quickfix"
)

func main() {
	// セッションIDを手動で作成
	sessionID := quickfix.SessionID{
		BeginString:  "FIX.4.4",
		SenderCompID: "SenderCompID",
		TargetCompID: "TargetCompID",
	}

	// FIXメッセージのバッファを作成（SOH文字で区切る）
	// rawMessage := bytes.NewBufferString(strings.ReplaceAll("8=FIX.4.4|35=0|49=SenderCompID|56=TargetCompID|34=1|52=20240123-15:30:45|", "|", "\x01"))
	// rawMessage := bytes.NewBufferString("8=FIX.4.4\x019=79\x0135=0\x0149=SenderCompID\x0156=TargetCompID\x0134=1\x0152=20240123-15:30:45\x0110=000\x01")
	// メッセージ本体（タグ35から始まる部分）
	messageBody := "35=0\x0149=SenderCompID\x0156=TargetCompID\x0134=1\x0152=20240123-15:30:45\x01"

	// メッセージ本体の長さを計算
	messageBodyLength := len(messageBody)

	// 完全なFIXメッセージを構築
	rawMessageString := fmt.Sprintf("8=FIX.4.4\x019=%d\x01%s10=000\x01", messageBodyLength, messageBody)
	rawMessage := bytes.NewBufferString(rawMessageString)

	// FIXメッセージをパース
	msg := quickfix.NewMessage()
	err := quickfix.ParseMessage(msg, rawMessage)
	if err != nil {
		fmt.Println("Error parsing FIX message:", err)
		return
	}

	// パースしたメッセージを表示
	fmt.Println("Parsed Message:", msg.String())

	// メッセージを送信
	err = quickfix.SendToTarget(msg, sessionID)
	if err != nil {
		fmt.Println("Error sending FIX message:", err)
		return
	}

	fmt.Println("Message sent successfully!")
}
