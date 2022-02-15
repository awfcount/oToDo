package entity

type TodoListFolder struct {
	Entity

	Name string `json:"name" gorm:"size:128"`

	UserID string `json:"user_id" gorm:"type:char(36);"`
	User   User   `json:"-"`

	TodoLists []TodoList `json:"-"`
}
