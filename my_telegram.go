package main

import (
    "encoding/xml"
    "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
    "net/http"
    "io/ioutil"
)

var moex = map[string]string {
    "/etf": "https://iss.moex.com/iss/engines/stock/markets/shares/boards/TQTF/securities.xml?iss.meta=off&iss.only=marketdata&marketdata.columns=SECID,LAST",
}

type MOEX struct {
    Items []Rows `xml:"data>rows>row"`
}

type Rows struct {
    SECID string `xml:"SECID,attr"`
    LAST  string `xml:"LAST,attr"`
}

func getMoex(url string) (*MOEX, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)

    moex := new(MOEX)
    err = xml.Unmarshal(body, moex)
    if err != nil {
        return nil, err
    }

    return moex, nil
}

func main() {
    bot, err := tgbotapi.NewBotAPI("XXXXXXXXXX")
    if err != nil {
        log.Panic(err)
    }

    bot.Debug = true
    log.Printf("Authorized on account %s", bot.Self.UserName)
    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60
    updates, err := bot.GetUpdatesChan(u)

    for update := range updates {
        if url, ok := moex[update.Message.Text]; ok {
            moex, err := getMoex(url)
            if err != nil {
                bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка",))
            } else {
                bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Текущие курсы ETF:",))
                for _, item := range moex.Items {
                    if (item.SECID == "FXUS" || item.SECID == "FXIM" || item.SECID == "FXDM") {
                        bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, item.SECID + " - " + item.LAST,))
                    }
                }
            }
        } else {
            bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, `Неизвестная команда`,))
        }
    }
}
