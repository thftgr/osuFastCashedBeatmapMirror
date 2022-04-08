package entity

type SEARCH_CACHE_STRING_INDEX struct {
	STRING string `gorm:"uniqueIndex,size:256,column:STRING"`
	ID     uint   `gorm:"primaryKey,autoIncrement:true,column:ID"`
}
type SEARCH_CACHE_TYPE struct {
	INDEX_KEY     int `gorm:"primaryKey,column:INDEX_KEY"`
	BEATMAPSET_ID int `gorm:"primaryKey,column:BEATMAPSET_ID"`
}
