package authorization

// CreateUserWithProfileRequest represents the payload needed to create a user and their profile in one call.
type CreateUserWithProfileRequest struct {
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