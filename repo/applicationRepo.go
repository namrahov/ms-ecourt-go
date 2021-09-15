package repo

import (
	"fmt"
	"github.com/namrahov/ms-ecourt-go/model"
	log "github.com/sirupsen/logrus"
)

type IApplicationRepo interface {
	GetApplications(offset int, count int) ([]*model.ApplicationResponse, error)
}

type ApplicationRepo struct {
}

func (r ApplicationRepo) GetApplications(offset int, count int) ([]*model.ApplicationResponse, error) {
	var applications []*model.ApplicationResponse

	/*	err := Db.Model(&applications).
		Limit(count).
		Offset(offset).
		Select()*/

	/*query := `SELECT a.id, a.court_name, c.description
	      FROM application a
	      JOIN comment c on a.id = c.application.id
	      where a.id = $1`

		var courtName, description string
		var id int

		row := conn.QueryRow(query, 1)

		err := row.Scan(&id, &courtName, &description)
		if err != nil {
			log.Fatal(err)
		}

		applicationResponse := &model.ApplicationResponse{}
		comment := &model.CommentDto{}

		fmt.Println("courtName=", courtName)
		if err != nil {
			log.Println(err)
		}

		applicationResponse.Comments = append(applicationResponse.Comments, *comment)
		applications = append(applications, applicationResponse)
		defer conn.Close()*/

	rows, err := Conn.Query("select id, court_name from application")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var courtName string
	var id int

	for rows.Next() {
		err := rows.Scan(&id, &courtName)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		fmt.Println("Record is", id, courtName)
	}

	if err = rows.Err(); err != nil {
		log.Fatal("Error scanning rows", err)
	}
	defer Conn.Close()

	return applications, err
}

/*
func QuestionByIdAndAnswers(id string) (*Question, []*Answer, error) {
  query := `
    SELECT q.id, q.body, a.id, a.question_id, a.body
    FROM question q
    JOIN answer a ON q.id = a.question_id
    WHERE q.id = ?
  `
  rows, err := db.Query(query, id)
  checkErr(err)

  question := &Question{}
  for rows.Next() {
    answer := &Answer{}
    err = rows.Scan(
      &question.ID,
      &question.Body,
      &answer.ID,
      &answer.QuestionID,
      &answer.Body,
    )
    checkErr(err)
    question.Answers = append(question.Answers, answer)
  }

  return question, question.Answers, nil
}
*/
