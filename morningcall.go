package morningbot

import (
  "strconv"

  "github.com/tucnak/telebot"
)

func (m *MorningBot) Start(msg *message) {
  m.log.Printf(`%s %s %s`, msg.Sender.ID, msg.Sender.FirstName, msg.Sender.LastName)
  m.SendMessage(msg.Chat, `Hi there! I can help you with the following things:

/subscribe - Subscribes you to morning call everyday at 7AM GMT+8`, nil)
  //`/feedback - Not on GMT+8? Want a different time? This lets you get back to my creators`
}

func (m *MorningBot) Subscribe(msg *message) {
  m.log.Printf(`%s %s %s`, msg.Sender.ID, msg.Sender.FirstName, msg.Sender.LastName)
  
  m.saveUserToDB(&msg.Sender)
  m.SendMessage(msg.Chat, `You're Subscribed! 7AM GMT+8`, nil)
}

func (m *MorningBot) MorningCall() {
  m.log.Printf(`Starting Morning Call`);

  userIDs, _ := m.getAllIDsForBroadcast()
  for _, userID := range userIDs {
    i, err := strconv.Atoi(userID)
    if (err != nil){
      m.log.Printf(`Problem with %s`, userID);
    } else {
      m.log.Printf(`Sending Morning Call to %s`, userID);
      m.SendMessage(telebot.User{ID: i}, "Good Morning!", nil)
    }
  }
}