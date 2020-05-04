package vksign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gorilla/schema"
	"net/url"
	"sort"
	"strings"
)

var decoder = schema.NewDecoder()
var ErrSignOrSecretNotValid = errors.New("sign or secret not valid")

func Parse(launchUrl string, secret string) (*Data, error) {
	var params url.Values
	if u, err := url.Parse(launchUrl); err == nil {
		if tempParams, err := url.ParseQuery(u.RawQuery); err == nil {
			params = tempParams
		} else {
			return &Data{}, err
		}
	} else {
		return &Data{}, err
	}

	return ParseWithUrlValues(params, secret)
}

func ParseWithUrlValues(params url.Values, secret string) (*Data, error) {
	if isValid(params, secret) {
		var data Data
		decoder.IgnoreUnknownKeys(true)
		err := decoder.Decode(&data, params)
		if err != nil {
			return &Data{}, err
		}

		return &data, nil
	}

	return &Data{}, ErrSignOrSecretNotValid
}

func isValid(params url.Values, secret string) bool {
	var sing string
	pair := make([]string, 0)
	for key, value := range params {
		if strings.Index(key, "vk_") == 0 {
			pair = append(pair, strings.Replace(fmt.Sprintf("%s=%s", key, value[0]), ",", "%2C", -1))
		}
		if key == "sign" {
			sing = value[0]
		}
	}
	sort.Strings(pair)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(strings.Join(pair, "&")))

	verifySign := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	verifySign = strings.Replace(verifySign, "+", "-", -1)
	verifySign = strings.Replace(verifySign, "/", "_", -1)
	verifySign = strings.TrimRight(verifySign, "=")

	return sing == verifySign
}

type Data struct {
	UserID                  uint32    `schema:"vk_user_id"`
	AppID                   int       `schema:"vk_app_id"`
	IsAppUser               bool      `schema:"vk_is_app_user"`
	AreNotificationsEnabled bool      `schema:"vk_are_notifications_enabled"`
	Language                Language  `schema:"vk_language"`
	Ref                     Ref       `schema:"vk_ref"`
	AccessTokenSettings     string    `schema:"vk_access_token_settings"`
	GroupID                 int       `schema:"vk_group_id"`
	ViewerGroupRole         GroupRole `schema:"vk_viewer_group_role"`
	Platform                Platform  `schema:"vk_platform"`
	IsFavorite              bool      `schema:"vk_is_favorite"`
}

type Language string

const (
	Ru Language = "ru"
	Uk Language = "uk"
	Be Language = "be"
	Kz Language = "kz"
	En Language = "en"
	Es Language = "es"
	Fi Language = "fi"
	De Language = "de"
	Id Language = "it"
)

type Ref string

const (
	FeaturingDiscover Ref = "featuring_discover"
	FeaturingMenu     Ref = "featuring_menu"
	FeaturingNew      Ref = "featuring_new"
	Other             Ref = "other"
)

type GroupRole string

const (
	None   GroupRole = "none"
	Member GroupRole = "member"
	Moder  GroupRole = "moder"
	Editor GroupRole = "editor"
	Admin  GroupRole = "admin"
)

type Platform string

const (
	MobileAndroid          Platform = "mobile_android"
	MobileIphone           Platform = "mobile_iphone"
	MobileWeb              Platform = "mobile_web"
	DesktopWeb             Platform = "desktop_web"
	MobileAndroidMessenger Platform = "mobile_android_messenger"
	MobileIphoneMessenger  Platform = "mobile_iphone_messenger"
)
