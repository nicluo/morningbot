package morningbot

//go:generate go-bindata -pkg $GOPACKAGE -o assets.go data/

import (
  "fmt"
  "log"
  "os"
  "path"
  "regexp"
  "strings"
  "time"

  "github.com/boltdb/bolt"
  "github.com/kardianos/osext"
  "github.com/tucnak/telebot"
)

var subscriptions_bucket_name = []byte("subscriptions")

type MorningBot struct {
  Name string // The name of the bot registered with Botfather
  bot  *telebot.Bot
  log  *log.Logger
  fmap FuncMap
  db   *bolt.DB
  keys map[string]string
}

// Wrapper struct for a message
type message struct {
  Cmd  string
  Args []string
  *telebot.Message
}

// GetArgs prints out the arguments for the message in one string.
func (m message) GetArgString() string {
  argString := ""
  for _, s := range m.Args {
    argString = argString + s + " "
  }
  return strings.TrimSpace(argString)
}

// A FuncMap is a map of command strings to response functions.
// It is use for routing comamnds to responses.
type FuncMap map[string]ResponseFunc

type ResponseFunc func(m *message)

// Initialise a MorningBot.
// lg is optional.
func InitMorningBot(name string, bot *telebot.Bot, lg *log.Logger, config map[string]string) *MorningBot {
  if lg == nil {
    lg = log.New(os.Stdout, "[morningbot] ", 0)
  }
  m := &MorningBot{Name: name, bot: bot, log: lg, keys: config}

  m.fmap = m.getDefaultFuncMap()

  // Setup database
  // Get current executing folder
  pwd, err := osext.ExecutableFolder()
  if err != nil {
    lg.Fatalf("cannot retrieve present working directory: %s", err)
  }

  db, err := bolt.Open(path.Join(pwd, "morningbot.db"), 0600, nil)
  if err != nil {
    lg.Fatal(err)
  }
  m.db = db
  createAllBuckets(db)

  return m
}

// Get the built-in, default FuncMap.
func (m *MorningBot) getDefaultFuncMap() FuncMap {
  return FuncMap{
    "/start":           m.Start,
    "/help":            m.Start,
    "/subscribe":       m.Subscribe,
    "/unsubscribe":     m.Unsubscribe,
  }
}

// Add a response function to the FuncMap
func (m *MorningBot) AddFunction(command string, resp ResponseFunc) error {
  if !strings.Contains(command, "/") {
    return fmt.Errorf("not a valid command string - it should be of the format /something")
  }
  m.fmap[command] = resp
  return nil
}

// Route received Telegram messages to the appropriate response functions.
func (m *MorningBot) Router(msg telebot.Message) {
  jmsg := m.parseMessage(&msg)
  if jmsg.Cmd != "" {
    m.log.Printf("[%s] command: %s, args: %s", time.Now().Format(time.RFC3339), jmsg.Cmd, jmsg.GetArgString())
  }
  execFn := m.fmap[jmsg.Cmd]

  if execFn != nil {
    m.GoSafely(func() { execFn(jmsg) })
  }
}

func (m *MorningBot) CloseDB() {
  m.db.Close()
}

// Ensure all buckets needed by morningbot are created.
func createAllBuckets(db *bolt.DB) error {
  // Check all buckets have been created
  err := db.Update(func(tx *bolt.Tx) error {
    _, err := tx.CreateBucketIfNotExists(subscriptions_bucket_name)
    if err != nil {
      return err
    }
    return nil
  })
  return err
}

// Helper to parse incoming messages and return MorningBot messages
func (m *MorningBot) parseMessage(msg *telebot.Message) *message {
  cmd := ""
  args := []string{}

  if msg.IsReply() {
    // We use a hack. All reply-to messages have the command it's replying to as the
    // part of the message.
    r := regexp.MustCompile(`\/\w*`)
    res := r.FindString(msg.ReplyTo.Text)
    for k, _ := range m.fmap {
      if res == k {
        cmd = k
        args = strings.Split(msg.Text, " ")
        break
      }
    }
  } else {
    msgTokens := strings.Split(msg.Text, " ")
    cmd, args = strings.ToLower(msgTokens[0]), msgTokens[1:]
  }

  return &message{Cmd: cmd, Args: args, Message: msg}
}

// SendMessage Shorthand
func (m *MorningBot) SendMessage(recipient telebot.Recipient, msg string, options *telebot.SendOptions) {
  m.bot.SendMessage(recipient, msg, options)
}
