package plugin

type (
	// defines the kind of the plugin e.g. Equinix, Scaphandre
	Kind string
	// defines the name of the account e.g. equinix-project-a
	AccountIdentifier string
	// defines the name of a resource which was found in an account e.g. equinix-server-b
	ResourceIdentifier string
)
