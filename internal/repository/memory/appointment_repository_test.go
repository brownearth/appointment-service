package memory

import (
	"appointment-service/internal/model"
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestMemoryAppointmentRepository tests the functionality of the in-memory appointment repository.
//
// It includes the following test cases:
//
// * Create appointment
// * List appointments
// * Delete appointment
// * Get booked appointments
//
// Note: There should be little to no business logic in the repository layer.
// Therefore, things like creating overlapping appointments for a single trainer
// or checking that appointments are 30 mins exactly are NOT tested here.
func TestMemoryAppointmentRepository(t *testing.T) {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Tests creating a new appointment in the memory storage:
	// - The appointment should be assigned the next available ID (1 for first appointment)
	// - The appointment should be stored in the repository's internal storage
	// - No error should be returned for a valid appointment
	t.Run("Create appointment", func(t *testing.T) {
		repo := New(logger)
		ctx := context.Background()

		appointment := model.Appointment{
			TrainerId: 1,
			StartTime: time.Now().Add(1 * time.Hour),
			EndTime:   time.Now().Add(2 * time.Hour),
		}

		newAppointment, err := repo.Create(ctx, appointment)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), newAppointment.Id)
		assert.Len(t, repo.appointments, 1)
	})

	// Tests listing appointments from the memory storage:
	// - Two appointments should be successfully created and stored
	// - The List function should return both appointments for the specified trainer
	// - Each appointment should maintain all its original data (ID, times, trainer)
	// - The appointments should be returned in order of creation
	t.Run("List appointments", func(t *testing.T) {
		repo := New(logger)
		ctx := context.Background()

		start1 := time.Now().Add(1 * time.Hour)
		end1 := time.Now().Add(2 * time.Hour)
		start2 := time.Now().Add(3 * time.Hour)
		end2 := time.Now().Add(4 * time.Hour)

		appointment1 := model.Appointment{
			TrainerId: 1,
			UserId:    100,
			StartTime: start1,
			EndTime:   end1,
		}

		appointment2 := model.Appointment{
			TrainerId: 1,
			UserId:    200,
			StartTime: start2,
			EndTime:   end2,
		}

		createdApp1, _ := repo.Create(ctx, appointment1)
		createdApp2, _ := repo.Create(ctx, appointment2)

		appointments, err := repo.List(ctx, 1)
		assert.NoError(t, err)
		assert.Len(t, appointments, 2)

		// Verify first appointment
		assert.Equal(t, createdApp1.Id, appointments[0].Id)
		assert.Equal(t, int64(1), appointments[0].TrainerId)
		assert.Equal(t, int64(100), appointments[0].UserId)
		assert.Equal(t, start1, appointments[0].StartTime)
		assert.Equal(t, end1, appointments[0].EndTime)

		// Verify second appointment
		assert.Equal(t, createdApp2.Id, appointments[1].Id)
		assert.Equal(t, int64(1), appointments[1].TrainerId)
		assert.Equal(t, int64(200), appointments[1].UserId)
		assert.Equal(t, start2, appointments[1].StartTime)
		assert.Equal(t, end2, appointments[1].EndTime)
	})

	// Tests deleting an appointment from the memory storage:
	// - An appointment should be successfully created first
	// - The Delete function should remove the appointment without error
	// - The repository should be empty after deletion
	// - No error should be returned for deleting an existing appointment
	t.Run("Delete appointment", func(t *testing.T) {
		repo := New(logger)
		ctx := context.Background()

		// Note: The Id is not set by the caller, but by the repository
		appointment := model.Appointment{
			TrainerId: 1,
			UserId:    100,
			StartTime: time.Now().Add(1 * time.Hour),
			EndTime:   time.Now().Add(2 * time.Hour),
		}

		// Note: requried to get the newly created instance in order to have the ID
		createdAppointment, _ := repo.Create(ctx, appointment)
		err := repo.Delete(ctx, createdAppointment.Id)
		assert.NoError(t, err)
		assert.Len(t, repo.appointments, 0)
	})

	// Tests deleteing appointment from memory storage when there are more than 1 appointments:
	// - Create 3 appointments, adding them to the repository
	// - Delete 1 of the appointments
	// - Ensure that the repository still contains the correct number of appointments
	// - Ensure that the repository still contains the 2 remaining appointments
	t.Run("Delete appointment when other appointments in storage", func(t *testing.T) {
		// Initialize an empty in-memory repository
		repo := New(logger)
		ctx := context.Background()

		// Create test appointments with different user IDs but same time slots
		appointment1 := model.Appointment{
			TrainerId: 1,
			UserId:    100,
			StartTime: time.Now().Add(1 * time.Hour),
			EndTime:   time.Now().Add(2 * time.Hour),
		}
		appointment2 := model.Appointment{
			TrainerId: 1,
			UserId:    200,
			StartTime: time.Now().Add(1 * time.Hour),
			EndTime:   time.Now().Add(2 * time.Hour),
		}
		appointment3 := model.Appointment{
			TrainerId: 1,
			UserId:    300,
			StartTime: time.Now().Add(1 * time.Hour),
			EndTime:   time.Now().Add(2 * time.Hour),
		}

		// Add all appointments to the repository
		createdAppointment1, _ := repo.Create(ctx, appointment1)
		createdAppointment2, _ := repo.Create(ctx, appointment2)
		createdAppointment3, _ := repo.Create(ctx, appointment3)

		// Delete the middle appointment (appointment2)
		err := repo.Delete(ctx, createdAppointment2.Id)

		assert.NoError(t, err)
		assert.Len(t, repo.appointments, 2)
		assert.Equal(t, createdAppointment1.Id, repo.appointments[0].Id)
		assert.Equal(t, createdAppointment3.Id, repo.appointments[1].Id)
	})
}
