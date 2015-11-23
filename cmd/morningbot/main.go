package main

import (
  "encoding/json"
  "io/ioutil"
  "log"
  "os"
  "path"
  "time"

  "github.com/nicluo/morningbot"
  "github.com/kardianos/osext"
  "github.com/tucnak/telebot"
)

func main() {
  // Grab current executing directory
  // In most cases it's the folder in which the Go binary is located.
  pwd, err := osext.ExecutableFolder()
  if err != nil {
    log.Fatalf("error getting executable folder: %s", err)
  }
  configJSON, err := ioutil.ReadFile(path.Join(pwd, "config.json"))
  if err != nil {
    log.Fatalf("error reading config file! Boo: %s", err)
  }

  var config map[string]string
  json.Unmarshal(configJSON, &config)

  telegramAPIKey, ok := config["telegram_api_key"]
  if !ok {
    log.Fatalf("config.json exists but doesn't contain a Telegram API Key! Read https://core.telegram.org/bots#3-how-do-i-create-a-bot on how to get one!")
  }
  botName, ok := config["name"]
  if !ok {
    log.Fatalf("config.json exists but doesn't contain a bot name. Set your botname when registering with The Botfather.")
  }

  bot, err := telebot.NewBot(telegramAPIKey)
  if err != nil {
    log.Fatalf("error creating new bot, %s", err)
  }

  logger := log.New(os.Stdout, "[morningbot] ", 0)

  logger.Printf("Args: %s %s %s", botName, bot, logger)

  mb := morningbot.InitMorningBot(botName, bot, logger, config)
  defer mb.CloseDB()

  mb.GoSafely(func() {
    logger.Println("Scheduling Time Check")
    for {
      nextHour := time.Now().Truncate(time.Hour).Add(time.Hour)
      timeToNextHour := nextHour.Sub(time.Now())
      time.Sleep(timeToNextHour)
      logger.Printf("[%s] [%s] !", time.Now().Format(time.RFC3339), time.Now().Hour())
      if (time.Now().Hour() == 7){
        mb.MorningCall()
      }
    }
  })

  messages := make(chan telebot.Message)
  bot.Listen(messages, 1*time.Second)

  for message := range messages {
    mb.Router(message)
  }
}
