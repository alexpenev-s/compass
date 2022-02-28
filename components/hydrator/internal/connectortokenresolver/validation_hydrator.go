package connectortokenresolver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
	"github.com/kyma-incubator/compass/components/hydrator/internal/director"

	"github.com/kyma-incubator/compass/components/connector/pkg/oathkeeper"
	"github.com/kyma-incubator/compass/components/director/pkg/httputils"
	"github.com/kyma-incubator/compass/components/director/pkg/log"
	"github.com/pkg/errors"
)

// ValidationHydrator missing godoc
type ValidationHydrator interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type validationHydrator struct {
	directorClient director.Client
}

// NewValidationHydrator missing godoc
func NewValidationHydrator(clientProvider director.Client) ValidationHydrator {
	return &validationHydrator{
		directorClient: clientProvider,
	}
}

// ResolveConnectorTokenHeader missing godoc
func (vh *validationHydrator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var authSession oathkeeper.AuthenticationSession
	if err := json.NewDecoder(r.Body).Decode(&authSession); err != nil {
		log.C(ctx).WithError(err).Errorf("Failed to decode request body: %v", err)
		httputils.RespondWithError(ctx, w, http.StatusBadRequest, errors.Wrap(err, "failed to decode Authentication Session from body"))
		return
	}
	defer httputils.Close(ctx, r.Body)

	connectorToken := r.Header.Get(oathkeeper.ConnectorTokenHeader)
	if connectorToken == "" {
		connectorToken = r.URL.Query().Get(oathkeeper.ConnectorTokenQueryParam)
	}

	if connectorToken == "" {
		log.C(ctx).Info("Token not provided")
		respondWithAuthSession(ctx, w, authSession)
		return
	}

	log.C(ctx).Info("Trying to resolve token...")

	systemAuth, err := vh.directorClient.GetSystemAuthByToken(ctx, connectorToken)
	if err != nil {
		log.C(ctx).WithError(err).Errorf("Invalid token provided: %s", err.Error())
		respondWithAuthSession(ctx, w, authSession)
		return
	}

	if authSession.Header == nil {
		authSession.Header = map[string][]string{}
	}

	var sysAuthID string

	switch v := systemAuth.(type) {
	case graphql.AppSystemAuth:
		sysAuthID = v.ID
	case graphql.RuntimeSystemAuth:
		sysAuthID = v.ID
	case graphql.IntSysSystemAuth:
		sysAuthID = v.ID
	default:
		respondWithAuthSession(ctx, w, authSession)
		return
	}

	authSession.Header.Add(oathkeeper.ClientIdFromTokenHeader, sysAuthID)

	if _, err := vh.directorClient.InvalidateSystemAuthOneTimeToken(ctx, sysAuthID); err != nil {
		log.C(ctx).WithError(err).Errorf("Failed to invalidate token: %v", err)
		httputils.RespondWithError(ctx, w, http.StatusInternalServerError, errors.New("could not invalidate token"))
		return
	}

	log.C(ctx).Infof("Token for %s resolved successfully", sysAuthID)
	respondWithAuthSession(ctx, w, authSession)
}

func respondWithAuthSession(ctx context.Context, w http.ResponseWriter, authSession oathkeeper.AuthenticationSession) {
	httputils.RespondWithBody(ctx, w, http.StatusOK, authSession)
}
