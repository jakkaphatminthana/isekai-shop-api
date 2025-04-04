package exception

import "fmt"

type InventoryFilling struct {
	PlayerID string
	ItemID   uint64
}

func (e *InventoryFilling) Error() string {
	return fmt.Sprintf("inventory filling for player %s with item %d failed", e.PlayerID, e.ItemID)
}
