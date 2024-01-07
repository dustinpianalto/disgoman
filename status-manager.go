package disgoman

import (
	"log"
	"math/rand"
	"time"

	"github.com/dustinpianalto/discordgo"
)

/* status-manager.go:
 * Built in status manager which cycles through a list of status at a specified interval
 *
 * Disgoman (c) 2020 Dusty.P/dustinpianalto
 */

// AddStatus to the manager
func (s *StatusManager) AddStatus(status string) {
	s.Values = append(s.Values, status)
}

// RemoveStatus from the manager
func (s *StatusManager) RemoveStatus(status string) []string {
	for i, v := range s.Values {
		if v == status {
			s.Values = append(s.Values[:i], s.Values[i+1:]...)
			break
		}
	}
	return s.Values
}

// SetInterval changes the update interval to new value
func (s *StatusManager) SetInterval(interval string) {
	s.Interval = interval
}

// UpdateStatus updates the status of the bot
func (s *StatusManager) UpdateStatus(session *discordgo.Session) error {
	i := rand.Intn(len(s.Values))
	err := session.UpdateGameStatus(0, s.Values[i])
	log.Println(err)
	return err
}

// OnReady is the default StatusManager ready function which updates the status at the specified interval
func (s *StatusManager) OnReady(session *discordgo.Session, _ *discordgo.Ready) {
	interval, err := time.ParseDuration(s.Interval)
	if err != nil {
		return
	}

	err = s.UpdateStatus(session)
	if err != nil {
		log.Println(err)
	}

	ticker := time.NewTicker(interval)
	for range ticker.C {
		err = s.UpdateStatus(session)
		if err != nil {
			log.Println(err)
		}
	}
}
