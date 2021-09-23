package mapper

import "github.com/namrahov/ms-ecourt-go/model"

func ApplicationsToLightApplications(applications *[]model.Application) []model.LightApplicationDto {
	if applications == nil {
		return nil
	}

	var lightApplicationDtos []model.LightApplicationDto
	var lightApplicationDto model.LightApplicationDto

	for _, application := range *applications {
		lightApplicationDto.Id = application.Id
		lightApplicationDto.CourtName = application.CourtName
		lightApplicationDto.JudgeName = application.JudgeName
		lightApplicationDtos = append(lightApplicationDtos, lightApplicationDto)
	}
	return lightApplicationDtos
}
