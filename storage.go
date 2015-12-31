package morningbot

import (
	"fmt"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	"github.com/tucnak/telebot"
)

// saveUserToDB saves the given user to the time bucket in Bolt.
func (m *MorningBot) saveUser(sender *telebot.User) error {
	userID := strconv.Itoa(sender.ID)

	err := m.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(subscriptions_bucket_name)

		// temporarily hardcode for GMT Offset
		gb, err := b.CreateBucketIfNotExists([]byte("+8"))
		if err != nil {
			return err
		}

		err = gb.Put([]byte(userID), []byte(time.Now().Format(time.RFC3339)))
		if err != nil {
			return err
		}

		return nil
	})
	return err
}

// saveUserToDB saves the given user to the time bucket in Bolt.
func (m *MorningBot) removeUser(sender *telebot.User) error {
	userID := strconv.Itoa(sender.ID)

	err := m.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(subscriptions_bucket_name)

		// temporarily hardcode for GMT Offset
		gb, err := b.CreateBucketIfNotExists([]byte("+8"))
		if err != nil {
			return err
		}

		err = gb.Delete([]byte(userID))
		if err != nil {
			return err
		}

		return nil
	})
	return err
}

// getAllIDsForBroadcast returns all user IDs for a morning broadcast
func (m *MorningBot) getAllIDsForBroadcast() ([]string, error) {
	uArray := []string{}

	err := m.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(subscriptions_bucket_name)

		// temporarily hardcode for GMT Offset
		gb := b.Bucket([]byte("+8"))
		if gb == nil {
			return fmt.Errorf("error retrieving bucket for broadcast time %s", 8)
		}

		gb.ForEach(func(k, v []byte) error {
			key, _ := string(k), string(v)
			uArray = append(uArray, key)
			return nil
		})
		return nil
	})

	if err != nil {
		return []string{}, err
	}

	return uArray, nil
}
