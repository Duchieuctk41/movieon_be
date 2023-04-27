package model

type User struct {
	BaseModel
	Name     string  `json:"name" gorm:"name;type:varchar(100)"`
	Avatar   string  `json:"avatar" gorm:"avatar"`
	Email    string  `json:"email" gorm:"email;type:varchar(250);not null;unique"`
	Password string  `json:"password" gorm:"password;type:varchar(250);not null"`
	Gender   string  `json:"gender" gorm:"column:gender"`
	Phone    string  `json:"phone" gorm:"phone;type:varchar(20)"`
	Address  string  `json:"address" gorm:"address;type:varchar(250)"`
	IsActive bool    `json:"is_active" gorm:"is_active"`
	Point    float64 `json:"point" gorm:"point"`
	RankName string  `json:"rank_name" gorm:"rank_name"`
	Role     int     `json:"role" gorm:"role"`
	IdOld    string  `json:"id_old" gorm:"id_old"`
}

func (User) TableName() string {
	return "users"
}

type UserRequest struct {
	ID       *string `json:"id"`
	Name     *string `json:"name"`
	Avatar   *string `json:"avatar"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
	Gender   *string `json:"gender"`
	Phone    *string `json:"phone"`
	Address  *string `json:"address"`
	IsActive *bool   `json:"is_active"`
}

type UserParams struct {
	BaseParam
}

type UserResponse struct {
	Data []User                 `json:"data"`
	Meta map[string]interface{} `json:"meta"`
}
