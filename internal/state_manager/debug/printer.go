package debug

import (
	"fmt"
	"telegram-bot/internal/state_manager/state"
)

// PrintState выводит состояние в консоль
func PrintState(state *state.UserState) {
	fmt.Println(DebugState(state))
}
