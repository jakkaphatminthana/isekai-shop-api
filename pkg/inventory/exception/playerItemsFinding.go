package exception

import "fmt"

type PlayerItemFinding struct {
	PlayerID string
}

func (e *PlayerItemFinding) Error() string {
	return fmt.Sprintf("finding player items for playerID: %s failed", e.PlayerID)
}
