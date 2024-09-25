package domain

// Text describes a text entity with associated metadata and text content.
type Text struct {
	Metadata string
	Text     string
}

// NewText creates a new Text object with provided metadata and text content.
func NewText(metadata, text string) *Text {
	return &Text{
		Metadata: metadata,
		Text:     text,
	}
}

// Clone creates a copy of the current Text instance, including its Metadata
// and Text fields, and returns a pointer to the new copy.
func (t *Text) Clone() *Text {
	return &Text{
		Metadata: t.Metadata,
		Text:     t.Text,
	}
}
