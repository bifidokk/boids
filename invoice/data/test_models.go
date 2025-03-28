package data

import (
	"database/sql"
	"time"
)

type UserTest struct {
	ID        int
	Email     string
	FirstName string
	LastName  string
	Password  string
	Active    int
	IsAdmin   int
	CreatedAt time.Time
	UpdatedAt time.Time
	Plan      *Plan
}

type PlanTest struct {
	ID                  int
	PlanName            string
	PlanAmount          int
	PlanAmountFormatted string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (p *PlanTest) GetAll() ([]*Plan, error) {
	var plans []*Plan

	plan := Plan{
		ID:                  1,
		PlanName:            "My plan",
		PlanAmount:          10,
		PlanAmountFormatted: "10",
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	plans = append(plans, &plan)

	return plans, nil
}

func (p *PlanTest) GetOne(id int) (*Plan, error) {
	plan := Plan{
		ID:                  1,
		PlanName:            "My plan",
		PlanAmount:          10,
		PlanAmountFormatted: "10",
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	return &plan, nil
}

func (p *PlanTest) SubscribeUserToPlan(user User, plan Plan) error {
	return nil
}

func (p *PlanTest) AmountForDisplay() string {
	return "10"
}

func (u *UserTest) GetAll() ([]*User, error) {
	var users []*User

	user := User{
		ID:        1,
		Email:     "test@test.com",
		FirstName: "test",
		LastName:  "test",
		Password:  "test",
		Active:    1,
		IsAdmin:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	users = append(users, &user)

	return users, nil
}

func (u *UserTest) GetByEmail(email string) (*User, error) {
	user := User{
		ID:        1,
		Email:     "test@test.com",
		FirstName: "test",
		LastName:  "test",
		Password:  "test",
		Active:    1,
		IsAdmin:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return &user, nil
}

func (u *UserTest) GetOne(id int) (*User, error) {
	return u.GetByEmail("")
}

func (u *UserTest) Update() error {
	return nil
}

func (u *UserTest) Delete() error {
	return nil
}

func (u *UserTest) DeleteByID(id int) error {
	return nil
}

func (u *UserTest) Insert(user User) (int, error) {
	return 2, nil
}

func (u *UserTest) ResetPassword(password string) error {
	return nil
}

func (u *UserTest) PasswordMatches(plainText string) (bool, error) {
	return true, nil
}

func TestNew(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		&UserTest{},
		&PlanTest{},
	}
}
