package sqlite3

import (
	"appointment-service/internal/model"
	"time"
)

type dbAppointment struct {
	ID        int64     `db:"id"`
	TrainerId int64     `db:"trainer_id"`
	UserId    int64     `db:"user_id"`
	StartTime time.Time `db:"start_time"`
	EndTime   time.Time `db:"end_time"`
}

func toDBModel(a model.Appointment) dbAppointment {
	return dbAppointment{
		ID:        a.Id,
		TrainerId: a.TrainerId,
		UserId:    a.UserId,
		StartTime: a.StartTime,
		EndTime:   a.EndTime,
	}
}

func toDomainModel(a dbAppointment) model.Appointment {
	return model.Appointment{
		Id:        a.ID,
		TrainerId: a.TrainerId,
		UserId:    a.UserId,
		StartTime: a.StartTime,
		EndTime:   a.EndTime,
	}
}

func toDomainModels(dbAppts []dbAppointment) []model.Appointment {
	appts := make([]model.Appointment, len(dbAppts))
	for i, a := range dbAppts {
		appts[i] = toDomainModel(a)
	}
	return appts
}
