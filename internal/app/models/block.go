package models

// Block is a struct that represents blocks saved in json format
// Specific block types: text, image, etc
type Block struct {
	Type string     `bson:"type" json:"type"`
	Data *BlockData `bson:"data" json:"data"`
}

// BlockData represents structure of blocks data field
type BlockData struct {
	Text  string    `bson:"text,omitempty" json:"text,omitempty"`
	Level int8      `bson:"level,omitempty" json:"level,omitempty"`
	File  *FileInfo `bson:"file,omitempty" json:"file,omitempty"`
}

// FileInfo represents basic file structure for editors "file" field
type FileInfo struct {
	URL    string `bson:"url,omitempty" json:"url,omitempty"`
	Width  int    `bson:"width,omitempty" json:"width,omitempty"`
	Height int    `bson:"height,omitempty" json:"height,omitempty"`
}
