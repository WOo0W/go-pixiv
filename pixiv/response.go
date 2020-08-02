package pixiv

import (
	"errors"
	"fmt"
	"net/http"
)

//
var (
	ErrEmptyNextURL = errors.New("pixiv: empty next_url field")
)

// Generated by https://quicktype.io

// RespAuth is the response from POST https://oauth.secure.pixiv.net/auth/token
type RespAuth struct {
	Response struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
		Scope        string `json:"scope"`
		RefreshToken string `json:"refresh_token"`
		User         struct {
			ProfileImageURLs struct {
				PX16X16   string `json:"px_16x16"`
				PX50X50   string `json:"px_50x50"`
				PX170X170 string `json:"px_170x170"`
			} `json:"profile_image_urls"`

			// The ID in original response is of the type string
			ID string `json:"id"`

			Name                   string `json:"name"`
			Account                string `json:"account"`
			MailAddress            string `json:"mail_address"`
			IsPremium              bool   `json:"is_premium"`
			XRestrict              int    `json:"x_restrict"`
			IsMailAuthorized       bool   `json:"is_mail_authorized"`
			RequirePolicyAgreement bool   `json:"require_policy_agreement"`
		} `json:"user"`
		DeviceToken string `json:"device_token"`
	} `json:"response"`
}

