

package main

import (
	"time"
)



type GRepo struct {
	ID	int     `gorm:"primary_key;index;AUTO_INCREMENT"`
	Path	string	`gorm:"not null;unque_index"`
	Owner	string  `gorm:"not null;index"`
	Name	string  `gorm:"not null"`
	Starred	int
	Forked	int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GUser struct {
	ID	int     `gorm:"primary_key;index;AUTO_INCREMENT"`
	Username string  `gorm:"not null;unque_index"`
	Fullname string
	Email   string  `gorm:"type:varchar(100);unique_index"`
	Link	string
	Starred	int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GOrga struct {
	ID	int     `gorm:"primary_key;index;AUTO_INCREMENT"`
	Name	string  `gorm:"not null;unque_index"`
	Fullname string
	Email   string  `gorm:"type:varchar(100);unique_index"`
	Link	string
	CreatedAt time.Time
	UpdatedAt time.Time
}

