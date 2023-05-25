package entity

import (
	"gorm.io/gorm"
	"time"
)

type BanchoBeatmapSetEntity struct {
	BeatmapsetId   int     `json:"beatmapset_id" gorm:"column:BEATMAPSET_ID"`
	Artist         *string `json:"artist" gorm:"column:ARTIST"`
	ArtistUnicode  *string `json:"artist_unicode" gorm:"column:ARTIST_UNICODE"`
	Creator        *string `json:"creator" gorm:"column:CREATOR"`
	FavouriteCount *int    `json:"favourite_count" gorm:"column:FAVOURITE_COUNT"`
	HypeCurrent    *int    `json:"-" gorm:"column:HYPE_CURRENT"`
	HypeRequired   *int    `json:"-" gorm:"column:HYPE_REQUIRED"`
	Hype           struct {
		Current  *int `json:"current"`
		Required *int `json:"required"`
	} `json:"hype" gorm:"-"`
	Nsfw                         *bool   `json:"nsfw" gorm:"column:NSFW"`
	PlayCount                    *int    `json:"play_count" gorm:"column:PLAY_COUNT"`
	Source                       *string `json:"source" gorm:"column:SOURCE"`
	Status                       *string `json:"status" gorm:"column:STATUS"`
	Title                        *string `json:"title" gorm:"column:TITLE"`
	TitleUnicode                 *string `json:"title_unicode" gorm:"column:TITLE_UNICODE"`
	UserId                       *int    `json:"user_id" gorm:"column:USER_ID"`
	Video                        *bool   `json:"video" gorm:"column:VIDEO"`
	AvailabilityDownloadDisabled *bool   `json:"-" gorm:"column:AVAILABILITY_DOWNLOAD_DISABLED"`
	AvailabilityMoreInformation  *string `json:"-" gorm:"column:AVAILABILITY_MORE_INFORMATION"`
	Availability                 struct {
		DownloadDisabled *bool   `json:"download_disabled"`
		MoreInformation  *string `json:"more_information"`
	} `json:"availability" gorm:"-"`
	Bpm               *float64 `json:"bpm" gorm:"column:BPM"`
	CanBeHyped        *bool    `json:"can_be_hyped" gorm:"column:CAN_BE_HYPED"`
	DiscussionEnabled *bool    `json:"-" gorm:"column:DISCUSSION_ENABLED"`
	DiscussionLocked  *bool    `json:"-" gorm:"column:DISCUSSION_LOCKED"`
	Discussion        struct {
		Enabled *bool `json:"enabled"`
		Locked  *bool `json:"locked"`
	} `json:"discussion" gorm:"-"`
	IsScoreable                *bool      `json:"is_scoreable" gorm:"column:IS_SCOREABLE"`
	LastUpdated                *time.Time `json:"last_updated" gorm:"column:LAST_UPDATED"`
	LegacyThreadUrl            *string    `json:"legacy_thread_url" gorm:"column:LEGACY_THREAD_URL"`
	NominationsSummaryCurrent  *int       `json:"-" gorm:"column:NOMINATIONS_SUMMARY_CURRENT"`
	NominationsSummaryRequired *int       `json:"-" gorm:"column:NOMINATIONS_SUMMARY_REQUIRED"`
	Nominations                struct {
		SummaryCurrent  *int `json:"summary_current"`
		SummaryRequired *int `json:"summary_required"`
	} `json:"nominations" gorm:"-"`
	Ranked        *int                  `json:"ranked" gorm:"column:RANKED"`
	RankedDate    *time.Time            `json:"ranked_date" gorm:"column:RANKED_DATE"`
	Storyboard    *bool                 `json:"storyboard" gorm:"column:STORYBOARD"`
	SubmittedDate *string               `json:"submitted_date" gorm:"column:SUBMITTED_DATE"`
	Tags          *string               `json:"tags" gorm:"column:TAGS"`
	HasFavourited *bool                 `json:"has_favourited" gorm:"column:HAS_FAVOURITED"`
	Beatmaps      []BanchoBeatmapEntity `json:"beatmaps" gorm:"foreignKey:BEATMAPSET_ID;references:BEATMAPSET_ID"`

	Description *string `json:"description" gorm:"column:DESCRIPTION"`
	GenreId     *string `json:"-" gorm:"column:GENRE_ID"`
	GenreName   *string `json:"-" gorm:"column:GENRE_NAME"`
	Genre       struct {
		Id   *string `json:"id"`
		Name *string `json:"name"`
	} `json:"genre" gorm:"-"`
	LanguageId   *string `json:"-" gorm:"column:LANGUAGE_ID"`
	LanguageName *string `json:"-" gorm:"column:LANGUAGE_NAME"`
	Language     struct {
		Id   *string `json:"id"`
		Name *string `json:"name"`
	} `json:"language" gorm:"-"`
	Ratings *string `json:"ratings" gorm:"column:RATINGS"`
}

func (v *BanchoBeatmapSetEntity) AfterFind(tx *gorm.DB) (err error) {
	v.Hype.Required = v.HypeRequired
	v.Hype.Current = v.HypeCurrent
	v.Availability.DownloadDisabled = v.AvailabilityDownloadDisabled
	v.Availability.MoreInformation = v.AvailabilityMoreInformation
	v.Discussion.Enabled = v.DiscussionEnabled
	v.Discussion.Locked = v.DiscussionLocked
	v.Nominations.SummaryCurrent = v.NominationsSummaryCurrent
	v.Nominations.SummaryRequired = v.NominationsSummaryRequired
	v.Genre.Id = v.GenreId
	v.Genre.Name = v.GenreName
	v.Language.Id = v.LanguageId
	v.Language.Name = v.LanguageName
	return
}
func (BanchoBeatmapSetEntity) TableName() string {
	return "BEATMAPSET"
}
