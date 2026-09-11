package custom

import (
	"encoding/json"
	_ "embed"
	"github.com/pgaskin/lithiumpatch/dict"
)

// --- Embed all 4 JSON files ---

//go:embed mw_collegiate.json
var mwCollegiateJSON []byte

//go:embed mw_learners.json
var mwLearnerJSON []byte

//go:embed oxford.json
var oxfordJSON []byte

//go:embed tolkien.json
var tolkienJSON []byte

func init() {
	// 1. Merriam-Webster Collegiate
	dict.Register("mw_collegiate", 100, func() ([]dict.Entry, error) {
		var entries []dict.Entry
		if err := json.Unmarshal(mwCollegiateJSON, &entries); err != nil {
			return nil, err
		}
		return entries, nil
	})

	// 2. Merriam-Webster Advanced Learner's
	dict.Register("mw_learners", 100, func() ([]dict.Entry, error) {
		var entries []dict.Entry
		if err := json.Unmarshal(mwLearnerJSON, &entries); err != nil {
			return nil, err
		}
		return entries, nil
	})

	// 3. Oxford Advanced Learner's
	dict.Register("oxford_learners", 100, func() ([]dict.Entry, error) {
		var entries []dict.Entry
		if err := json.Unmarshal(oxfordJSON, &entries); err != nil {
			return nil, err
		}
		return entries, nil
	})

	// 4. Tolkien Gateway
	dict.Register("tolkien_gateway", 100, func() ([]dict.Entry, error) {
		var entries []dict.Entry
		if err := json.Unmarshal(tolkienJSON, &entries); err != nil {
			return nil, err
		}
		return entries, nil
	})
}