package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"myapp/database"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Connect to the database
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Successfully connected to the database")

	// Обработчик для статических файлов из папки frontend
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./frontend"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./backend/images"))))

	// Создаем новый экземпляр маршрутизатора mux
	router := mux.NewRouter()

	// Настройка маршрутов для API-запросов
	router.HandleFunc("/category", func(w http.ResponseWriter, r *http.Request) {
		allCategory(w, r, db)
	})

	router.HandleFunc("/news", func(w http.ResponseWriter, r *http.Request) {
		allTasks(w, r, db)
	})
	router.HandleFunc("/news-page", func(w http.ResponseWriter, r *http.Request) {
		// Загрузка и отправка HTML шаблона для отображения новостей в виде карточек
		newsHTML, err := ioutil.ReadFile("../frontend/news.html")
		if err != nil {
			http.Error(w, "Failed to read news.html", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write(newsHTML)
	})

	router.HandleFunc("/newsByCategory/{categoryID}", func(w http.ResponseWriter, r *http.Request) {
		tasksByCategory(w, r, db)
	}).Methods("GET")

	// Настройка обработки статических файлов из папки frontend
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Загружаем шаблон
		tmpl, err := template.ParseFiles("../frontend/startPage.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Выводим шаблон в ответ
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Start the HTTP server
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func allCategory(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	enableCors(&w)
	// Выполняем запрос к базе данных для получения всех категорий
	rows, err := db.Query("SELECT name FROM categories")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Формируем слайс категорий
	var categories []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		categories = append(categories, name)
	}

	// Преобразуем слайс категорий в JSON и отправляем в ответ на запрос
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(categories); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Task структура представляет данные о задаче
type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	CategoryID  int    `json:"category_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func allTasks(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	enableCors(&w)
	// Выполняем запрос к базе данных для получения всех задач
	rows, err := db.Query("SELECT * FROM tasks")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Формируем слайс структур для хранения данных о задачах
	var tasks []Task

	// Итерируем по строкам результата запроса
	for rows.Next() {
		// Создаем переменную для хранения данных об отдельной задаче
		var task Task

		// Извлекаем значения всех столбцов из строки результата запроса
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Image, &task.CategoryID, &task.CreatedAt, &task.UpdatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Формируем абсолютный URL-адрес для изображения
		imageURL := fmt.Sprintf("http://D:/praktika/backend/images/%s", task.Image)
		task.Image = imageURL

		// Добавляем структуру с данными об отдельной задаче в слайс задач
		tasks = append(tasks, task)
	}

	// Преобразуем слайс структур в JSON и отправляем в ответ на запрос
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func tasksByCategory(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	enableCors(&w)
	// Извлекаем значение переменной маршрута "categoryID"
	vars := mux.Vars(r)
	categoryID := vars["categoryID"]

	// Выполняем запрос к базе данных для получения задач по указанной категории
	rows, err := db.Query("SELECT * FROM tasks WHERE category_id = ?", categoryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Формируем слайс структур для хранения данных о задачах
	var tasks []Task

	// Итерируем по строкам результата запроса
	for rows.Next() {
		var task Task

		// Извлекаем значения всех столбцов из строки результата запроса
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Image, &task.CategoryID, &task.CreatedAt, &task.UpdatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Добавляем структуру с данными об отдельной задаче в слайс задач
		tasks = append(tasks, task)
	}

	// Преобразуем слайс структур в JSON и отправляем в ответ на запрос
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
