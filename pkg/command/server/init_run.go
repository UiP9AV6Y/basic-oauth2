package server

import (
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	AuthorizationExpiration int32 = 250
	AccessExpiration        int32 = 3600
	ErrorStatusCode         int   = 200
)

func initRunAllowedAuthorize(config *viper.Viper, flags *flag.FlagSet) {
	f := "allow-authorize"
	c := "allow_authorize"

	flags.StringSlice(f, allowedAuthorizeTypes.FlagValues(), "Allowed authorization request types")
	_ = config.BindPFlag(c, flags.Lookup(f))
	_ = config.BindEnv(c)

	for key, allowed := range allowedAuthorizeTypes {
		aF := f + "-" + key
		aC := c + "_" + string(allowed)

		flags.Bool(aF, false, "Allow authorization selectively. See --allow-authorize")
		_ = config.BindPFlag(aC, flags.Lookup(aF))
		_ = config.BindEnv(aC)
	}
}

func initRunAllowedAccess(config *viper.Viper, flags *flag.FlagSet) {
	f := "allow-access"
	c := "allow_access"

	flags.StringSlice(f, allowedAccessTypes.FlagValues(), "Allowed access request types")
	_ = config.BindPFlag(c, flags.Lookup(f))
	_ = config.BindEnv(c)

	for key, allowed := range allowedAccessTypes {
		aF := f + "-" + key
		aC := c + "_" + string(allowed)

		flags.Bool(aF, false, "Allow access request selectively. See --allow-access")
		_ = config.BindPFlag(aC, flags.Lookup(aF))
		_ = config.BindEnv(aC)
	}
}

func initRunAuthorizationExpiration(config *viper.Viper, flags *flag.FlagSet) {
	flags.Int32("authorization-expiration", AuthorizationExpiration, "Authorization token expiration in seconds")
	_ = config.BindPFlag("authorization_expiration", flags.Lookup("authorization-expiration"))
	_ = config.BindEnv("authorization_expiration")
}

func initRunAccessExpiration(config *viper.Viper, flags *flag.FlagSet) {
	flags.Int32("access-expiration", AccessExpiration, "Access token expiration in seconds")
	_ = config.BindPFlag("access_expiration", flags.Lookup("access-expiration"))
	_ = config.BindEnv("access_expiration")
}

func initRunErrorStatusCode(config *viper.Viper, flags *flag.FlagSet) {
	flags.Int("error-status-code", ErrorStatusCode, "HTTP status code to return for errors")
	_ = config.BindPFlag("error_status_code", flags.Lookup("error-status-code"))
	_ = config.BindEnv("error_status_code")
}

func initRunAllowClientSecretInParams(config *viper.Viper, flags *flag.FlagSet) {
	flags.Bool("allow-params-secret", false, "Allows client secret also in request parameters in addition to Authorization header")
	_ = config.BindPFlag("allow_params_secret", flags.Lookup("allow-params-secret"))
	_ = config.BindEnv("allow_params_secret")
}

func initRunAllowGetAccessRequest(config *viper.Viper, flags *flag.FlagSet) {
	flags.Bool("allow-get-access-request", false, "Allow access request using GET and POST")
	_ = config.BindPFlag("allow_get_access_request", flags.Lookup("allow-get-access-request"))
	_ = config.BindEnv("allow_get_access_request")
}

func initRunRequirePKCEForPublicClients(config *viper.Viper, flags *flag.FlagSet) {
	flags.Bool("require-pkce", false, "Require PKCE for code flows for public OAuth clients")
	_ = config.BindPFlag("require_pkce", flags.Lookup("require-pkce"))
	_ = config.BindEnv("require_pkce")
}

func initRunRedirectUriSeparator(config *viper.Viper, flags *flag.FlagSet) {
	flags.String("redirect-uri-separator", "", "Delimiter for clients to specify multiple redirect URIs")
	_ = config.BindPFlag("redirect_uri_separator", flags.Lookup("redirect-uri-separator"))
	_ = config.BindEnv("redirect_uri_separator")
}

func initRunRetainTokenAfterRefresh(config *viper.Viper, flags *flag.FlagSet) {
	flags.Bool("retain-token-after-refresh", false, "Retain the access and refresh token for re-use")
	_ = config.BindPFlag("retain_token_after_refresh", flags.Lookup("retain-token-after-refresh"))
	_ = config.BindEnv("retain_token_after_refresh")
}

func initRunIssuerUrl(config *viper.Viper, flags *flag.FlagSet) {
	flags.String("issuer-url", "", "URL of this server")
	_ = config.BindPFlag("issuer_url", flags.Lookup("issuer-url"))
	_ = config.BindEnv("issuer_url")
}

func initRunAuthBaseUrl(config *viper.Viper, flags *flag.FlagSet) {
	flags.String("auth-base-url", "", "Base URL for authentication redirects")
	_ = config.BindPFlag("auth_base_url", flags.Lookup("auth-base-url"))
	_ = config.BindEnv("auth_base_url")
}

func initRunInfoBaseUrl(config *viper.Viper, flags *flag.FlagSet) {
	flags.String("info-base-url", "", "Base URL for the user information endpoint")
	_ = config.BindPFlag("info_base_url", flags.Lookup("info-base-url"))
	_ = config.BindEnv("info_base_url")
}

func initRunTokenBaseUrl(config *viper.Viper, flags *flag.FlagSet) {
	flags.String("token-base-url", "", "Base URL for JWT issuing")
	_ = config.BindPFlag("token_base_url", flags.Lookup("token-base-url"))
	_ = config.BindEnv("token_base_url")
}

func initRunKeysBaseUrl(config *viper.Viper, flags *flag.FlagSet) {
	flags.String("keys-base-url", "", "Base URL for public key information")
	_ = config.BindPFlag("keys_base_url", flags.Lookup("keys-base-url"))
	_ = config.BindEnv("keys_base_url")
}

func initRunEmailDomain(config *viper.Viper, flags *flag.FlagSet) {
	flags.String("email-domain", "localhost", "Email domain for auto-generated principal addresses")
	_ = config.BindPFlag("email_domain", flags.Lookup("email-domain"))
	_ = config.BindEnv("email_domain")
}

func initRunEmailVerified(config *viper.Viper, flags *flag.FlagSet) {
	flags.Bool("email-verified", false, "Mark auto-generated principal addresses as verified")
	_ = config.BindPFlag("email_verified", flags.Lookup("email-verified"))
	_ = config.BindEnv("email_verified")
}

func initRunKeySeed(config *viper.Viper, flags *flag.FlagSet) {
	flags.Int64("key-seed", 0, "Seed for generated key. If not set, the OS crypto source will be used.")
	_ = config.BindPFlag("key_seed", flags.Lookup("key-seed"))
	_ = config.BindEnv("key_seed")
}

func initRunKeySize(config *viper.Viper, flags *flag.FlagSet) {
	flags.Uint16("key-size", 2048, "Number of bits to use for the generated key.")
	_ = config.BindPFlag("key_size", flags.Lookup("key-size"))
	_ = config.BindEnv("key_size")
}

func initRunKeyGenerate(config *viper.Viper, flags *flag.FlagSet) {
	flags.Bool("key-generate", false, "Generate a key instead of reading from a file.")
	_ = config.BindPFlag("key_generate", flags.Lookup("key-generate"))
	_ = config.BindEnv("key_generate")
}

func initRunLogLevel(config *viper.Viper, flags *flag.FlagSet) {
	flags.String("log-level", "info", "Output verbosity")
	_ = config.BindPFlag("log_level", flags.Lookup("log-level"))
	_ = config.BindEnv("log_level")
}
