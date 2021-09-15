package repo

import (
	"database/sql"
	"github.com/namrahov/ms-ecourt-go/model"
	log "github.com/sirupsen/logrus"
)

var conn *sql.DB

type IApplicationRepo interface {
	GetApplications(offset int, count int) ([]*model.ApplicationResponse, error)
}

type ApplicationRepo struct{}

func (r ApplicationRepo) GetApplications(offset int, count int) ([]*model.ApplicationResponse, error) {
	var applications []*model.ApplicationResponse

	/*	err := Db.Model(&applications).
		Limit(count).
		Offset(offset).
		Select()*/

	query := `
      SELECT a.id, a.court_name, c.id, c.description
      FROM application a
      JOIN comment c on a.id = c.application.id
      where a.id = $1
    `

	rows, err := conn.Query(query, 4)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	application := &model.ApplicationResponse{}

	for rows.Next() {
		comment := &model.CommentDto{}
		err = rows.Scan(
			&application.Id,
			&application.CourtName,
			&comment.Id,
			&comment.Description,
		)

		application.Comments = append(application.Comments, *comment)
	}

	applications = append(applications, application)

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
