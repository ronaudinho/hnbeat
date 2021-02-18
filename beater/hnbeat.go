package beater

import (
	"errors"
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
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
			_, wordFreq, err := bt.countWords()
			if err != nil {
				logp.Err(err.Error())
				continue
			}

			event := beat.Event{
				Timestamp: time.Now(),
				Fields: common.MapStr{
					"type": b.Info.Name,
					// "words":  wordCloud,
					"sorted": wordFreq,
				},
			}
			bt.client.Publish(event)
			logp.Info("Event sent")
		}
	}
}

// Stop stops hnbeat.
func (bt *hnbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}

// countWords counts words that appear since last max item.
func (bt *hnbeat) countWords() (map[string]int, map[int][]string, error) {
	wordCloud := make(map[string]int)
	x, err := bt.hnClient.GetMaxItem()
	if err != nil {
		return nil, nil, err
	}

	if bt.hnClient.MaxItem == 0 {
		bt.hnClient.MaxItem = x - int64(bt.config.ItemCount)
	} else if bt.hnClient.MaxItem == x {
		return nil, nil, errors.New("no new item")
	}

	var words []string
	for bt.hnClient.MaxItem < x {
		i, err := bt.hnClient.GetItem(x)
		if err != nil {
			logp.Err(err.Error())
			continue
		}
		r := strings.NewReplacer(".", " ", ",", " ", ";", " ", ":", " ", "!", " ", "?", " ", "-", " ", "_", " ")
		if i.Text != "" {
			words = append(words, strings.Fields(
				r.Replace(
					strings.ToLower(
						html.UnescapeString(i.Text),
					),
				),
			)...)
		}
		if i.Title != "" {
			words = append(words, strings.Fields(r.Replace(strings.ToLower(html.UnescapeString(i.Title))))...)
		}
		x--
	}

	bt.hnClient.MaxItem = x
	for _, word := range words {
		_, ok := wordCloud[word]
		if ok {
			wordCloud[word]++
		} else {
			wordCloud[word] = 1
		}
	}

	wordFreq := make(map[int][]string)
	for k, v := range wordCloud {
		wordFreq[v] = append(wordFreq[v], k)
	}

	return wordCloud, wordFreq, nil
}
