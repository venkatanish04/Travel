package booking

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	pb "travelraft/api/proto"
	"travelraft/internal/raft"
	"travelraft/internal/state"
	"travelraft/internal/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCService struct {
	pb.UnimplementedBookingServiceServer

	Service *DatabaseService
	DB      *sql.DB
	Raft    *raft.Node
}

// GRPCServer is kept as a compatibility alias for the older call sites in this repo.
type GRPCServer = GRPCService

func NewGRPCServer(service *DatabaseService, db *sql.DB, nodes ...*raft.Node) *GRPCService {
	server := &GRPCService{Service: service, DB: db}
	if len(nodes) > 0 {
		server.Raft = nodes[0]
	}
	return server
}

func NewGRPCService(service *DatabaseService, db *sql.DB, nodes ...*raft.Node) *GRPCService {
	return NewGRPCServer(service, db, nodes...)
}

func (s *GRPCService) Search(_ context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	if req.GetVehicleType() == "" {
		return nil, status.Error(codes.InvalidArgument, "vehicle type is required")
	}
	vehicles, err := s.Service.Search(req.GetVehicleType(), req.GetSource(), req.GetDestination())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	response := &pb.SearchResponse{}
	for _, vehicle := range vehicles {
		response.Vehicles = append(response.Vehicles, &pb.Vehicle{
			Id:            int32(vehicle.ID),
			VehicleType:   vehicle.VehicleType,
			VehicleNumber: vehicle.VehicleNumber,
			Name:          vehicle.Name,
			Source:        vehicle.Source,
			Destination:   vehicle.Destination,
			DepartureTime: vehicle.DepartureTime,
		})
	}
	return response, nil
}

func (s *GRPCService) GetSeats(_ context.Context, req *pb.GetSeatsRequest) (*pb.GetSeatsResponse, error) {
	if req.GetVehicleId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "vehicle id is required")
	}
	seats, err := s.Service.Seats(int(req.GetVehicleId()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	response := &pb.GetSeatsResponse{}
	for _, seat := range seats {
		response.Seats = append(response.Seats, &pb.Seat{
			Id:         int32(seat.ID),
			VehicleId:  int32(seat.VehicleID),
			SeatNumber: seat.SeatNumber,
			Status:     seat.Status,
		})
	}
	return response, nil
}

func (s *GRPCService) Book(_ context.Context, req *pb.BookRequest) (*pb.BookResponse, error) {
	if req.GetPassengerName() == "" {
		return nil, status.Error(codes.InvalidArgument, "passenger name is required")
	}
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.GetVehicleId() <= 0 || req.GetSeatId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "vehicle id and seat id are required")
	}

	userID, err := s.Service.RegisterUser(req.GetPassengerName(), req.GetEmail())
	if err != nil {
		return &pb.BookResponse{Message: err.Error()}, nil
	}

	if s.Raft != nil {
		if !s.Raft.IsLeader() {
			return &pb.BookResponse{Success: false, Message: "NOT_LEADER"}, nil
		}
		seats, err := s.Service.Seats(int(req.GetVehicleId()))
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		seatAvailable := false
		for _, seat := range seats {
			if seat.ID == int(req.GetSeatId()) {
				seatAvailable = seat.Status == "AVAILABLE"
				break
			}
		}
		if !seatAvailable {
			return &pb.BookResponse{Success: false, Message: "seat is already booked or does not belong to vehicle"}, nil
		}

		pnr := s.Service.NewPNR()
		command := state.BookingCommand{
			Type:          state.CommandBook,
			PNR:           pnr,
			UserID:        int(userID),
			VehicleID:     int(req.GetVehicleId()),
			SeatID:        int(req.GetSeatId()),
			PassengerName: req.GetPassengerName(),
		}
		commandString, err := state.EncodeBookingCommand(command)
		if err != nil {
			return nil, err
		}

		accepted, index, err := s.Raft.ProposeCommand(string(commandString))
		if err != nil {
			return nil, err
		}
		if !accepted {
			return &pb.BookResponse{Success: false, Message: "NOT_LEADER"}, nil
		}

		s.Raft.ReplicateNow()
		if !s.Raft.WaitForCommit(index, 5*time.Second) {
			return &pb.BookResponse{Success: false, Message: "BOOKING_COMMIT_TIMEOUT"}, nil
		}
		if !s.Raft.WaitForApply(index, 5*time.Second) {
			return &pb.BookResponse{Success: false, Message: "BOOKING_APPLY_TIMEOUT"}, nil
		}
		return &pb.BookResponse{Success: true, Pnr: pnr, Message: fmt.Sprintf("Booking confirmed. PNR: %s", pnr)}, nil
	}

	pnr, err := s.Service.Book(int(userID), int(req.GetVehicleId()), int(req.GetSeatId()), req.GetPassengerName())
	if err != nil {
		return &pb.BookResponse{Message: err.Error()}, nil
	}
	return &pb.BookResponse{Success: true, Pnr: pnr, Message: "booking confirmed"}, nil
}

