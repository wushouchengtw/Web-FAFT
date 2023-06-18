package result

import "backend/lib/models"

type ResultMem struct {
	id   int
	data map[int]*models.Result
}

func (r *ResultMem) Save(v *models.Result) (int, error) {
	r.id += 1
	r.data[r.id] = v
	return r.id, nil
}
func (r *ResultMem) SearchTestHaus(search models.SearchParams)
