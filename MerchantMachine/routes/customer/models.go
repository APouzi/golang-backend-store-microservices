package customer

type CustomerRegistration struct {
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	PhotoURL    string `json:"photo_url"`
	ProviderID  string `json:"provider_id"`
	UID         string `json:"uid"`
}

// DBUserProfileRequest matches the CreateUserWithProfileRequest in DBLayer
type DBUserProfileRequest struct {
	Email                    string  `json:"email"`
	FirstName                string  `json:"first_name"`
	LastName                 string  `json:"last_name"`
	PhoneNumberMobileE164    string  `json:"phone_number_mobile_e164"`
	PhoneNumberHomeE164      string  `json:"phone_number_home_e164"`
	PrimaryShippingAddressID *int    `json:"primary_shipping_address_id,omitempty"`
	PrimaryBillingAddressID  *int    `json:"primary_billing_address_id,omitempty"`
	PreferredLocale          *string `json:"preferred_locale,omitempty"`
	PreferredTimeZone        *string `json:"preferred_timezone,omitempty"`
}