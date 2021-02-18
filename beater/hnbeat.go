package beater

import (
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/logp"

	"github.com/ronaudinho/hnbeat/config"
	"github.com/ronaudinho/hnbeat/hn"
)

// hnbeat configuration.
type hnbeat struct {
	done     chan struct{}
	config   config.Config
	client   beat.Client
	hnClient *hn.Client
}

// New creates an instance of hnbeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &hnbeat{
		done:     make(chan struct{}),
		config:   c,
		hnClient: hn.NewClient(),
	}
	return bt, nil
}

// Run starts hnbeat.
func (bt *hnbeat) Run(b *beat.Beat) error {
	logp.Info("hnbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(bt.config.Period)
	counter := 1
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
			wordMap, err := bt.countWords()
			if err != nil {
				logp.Err(err.Error())
				continue
			}

			event := beat.Event{
				Timestamp: time.Now(),
				Fields: common.MapStr{
					"type":  b.Info.Name,
					"words": wordMap,
				},
			}
			bt.client.Publish(event)
			logp.Info("Event sent")
			counter++
		}
	}
}

// Stop stops hnbeat.
func (bt *hnbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}

// countWords counts words that appear since last max item.
func (bt *hnbeat) countWords() (map[string]int, error) {
	wordMap := make(map[string]int)
	x, err := bt.hnClient.GetMaxItem()
	if err != nil {
		return nil, err
	}
	if bt.hnClient.MaxItem == 0 {
		bt.hnClient.MaxItem = x - 100
	}
	var words []string
	for bt.hnClient.MaxItem < x {
		i, err := bt.hnClient.GetItem(x)
		if err != nil {
			logp.Err(err.Error())
			continue
		}
		if i.Text != "" {
			r := strings.NewReplacer(".", " ", ",", " ", ";", " ", "!", " ", "?", " ", "-", " ", "_", " ")
			words = append(words, strings.Fields(
				r.Replace(
					strings.ToLower(
						html.UnescapeString(i.Text),
					),
				),
			)...)
		}
		if i.Title != "" {
			r := strings.NewReplacer(".", " ", ",", " ", ";", " ", "!", " ", "?", " ", "-", " ", "_", " ")
			words = append(words, strings.Fields(r.Replace(strings.ToLower(html.UnescapeString(i.Title))))...)
		}
		x--
	}

	bt.hnClient.MaxItem = x
	for _, word := range words {
		_, ok := wordMap[word]
		if ok {
			wordMap[word]++
		} else {
			wordMap[word] = 1
		}
	}
	return wordMap, nil
}
