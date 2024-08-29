package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/faraz-wq/portfolio-service/models"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

var db *sql.DB

// Initialize handlers with database connection
func Init(dbInstance *sql.DB) {
	db = dbInstance
}

// GetProjects handles GET /projects
func GetProjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query("SELECT id, title, description, image, tag, giturl, previewurl FROM projects")
	if err != nil {
		http.Error(w, "Unable to retrieve projects", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var project models.Project
		var tags pq.StringArray
		if err := rows.Scan(&project.ID, &project.Title, &project.Description, &project.Image, &tags, &project.GitURL, &project.PreviewURL); err != nil {
			http.Error(w, "Unable to scan project", http.StatusInternalServerError)
			return
		}
		project.Tag = tags
		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error occurred while processing rows", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(projects)
}

// GetProject handles GET /projects/{id}
func GetProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	var project models.Project
	var tags pq.StringArray
	row := db.QueryRow("SELECT id, title, description, image, tag, giturl, previewurl FROM projects WHERE id=$1", id)
	err = row.Scan(&project.ID, &project.Title, &project.Description, &project.Image, &tags, &project.GitURL, &project.PreviewURL)
	if err == sql.ErrNoRows {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Unable to retrieve project", http.StatusInternalServerError)
		return
	}
	project.Tag = tags

	json.NewEncoder(w).Encode(project)
}

// CreateProject handles POST /projects
func CreateProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var project models.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Insert the project into the database
	query := `INSERT INTO projects (title, description, image, tag, giturl, previewurl) 
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	var id int
	err := db.QueryRow(query, project.Title, project.Description, project.Image, pq.Array(project.Tag), project.GitURL, project.PreviewURL).Scan(&id)
	if err != nil {
		http.Error(w, "Unable to create project", http.StatusInternalServerError)
		return
	}

	// Return the created project with its ID
	project.ID = id
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(project)
}

// DeleteProject handles DELETE /projects/{id}
func DeleteProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	// Delete the project from the database
	query := `DELETE FROM projects WHERE id = $1`
	res, err := db.Exec(query, id)
	if err != nil {
		http.Error(w, "Unable to delete project", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, "Unable to check affected rows", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent) // No content to return
}
