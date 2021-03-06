package format

import (
	"fmt"

	"github.com/kapitan123/telegrofler/internal/sortedmap"
)

// AK TODO this is not really apparent that format is a formattable string using two arguments
func AsDescendingList(m map[string]int, format string) string {
	listMeassge := ""

	sm := sortedmap.Sort(m)

	for _, pair := range sm {
		listMeassge += formatLine(format, pair.Key, pair.Value)
	}
	return listMeassge
}

func AsUserMention(userId int64, username string) string {
	return fmt.Sprintf("<a href=\"tg://user?id=%d\">%s</a>", userId, username)
}

func formatLine(format string, username string, score int) string {
	return fmt.Sprintf(format, username, score) + "\n"
}
