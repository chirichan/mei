package entities

import "encoding/json"

type NoteFullPageRespDataEntry struct {
	Snippet    string `json:"snippet"`
	ModifyDate int64  `json:"modifyDate"`
	ColorId    int    `json:"colorId"`
	Subject    string `json:"subject"`
	AlertDate  int    `json:"alertDate"`
	Id         any    `json:"id"`
	Tag        any    `json:"tag"`
	Type       string `json:"type"`
	FolderId   any    `json:"folderId"`
	CreateDate int64  `json:"createDate"`
	Status     string `json:"status"`
	Setting    struct {
		StickyTime int `json:"stickyTime"`
		Version    int `json:"version"`
		ThemeId    int `json:"themeId,omitempty"`
		Call       struct {
			Phone string `json:"phone"`
			Time  int64  `json:"time"`
		} `json:"call,omitempty"`
		Data []struct {
			Digest   string `json:"digest"`
			MimeType string `json:"mimeType"`
			FileId   string `json:"fileId"`
		} `json:"data,omitempty"`
	} `json:"setting"`
	AlertTag   int    `json:"alertTag,omitempty"`
	ExtraInfo  string `json:"extraInfo,omitempty"`
	DeleteTime int    `json:"deleteTime,omitempty"`
}

type NoteFullPageRespData struct {
	Entries    []NoteFullPageRespDataEntry `json:"entries"`
	Folders    []interface{}               `json:"folders"`
	LastNoteId string                      `json:"lastNoteId"`
	LastPage   bool                        `json:"lastPage"`
	SyncTag    string                      `json:"syncTag"`
}

type MiNoteDetail struct {
	Entry struct {
		Snippet    string `json:"snippet"`
		ModifyDate int64  `json:"modifyDate"`
		ColorId    int    `json:"colorId"`
		Subject    string `json:"subject"`
		AlertDate  int    `json:"alertDate"`
		Type       string `json:"type"`
		FolderId   any    `json:"folderId"`
		Content    string `json:"content"`
		Setting    struct {
			Data []struct {
				Digest   string `json:"digest"`
				MimeType string `json:"mimeType"` // image/jpeg
				FileId   string `json:"fileId"`
			} `json:"data"`
			ThemeId    int `json:"themeId"`
			StickyTime int `json:"stickyTime"`
			Version    int `json:"version"`
		} `json:"setting"`
		DeleteTime int       `json:"deleteTime"`
		AlertTag   int       `json:"alertTag"`
		Id         any       `json:"id"`
		Tag        any       `json:"tag"`
		CreateDate int64     `json:"createDate"`
		Status     string    `json:"status"`
		ExtraInfo  ExtraInfo `json:"extraInfo"`
	} `json:"entry"`
}

type MiNoteResult[T any] struct {
	Result      string `json:"result"`
	Retriable   bool   `json:"retriable"`
	Code        int    `json:"code"`
	Data        T      `json:"data"`
	Description string `json:"description"`
	Ts          int64  `json:"ts"`
}

func (r MiNoteResult[T]) Success() bool {
	return r.Code == 0 && r.Result == "ok"
}

type ExtraInfo struct {
	NoteContentType      string `json:"note_content_type"`
	MindContentPlainText string `json:"mind_content_plain_text"`
	Title                string `json:"title"`
	MindContent          string `json:"mind_content"`
}

func (e *ExtraInfo) UnmarshalJSON(data []byte) error {
	var rawString string
	if err := json.Unmarshal(data, &rawString); err == nil {
		if rawString == "" || rawString == "{}" {
			*e = ExtraInfo{}
			return nil
		}
		return json.Unmarshal([]byte(rawString), e)
	}

	type Alias ExtraInfo // 避免递归调用
	var alias Alias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*e = ExtraInfo(alias)
	return nil
}

type MiNoteCSV struct {
	Title      string `json:"title" csv:"title"`
	Content    string `json:"content" csv:"content"`
	Snippet    string `json:"snippet" csv:"snippet"`
	Id         string `json:"id" csv:"id"`
	Tag        string `json:"tag" csv:"tag"`
	CreateDate int64  `json:"createDate" csv:"createDate"`
	ModifyDate int64  `json:"modifyDate" csv:"modifyDate"`
	Type       string `json:"type" csv:"type"` // note
	Data       []struct {
		Digest   string `json:"digest"`
		MimeType string `json:"mimeType"` // image/jpeg
		FileId   string `json:"fileId"`
	} `json:"data" csv:"-"`
	FolderId string `json:"folderId"`
}
