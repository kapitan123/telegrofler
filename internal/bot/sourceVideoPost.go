package bot

import "fmt"

type SourceVideoPost struct {
	Sender            string
	ChatId            int64
	OriginalMessageId int
	Url               string
	VideoData         VideoData
}
type VideoData struct {
	Id         string
	Duration   int
	Title      string
	Payload    []byte
	LikesCount int
}

func (tp *SourceVideoPost) GetCaption() string {
	return fmt.Sprintf("<b>Rofler:</b> 🔥@%s🔥\n<b>Title</b>: %s", tp.Sender, tp.VideoData.Title)
}