package osu

// ===========
// TIME FORMAT: "2017-04-12T04:23:36Z"
// CheesegullBeatmapSet > LastChecked: osu.BEATMAPSET SYSTEM_UPDATE_TIMESTAMP
// ===========

type CheesegullBeatmapSet struct {
	SetId				int 				`json:"SetID"`
	ChildrenBeatmaps 	[]CheesegullBeatmap `json:"ChildrenBeatmaps"`
	RankedStatus        int          		`json:"RankedStatus"`
	ApprovedDate    	*string      		`json:"ApprovedDate"`
	LastUpdate        	*string 			`json:"LastUpdate"`
	LastChecked			*string				`json:"LastChecked"`
	Artist  			*string 			`json:"Artist"`
	Title 				*string 			`json:"Title"`
	Creator        		*string 			`json:"Creator"`
	CreatorId        	*string 			`json:"CreatorID"`
	Source       		*string 			`json:"Source"`
	Tags          		*string      		`json:"Tags"`
	HasVideo        	*bool   			`json:"HasVideo"`
	Genre				*int				`json:"Genre"`
	Language			*int				`json:"Language"`
	Favourites			*int				`json:"Favourites"`
}

type CheesegullBeatmap struct {
	BeatmapId			*int 		`json:"BeatmapID"`
	ParentSetId			*int 		`json:"ParentSetID"`
	DiffName			*string 	`json:"DiffName"`
	FileMd5				*string		`json:"FileMD5"`
	Mode          		*int     	`json:"Mode"`
	Bpm             	*float64  	`json:"BPM"`
	Ar              	*float64 	`json:"AR"`
	Od         			*float64 	`json:"OD"`
	Cs              	*float64 	`json:"CS"`
	Hp            		*float64 	`json:"HP"`
	TotalLength     	*int     	`json:"TotalLength"`
	HitLength       	*int     	`json:"HitLength"`
	Playcount       	*int     	`json:"Playcount"`
	Passcount       	*int     	`json:"Passcount"`
	MaxCombo         	*int     	`json:"MaxCombo"`
	DifficultyRating 	*float64 	`json:"DifficultyRating"`
}