func (s *GRPCService) GetBooking(_ context.Context, req *pb.GetBookingRequest) (*pb.GetBookingResponse, error) {
	if req.GetPnr() == "" {
		return nil, status.Error(codes.InvalidArgument, "pnr is required")
	}
	booking, err := s.Service.GetBooking(req.GetPnr())
	if err != nil {
		return &pb.GetBookingResponse{Message: "booking not found"}, nil
	}
	return &pb.GetBookingResponse{
		Found: true,
		Booking: &pb.Booking{
			Id:            int32(booking.ID),
			Pnr:           booking.PNR,
			UserId:        int32(booking.UserID),
			VehicleId:     int32(booking.VehicleID),
			SeatId:        int32(booking.SeatID),
			PassengerName: booking.PassengerName,
			Status:        booking.Status,
			VehicleName:   booking.VehicleName,
			VehicleNumber: booking.VehicleNumber,
			SeatNumber:    booking.SeatNumber,
		},
	}, nil
}

func (s *GRPCService) Cancel(_ context.Context, req *pb.CancelRequest) (*pb.CancelResponse, error) {
	if req.GetPnr() == "" {
		return nil, status.Error(codes.InvalidArgument, "pnr is required")
	}
	if s.Raft != nil {
		if !s.Raft.IsLeader() {
			return &pb.CancelResponse{Success: false, Message: "NOT_LEADER"}, nil
		}
		command, err := state.EncodeBookingCommand(state.BookingCommand{Type: state.CommandCancel, PNR: req.GetPnr()})
		if err != nil {
			return &pb.CancelResponse{Message: err.Error()}, nil
		}
		accepted, index, err := s.Raft.ProposeCommand(string(command))
		if err != nil {
			return &pb.CancelResponse{Message: err.Error()}, nil
		}
		if !accepted {
			return &pb.CancelResponse{Success: false, Message: "NOT_LEADER"}, nil
		}
		s.Raft.ReplicateNow()
		if !s.Raft.WaitForCommit(index, 5*time.Second) {
			return &pb.CancelResponse{Success: false, Message: "BOOKING_COMMIT_TIMEOUT"}, nil
		}
		if !s.Raft.WaitForApply(index, 5*time.Second) {
			return &pb.CancelResponse{Success: false, Message: "BOOKING_APPLY_TIMEOUT"}, nil
		}
		return &pb.CancelResponse{Success: true, Message: "booking cancelled successfully"}, nil
	}
	if err := s.Service.Cancel(req.GetPnr()); err != nil {
		return &pb.CancelResponse{Message: err.Error()}, nil
	}
	return &pb.CancelResponse{Success: true, Message: "booking cancelled successfully"}, nil
}

func GetExistingUser(db *sql.DB, email string) (*storage.User, error) {
	return storage.GetUserByEmail(db, email)
}
