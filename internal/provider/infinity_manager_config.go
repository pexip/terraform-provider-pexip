package provider

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"net"
	"net/mail"
	"strings"
)

type InfinityManagerConfig struct {
	Hostname            string `json:"hostname"`
	Domain              string `json:"domain"`
	IP                  string `json:"ip"`
	Mask                string `json:"mask"`
	GW                  string `json:"gw"`
	DNS                 string `json:"dns"`
	NTP                 string `json:"ntp"`
	User                string `json:"user"`
	Pass                string `json:"pass"`
	AdminPassword       string `json:"admin_password"`
	ErrorReports        bool   `json:"error_reports"`
	EnableAnalytics     bool   `json:"enable_analytics"`
	ContactEmailAddress string `json:"contact_email_address"`
}

func (c *InfinityManagerConfig) Validate() error {
	var errs []string

	checkRequired := func(value, fieldName string) {
		if value == "" {
			errs = append(errs, fmt.Sprintf("%s is required", fieldName))
		}
	}

	checkRequiredIP := func(value, fieldName string) {
		if value == "" {
			errs = append(errs, fmt.Sprintf("%s is required", fieldName))
		} else if net.ParseIP(value) == nil {
			errs = append(errs, fmt.Sprintf("invalid %s: %s", fieldName, value))
		}
	}

	checkRequiredEmail := func(value, fieldName string) {
		if value == "" {
			errs = append(errs, fmt.Sprintf("%s is required", fieldName))
		} else {
			_, err := mail.ParseAddress(value)
			if err != nil {
				errs = append(errs, fmt.Sprintf("invalid %s '%s': %v", fieldName, value, err))
			}
		}
	}

	checkRequired(c.Hostname, "hostname")
	checkRequired(c.Domain, "domain")

	checkRequiredIP(c.IP, "IP address")
	checkRequiredIP(c.Mask, "subnet mask")
	checkRequiredIP(c.GW, "gateway")

	checkRequired(c.DNS, "DNS server")
	checkRequired(c.NTP, "NTP server")

	checkRequired(c.User, "user")
	checkRequired(c.Pass, "password")
	checkRequired(c.AdminPassword, "admin password")

	checkRequiredEmail(c.ContactEmailAddress, "contact email address")

	if len(errs) > 0 {
		return fmt.Errorf("validation failed: %s", strings.Join(errs, "; "))
	}

	return nil
}

func (c *InfinityManagerConfig) String() string {
	b, err := json.Marshal(c)
	if err != nil {
		log.Error().Err(err).Msgf("error marshalling InfinityManagerConfig: %v", err)
		return "{}"
	}
	return string(b)
}
