package main

// Modify to include target specific words, which can be normalized to reduce results
// left side must contain a unique string like FOO, not too long, not too short. This is used internally as replacement
// words are built into a regex, don't break it, but you can use that to your advantage
// This isn't tested well!
var equivalences = map[string][]string{
	//"TESLA": {"model-3", "model-y", ...},
	// langcodes are too small and are treated seperately
}
