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
	config.BindPFlag(c, flags.Lookup(f))
	config.BindEnv(c)

	for key, allowed := range allowedAuthorizeTypes {
		aF := f + "-" + key
		aC := c + "_" + string(allowed)

		flags.Bool(aF, false, "Allow authorization selectively. See --allow-authorize")
		config.BindPFlag(aC, flags.Lookup(aF))
		config.BindEnv(aC)
	}
}

func initRunAllowedAccess(config *viper.Viper, flags *flag.FlagSet) {
	f := "allow-access"
	c := "allow_access"

	flags.StringSlice(f, allowedAccessTypes.FlagValues(), "Allowed access request types")
	config.BindPFlag(f, flags.Lookup(f))
	config.BindEnv(f)

	for key, allowed := range allowedAccessTypes {
		aF := f + "-" + key
		aC := c + "_" + string(allowed)

		flags.Bool(aF, false, "Allow access request selectively. See --allow-access")
		config.BindPFlag(aC, flags.Lookup(aF))
		config.BindEnv(aC)
	}
}

func initRunAuthorizationExpiration(config *viper.Viper, flags *flag.FlagSet) {
	flags.Int32("authorization-expiration", AuthorizationExpiration, "Authorization token expiration in seconds")
	config.BindPFlag("authorization_expiration", flags.Lookup("authorization-expiration"))
	config.BindEnv("authorization_expiration")
}

func initRunAccessExpiration(config *viper.Viper, flags *flag.FlagSet) {
	flags.Int32("access-expiration", AccessExpiration, "Access token expiration in seconds")
	config.BindPFlag("access_expiration", flags.Lookup("access-expiration"))
	config.BindEnv("access_expiration")
}

func initRunErrorStatusCode(config *viper.Viper, flags *flag.FlagSet) {
	flags.Int("error-status-code", ErrorStatusCode, "HTTP status code to return for errors")
	config.BindPFlag("error_status_code", flags.Lookup("error-status-code"))
	config.BindEnv("error_status_code")
}

func initRunAllowClientSecretInParams(config *viper.Viper, flags *flag.FlagSet) {
	flags.Bool("allow-params-secret", false, "Allows client secret also in request parameters in addition to Authorization header")
	config.BindPFlag("allow_params_secret", flags.Lookup("allow-params-secret"))
	config.BindEnv("allow_params_secret")
}

func initRunAllowGetAccessRequest(config *viper.Viper, flags *flag.FlagSet) {
	flags.Bool("allow-get-access-request", false, "Allow access request using GET and POST")
	config.BindPFlag("allow_get_access_request", flags.Lookup("allow-get-access-request"))
	config.BindEnv("allow_get_access_request")
}

func initRunRequirePKCEForPublicClients(config *viper.Viper, flags *flag.FlagSet) {
	flags.Bool("require-pkce", false, "Require PKCE for code flows for public OAuth clients")
	config.BindPFlag("require_pkce", flags.Lookup("require-pkce"))
	config.BindEnv("require_pkce")
}

func initRunRedirectUriSeparator(config *viper.Viper, flags *flag.FlagSet) {
	flags.String("redirect-uri-separator", "", "Delimiter for clients to specify multiple redirect URIs")
	config.BindPFlag("redirect_uri_separator", flags.Lookup("redirect-uri-separator"))
	config.BindEnv("redirect_uri_separator")
}

func initRunRetainTokenAfterRefresh(config *viper.Viper, flags *flag.FlagSet) {
	flags.Bool("retain-token-after-refresh", false, "Retain the access and refresh token for re-use")
	config.BindPFlag("retain_token_after_refresh", flags.Lookup("retain-token-after-refresh"))
	config.BindEnv("retain_token_after_refresh")
}

func initRunIssuerUrl(config *viper.Viper, flags *flag.FlagSet) {
	flags.String("issuer-url", "", "URL of this server")
	config.BindPFlag("issuer_url", flags.Lookup("issuer-url"))
	config.BindEnv("issuer_url")
}

func initRunAuthBaseUrl(config *viper.Viper, flags *flag.FlagSet) {
	flags.String("auth-base-url", "", "Base URL for authentication redirects")
	config.BindPFlag("auth_base_url", flags.Lookup("auth-base-url"))
	config.BindEnv("auth_base_url")
}

func initRunInfoBaseUrl(config *viper.Viper, flags *flag.FlagSet) {
	flags.String("info-base-url", "", "Base URL for the user information endpoint")
	config.BindPFlag("info_base_url", flags.Lookup("info-base-url"))
	config.BindEnv("info_base_url")
}

func initRunTokenBaseUrl(config *viper.Viper, flags *flag.FlagSet) {
	flags.String("token-base-url", "", "Base URL for JWT issuing")
	config.BindPFlag("token_base_url", flags.Lookup("token-base-url"))
	config.BindEnv("token_base_url")
}

func initRunKeysBaseUrl(config *viper.Viper, flags *flag.FlagSet) {
	flags.String("keys-base-url", "", "Base URL for public key information")
	config.BindPFlag("keys_base_url", flags.Lookup("keys-base-url"))
	config.BindEnv("keys_base_url")
}

func initRunEmailDomain(config *viper.Viper, flags *flag.FlagSet) {
	flags.String("email-domain", "localhost", "Email domain for auto-generated principal addresses")
	config.BindPFlag("email_domain", flags.Lookup("email-domain"))
	config.BindEnv("email_domain")
}

func initRunEmailVerified(config *viper.Viper, flags *flag.FlagSet) {
	flags.Bool("email-verified", false, "Mark auto-generated principal addresses as verified")
	config.BindPFlag("email_verified", flags.Lookup("email-verified"))
	config.BindEnv("email_verified")
}

func initRunKeySeed(config *viper.Viper, flags *flag.FlagSet) {
	flags.Int64("key-seed", 0, "Seed for generated key. If not set, the OS crypto source will be used.")
	config.BindPFlag("key_seed", flags.Lookup("key-seed"))
	config.BindEnv("key_seed")
}

func initRunKeySize(config *viper.Viper, flags *flag.FlagSet) {
	flags.Uint16("key-size", 2048, "Number of bits to use for the generated key.")
	config.BindPFlag("key_size", flags.Lookup("key-size"))
	config.BindEnv("key_size")
}

func initRunKeyGenerate(config *viper.Viper, flags *flag.FlagSet) {
	flags.Bool("key-generate", false, "Generate a key instead of reading from a file.")
	config.BindPFlag("key_generate", flags.Lookup("key-generate"))
	config.BindEnv("key_generate")
}

func initRunLogLevel(config *viper.Viper, flags *flag.FlagSet) {
	flags.String("log-level", "info", "Output verbosity")
	config.BindPFlag("log_level", flags.Lookup("log-level"))
	config.BindEnv("log_level")
}
