package deviceworker

import (
	"math/rand"
	"time"
)

// randomInterval returns a random duration between 1-10 seconds
func randomInterval() time.Duration {
	seconds := rand.Intn(10) + 1
	return time.Duration(seconds) * time.Second
}

// randomUsername returns a random username from a predefined list
func randomUsername() string {
	usernames := []string{
		"john.doe", "jane.smith", "mike.johnson", "sarah.williams",
		"david.brown", "emily.jones", "chris.davis", "lisa.miller",
		"tom.wilson", "anna.moore", "james.taylor", "mary.anderson",
		"robert.thomas", "linda.jackson", "william.white", "jennifer.harris",
	}
	return usernames[rand.Intn(len(usernames))]
}

func init() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())
}
