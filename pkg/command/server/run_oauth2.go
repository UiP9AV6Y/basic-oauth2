package server

import (
	cr "crypto/rand"
	"errors"
	"fmt"
	"io"
	mr "math/rand"
	"time"

	"github.com/openshift/osin"
	"github.com/spf13/viper"

	"github.com/UiP9AV6Y/basic-oauth2/pkg/client"
	"github.com/UiP9AV6Y/basic-oauth2/pkg/jwt"
	"github.com/UiP9AV6Y/basic-oauth2/pkg/log"
	"github.com/UiP9AV6Y/basic-oauth2/pkg/storage"
	"github.com/UiP9AV6Y/basic-oauth2/pkg/utils"
)

var ErrNoKey = errors.New("no key provided")

func newOauth2Secret(config *viper.Viper, logger *log.Controller) (file string, secret *jwt.Secret, err error) {
	var source io.Reader
	if file = config.GetString("key_file"); file != "" {
		logger.Info().Println("using private key", file)
		secret, err = jwt.ParseSecretFile(file)
		return
	} else if config.GetBool("key_generate") {
		if seed := config.GetInt64("key_seed"); seed > 0 {
			logger.Info().Println("generating seeded private key")
			source = mr.New(mr.NewSource(seed))
		} else {
			logger.Info().Println("generating random private key")
			source = cr.Reader
		}

		file = "generated"
		size := config.GetInt("key_size")
		secret, err = jwt.GenerateSecret(source, size)
		return
	}

	file = "nil"
	err = ErrNoKey
	return
}

func newOauth2Config(config *viper.Viper) (*osin.ServerConfig, error) {
	cfg := osin.NewServerConfig()

	cfg.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{
		osin.CODE,
		osin.TOKEN,
	}
	cfg.AllowedAccessTypes = osin.AllowedAccessType{
		osin.AUTHORIZATION_CODE,
		osin.REFRESH_TOKEN,
		osin.PASSWORD,
		osin.CLIENT_CREDENTIALS,
		osin.ASSERTION,
	}
	cfg.AuthorizationExpiration = config.GetInt32("authorization_expiration")
	cfg.AccessExpiration = config.GetInt32("access_expiration")
	cfg.ErrorStatusCode = config.GetInt("error_status_code")
	cfg.RequirePKCEForPublicClients = config.GetBool("require_pkce")
	cfg.AllowClientSecretInParams = config.GetBool("allow_params_secret")
	cfg.AllowGetAccessRequest = config.GetBool("allow_get_access_request")
	cfg.RetainTokenAfterRefresh = config.GetBool("retain_token_after_refresh")
	cfg.RedirectUriSeparator = config.GetString("redirect_uri_separator")

	if cfg.AuthorizationExpiration <= 0 {
		return nil, fmt.Errorf("invalid value %q for authorization expiration", cfg.AuthorizationExpiration)
	}

	if cfg.AccessExpiration <= 0 {
		return nil, fmt.Errorf("invalid value %q for access expiration", cfg.AccessExpiration)
	}

	if cfg.ErrorStatusCode <= 0 {
		return nil, fmt.Errorf("invalid value %q for error status code", cfg.ErrorStatusCode)
	}

	return cfg, nil
}

func newOauth2Storage(config *viper.Viper) (osin.Storage, error) {
	vars := newStringProcessor(config)
	store := storage.NewClientsMemory()
	clientFile := config.GetString("clients_file")
	clientMap := config.GetStringMap("clients")
	clientLoader := client.NewLoader(vars)
	clients := []osin.Client{}

	if clientFile != "" {
		clientFile = utils.ResolveRelativePath(config.ConfigFileUsed(), clientFile)
		file, err := vars.Process(clientFile)
		if err != nil {
			return nil, fmt.Errorf("unable to load clients from %q: %w", clientFile, err)
		} else if file == "" {
			return nil, fmt.Errorf("unable to load clients from %q", clientFile)
		}

		parsed, err := clientLoader.ParseFile(file)
		if err != nil {
			return nil, err
		}

		clients = append(clients, parsed...)
	}

	if clientMap != nil {
		parsed, err := clientLoader.ParseMap(clientMap)
		if err != nil {
			return nil, fmt.Errorf("unable to load clients from config: %w", err)
		}

		clients = append(clients, parsed...)
	}

	for _, c := range clients {
		if ok := store.AddClient(c); !ok {
			return nil, fmt.Errorf("client %q already registered", c.GetId())
		}
	}

	return store, nil
}

func newOauth2Server(config *viper.Viper, logger *log.Controller) (*osin.Server, error) {
	cfg, err := newOauth2Config(config)
	if err != nil {
		return nil, err
	}

	store, err := newOauth2Storage(config)
	if err != nil {
		return nil, err
	}

	access := &osin.AccessTokenGenDefault{}
	authorize := &osin.AuthorizeTokenGenDefault{}
	server := &osin.Server{
		Config:            cfg,
		Storage:           store,
		AuthorizeTokenGen: authorize,
		AccessTokenGen:    access,
		Now:               time.Now,
		Logger:            logger.Debug().Logger(),
	}

	return server, nil
}
