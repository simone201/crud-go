package db

import (
	"HTTPChiSqlite/model"
	"context"
)

type PeopleRepository interface {
	FindAllPeople() ([]model.Person, error)
	FindPerson(id int) (model.Person, error)
	SavePerson(ctx context.Context, person model.Person) (int64, error)
	UpdatePerson(ctx context.Context, id int, patch model.Person) (model.Person, error)
	DeletePerson(ctx context.Context, id int) error
}

type PeopleRepositoryImpl struct{}

func (pr PeopleRepositoryImpl) FindAllPeople() ([]model.Person, error) {
	rows, err := Instance.GetDB().Query("SELECT * FROM people")
	if err != nil {
		return []model.Person{}, err
	}

	var people []model.Person
	for rows.Next() {
		var person model.Person
		if err := rows.Scan(&person.Id, &person.Name, &person.Birth,
			&person.CreatedAt, &person.UpdatedAt); err != nil {
			return people, err
		}
		people = append(people, person)
	}

	err = rows.Close()
	if err != nil {
		return people, err
	}

	return people, nil
}

func (pr PeopleRepositoryImpl) FindPerson(id int) (model.Person, error) {
	var person model.Person
	if err := Instance.GetDB().QueryRow("SELECT * FROM people WHERE id = ?", id).Scan(
		&person.Id, &person.Name, &person.Birth, &person.CreatedAt, &person.UpdatedAt); err != nil {
		return person, err
	}
	return person, nil
}

func (pr PeopleRepositoryImpl) SavePerson(ctx context.Context, person model.Person) (int64, error) {
	tx, ctx, cancelFunc, rollbackFunc, err := initBackgroundTransaction(Instance, ctx, 5)
	if err != nil {
		return -1, err
	}
	defer cancelFunc()
	defer rollbackFunc(tx)

	res, err := tx.ExecContext(ctx, "INSERT INTO people (name, birth) VALUES (?, ?)", person.Name, person.Birth)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	if err = tx.Commit(); err != nil {
		return -1, err
	}

	return id, nil
}

func (pr PeopleRepositoryImpl) UpdatePerson(ctx context.Context, id int, patch model.Person) (model.Person, error) {
	var person model.Person
	if id <= 0 {
		return person, &ParamsNotValidError{params: []string{"id"}}
	}

	person, err := pr.FindPerson(id)
	if err != nil {
		return person, err
	}

	if !person.Update(patch) {
		return person, &NoRowsAffectedError{}
	}

	tx, ctx, cancelFunc, rollbackFunc, err := initBackgroundTransaction(Instance, ctx, 5)
	if err != nil {
		return person, err
	}
	defer cancelFunc()
	defer rollbackFunc(tx)

	res, err := tx.ExecContext(ctx,
		"UPDATE people SET name = ?, birth = ?, updatedAt = ? WHERE id = ?",
		person.Name, person.Birth, person.UpdatedAt, person.Id)
	if err != nil {
		return person, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return person, err
	}
	if count == 0 {
		return person, &NoRowsAffectedError{}
	}

	if err = tx.Commit(); err != nil {
		return person, err
	}

	return person, nil
}

func (pr PeopleRepositoryImpl) DeletePerson(ctx context.Context, id int) error {
	tx, ctx, cancelFunc, rollbackFunc, err := initBackgroundTransaction(Instance, ctx, 5)
	if err != nil {
		return err
	}
	defer cancelFunc()
	defer rollbackFunc(tx)

	_, err = tx.ExecContext(ctx, "DELETE FROM people WHERE id = ?", id)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
