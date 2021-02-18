package hn

// ItemType is the type of HN item.
type ItemType string

var (
	// Job ItemType
	Job ItemType = "job"
	// Story ItemType
	Story = "story"
	// Comment ItemType
	Comment = "comment"
	// Poll ItemType
	Poll = "poll"
	// PollOpt ItemType
	PollOpt = "pollopt"
)

// Item defines HN Item.
type Item struct {
	ID          int64    `json:"id"`                    // The item's unique id.
	Deleted     bool     `json:"deleted,omitempty"`     // true if the item is deleted.
	Type        ItemType `json:"type,omitempty"`        // The type of item. One of "job", "story", "comment", "poll", or "pollopt".
	By          string   `json:"by,omitempty"`          // The username of the item's author.
	Time        int64    `json:"time,omitempty"`        // Creation date of the item, in Unix Time.
	Text        string   `json:"text,omitempty"`        // The comment, story or poll text. HTML.
	Dead        bool     `json:"dead,omitempty"`        // true if the item is dead.
	Parent      int64    `json:"parent,omitempty"`      // The comment's parent: either another comment or the relevant story.
	Poll        string   `json:"poll,omitempty"`        // The pollopt's associated poll.
	Kids        []int64  `json:"kids,omitempty"`        // The ids of the item's comments, in ranked display order.
	Url         string   `json:"url,omitempty"`         // The URL of the story.
	Score       int64    `json:"score,omitempty"`       // The story's score, or the votes for a pollopt.
	Title       string   `json:"title,omitempty"`       // The title of the story, poll or job. HTML.
	Parts       []int64  `json:"parts,omitempty"`       // A list of related pollopts, in display order.
	Descendants int64    `json:"descendants,omitempty"` // In the case of stories or polls, the total comment count.
}

// MaxItem is the current largest item ID.
type MaxItem int64
