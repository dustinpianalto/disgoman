package disgoman

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"math/rand"
	"time"
)

/* status-manager.go:
 * Built in status manager which cycles through a list of status at a specified interval
 *
 * Disgoman (c) 2020 Dusty.P/dustinpianalto
 */

// Add a status to the manager
func (s *StatusManager) AddStatus(status string) {
	s.Values = append(s.Values, status)
}

// Remove a status from the manager
func (s *StatusManager) RemoveStatus(status string) []string {
	for i, v := range s.Values {
		if v == status {
			s.Values = append(s.Values[:i], s.Values[i+1:]...)
			break
		}
	}
	return s.Values
}

// Sets interval to new value
func (s *StatusManager) SetInterval(interval string) {
	s.Interval = interval
}

// Update the status now
func (s *StatusManager) UpdateStatus(session *discordgo.Session) error {
	i := rand.Intn(len(s.Values))
	err := session.UpdateStatus(0, s.Values[i])
	return err
}

// Default StatusManager ready function which updates the status at the specified interval
func (s *StatusManager) OnReady(session *discordgo.Session, _ *discordgo.Ready) {
	interval, err := time.ParseDuration(s.Interval)
	if err != nil {
		return
	}

	err = s.UpdateStatus(session)
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(interval)
	for range ticker.C {
		err = s.UpdateStatus(session)
		if err != nil {
			log.Fatal(err)
		}
	}
}