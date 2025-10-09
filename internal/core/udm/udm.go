package udm

import (
	"encoding/json"
	"fmt"
	"os"
)

type Subscriber struct {
	IMSI               string `json:"imsi"`
	PhoneNumber        string `json:"phone_number"`
	SubscriptionStatus string `json:"subscription_status"`
	MaxDataRate        int    `json:"max_data_rate"`
}

type UDM struct {
	Subscribers map[string]*Subscriber // key = IMSI
}

type subscriberFile struct {
	Subscribers []Subscriber `json:"subscribers"`
}

func NewUDM() (*UDM, error) {
	file, err := os.Open("internal/configs/subscribers.json")
	if err != nil {
		return nil, fmt.Errorf("failed to open subscribers file: %w", err)
	}
	defer file.Close()

	var data subscriberFile
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to parse subscribers JSON: %w", err)
	}

	// Convert slice to map for O(1) lookups
	subscribers := make(map[string]*Subscriber)
	for i := range data.Subscribers {
		sub := &data.Subscribers[i]
		subscribers[sub.IMSI] = sub
	}

	return &UDM{
		Subscribers: subscribers,
	}, nil
}

// GetSubscriber checks if a subscriber exists and returns it
func (u *UDM) GetSubscriber(imsi string) (*Subscriber, error) {
	sub, exists := u.Subscribers[imsi]
	if !exists {
		return nil, fmt.Errorf("subscriber not found: %s", imsi)
	}
	return sub, nil
}
