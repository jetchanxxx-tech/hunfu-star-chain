package ai

type EmergencyDetector struct {
	trie     *Trie
	keywords []string
}

func NewEmergencyDetector(keywords []string) *EmergencyDetector {
	return &EmergencyDetector{
		trie:     NewTrie(keywords),
		keywords: keywords,
	}
}

type EmergencyResult struct {
	IsEmergency bool     `json:"is_emergency"`
	Matched     []string `json:"matched"`
}

func (d *EmergencyDetector) Check(text string) EmergencyResult {
	matched := d.trie.Search(text)
	return EmergencyResult{
		IsEmergency: len(matched) > 0,
		Matched:     matched,
	}
}
