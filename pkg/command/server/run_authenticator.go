package server

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/UiP9AV6Y/basic-oauth2/pkg/principal"
	"github.com/UiP9AV6Y/basic-oauth2/pkg/utils"
)

func newAuthenticator(config *viper.Viper) (principal.Authenticator, error) {
	vars := newStringProcessor(config)
	login := principal.Principals{}
	principalHtp := config.GetString("principals_htpasswd")
	principalFile := config.GetString("principals_file")
	principalMap := config.GetStringMap("principals")
	principalLoader := principal.NewLoader(vars)
	principals := []principal.Principal{}

	if domain := config.GetString("email_domain"); domain != "" {
		principalLoader.EmailDomain = domain
	}

	if principalHtp != "" {
		principalHtp = utils.ResolveRelativePath(config.ConfigFileUsed(), principalHtp)
		file, err := vars.Process(principalHtp)
		if err != nil {
			return nil, fmt.Errorf("unable to load principals from %q: %w", principalHtp, err)
		} else if file == "" {
			return nil, fmt.Errorf("unable to load principals from %q", principalHtp)
		}

		parsed, err := principalLoader.ParseHtpasswd(file)
		if err != nil {
			return nil, err
		}

		principals = append(principals, parsed...)
	}

	if principalFile != "" {
		principalFile = utils.ResolveRelativePath(config.ConfigFileUsed(), principalFile)
		file, err := vars.Process(principalFile)
		if err != nil {
			return nil, fmt.Errorf("unable to load principals from %q: %w", principalFile, err)
		} else if file == "" {
			return nil, fmt.Errorf("unable to load principals from %q", principalFile)
		}

		parsed, err := principalLoader.ParseFile(file)
		if err != nil {
			return nil, err
		}

		principals = append(principals, parsed...)
	}

	if principalMap != nil {
		parsed, err := principalLoader.ParseMap(principalMap)
		if err != nil {
			return nil, fmt.Errorf("unable to load principals from config: %w", err)
		}

		principals = append(principals, parsed...)
	}

	for _, p := range principals {
		if ok := login.AddPrincipal(p); !ok {
			return nil, fmt.Errorf("principal %q already registered", p.GetIdent())
		}
	}

	return login, nil
}
