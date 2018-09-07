package console

import (
	prompt "github.com/c-bata/go-prompt"
)

var initialSuggestions = []prompt.Suggest{
	{Text: ".exit", Description: "Exit the console"},
	{Text: ".help", Description: "Print the help message"},
	{Text: "methods", Description: "List all RPC methods"},
	{Text: "login()", Description: "Login to begin an authenticated session"},
}

var commonFunc = [][]string{
	{".help", "Print the help message"},
	{".exit", "Exit the console"},
}

// SuggestionManager manages suggestions
type SuggestionManager struct {
	suggestions []prompt.Suggest
}

// NewSuggestionManager creates a suggestion manager. Initialize suggestions
// by providing initial suggestion as argument
func newSuggestionManager(initialSuggestions []prompt.Suggest) *SuggestionManager {
	sm := new(SuggestionManager)
	sm.suggestions = initialSuggestions
	return sm
}

func (sm *SuggestionManager) completer(d prompt.Document) []prompt.Suggest {
	if words := d.GetWordBeforeCursor(); len(words) > 1 {
		return prompt.FilterHasPrefix(sm.suggestions, words, true)
	}
	return nil
}
