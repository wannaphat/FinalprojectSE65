package booking

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"se/jwt-api/orm"
	"time"
)

// สร้าง structure เพื่อรองรับ json
type BookingBody struct {
	UserID string
	CarID  string
	Start  time.Time
	End    time.Time
}

func BookingCar(c *gin.Context) {
	var json BookingBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// การแปลงค่า
	layout := "2006-01-02"
	start, err := time.Parse(layout, json.Start.Format(layout))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start time"})
		return
	}
	if start.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start time must be in the future"})
		return
	}
	end, err := time.Parse(layout, json.End.Format(layout))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end time"})
		return
	}
	if end.Before(start) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end time must be after start time"})
		return
	}

	// Check if the start time is before the end time
	if json.Start.After(json.End) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking period"})
		return
	}
	// Query the database using Gorm
	var results []orm.Booking
	orm.Db.Where("car_id = ? AND ((start <= ? AND end > ?) OR (start >= ? AND start < ?))", json.CarID, end, start, start, end).Find(&results)
	if len(results) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Booking already exists"})
		return
	}
	// Check if the booking already exists
	if len(results) > 0 {
		c.JSON(400, gin.H{"status": "error", "message": "Booking Exists"})
		return
	}
	// Create the booking
	booking := orm.Booking{UserID: json.UserID, CarID: json.CarID, Start: start,
		End: end}
	orm.Db.Create(&booking)
	c.JSON(200, gin.H{"status": "success", "data": booking})
}
