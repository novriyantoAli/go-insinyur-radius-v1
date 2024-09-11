package domain

import "context"

type DashboardQuick struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	Value           int    `json:"value"`
	Icon            string `json:"icon"`
	ChangeText      string `json:"changeText"`
	ChangeDirection string `json:"changeDirection"`
}

type DashboardTopLocation struct {
	Name          string `json:"name"`
	CountingToday int    `json:"counting_today"`
	CountingWeek  int    `json:"counting_week"`
	CountingMonth int    `json:"counting_month"`
}

type DashboardUsecase interface {
	Timeline(c context.Context) (res []Timeline, err error)
	QuickCount(c context.Context) (res []DashboardQuick, err error)
	TopLocation(c context.Context) (res []DashboardTopLocation, err error)
}
