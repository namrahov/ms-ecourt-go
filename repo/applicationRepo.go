package repo

import (
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

	query := `SELECT id, court_name FROM application limit $1 offset $2`

	rows, err := Conn.Query(query, count, offset)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var application model.ApplicationResponse
		err := rows.Scan(&application.Id, &application.CourtName)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		var comments []*model.CommentDto
		commentsRows, err2 := Conn.Query("SELECT description FROM comment where application_id in ($1)", application.Id)
		if err2 != nil {
			log.Println(err2)
			return nil, err2
		}
		for commentsRows.Next() {
			var comment model.CommentDto
			err2 := commentsRows.Scan(&comment.Description)
			if err2 != nil {
				log.Println(err2)
				return nil, err2
			}
			comments = append(comments, &comment)
		}

		application.Comments = comments
		applications = append(applications, &application)
	}

	defer Conn.Close()

	return applications, err
}
