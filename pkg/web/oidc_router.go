package web

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/openshift/osin"
	"gopkg.in/square/go-jose.v2"

	"github.com/UiP9AV6Y/basic-oauth2/pkg/jwt"
	"github.com/UiP9AV6Y/basic-oauth2/pkg/log"
	"github.com/UiP9AV6Y/basic-oauth2/pkg/principal"
	principals "github.com/UiP9AV6Y/basic-oauth2/pkg/principal"
	"github.com/UiP9AV6Y/basic-oauth2/pkg/utils"
)

var (
	DiscoEndpoint = "/.well-known/openid-configuration"
	KeysEndpoint  = "/.well-known/jwks.json"
	AuthEndpoint  = "/authorize"
	InfoEndpoint  = "/userinfo"
	TokenEndpoint = "/token"
)

type OIDCRouterOptions struct {
	PubKeys []jose.JSONWebKey
	Signer  jose.Signer
	Server  *osin.Server

	Logger *log.Controller

	Login principals.Authenticator

	IssuerUrl    string
	AuthBaseUrl  string
	InfoBaseUrl  string
	TokenBaseUrl string
	KeysBaseUrl  string
}

func NewOIDCRouterOptions(server *osin.Server, signer jose.Signer, publicKeys ...jose.JSONWebKey) *OIDCRouterOptions {
	options := &OIDCRouterOptions{
		Signer:  signer,
		Server:  server,
		PubKeys: publicKeys,
		Logger:  log.NewOffController(),
	}

	return options
}

func (o *OIDCRouterOptions) AuthorizeTypes() []string {
	types := make([]string, len(o.Server.Config.AllowedAuthorizeTypes))
	for i, t := range o.Server.Config.AllowedAuthorizeTypes {
		types[i] = string(t)
	}

	return types
}

func (o *OIDCRouterOptions) AccessTypes() []string {
	types := make([]string, len(o.Server.Config.AllowedAccessTypes))
	for i, t := range o.Server.Config.AllowedAccessTypes {
		types[i] = string(t)
	}

	return types
}

func (o *OIDCRouterOptions) DiscoveryData() map[string]interface{} {
	data := map[string]interface{}{
		"issuer":                                strings.TrimRight(o.IssuerUrl, "/"),
		"authorization_endpoint":                UrlJoin(o.AuthBaseUrl, AuthEndpoint),
		"userinfo_endpoint":                     UrlJoin(o.InfoBaseUrl, InfoEndpoint),
		"token_endpoint":                        UrlJoin(o.TokenBaseUrl, TokenEndpoint),
		"jwks_uri":                              UrlJoin(o.KeysBaseUrl, KeysEndpoint),
		"scopes_supported":                      []string{"openid", "email"},
		"subject_types_supported":               []string{"public"},
		"response_types_supported":              o.AuthorizeTypes(),
		"token_endpoint_auth_methods_supported": o.AccessTypes(),
		"id_token_signing_alg_values_supported": o.SigningAlgorithms(),
		"claims_supported": []string{
			"aud", "exp", "iat", "iss", "sub",
			"email", "email_verified",
		},
	}

	return data
}

func (o *OIDCRouterOptions) SigningAlgorithms() []string {
	algsUsed := make(map[string]bool)
	signingAlgs := make([]string, 0, len(o.PubKeys))

	for _, key := range o.PubKeys {
		if _, used := algsUsed[key.Algorithm]; !used {
			algsUsed[key.Algorithm] = true
			signingAlgs = append(signingAlgs, key.Algorithm)
		}
	}

	return signingAlgs
}

func (o *OIDCRouterOptions) Router() (*OIDCRouter, error) {
	handler := &OIDCRouter{
		login:   o.Login,
		logger:  o.Logger,
		signer:  o.Signer,
		server:  o.Server,
		pubKeys: o.PubKeys,
		disco:   o.DiscoveryData(),
	}
	return handler, nil
}

func UrlJoin(base, path string) string {
	b := strings.TrimRight(base, "/")
	p := strings.TrimLeft(path, "/")
	s := ""

	if !strings.Contains(base, "://") {
		s = "http://"
	}

	if base == "" {
		b = "127.0.0.1"
	}

	return fmt.Sprintf("%s%s/%s", s, b, p)
}

type OIDCRouter struct {
	logger  *log.Controller
	pubKeys []jose.JSONWebKey
	signer  jose.Signer
	server  *osin.Server
	disco   map[string]interface{}
	login   principals.Authenticator
}

func (h *OIDCRouter) Handler() http.Handler {
	mux := http.NewServeMux()

	h.DecorateHandler(mux)

	return mux
}

func (h *OIDCRouter) DecorateHandler(mux *http.ServeMux) {
	mux.HandleFunc(DiscoEndpoint, h.HandleDiscovery)
	mux.HandleFunc(KeysEndpoint, h.HandlePublicKeys)
	mux.HandleFunc(TokenEndpoint, h.HandleToken)
	mux.HandleFunc(InfoEndpoint, h.HandleInfo)
	mux.HandleFunc(AuthEndpoint, h.HandleAuth)
}

