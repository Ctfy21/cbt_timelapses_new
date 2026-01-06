package database

import (
	"cbt_timelapses_backend/m/v2/configs"
	"cbt_timelapses_backend/m/v2/internal/folders"
	order "cbt_timelapses_backend/m/v2/internal/instances"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func StartClient() *DB {
	dbPath := "/home/blunder/bin/cbt_timelapses_new/cbt_timelapses_backend/orders.db"
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	
	// Проверка прав на запись
	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		log.Fatal("Database is read-only or inaccessible:", err)
	}

	// Создание таблицы для заказов
	createTableSQL := `CREATE TABLE IF NOT EXISTS orders (
		id INTEGER PRIMARY KEY,
		json_data TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	// Создание таблицы для счетчика ID
	createCounterSQL := `CREATE TABLE IF NOT EXISTS counters (
		name TEXT PRIMARY KEY,
		value INTEGER NOT NULL DEFAULT 0
	);`

	_, err = db.Exec(createCounterSQL)
	if err != nil {
		log.Fatal("Failed to create counters table:", err)
	}

	// Инициализация счетчика OrderID если его нет
	_, err = db.Exec(`INSERT OR IGNORE INTO counters (name, value) VALUES ('OrderID', 0)`)
	if err != nil {
		log.Fatal("Failed to initialize counter:", err)
	}

	// Запуск фоновой задачи для очистки старых записей
	go cleanupOldOrders(db)

	return &DB{db}
}

// Очистка заказов старше 7 дней
func cleanupOldOrders(db *sql.DB) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		cutoffTime := time.Now().Add(-configs.ORDER_TTL)
		_, err := db.Exec(`DELETE FROM orders WHERE created_at < ?`, cutoffTime)
		if err != nil {
			log.Println("Failed to cleanup old orders:", err)
		} else {
			log.Println("Cleaned up old orders")
		}
	}
}

func SetJSONData(db *DB, key string, data []byte) {
	// Извлечение ID из ключа (например "Order:123" -> 123)
	var id int
	_, err := fmt.Sscanf(key, "Order:%d", &id)
	if err != nil {
		log.Println("Failed to parse ID from key:", err)
		return
	}

	_, err = db.Exec(`INSERT OR REPLACE INTO orders (id, json_data, created_at) VALUES (?, ?, CURRENT_TIMESTAMP)`,
		id, string(data))
	if err != nil {
		log.Println("SQLite set JSON error:", err)
	}
}

func GetJSONData(db *DB, key string) []byte {
	var id int
	_, err := fmt.Sscanf(key, "Order:%d", &id)
	if err != nil {
		log.Println("Failed to parse ID from key:", err)
		return nil
	}

	var jsonData string
	err = db.QueryRow(`SELECT json_data FROM orders WHERE id = ?`, id).Scan(&jsonData)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Println("SQLite get JSON error:", err)
		}
		return nil
	}
	return []byte(jsonData)
}

func GetIncrId(db *DB, key string) int {
	tx, err := db.Begin()
	if err != nil {
		log.Println("SQLite error starting transaction:", err)
		return 0
	}
	defer tx.Rollback()

	_, err = tx.Exec(`UPDATE counters SET value = value + 1 WHERE name = ?`, key)
	if err != nil {
		log.Println("SQLite error incrementing ID:", err)
		return 0
	}

	var val int
	err = tx.QueryRow(`SELECT value FROM counters WHERE name = ?`, key).Scan(&val)
	if err != nil {
		log.Println("SQLite error getting incremented ID:", err)
		return 0
	}

	if err = tx.Commit(); err != nil {
		log.Println("SQLite error committing transaction:", err)
		return 0
	}

	log.Println("SQLite get incremented ID val:", val)
	return val
}

func GetJSONArrayValuesFromKeyPattern(db *DB, key string, needFolders bool) []byte {
	var newOrdersJson order.OrdersJSONType
	if needFolders {
		log.Println("GetJSONArrayValuesFromKeyPattern: Getting folders...")
		jsonFolders := folders.GetScreenshotsFolders()
		log.Printf("GetJSONArrayValuesFromKeyPattern: Got %d folders: %+v\n", len(jsonFolders), jsonFolders)
		newOrdersJson = order.OrdersJSONType{OrdersJSON: []string{}, Folders: jsonFolders}
	} else {
		newOrdersJson = order.OrdersJSONType{OrdersJSON: []string{}}
	}

	rows, err := db.Query(`SELECT json_data FROM orders ORDER BY id`)
	if err != nil {
		log.Println("Error during get orders operation:", err)
		return []byte("{\"Orders\":[],\"Folders\":{}}")
	}
	defer rows.Close()

	for rows.Next() {
		var jsonData string
		if err := rows.Scan(&jsonData); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		newOrdersJson.OrdersJSON = append(newOrdersJson.OrdersJSON, jsonData)
	}

	j, err := json.Marshal(newOrdersJson)
	if err != nil {
		log.Println("Error during json operation in GetJSONArrayValuesFromKeyPattern:", err)
		return []byte("{\"Orders\":[],\"Folders\":{}}")
	}

	log.Printf("GetJSONArrayValuesFromKeyPattern: Returning JSON: %s\n", string(j))
	return j
}