// ErrAuth is the error from POST https://oauth.secure.pixiv.net/auth/token
type ErrAuth struct {
	HasError bool `json:"has_error"`
	Errors   struct {
		System struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"system"`
	} `json:"errors"`

	response *http.Response
}

func (e *ErrAuth) Error() string {
	return fmt.Sprintf("pixiv auth: http %d: code %d: %s", e.response.StatusCode, e.Errors.System.Code, e.Errors.System.Message)
}

// ErrAppAPI is the error from app-api.pixiv.net
type ErrAppAPI struct {
	Errors struct {
		UserMessage string `json:"user_message"`
		Message     string `json:"message"`
		Reason      string `json:"reason"`
		// UserMessageDetails struct {
		// } `json:"user_message_details"`
	} `json:"error"`

	response *http.Response
}

func (e *ErrAppAPI) Error() string {
	return fmt.Sprintf("pixiv: %s %q %d: %s %s %s", e.response.Request.Method, e.response.Request.URL, e.response.StatusCode, e.Errors.Message, e.Errors.Reason, e.Errors.UserMessage)
}

// RespComments is the response from:
//  /v2/illust/comments?illust_id=...
//  /v2/novel/comments?novel_id=...
//  /v1/illust/comment/replies?comment_id=...
type RespComments struct {
	Comments []*Comment `json:"comments"`
	NextURL  string     `json:"next_url"`

	api *AppAPI
}

// NextComments fetches NextURL with API.
func (r *RespComments) NextComments() (*RespComments, error) {
	if r.NextURL == "" {
		return nil, ErrEmptyNextURL
	}
	rn := &RespComments{api: r.api}
	err := r.api.get(rn, r.NextURL, nil)
	if err != nil {
		return nil, err
	}
	return rn, nil
}

// RespNovel is the response from:
//
//  /v2/novel/detail?novel_id=...
type RespNovel struct {
	Novel Novel `json:"novel"`
}

// RespNovels is the response from:
//
//  /v1/user/novels?user_id=...
//  /v1/user/bookmarks/novel?user_id=...&restrict=...
type RespNovels struct {
	Novels  []*Novel `json:"novels"`
	NextURL string   `json:"next_url"`

	RankingNovels []*Novel `json:"ranking_novels"`

	SearchSpanLimit int `json:"search_span_limit"`

	api *AppAPI
}

// NextNovels fetches NextURL with API.
func (r *RespNovels) NextNovels() (*RespNovels, error) {
	if r.NextURL == "" {
		return nil, ErrEmptyNextURL
	}
	rn := &RespNovels{api: r.api}
	err := r.api.get(rn, r.NextURL, nil)
	if err != nil {
		return nil, err
	}
	return rn, nil
}

// RespNovelText is the response from:
//
//  /v1/novel/text?novel_id=...
type RespNovelText struct {
	NovelMarker NovelMarker `json:"novel_marker"`

	NovelText  string `json:"novel_text"`
	SeriesPrev *Novel `json:"series_prev"`
	SeriesNext *Novel `json:"series_next"`
}

// RespIllust is the response from:
//
//  /v1/illust/detail?illust_id=...
type RespIllust struct {
	Illust Illust `json:"illust"`
}

// RespIllusts is the response from:
//
//  /v2/illust/mypixiv
//  /v1/illust/new?content_type=...
//  /v1/user/illusts?user_id=...&type=...
type RespIllusts struct {
	Illusts []*Illust `json:"illusts"`
	NextURL string    `json:"next_url"`

	// For queries of recommended illusts and manga, RankingIllusts contains ranking illusts.
	RankingIllusts []*Illust `json:"ranking_illusts"`

	SearchSpanLimit int `json:"search_span_limit"`

	api *AppAPI
}

// NextIllusts fetches NextURL with API.
func (r *RespIllusts) NextIllusts() (*RespIllusts, error) {
	if r.NextURL == "" {
		return nil, ErrEmptyNextURL
	}
	rn := &RespIllusts{api: r.api}
	err := r.api.get(rn, r.NextURL, nil)
	if err != nil {
		return nil, err
	}
	return rn, nil
}

// RespUserDetail is the response from:
//
//  /v1/user/detail?user_id=...
type RespUserDetail struct {
	User    User    `json:"user"`
	Profile Profile `json:"profile"`

	// All fields here except Pawoo are all "private" or "public"
	ProfilePublicity struct {
		Gender    string `json:"gender"`
		Region    string `json:"region"`
		BirthDay  string `json:"birth_day"`
		BirthYear string `json:"birth_year"`
		Job       string `json:"job"`
		Pawoo     bool   `json:"pawoo"`
	} `json:"profile_publicity"`
	Workspace map[string]string `json:"workspace"`
}

// RespUserPreviews is the response from:
//
//  /v1/user/following?restrict=...&user_id=...
type RespUserPreviews struct {
	UserPreviews []*UserPreview `json:"user_previews"`
	NextURL      string         `json:"next_url"`

	api *AppAPI
}

// UserPreview contains last 3 illusts and novels of a user.
type UserPreview struct {
	User    User      `json:"user"`
	Illusts []*Illust `json:"illusts"`
	Novels  []*Novel  `json:"novels"`
	IsMuted bool      `json:"is_muted"`
}

// NextFollowing fetches NextURL with API.
func (r *RespUserPreviews) NextFollowing() (*RespUserPreviews, error) {
	if r.NextURL == "" {
		return nil, ErrEmptyNextURL
	}
	rn := &RespUserPreviews{api: r.api}
	err := r.api.get(rn, r.NextURL, nil)
	if err != nil {
		return nil, err
	}
	return rn, nil
}

// RespBookmarkTags is the response from:
//
//  /v1/user/bookmark-tags/illust
type RespBookmarkTags struct {
	BookmarkTags []struct {
		Count int    `json:"count"`
		Name  string `json:"name"`
	} `json:"bookmark_tags"`
	NextURL string

	api *AppAPI
}

// RespUgoiraMetadata is the response from:
//
//  /v1/ugoira/metadata?illust_id=...
type RespUgoiraMetadata struct {
	UgoiraMetadata struct {
		ZipURLs struct {
			Medium string `json:"medium"`
		} `json:"zip_urls"`
		Frames []struct {
			File  string `json:"file"`
			Delay int    `json:"delay"`
		} `json:"frames"`
	} `json:"ugoira_metadata"`
}

// RespTrendingTags is the response from:
//
//  /v1/trending-tags/illust
type RespTrendingTags struct {
	TrendTags []struct {
		Name           string `json:"tag"`
		TranslatedName string `json:"translated_name"`
		Illust         Illust `json:"illust"`
	} `json:"trend_tags"`
}

// RespTags is the response from:
//
//  /v2/search/autocomplete?word=...
type RespTags struct {
	Tags []Tag `json:"tags"`
}

// RespComment is the response from:
//
//  POST /v1/illust/comment/add
type RespComment struct {
	Comment Comment `json:"comment"`
}
