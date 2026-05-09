// Package main seeds the database with initial data.
// Layer 5: Entry point — depends on all internal packages.
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/example/dormitory-management/pkg/config"
	"github.com/example/dormitory-management/pkg/database"
	"github.com/example/dormitory-management/pkg/logger"
)

func main() {
	// Load .env file if present
	godotenv.Load()

	// Initialize logger first (before any logging)
	zapLogger, err := logger.NewProduction()
	if err != nil {
		// No logger yet — use raw stderr
		os.Stderr.WriteString("Failed to initialize logger: " + err.Error() + "\n")
		os.Exit(1)
	}
	defer zapLogger.Sync()

	// Load configuration
	cfg, err := config.Load(".")
	if err != nil {
		zapLogger.Fatal("Failed to load configuration", logger.Error(err))
	}

	// Connect to database
	ctx := context.Background()
	db, err := database.NewPostgres(ctx, cfg.DatabaseURL)
	if err != nil {
		zapLogger.Fatal("Failed to connect to database", logger.Error(err))
	}
	defer db.Close(ctx)

	zapLogger.Info("Seeding database...")

	// Seed buildings: 1 male, 1 female building
	maleBuildingID := uuid.New().String()
	femaleBuildingID := uuid.New().String()

	buildings := []struct {
		id          string
		name        string
		gender      string
		floorCount  int
		roomPerFloor int
		description string
	}{
		{maleBuildingID, "男生宿舍楼1号楼", "male", 4, 10, "男生宿舍楼，共4层，每层10间房"},
		{femaleBuildingID, "女生宿舍楼1号楼", "female", 4, 10, "女生宿舍楼，共4层，每层10间房"},
	}

	for _, b := range buildings {
		err := db.Exec(ctx, `
			INSERT INTO buildings (id, name, gender, floor_count, room_per_floor, description, status, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, 'active', NOW(), NOW())
			ON CONFLICT (id) DO NOTHING
		`, b.id, b.name, b.gender, b.floorCount, b.roomPerFloor, b.description)
		if err != nil {
			zapLogger.Error("Failed to seed building", logger.String("name", b.name), logger.Error(err))
		}
	}
	zapLogger.Info("Buildings seeded", logger.Int("count", len(buildings)))

	// Seed rooms: 4 floors * 10 rooms per floor for each building
	type roomSeed struct {
		id         string
		buildingID string
		floor      int
		number     string
		capacity   int
		status     string
	}

	var allRooms []roomSeed

	for _, b := range buildings {
		for floor := 1; floor <= b.floorCount; floor++ {
			for roomNum := 1; roomNum <= b.roomPerFloor; roomNum++ {
				roomNumber := fmt.Sprintf("%d%02d", floor, roomNum) // e.g., 101, 102, ...
				roomID := uuid.New().String()
				capacity := 4 // default capacity
				if b.gender == "female" && roomNum%3 == 0 {
					capacity = 6 // some rooms have 6 beds
				}

			room := roomSeed{
				id:         roomID,
				buildingID: b.id,
				floor:      floor,
				number:     roomNumber,
				capacity:   capacity,
				status:     "available",
			}
			allRooms = append(allRooms, room)

			err := db.Exec(ctx, `
				INSERT INTO rooms (id, building_id, floor, number, capacity, status, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
				ON CONFLICT (building_id, floor, number) DO NOTHING
			`, room.id, room.buildingID, room.floor, room.number, room.capacity, room.status)
				if err != nil {
					zapLogger.Error("Failed to seed room", logger.String("room", roomNumber), logger.Error(err))
				}
			}
		}
	}
	zapLogger.Info("Rooms seeded", logger.Int("count", len(allRooms)))

	// Seed students: 20 students
	type studentSeed struct {
		id                string
		studentNo         string
		name              string
		gender            string
		enrollmentYear    int
		major             string
		phone             string
		emergencyContact  string
		status            string
	}

	students := []studentSeed{
		{uuid.New().String(), "2021001", "张三", "male", 2021, "计算机科学与技术", "13800001001", "13800001002", "active"},
		{uuid.New().String(), "2021002", "李四", "male", 2021, "软件工程", "13800001003", "13800001004", "active"},
		{uuid.New().String(), "2021003", "王五", "male", 2021, "信息安全", "13800001005", "13800001006", "active"},
		{uuid.New().String(), "2022001", "赵六", "male", 2022, "计算机科学与技术", "13800001007", "13800001008", "active"},
		{uuid.New().String(), "2022002", "钱七", "male", 2022, "数据科学", "13800001009", "13800001010", "active"},
		{uuid.New().String(), "2022003", "孙八", "male", 2022, "人工智能", "13800001011", "13800001012", "active"},
		{uuid.New().String(), "2023001", "周九", "male", 2023, "计算机科学与技术", "13800001013", "13800001014", "active"},
		{uuid.New().String(), "2023002", "吴十", "male", 2023, "软件工程", "13800001015", "13800001016", "active"},
		{uuid.New().String(), "2023003", "郑一", "male", 2023, "网络工程", "13800001017", "13800001018", "active"},
		{uuid.New().String(), "2023004", "冯二", "male", 2023, "信息安全", "13800001019", "13800001020", "active"},
		{uuid.New().String(), "2021004", "陈雪", "female", 2021, "计算机科学与技术", "13900001001", "13900001002", "active"},
		{uuid.New().String(), "2021005", "李梅", "female", 2021, "软件工程", "13900001003", "13900001004", "active"},
		{uuid.New().String(), "2022004", "王芳", "female", 2022, "数据科学", "13900001005", "13900001006", "active"},
		{uuid.New().String(), "2022005", "赵丽", "female", 2022, "人工智能", "13900001007", "13900001008", "active"},
		{uuid.New().String(), "2022006", "钱敏", "female", 2022, "计算机科学与技术", "13900001009", "13900001010", "active"},
		{uuid.New().String(), "2023005", "孙燕", "female", 2023, "软件工程", "13900001011", "13900001012", "active"},
		{uuid.New().String(), "2023006", "周艳", "female", 2023, "网络工程", "13900001013", "13900001014", "active"},
		{uuid.New().String(), "2023007", "吴娟", "female", 2023, "信息安全", "13900001015", "13900001016", "active"},
		{uuid.New().String(), "2023008", "郑萍", "female", 2023, "计算机科学与技术", "13900001017", "13900001018", "active"},
		{uuid.New().String(), "2023009", "冯婷", "female", 2023, "数据科学", "13900001019", "13900001020", "active"},
	}

	for _, s := range students {
		err := db.Exec(ctx, `
			INSERT INTO students (id, student_no, name, gender, enrollment_year, major, phone, emergency_contact, status, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
			ON CONFLICT (student_no) DO NOTHING
		`, s.id, s.studentNo, s.name, s.gender, s.enrollmentYear, s.major, s.phone, s.emergencyContact, s.status)
		if err != nil {
			zapLogger.Error("Failed to seed student", logger.String("student_no", s.studentNo), logger.Error(err))
		}
	}
	zapLogger.Info("Students seeded", logger.Int("count", len(students)))

	// Seed allocations: 15 allocation records
	// Get male and female rooms
	var maleRooms []roomSeed
	var femaleRooms []roomSeed
	for _, r := range allRooms {
		if r.buildingID == maleBuildingID {
			maleRooms = append(maleRooms, r)
		} else {
			femaleRooms = append(femaleRooms, r)
		}
	}

	// Assign male students to male building rooms
	maleStudents := []studentSeed{}
	femaleStudents := []studentSeed{}
	for _, s := range students {
		if s.gender == "male" {
			maleStudents = append(maleStudents, s)
		} else {
			femaleStudents = append(femaleStudents, s)
		}
	}

	// Create 15 allocations
	allocCount := 0
	usedMaleRooms := make(map[string]bool)
	usedFemaleRooms := make(map[string]bool)
	usedMaleStudents := make(map[string]bool)
	usedFemaleStudents := make(map[string]bool)

	// Allocate 8 male students to male rooms
	for i, ms := range maleStudents {
		if i >= 8 || i >= len(maleRooms) {
			break
		}
		// Find an available room
		var assignedRoom *roomSeed
		for j, mr := range maleRooms {
			if !usedMaleRooms[mr.id] && mr.capacity > allocCount%4 {
				assignedRoom = &maleRooms[j]
				usedMaleRooms[mr.id] = true
				break
			}
		}
		if assignedRoom == nil {
			continue
		}

		allocationID := uuid.New().String()
		checkInAt := time.Now().AddDate(0, -6, 0) // 6 months ago

		err := db.Exec(ctx, `
			INSERT INTO allocations (id, student_id, room_id, bed_number, status, check_in_at, created_at)
			VALUES ($1, $2, $3, $4, 'active', $5, NOW())
			ON CONFLICT DO NOTHING
		`, allocationID, ms.id, assignedRoom.id, i%4+1, checkInAt)
		if err != nil {
			zapLogger.Error("Failed to seed allocation", logger.String("student", ms.studentNo), logger.Error(err))
		} else {
			allocCount++
			usedMaleStudents[ms.id] = true
		}
	}

	// Allocate 7 female students to female rooms
	for i, fs := range femaleStudents {
		if i >= 7 || i >= len(femaleRooms) {
			break
		}
		// Find an available room
		var assignedRoom *roomSeed
		for j, fr := range femaleRooms {
			if !usedFemaleRooms[fr.id] && fr.capacity > allocCount%6 {
				assignedRoom = &femaleRooms[j]
				usedFemaleRooms[fr.id] = true
				break
			}
		}
		if assignedRoom == nil {
			continue
		}

		allocationID := uuid.New().String()
		checkInAt := time.Now().AddDate(0, -6, 0) // 6 months ago

		err := db.Exec(ctx, `
			INSERT INTO allocations (id, student_id, room_id, bed_number, status, check_in_at, created_at)
			VALUES ($1, $2, $3, $4, 'active', $5, NOW())
			ON CONFLICT DO NOTHING
		`, allocationID, fs.id, assignedRoom.id, i%6+1, checkInAt)
		if err != nil {
			zapLogger.Error("Failed to seed allocation", logger.String("student", fs.studentNo), logger.Error(err))
		} else {
			allocCount++
			usedFemaleStudents[fs.id] = true
		}
	}

	zapLogger.Info("Allocations seeded", logger.Int("count", allocCount))
	zapLogger.Info("Database seeded successfully")
}
