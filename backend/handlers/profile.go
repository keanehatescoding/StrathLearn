package handlers

import (
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"strathlearn/backend/auth"
	"strathlearn/backend/db"
)

type DailySubmission struct {
	Date        string `json:"date"`
	Count       int    `json:"count"`
	ChallengeID string `json:"challengeId,omitempty"`
}

// Enhanced GetUserSubmissions to include streak information
func (h *APIHandler) GetUserSubmissions(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.GetUserFromContext(r)
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	startDateStr := r.URL.Query().Get("startDate")
	endDateStr := r.URL.Query().Get("endDate")

	var startDate, endDate time.Time
	var err error

	if startDateStr == "" {
		startDate = time.Now().AddDate(-1, 0, 0) // Default to 1 year ago
	} else {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			http.Error(w, "Invalid start date format", http.StatusBadRequest)
			return
		}
	}

	if endDateStr == "" {
		endDate = time.Now()
	} else {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			http.Error(w, "Invalid end date format", http.StatusBadRequest)
			return
		}
	}

	var submissions []struct {
		Date        time.Time
		Count       int
		ChallengeID string
	}

	result := db.DB.Model(&db.Submission{}).
		Select("DATE(created_at) as date, COUNT(*) as count, challenge_id").
		Where("user_id = ? AND created_at BETWEEN ? AND ?", user.ID, startDate, endDate).
		Group("DATE(created_at), challenge_id").
		Order("date ASC").
		Find(&submissions)

	if result.Error != nil {
		http.Error(w, "Failed to fetch submissions: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Calculate streaks
	var profile db.Profile
	db.DB.Where("user_id = ?", user.ID).First(&profile)

	activeStreak, longestStreak := calculateStreaks(submissions)

	dailySubmissions := make([]DailySubmission, 0, len(submissions))
	for _, sub := range submissions {
		dailySubmissions = append(dailySubmissions, DailySubmission{
			Date:        sub.Date.Format("2006-01-02"),
			Count:       sub.Count,
			ChallengeID: sub.ChallengeID,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"submissions":   dailySubmissions,
		"activeStreak":  activeStreak,
		"longestStreak": longestStreak,
	})
}

func calculateStreaks(submissions []struct {
	Date        time.Time
	Count       int
	ChallengeID string
}) (activeStreak int, longestStreak int) {
	if len(submissions) == 0 {
		return 0, 0
	}

	sort.Slice(submissions, func(i, j int) bool {
		return submissions[i].Date.Before(submissions[j].Date)
	})

	currentStreak := 1
	maxStreak := 1
	lastDate := submissions[0].Date

	for i := 1; i < len(submissions); i++ {
		dayDiff := submissions[i].Date.Sub(lastDate).Hours() / 24
		if dayDiff == 1 {

			currentStreak++
		} else if dayDiff > 1 {

			maxStreak = max(maxStreak, currentStreak)
			currentStreak = 1
		}
		lastDate = submissions[i].Date
	}

	now := time.Now()
	daysSinceLastSubmission := now.Sub(submissions[len(submissions)-1].Date).Hours() / 24

	if daysSinceLastSubmission > 1 {

		maxStreak = max(maxStreak, currentStreak)
		currentStreak = 0
	}

	return currentStreak, maxStreak
}
func (h *APIHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.GetUserFromContext(r)
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	var stats struct {
		TotalSubmissions int64
		ChallengeSolved  int64
		FirstSubmission  time.Time
		LastSubmission   time.Time
	}

	db.DB.Model(&db.Submission{}).
		Where("user_id = ?", user.ID).
		Count(&stats.TotalSubmissions)

	db.DB.Model(&db.Submission{}).
		Where("user_id = ? AND status_code = 3", user.ID).
		Distinct("challenge_id").
		Count(&stats.ChallengeSolved)

	var firstSub, lastSub db.Submission
	db.DB.Where("user_id = ?", user.ID).Order("created_at ASC").First(&firstSub)
	db.DB.Where("user_id = ?", user.ID).Order("created_at DESC").First(&lastSub)

	stats.FirstSubmission = firstSub.CreatedAt
	stats.LastSubmission = lastSub.CreatedAt

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user": map[string]interface{}{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
		"stats": map[string]interface{}{
			"totalSubmissions": stats.TotalSubmissions,
			"challengesSolved": stats.ChallengeSolved,
			"firstSubmission":  stats.FirstSubmission,
			"lastSubmission":   stats.LastSubmission,
		},
	})
}
