package exception

type ItemListing struct {
}

// implement form Error interface
func (e *ItemListing) Error() string {
	return "Item listing failed"
}
