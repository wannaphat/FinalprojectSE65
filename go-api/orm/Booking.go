package orm
import (
"time"
"gorm.io/gorm" // framwork ต่อกับ database ภาษา GO
)
type Booking struct { // สร้าง ตารางใน database ชื่อ User
gorm.Model
UserID string
CarID string
Start time.Time
End time.Time
}