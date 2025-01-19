package search

import (
	"suparna/internal/database"
)

// SearchFiles searches files by keyword and filters
func SearchFiles(keyword string) ([]map[string]interface{}, error) {
	query := `
	SELECT name, path, size, modified_time 
	FROM files 
	WHERE name LIKE ? 
	ORDER BY modified_time DESC;
	`
	rows, err := database.GetDB().Query(query, "%"+keyword+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var name, path string
		var size int64
		var modifiedTime string
		if err := rows.Scan(&name, &path, &size, &modifiedTime); err != nil {
			return nil, err
		}
		results = append(results, map[string]interface{}{
			"name":          name,
			"path":          path,
			"size":          size,
			"modified_time": modifiedTime,
		})
	}
	return results, nil
}
