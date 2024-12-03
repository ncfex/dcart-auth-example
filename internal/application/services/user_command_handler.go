package services

import (
	"context"

	"github.com/ncfex/dcart-auth/internal/application/ports/outbound"
	userDomain "github.com/ncfex/dcart-auth/internal/domain/user"
)

type userCommandHandler struct {
	eventStore outbound.EventStore // separate read/write
}

func NewUserCommandHandler(eventStore outbound.EventStore) *userCommandHandler {
	return &userCommandHandler{
		eventStore: eventStore,
	}
}

func (h *userCommandHandler) HandleRegisterUser(ctx context.Context, cmd userDomain.RegisterUserCommand) (*userDomain.User, error) {
	user, err := userDomain.NewUser(cmd.Username, cmd.Password)
	if err != nil {
		return nil, err
	}

	if err := h.eventStore.SaveEvents(ctx, user.GetID(), user.GetUncommittedChanges()); err != nil {
		return nil, err
	}

	user.ClearUncommittedChanges()
	return user, nil
}