// handleDiscovery returns the OpenID Connect discovery object, allowing clients
// to discover OAuth2 resources.
func (h *OIDCRouter) HandleDiscovery(w http.ResponseWriter, r *http.Request) {
	resp := h.server.NewResponse()
	defer resp.Close()

	resp.Output = osin.ResponseData(h.disco)

	err := osin.OutputJSON(resp, w, r)
	if err != nil {
		h.logger.Error().Printf("%s: %v", DiscoEndpoint, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handlePublicKeys publishes the public part of this server's signing keys.
// This allows clients to verify the signature of ID Tokens.
func (h *OIDCRouter) HandlePublicKeys(w http.ResponseWriter, r *http.Request) {
	resp := h.server.NewResponse()
	defer resp.Close()

	resp.Output["keys"] = h.pubKeys

	err := osin.OutputJSON(resp, w, r)
	if err != nil {
		h.logger.Error().Printf("%s: %v", KeysEndpoint, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *OIDCRouter) HandleInfo(w http.ResponseWriter, r *http.Request) {
	resp := h.server.NewResponse()
	defer resp.Close()

	if ir := h.server.HandleInfoRequest(resp, r); ir != nil {
		h.logger.Debug().Printf("%s: received information request: %v", InfoEndpoint, ir.Code)
		h.server.FinishInfoRequest(resp, r, ir)
	}

	err := osin.OutputJSON(resp, w, r)
	if err != nil {
		h.logger.Error().Printf("%s: %v", InfoEndpoint, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *OIDCRouter) HandleAuth(w http.ResponseWriter, r *http.Request) {
	resp := h.server.NewResponse()
	defer resp.Close()

	if ar := h.server.HandleAuthorizeRequest(resp, r); ar != nil {
		if userData, ok := h.authorize(ar, r); ok {
			h.logger.Debug().Printf("%s: succeeded via %s", AuthEndpoint, ar.Client.GetId())
			ar.Authorized = true
			ar.UserData = userData
		} else {
			h.logger.Info().Printf("%s: authentication failed", AuthEndpoint)
			w.Header().Set("WWW-Authenticate",
				fmt.Sprintf("Basic realm=%q", ar.Client.GetId()))
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		h.server.FinishAuthorizeRequest(resp, r, ar)
	}

	if resp.IsError && resp.InternalError != nil {
		h.logger.Info().Printf("%s: authorization failed: %v", AuthEndpoint, resp.InternalError)
		http.Error(w, resp.InternalError.Error(), http.StatusInternalServerError)
		return
	}

	err := osin.OutputJSON(resp, w, r)
	if err != nil {
		h.logger.Error().Printf("%s: %v", AuthEndpoint, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *OIDCRouter) HandleToken(w http.ResponseWriter, r *http.Request) {
	resp := h.server.NewResponse()
	defer resp.Close()

	if ar := h.server.HandleAccessRequest(resp, r); ar != nil {
		if idToken, ok := h.access(ar, r); ok {
			h.logger.Debug().Printf("%s: access granted via %s", TokenEndpoint, ar.Client.GetId())
			ar.Authorized = true
			h.encodeToken(resp, idToken)
		} else {
			h.logger.Debug().Printf("%s: access denied via %s", TokenEndpoint, ar.Client.GetId())
		}

		h.server.FinishAccessRequest(resp, r, ar)
	}

	if resp.IsError && resp.InternalError != nil {
		h.logger.Info().Printf("%s: access failed: %v", TokenEndpoint, resp.InternalError)
		http.Error(w, resp.InternalError.Error(), http.StatusInternalServerError)
		return
	}

	err := osin.OutputJSON(resp, w, r)
	if err != nil {
		h.logger.Error().Printf("%s: %v", TokenEndpoint, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *OIDCRouter) access(ar *osin.AccessRequest, r *http.Request) (*jwt.IDToken, bool) {
	idToken, _ := ar.UserData.(*jwt.IDToken)

	return idToken, true
}

func (h *OIDCRouter) encodeToken(w *osin.Response, token *jwt.IDToken) {
	if token == nil {
		return
	}

	raw, err := token.Encode(h.signer)
	if err != nil {
		w.IsError = true
		w.ErrorId = osin.E_SERVER_ERROR
		w.InternalError = err
	} else {
		h.logger.Debug().Printf("adding token to response for %s/%s", token.ClientID, token.UserID)
		w.Output["id_token"] = raw
	}
}

func (h *OIDCRouter) authorize(ar *osin.AuthorizeRequest, r *http.Request) (*jwt.IDToken, bool) {
	principal, err := h.authenticate(ar, r)
	if err != nil {
		h.logger.Debug().Printf("authentication failed: %v", err)
		return nil, false
	}

	scopes := utils.MapFields(ar.Scope)

	if scopes["openid"] {
		issuer, _ := h.disco["issuer"]
		idToken := jwt.NewIDToken(principal.GetIdent(), time.Hour)

		idToken.Issuer = issuer.(string)
		idToken.ClientID = ar.Client.GetId()
		idToken.Nonce = r.URL.Query().Get("nonce")

		if scopes["email"] {
			e, v := principal.GetEmail()
			idToken.Email, idToken.EmailVerified = e, &v
		}

		return idToken, true
	}

	return nil, true
}

func (h *OIDCRouter) authenticate(ar *osin.AuthorizeRequest, r *http.Request) (principal.Principal, error) {
	u, p, ok := r.BasicAuth()

	if !ok {
		return nil, fmt.Errorf("unable to authenticate b/c no credentials where provided")
	}

	if h.login == nil {
		return nil, fmt.Errorf("unable to authenticate due to lack of authenticator")
	}

	return h.login.Authenticate(u, p)
}
