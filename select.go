package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:qwerty123)@/observatory")

	if err != nil {
		fmt.Println("error")
		panic(err)
	}
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		table := r.URL.Query().Get("table")
		viewSelect(w, db, table)
	})

	http.HandleFunc("/postform", func(w http.ResponseWriter, r *http.Request) {
		addSector(w, r, db)
	})

	http.HandleFunc("/edit", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		renderEditForm(w, db, id)
	})

	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		updateSector(w, r, db)
	})

	fmt.Println("Server is listening on http://localhost:8181/")
	http.ListenAndServe(":8181", nil)
}

func viewSelect(w http.ResponseWriter, db *sql.DB, table string) {
	file, err := os.Open("select.html")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() != "@tr" && scanner.Text() != "@ver" {
			fmt.Fprintf(w, scanner.Text())
		}
		if scanner.Text() == "@tr" {
			if table == "" {
			viewHeadQuery(w, db, "SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_NAME = 'sector' ORDER BY ORDINAL_POSITION")

				viewSelectQuery(w, db, "SELECT * FROM sector ORDER BY id DESC")
			} else {
				tables := splitTables(table)
				viewSelectMultiQuery(w, db, tables[0], tables[1])
			}
		}
		if scanner.Text() == "@ver" {
			viewSelectVerQuery(w, db, "SELECT VERSION() AS ver")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func viewHeadQuery(w http.ResponseWriter, db *sql.DB, sShow string) {
	type sHead struct {
		clnme string
	}
	rows, err := db.Query(sShow)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	fmt.Fprintf(w, "<tr>")
	i := 0
	for i < 9 {
		rows.Next()
		p := sHead{}
		err := rows.Scan(&p.clnme)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Fprintf(w, "<td>"+p.clnme+"</td>")
		i++
	}
	fmt.Fprintf(w, "</tr>")
}

func viewSelectQuery(w http.ResponseWriter, db *sql.DB, sSelect string) {
	type sector struct {
		id              int
		coordinates     sql.NullString
		lightIntensity  sql.NullFloat64
		foreignObjects  sql.NullInt64
		starObjects     sql.NullInt64
		unknownObjects  sql.NullInt64
		knownObjects    sql.NullInt64
		notes           sql.NullString
		dateUpdate      sql.NullString
	}
	sectors := []sector{}

	rows, err := db.Query(sSelect)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		p := sector{}
		err := rows.Scan(&p.id, &p.coordinates, &p.lightIntensity, &p.foreignObjects, &p.starObjects, &p.unknownObjects, &p.knownObjects, &p.notes, &p.dateUpdate)
		if err != nil {
			fmt.Println(err)
			continue
		}
		sectors = append(sectors, p)
	}

	for _, p := range sectors {
		fmt.Fprintf(w, "<tr>")
		fmt.Fprintf(w, "<td>%d</td>", p.id)

		if p.coordinates.Valid {
			fmt.Fprintf(w, "<td>%s</td>", p.coordinates.String)
		} else {
			fmt.Fprintf(w, "<td></td>")
		}

		if p.lightIntensity.Valid {
			fmt.Fprintf(w, "<td>%.2f</td>", p.lightIntensity.Float64)
		} else {
			fmt.Fprintf(w, "<td></td>")
		}

		if p.foreignObjects.Valid {
			fmt.Fprintf(w, "<td>%d</td>", p.foreignObjects.Int64)
		} else {
			fmt.Fprintf(w, "<td></td>")
		}

		if p.starObjects.Valid {
			fmt.Fprintf(w, "<td>%d</td>", p.starObjects.Int64)
		} else {
			fmt.Fprintf(w, "<td></td>")
		}

		if p.unknownObjects.Valid {
			fmt.Fprintf(w, "<td>%d</td>", p.unknownObjects.Int64)
		} else {
			fmt.Fprintf(w, "<td></td>")
		}

		if p.knownObjects.Valid {
			fmt.Fprintf(w, "<td>%d</td>", p.knownObjects.Int64)
		} else {
			fmt.Fprintf(w, "<td></td>")
		}

		if p.notes.Valid {
			fmt.Fprintf(w, "<td>%s</td>", p.notes.String)
		} else {
			fmt.Fprintf(w, "<td></td>")
		}

		if p.dateUpdate.Valid {
			fmt.Fprintf(w, "<td>%s</td>", p.dateUpdate.String)
		} else {
			fmt.Fprintf(w, "<td></td>")
		}

		fmt.Fprintf(w, "<td><a href=\"/edit?id=%d\">Edit</a></td>", p.id)
		fmt.Fprintf(w, "</tr>")
	}
}


// func viewSelectQuery(w http.ResponseWriter, db *sql.DB, sSelect string) {
// 	type sector struct {
// 		id              int
// 		coordinates     string
// 		lightIntensity  float64
// 		foreignObjects  int
// 		starObjects     int
// 		unknownObjects  int
// 		knownObjects    int
// 		notes           string
// 	}
// 	sectors := []sector{}

// 	rows, err := db.Query(sSelect)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		p := sector{}
// 		err := rows.Scan(&p.id, &p.coordinates, &p.lightIntensity, &p.foreignObjects, &p.starObjects, &p.unknownObjects, &p.knownObjects, &p.notes)
// 		if err != nil {
// 			fmt.Println(err)
// 			continue
// 		}
// 		sectors = append(sectors, p)
// 	}

// 	// Теперь все операции с rows завершены, поэтому можно использовать данные из sectors

// 	for _, p := range sectors {
// 		fmt.Fprintf(w, "<tr><td>"+strconv.Itoa(p.id)+"</td><td>"+p.coordinates+"</td><td>"+strconv.FormatFloat(p.lightIntensity, 'f', 2, 64)+"</td><td>"+strconv.Itoa(p.foreignObjects)+"</td><td>"+strconv.Itoa(p.starObjects)+"</td><td>"+strconv.Itoa(p.unknownObjects)+"</td><td>"+strconv.Itoa(p.knownObjects)+"</td><td>"+p.notes+"</td>")
// 		fmt.Fprintf(w, "<td><a href=\"/edit?id="+strconv.Itoa(p.id)+"\">Edit</a></td></tr>")
// 	}
// }

func viewSelectVerQuery(w http.ResponseWriter, db *sql.DB, sSelect string) {
	type sVer struct {
		ver string
	}
	rows, err := db.Query(sSelect)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		p := sVer{}
		err := rows.Scan(&p.ver)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Fprintf(w, p.ver)
	}
}

