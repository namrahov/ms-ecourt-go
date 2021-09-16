package repo

import (
	"github.com/namrahov/ms-ecourt-go/model"
	log "github.com/sirupsen/logrus"
)

type IApplicationRepo interface {
	GetApplications(offset int, count int, applicationCriteria model.ApplicationCriteria) ([]*model.Application, error)
	GetTotalCount() (int, error)
}

type ApplicationRepo struct {
}

func (r ApplicationRepo) GetApplications(offset int, count int, applicationCriteria model.ApplicationCriteria) ([]*model.Application, error) {
	var applications []*model.Application

	query := `SELECT id, request_id, checked_id, person, customer_type, customer_name, file_path,
                     court_name, judge_name, decision_number, note, status, deadline, assignee_id, 
                     priority, assignee_name, begin_date, end_date, created_at
              FROM application
              WHERE court_name like $1 
                and judge_name like $2
                and person like $3 
                and created_at::DATE >= $4 and created_at::DATE <= $5               
              limit $6 offset $7`

	rows, err := Conn.Query(query,
		"%"+applicationCriteria.CourtName+"%", "%"+applicationCriteria.JudgeName+"%",
		"%"+applicationCriteria.Person+"%", applicationCriteria.CreateDateFrom,
		applicationCriteria.CreateDateTo, count, offset)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var application model.Application
		err := rows.Scan(&application.Id, &application.RequestId, &application.CheckedId,
			&application.Person, &application.CustomerType, &application.CustomerName,
			&application.FilePath, &application.CourtName, &application.JudgeName,
			&application.DecisionNumber, &application.Note, &application.Status,
			&application.Deadline, &application.AssigneeId, &application.Priority,
			&application.AssigneeName, &application.BeginDate, &application.EndDate,
			&application.CreatedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		var comments []*model.Comment
		commentsRows, err := Conn.Query("SELECT id, commentator, description, comment_type, created_at FROM comment where application_id in ($1)", application.Id)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		for commentsRows.Next() {
			var comment model.Comment
			err := commentsRows.Scan(&comment.Id, &comment.Commentator, &comment.Description, &comment.CommentType, &comment.CreatedAt)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			comments = append(comments, &comment)
		}

		application.Comments = comments
		applications = append(applications, &application)
	}

	return applications, err
}

func (r ApplicationRepo) GetTotalCount() (int, error) {
	var totalCount int
	query := `SELECT count(*) FROM application`
	row := Conn.QueryRow(query)
	err = row.Scan(&totalCount)
	if err != nil {
		log.Fatal(err)
	}
	return totalCount, nil
}
