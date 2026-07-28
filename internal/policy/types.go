package policy

type Rule struct {
	From, To, Submodule string
	Strategy            Strategy
	Description         string
}

type Policy struct {
	Version int
	Rules   []Rule
}
