package model

import (
	"fmt"
	"net/http"
	"time"
)

type Person struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Birth     time.Time `json:"birth"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (p Person) String() string {
	return fmt.Sprintf(
		"{Id: %d, Name: %s, Birth: %s, CreatedAt: %s, UpdatedAt: %s}",
		p.Id, p.Name, p.Birth, p.CreatedAt, p.UpdatedAt,
	)
}

func (p Person) Bind(r *http.Request) error {
	return nil
}

func (p Person) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (p *Person) Update(np Person) bool {
	var isUpdated = false
	if np.Name != "" {
		p.Name = np.Name
		isUpdated = true
	}
	if !np.Birth.IsZero() {
		p.Birth = np.Birth
		isUpdated = true
	}
	if isUpdated {
		p.UpdatedAt = time.Now()
	}
	return isUpdated
}
