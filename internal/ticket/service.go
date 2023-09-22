package ticket

import (
	"context"
	"errors"
	"fmt"

	"github.com/nkyizbay/ticket_store/internal/auth"
	"github.com/nkyizbay/ticket_store/internal/trip"
)

const (
	CorporatedLimit       = 20
	IndividualLimit       = 5
	LeastMaleTicketNumber = 2
)

var (
	ErrNoCapacity   = errors.New("capacity is full")
	ErrTripNotFound = errors.New("this trip does not exist")

	ErrExceedAllowedTicketToPurchase = func(limit int) error {
		return fmt.Errorf("exceed number of tickets allowed to be purchased(%d)", limit)
	}

	ErrExceedMaleTicketNumber = errors.New("exceed number of male ticket allowed to be purchased")
)

type Service interface {
	Purchase(ctx context.Context, tickets []Ticket, claims auth.Claims) error
}

type defaultService struct {
	ticketRepo Repository
	tripRepo   trip.Repository
}

func NewService(ticketRepo Repository, tripRepo trip.Repository) Service {
	return &defaultService{ticketRepo: ticketRepo, tripRepo: tripRepo}
}

func (s *defaultService) Purchase(ctx context.Context, tickets []Ticket, claims auth.Claims) error {
	if err := checkCorporatedLimit(claims, tickets); err != nil {
		return err
	}

	if err := checkIndividualLimit(claims, tickets); err != nil {
		return err
	}

	if err := checkMaleTicketLimit(tickets, claims); err != nil {
		return err
	}

	for i := range tickets {
		ticket := tickets[i]

		requestedTrip, err := s.tripRepo.FindByTripID(ctx, ticket.TripID)
		if err != nil {
			if errors.Is(err, trip.ErrTripNotFound) {
				return ErrTripNotFound
			}
			return err
		}

		purchasedTicket := Ticket{
			TripID: requestedTrip.ID,
			UserID: claims.UserID,
			Passenger: Passenger{
				Gender:   ticket.Gender,
				FullName: ticket.FullName,
				Email:    ticket.Email,
				Phone:    ticket.Phone,
			},
		}

		if ok := requestedTrip.CheckAvailableSeat(len(tickets)); !ok {
			return ErrNoCapacity
		}

		if err = s.tripRepo.UpdateAvailableSeat(ctx, requestedTrip.ID, len(tickets)); err != nil {
			return ErrNoCapacity
		}

		if err = s.ticketRepo.CreateTicketWithDetails(ctx, &purchasedTicket); err != nil {
			return err
		}
	}

	return nil
}

func checkIndividualLimit(claims auth.Claims, tickets []Ticket) error {
	if claims.IsIndividualUser() && len(tickets) > IndividualLimit {
		return ErrExceedAllowedTicketToPurchase(IndividualLimit)
	}
	return nil
}

func checkCorporatedLimit(claims auth.Claims, tickets []Ticket) error {
	if claims.IsCorporatedUser() && len(tickets) > CorporatedLimit {
		return ErrExceedAllowedTicketToPurchase(CorporatedLimit)
	}
	return nil
}

func checkMaleTicketLimit(tickets []Ticket, claims auth.Claims) error {
	var maleNum int

	for i := range tickets {
		if tickets[i].Gender == Male {
			maleNum++
		}
	}

	if claims.IsIndividualUser() && maleNum > LeastMaleTicketNumber {
		return ErrExceedMaleTicketNumber
	}

	return nil
}
