// Copyright 2024 eve.  All rights reserved.

package options

import (
	"github.com/spf13/pflag"
)

var _ IOptions = (*ClientCertAuthenticationOptions)(nil)

// ClientCertAuthenticationOptions provides different options for client cert auth.
type ClientCertAuthenticationOptions struct {
	// ClientCA is the certificate bundle for all the signers that you'll recognize for incoming client certificates
	ClientCA string `json:"client-ca-file" mapstructure:"client-ca-file"`
}

// NewClientCertAuthenticationOptions creates a ClientCertAuthenticationOptions object with default parameters.
func NewClientCertAuthenticationOptions() *ClientCertAuthenticationOptions {
	return &ClientCertAuthenticationOptions{
		ClientCA: "",
	}
}

// Validate is used to parse and validate the parameters entered by the user at
// the command line when the program starts.
func (o *ClientCertAuthenticationOptions) Validate() []error {
	return []error{}
}

// AddFlags adds flags related to ClientCertAuthenticationOptions for a specific server to the
// specified FlagSet.
func (o *ClientCertAuthenticationOptions) AddFlags(fs *pflag.FlagSet, prefixs ...string) {
	fs.StringVar(&o.ClientCA, "client-ca-file", o.ClientCA, ""+
		"If set, any request presenting a client certificate signed by one of "+
		"the authorities in the client-ca-file is authenticated with an identity "+
		"corresponding to the CommonName of the client certificate.")
}
