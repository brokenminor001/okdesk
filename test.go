package main

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

type groupe struct {
	sequential_id string
	title         string
	last_name     string
	start_at      string
	end_at        string
	name          string
	worktype      string
}

func cut(text string) string {
	i := strings.Index(text, "(")
	runes := []rune(text)
	if len(runes) >= i {
		return string(runes[i:])
	}
	return text
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "csv.html")

	})
	http.HandleFunc("/csv", func(w http.ResponseWriter, r *http.Request) {

		name := r.FormValue("select")
		name2 := r.FormValue("select2")
		csvfile, err := os.Create("/home/andrey/lalalalalal.csv")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer csvfile.Close()

		connStr := "user=bi_soyuzintegro password=30a08bc2d9d782eb host=188.42.59.228 port=10100 dbname=bi_db_soyuzintegro sslmode=disable"
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			panic(err)
		}
		rows, err := db.Query("select issues.sequential_id, issues.title, users.last_name, issue_status_times.start_at, issue_status_times.end_at, equipment_models.name, issue_work_types.name as worktype from issues left join issue_status_times on issue_status_times.issue_id = issues.id left join users on users.id = issues.assignee_id left join issue_work_types on issues.work_type_id = issue_work_types.id left join issue_equipments on issue_equipments.issue_id = issues.id left join equipments on issue_equipments.equipment_id = equipments.id left join equipment_models on equipments.equipment_model_id = equipment_models.id  where issue_status_times.end_at between $1 and $2 and issue_status_times.status_id = '44278'", name, name2)
		groupes := []groupe{}
		for rows.Next() {
			p := groupe{}
			err := rows.Scan(&p.sequential_id, &p.title, &p.last_name, &p.start_at, &p.end_at, &p.name, &p.worktype)
			if err != nil {
				fmt.Println(err)
				continue
			}

			groupes = append(groupes, p)
		}
		for _, p := range groupes {
			m := map[string]string{
				"(МФУ-А3-Лазер-ЧБ-23ppm)":     "12",
				"(МФУ-А4-Лазер-ЦВ-30ppm)":     "13",
				"(МФУ-А3-Лазер-ЦВ-46ppm)":     "14",
				"(Принтер-А3-Лазер-ЧБ-50ppm)": "15",
			}
			text := p.name

			f := cut(text)
			ooo := m[f]
			records := [][]string{

				{p.sequential_id, p.title, p.last_name, p.start_at, p.end_at, p.name, p.worktype, ooo},
			}
			b := &bytes.Buffer{}
			writer := csv.NewWriter(b)
			for _, record := range records {
				err := writer.Write(record)
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
			}

			writer.Flush()
			w.Header().Set("Content-Type", "text/csv")
			w.Header().Set("Content-Disposition", "attachment;filename=TheCSVFileName.csv")

			w.Write(b.Bytes())
		}

	})

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8181", nil)
}