func renderEditForm(w http.ResponseWriter, db *sql.DB, id string) {
	type sector struct {
		id              int
		coordinates     string
		lightIntensity  float64
		foreignObjects  int
		starObjects     int
		unknownObjects  int
		knownObjects    int
		notes           string
	}
	var p sector
	err := db.QueryRow("SELECT id, coordinates, light_intensity, foreign_objects, star_objects, unknown_objects, known_objects, notes FROM sector WHERE id = ?", id).Scan(&p.id, &p.coordinates, &p.lightIntensity, &p.foreignObjects, &p.starObjects, &p.unknownObjects, &p.knownObjects, &p.notes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Fprintf(w, `<form action="/update" method="post">`)
	fmt.Fprintf(w, `<input type="hidden" name="id" value="%d">`, p.id)
	fmt.Fprintf(w, `Coordinates: <input type="text" name="coordinates" value="%s"><br>`, p.coordinates)
	fmt.Fprintf(w, `Light Intensity: <input type="text" name="light_intensity" value="%f"><br>`, p.lightIntensity)
	fmt.Fprintf(w, `Foreign Objects: <input type="text" name="foreign_objects" value="%d"><br>`, p.foreignObjects)
	fmt.Fprintf(w, `Star Objects: <input type="text" name="star_objects" value="%d"><br>`, p.starObjects)
	fmt.Fprintf(w, `Unknown Objects: <input type="text" name="unknown_objects" value="%d"><br>`, p.unknownObjects)
	fmt.Fprintf(w, `Known Objects: <input type="text" name="known_objects" value="%d"><br>`, p.knownObjects)
	fmt.Fprintf(w, `Notes: <input type="text" name="notes" value="%s"><br>`, p.notes)
	fmt.Fprintf(w, `<input type="submit" value="Update">`)
	fmt.Fprintf(w, `</form>`)
	fmt.Fprintf(w, `<br><a style="background-color: #33bee1; color: #fff; padding: 2px 3px; text-decoration: none" href="/">Return to main menu</a>`)
}

func splitTables(table string) []string {
	return strings.Split(table, ",")
}

func viewSelectMultiQuery(w http.ResponseWriter, db *sql.DB, table1, table2 string) {
	rows, err := db.Query(fmt.Sprintf("CALL multi_select('%s', '%s')", table1, table2))
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	fmt.Fprintf(w, "<tr>")
	for _, col := range columns {
		fmt.Fprintf(w, "<td>"+col+"</td>")
	}
	fmt.Fprintf(w, "</tr>")

	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, "<tr>")
		for _, col := range values {
			if col == nil {
				fmt.Fprintf(w, "<td></td>")
			} else {
				fmt.Fprintf(w, "<td>%s</td>", col)
			}
		}
		fmt.Fprintf(w, "</tr>")
	}
}

func updateSector(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	coordinates := r.FormValue("coordinates")
	lightIntensity, _ := strconv.ParseFloat(r.FormValue("light_intensity"), 64)
	foreignObjects, _ := strconv.Atoi(r.FormValue("foreign_objects"))
	starObjects, _ := strconv.Atoi(r.FormValue("star_objects"))
	unknownObjects, _ := strconv.Atoi(r.FormValue("unknown_objects"))
	knownObjects, _ := strconv.Atoi(r.FormValue("known_objects"))
	notes := r.FormValue("notes")

	_, err := db.Exec("UPDATE sector SET coordinates = ?, light_intensity = ?, foreign_objects = ?, star_objects = ?, unknown_objects = ?, known_objects = ?, notes = ? WHERE id = ?",
		coordinates, lightIntensity, foreignObjects, starObjects, unknownObjects, knownObjects, notes, id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func addSector(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	coordinates := r.FormValue("coordinates")
	lightIntensity, _ := strconv.ParseFloat(r.FormValue("light_intensity"), 64)
	foreignObjects, _ := strconv.Atoi(r.FormValue("foreign_objects"))
	starObjects, _ := strconv.Atoi(r.FormValue("star_objects"))
	unknownObjects, _ := strconv.Atoi(r.FormValue("unknown_objects"))
	knownObjects, _ := strconv.Atoi(r.FormValue("known_objects"))
	notes := r.FormValue("notes")

	query := "INSERT INTO sector (coordinates, light_intensity, foreign_objects, star_objects, unknown_objects, known_objects, notes) VALUES (?, ?, ?, ?, ?, ?, ?)"
    fmt.Printf("Executing query: %s\nWith parameters: %s, %f, %d, %d, %d, %d, %s\n",
        query, coordinates, lightIntensity, foreignObjects, starObjects, unknownObjects, knownObjects, notes)

	_, err := db.Exec(query, coordinates, lightIntensity, foreignObjects, starObjects, unknownObjects, knownObjects, notes)

	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}