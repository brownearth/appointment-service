package dto

import "appointment-service/internal/model"

func ToAppointmentModel(r *CreateAppointmentRequest) model.Appointment {
	return model.Appointment{
		StartTime: r.StartTime,
		EndTime:   r.EndTime,
		TrainerId: r.TrainerId,
		UserId:    r.UserId,
	}
}

func ToAppointmentResponse(m *model.Appointment) AppointmentResponse {
	return AppointmentResponse{
		Id:        m.Id,
		StartTime: m.StartTime,
		EndTime:   m.EndTime,
		TrainerId: m.TrainerId,
		UserId:    m.UserId,
	}
}

// ToListAppointmentsResponse converts model appointments to a response DTO
func ToListAppointmentsResponse(appointments []model.Appointment) []AppointmentResponse {
	response := make([]AppointmentResponse, len(appointments))

	for i, apt := range appointments {
		response[i] = AppointmentResponse{
			Id:        apt.Id,
			StartTime: apt.StartTime,
			EndTime:   apt.EndTime,
			TrainerId: apt.TrainerId,
			UserId:    apt.UserId,
		}
	}

	return response
}
