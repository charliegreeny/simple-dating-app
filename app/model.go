package app

type ErrorOutput struct {
	Message string `json:"error"`
}

type User struct {
	*UserOutput
	Pref         *Preference `json:"-"`
	Loc          *Location   `json:"-"`
	DistanceFrom int         `json:"distanceFrom,omitempty"`
}

type UserOutput struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	Age      int    `json:"age"`
	Password string `json:"-"`
}

type Preference struct {
	UserID      string `gorm:"column:user_id; primaryKey" json:"-"`
	Gender      string `gorm:"column:preference_gender" json:"gender"`
	MinAge      int    `gorm:"column:min_age" json:"minAge"`
	MaxAge      *int   `gorm:"column:max_age" json:"maxAge"`
	MaxDistance int    `gorm:"column:max_distance" json:"maxDistance"`
}

type Location struct {
	UserID string  `gorm:"column:user_id; primaryKey" json:"-"`
	Lat    float64 `gorm:"column:lat" json:"lat"`
	Long   float64 `gorm:"column:long" json:"long"`
}
