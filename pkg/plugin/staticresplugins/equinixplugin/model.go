package equinixplugin

import (
	"strings"
	"time"

	"carbonaut.dev/pkg/provider/data/account/project/resource"
)

func EquinixDataIntegration(device *EquinixDevice) *resource.StaticResData {
	staticresData := resource.StaticResData{
		ID:       device.ID,
		User:     device.User,
		OS:       &resource.OS{Version: device.OperatingSystem.Version, Distro: device.OperatingSystem.Distro, Name: device.OperatingSystem.Name},
		IPv4:     "",
		CPUs:     []*resource.CPU{},
		GPUs:     []*resource.GPU{},
		NICs:     []*resource.NIC{},
		Drives:   []*resource.DRIVE{},
		MemoryGB: device.Plan.Specs.Memory.Total,
		Location: &resource.Location{City: device.Facility.Address.City, Country: device.Facility.Address.Country, Address: device.Facility.Address.Address, ZipCode: device.Facility.Address.ZipCode, Code: device.Facility.Code},
	}
	// Add CPU Information
	for i := range device.Plan.Specs.Cpus {
		staticresData.CPUs = append(staticresData.CPUs, &resource.CPU{
			Count:        device.Plan.Specs.Cpus[i].Count,
			Type:         device.Plan.Specs.Cpus[i].Type,
			Cores:        device.Plan.Specs.Cpus[i].Cores,
			Threads:      device.Plan.Specs.Cpus[i].Threads,
			Speed:        device.Plan.Specs.Cpus[i].Speed,
			Arch:         device.Plan.Specs.Cpus[i].Arch,
			Model:        device.Plan.Specs.Cpus[i].Model,
			Manufacturer: device.Plan.Specs.Cpus[i].Manufacturer,
			Name:         device.Plan.Specs.Cpus[i].Name,
		})
	}
	// Add GPU Information
	for i := range device.Plan.Specs.Gpu {
		staticresData.GPUs = append(staticresData.GPUs, &resource.GPU{
			Count: device.Plan.Specs.Gpu[i].Count,
			Type:  device.Plan.Specs.Gpu[i].Type,
		})
	}
	// Add Drive Information
	for i := range device.Plan.Specs.Drives {
		staticresData.Drives = append(staticresData.Drives, &resource.DRIVE{
			Count: device.Plan.Specs.Drives[i].Count,
			Type:  device.Plan.Specs.Drives[i].Type,
			Size:  device.Plan.Specs.Drives[i].Size,
		})
	}
	// Add NIC Information
	for i := range device.Plan.Specs.Nics {
		staticresData.NICs = append(staticresData.NICs, &resource.NIC{
			Count: device.Plan.Specs.Nics[i].Count,
			Type:  device.Plan.Specs.Nics[i].Type,
		})
	}
	// Add IPv4
	for i := range device.IPAddresses {
		if device.IPAddresses[i].AddressFamily == 4 && !strings.HasPrefix(device.IPAddresses[i].Address, "10.") {
			staticresData.IPv4 = device.IPAddresses[i].Address
		}
	}
	return &staticresData
}

// TODO: implement supporting multiple pages Meta
type ProjectDevicesResponse struct {
	Devices []EquinixDevice `json:"devices"`
	// Meta struct {
	// 	First struct {
	// 		Href string `json:"href"`
	// 	} `json:"first"`
	// 	Previous interface{} `json:"previous"`
	// 	Self     struct {
	// 		Href string `json:"href"`
	// 	} `json:"self"`
	// 	Next interface{} `json:"next"`
	// 	Last struct {
	// 		Href string `json:"href"`
	// 	} `json:"last"`
	// 	CurrentPage int `json:"current_page"`
	// 	LastPage    int `json:"last_page"`
	// 	Total       int `json:"total"`
	// } `json:"meta"`
}

// TODO: implement supporting multiple pages Meta
type ProjectsResponse struct {
	Projects []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		// Description            interface{} `json:"description"`
		// CreatedAt              time.Time   `json:"created_at"`
		// UpdatedAt              time.Time   `json:"updated_at"`
		// BackendTransferEnabled bool        `json:"backend_transfer_enabled"`
		// Customdata             struct {
		// } `json:"customdata"`
		// Tags                    []interface{} `json:"tags"`
		// EventAlertConfiguration interface{}   `json:"event_alert_configuration"`
		// Memberships             []interface{} `json:"memberships"`
		// Invitations             []interface{} `json:"invitations"`
		// Devices                 []struct {
		// 	Href string `json:"href"`
		// } `json:"devices"`
		// SSHKeys []struct {
		// 	Href string `json:"href"`
		// } `json:"ssh_keys"`
		// Volumes []interface{} `json:"volumes"`
		// Members []struct {
		// 	Href string `json:"href"`
		// } `json:"members"`
		// PaymentMethod struct {
		// 	Href string `json:"href"`
		// } `json:"payment_method"`
		// Transfers    []interface{} `json:"transfers"`
		// Organization struct {
		// 	Href string `json:"href"`
		// } `json:"organization"`
		// Favorite bool   `json:"favorite"`
		// Type     string `json:"type"`
		// Href     string `json:"href"`
	} `json:"projects"`
	// Meta struct {
	// 	First struct {
	// 		Href string `json:"href"`
	// 	} `json:"first"`
	// 	Previous interface{} `json:"previous"`
	// 	Self     struct {
	// 		Href string `json:"href"`
	// 	} `json:"self"`
	// 	Next interface{} `json:"next"`
	// 	Last struct {
	// 		Href string `json:"href"`
	// 	} `json:"last"`
	// 	CurrentPage int `json:"current_page"`
	// 	LastPage    int `json:"last_page"`
	// 	Total       int `json:"total"`
	// } `json:"meta"`
}

type EquinixDeviceResponse struct {
	Devices []EquinixDevice `json:"devices"`
	Meta    struct {
		First struct {
			Href string `json:"href"`
		} `json:"first"`
		Previous interface{} `json:"previous"`
		Self     struct {
			Href string `json:"href"`
		} `json:"self"`
		Next interface{} `json:"next"`
		Last struct {
			Href string `json:"href"`
		} `json:"last"`
		CurrentPage int `json:"current_page"`
		LastPage    int `json:"last_page"`
		Total       int `json:"total"`
	} `json:"meta"`
}

type EquinixDevice struct {
	ID              string        `json:"id"`
	ShortID         string        `json:"short_id"`
	Hostname        string        `json:"hostname"`
	Description     interface{}   `json:"description"`
	Tags            []interface{} `json:"tags"`
	ImageURL        interface{}   `json:"image_url"`
	BillingCycle    string        `json:"billing_cycle"`
	User            string        `json:"user"`
	Iqn             string        `json:"iqn"`
	Locked          bool          `json:"locked"`
	BondingMode     int           `json:"bonding_mode"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	IpxeScriptURL   interface{}   `json:"ipxe_script_url"`
	AlwaysPxe       bool          `json:"always_pxe"`
	Storage         interface{}   `json:"storage"`
	Customdata      struct{}      `json:"customdata"`
	TerminationTime interface{}   `json:"termination_time"`
	CreatedBy       struct {
		ID             string    `json:"id"`
		ShortID        string    `json:"short_id"`
		FirstName      string    `json:"first_name"`
		LastName       string    `json:"last_name"`
		FullName       string    `json:"full_name"`
		Email          string    `json:"email"`
		CreatedAt      time.Time `json:"created_at"`
		UpdatedAt      time.Time `json:"updated_at"`
		Level          string    `json:"level"`
		AvatarThumbURL string    `json:"avatar_thumb_url"`
		Href           string    `json:"href"`
	} `json:"created_by"`
	OperatingSystem struct {
		ID                     string   `json:"id"`
		Slug                   string   `json:"slug"`
		Name                   string   `json:"name"`
		Version                string   `json:"version"`
		Preinstallable         bool     `json:"preinstallable"`
		Pricing                struct{} `json:"pricing"`
		DistroLabel            string   `json:"distro_label"`
		Distro                 string   `json:"distro"`
		DefaultOperatingSystem bool     `json:"default_operating_system"`
		ProvisionableOn        []string `json:"provisionable_on"`
		Licensed               bool     `json:"licensed"`
	} `json:"operating_system"`
	Facility struct {
		ID       string   `json:"id"`
		Name     string   `json:"name"`
		Code     string   `json:"code"`
		Features []string `json:"features"`
		Metro    struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Code    string `json:"code"`
			Country string `json:"country"`
		} `json:"metro"`
		IPRanges []interface{} `json:"ip_ranges"`
		Address  struct {
			ID          string      `json:"id"`
			Address     string      `json:"address"`
			Address2    interface{} `json:"address2"`
			City        string      `json:"city"`
			State       interface{} `json:"state"`
			ZipCode     string      `json:"zip_code"`
			Country     string      `json:"country"`
			Coordinates struct{}    `json:"coordinates"`
		} `json:"address"`
	} `json:"facility"`
	Project struct {
		Href string `json:"href"`
	} `json:"project"`
	SSHKeys []struct {
		Href string `json:"href"`
	} `json:"ssh_keys"`
	ProjectLite struct {
		Href string `json:"href"`
	} `json:"project_lite"`
	Volumes     []interface{} `json:"volumes"`
	IPAddresses []struct {
		ID            string        `json:"id"`
		AddressFamily int           `json:"address_family"`
		Netmask       string        `json:"netmask"`
		CreatedAt     time.Time     `json:"created_at"`
		Details       interface{}   `json:"details"`
		Tags          []interface{} `json:"tags"`
		Public        bool          `json:"public"`
		Cidr          int           `json:"cidr"`
		Management    bool          `json:"management"`
		Manageable    bool          `json:"manageable"`
		Enabled       bool          `json:"enabled"`
		GlobalIP      bool          `json:"global_ip"`
		Customdata    struct{}      `json:"customdata"`
		Project       struct {
			Href string `json:"href"`
		} `json:"project"`
		ProjectLite struct {
			Href string `json:"href"`
		} `json:"project_lite"`
		AssignedTo struct {
			Href string `json:"href"`
		} `json:"assigned_to"`
		Interface struct {
			Href string `json:"href"`
		} `json:"interface"`
		Network  string      `json:"network"`
		Address  string      `json:"address"`
		Gateway  string      `json:"gateway"`
		Href     string      `json:"href"`
		Facility interface{} `json:"facility"`
		Metro    struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Code    string `json:"code"`
			Country string `json:"country"`
		} `json:"metro"`
	} `json:"ip_addresses"`
	Favorite bool `json:"favorite"`
	Metro    struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Code    string `json:"code"`
		Country string `json:"country"`
	} `json:"metro"`
	Plan struct {
		ID              string   `json:"id"`
		Slug            string   `json:"slug"`
		Name            string   `json:"name"`
		Description     string   `json:"description"`
		Line            string   `json:"line"`
		DeploymentTypes []string `json:"deployment_types"`
		Categories      []string `json:"categories"`
		Type            string   `json:"type"`
		Class           string   `json:"class"`
		Legacy          bool     `json:"legacy"`
		Specs           struct {
			Cpus []struct {
				Count        int    `json:"count"`
				Type         string `json:"type"`
				Cores        string `json:"cores"`
				Threads      string `json:"threads"`
				Speed        string `json:"speed"`
				Arch         string `json:"arch"`
				Model        string `json:"model"`
				Manufacturer string `json:"manufacturer"`
				Name         string `json:"name"`
			} `json:"cpus"`
			Memory struct {
				Total string `json:"total"`
			} `json:"memory"`
			Drives []struct {
				Count    int    `json:"count"`
				Size     string `json:"size"`
				Type     string `json:"type"`
				Category string `json:"category"`
			} `json:"drives"`
			Nics []struct {
				Count int    `json:"count"`
				Type  string `json:"type"`
			} `json:"nics"`
			Gpu []struct {
				Count int    `json:"count"`
				Type  string `json:"type"`
			} `json:"gpu"`
			Features struct {
				Raid bool `json:"raid"`
				Txt  bool `json:"txt"`
			} `json:"features"`
		} `json:"specs"`
		Pricing struct {
			Hour float64 `json:"hour"`
		} `json:"pricing"`
		ReservationPricing struct {
			OneMonth struct {
				Month float64 `json:"month"`
			} `json:"one_month"`
			OneYear struct {
				Month float64 `json:"month"`
			} `json:"one_year"`
			Ml struct {
				OneMonth struct {
					Month float64 `json:"month"`
				} `json:"one_month"`
				OneYear struct {
					Month float64 `json:"month"`
				} `json:"one_year"`
			} `json:"ml"`
			Hk struct {
				OneMonth struct {
					Month float64 `json:"month"`
				} `json:"one_month"`
				OneYear struct {
					Month float64 `json:"month"`
				} `json:"one_year"`
			} `json:"hk"`
			Ty struct {
				OneMonth struct {
					Month float64 `json:"month"`
				} `json:"one_month"`
				OneYear struct {
					Month float64 `json:"month"`
				} `json:"one_year"`
			} `json:"ty"`
			Md struct {
				OneMonth struct {
					Month float64 `json:"month"`
				} `json:"one_month"`
				OneYear struct {
					Month float64 `json:"month"`
				} `json:"one_year"`
			} `json:"md"`
			La struct {
				OneMonth struct {
					Month float64 `json:"month"`
				} `json:"one_month"`
				OneYear struct {
					Month float64 `json:"month"`
				} `json:"one_year"`
			} `json:"la"`
			Ny struct {
				OneMonth struct {
					Month float64 `json:"month"`
				} `json:"one_month"`
				OneYear struct {
					Month float64 `json:"month"`
				} `json:"one_year"`
			} `json:"ny"`
			Mx struct {
				OneMonth struct {
					Month float64 `json:"month"`
				} `json:"one_month"`
				OneYear struct {
					Month float64 `json:"month"`
				} `json:"one_year"`
			} `json:"mx"`
			Ch struct {
				OneMonth struct {
					Month float64 `json:"month"`
				} `json:"one_month"`
				OneYear struct {
					Month float64 `json:"month"`
				} `json:"one_year"`
			} `json:"ch"`
			Sp struct {
				OneMonth struct {
					Month float64 `json:"month"`
				} `json:"one_month"`
				OneYear struct {
					Month float64 `json:"month"`
				} `json:"one_year"`
			} `json:"sp"`
			Fr struct {
				OneMonth struct {
					Month float64 `json:"month"`
				} `json:"one_month"`
				OneYear struct {
					Month float64 `json:"month"`
				} `json:"one_year"`
			} `json:"fr"`
			Da struct {
				OneMonth struct {
					Month float64 `json:"month"`
				} `json:"one_month"`
				OneYear struct {
					Month float64 `json:"month"`
				} `json:"one_year"`
			} `json:"da"`
			Pa struct {
				OneMonth struct {
					Month float64 `json:"month"`
				} `json:"one_month"`
				OneYear struct {
					Month float64 `json:"month"`
				} `json:"one_year"`
			} `json:"pa"`
			DB struct {
				OneMonth struct {
					Month float64 `json:"month"`
				} `json:"one_month"`
				OneYear struct {
					Month float64 `json:"month"`
				} `json:"one_year"`
			} `json:"db"`
			Ma struct {
				OneMonth struct {
					Month float64 `json:"month"`
				} `json:"one_month"`
				OneYear struct {
					Month float64 `json:"month"`
				} `json:"one_year"`
			} `json:"ma"`
			Sl struct {
				OneMonth struct {
					Month float64 `json:"month"`
				} `json:"one_month"`
				OneYear struct {
					Month float64 `json:"month"`
				} `json:"one_year"`
			} `json:"sl"`
			Ld struct {
				OneMonth struct {
					Month float64 `json:"month"`
				} `json:"one_month"`
				OneYear struct {
					Month float64 `json:"month"`
				} `json:"one_year"`
			} `json:"ld"`
			Mt struct {
				OneYear struct {
					Month float64 `json:"month"`
				} `json:"one_year"`
			} `json:"mt"`
			Sv struct {
				OneMonth struct {
					Month float64 `json:"month"`
				} `json:"one_month"`
				OneYear struct {
					Month float64 `json:"month"`
				} `json:"one_year"`
			} `json:"sv"`
			Se struct {
				OneMonth struct {
					Month float64 `json:"month"`
				} `json:"one_month"`
				OneYear struct {
					Month float64 `json:"month"`
				} `json:"one_year"`
			} `json:"se"`
			Sy struct {
				OneMonth struct {
					Month float64 `json:"month"`
				} `json:"one_month"`
				OneYear struct {
					Month float64 `json:"month"`
				} `json:"one_year"`
			} `json:"sy"`
			Am struct {
				OneMonth struct {
					Month float64 `json:"month"`
				} `json:"one_month"`
				OneYear struct {
					Month float64 `json:"month"`
				} `json:"one_year"`
			} `json:"am"`
			Tr struct {
				OneMonth struct {
					Month float64 `json:"month"`
				} `json:"one_month"`
				OneYear struct {
					Month float64 `json:"month"`
				} `json:"one_year"`
			} `json:"tr"`
			Dc struct {
				OneMonth struct {
					Month float64 `json:"month"`
				} `json:"one_month"`
			} `json:"dc"`
			Sg struct {
				OneMonth struct {
					Month float64 `json:"month"`
				} `json:"one_month"`
				OneYear struct {
					Month float64 `json:"month"`
				} `json:"one_year"`
			} `json:"sg"`
			Mb struct {
				OneMonth struct {
					Month float64 `json:"month"`
				} `json:"one_month"`
				OneYear struct {
					Month float64 `json:"month"`
				} `json:"one_year"`
			} `json:"mb"`
		} `json:"reservation_pricing"`
		AvailableIn []struct {
			Href  string `json:"href"`
			Price struct {
				Hour float64 `json:"hour"`
			} `json:"price"`
		} `json:"available_in"`
		AvailableInMetros []struct {
			Href  string `json:"href"`
			Price struct {
				Hour float64 `json:"hour"`
			} `json:"price"`
		} `json:"available_in_metros"`
	} `json:"plan"`
	DeviceType string `json:"device_type"`
	Actions    []struct {
		Type string `json:"type"`
		Name string `json:"name"`
	} `json:"actions"`
	NetworkFrozen bool   `json:"network_frozen"`
	Userdata      string `json:"userdata"`
	RootPassword  string `json:"root_password"`
	SwitchUUID    string `json:"switch_uuid"`
	NetworkPorts  []struct {
		ID   string `json:"id"`
		Type string `json:"type"`
		Name string `json:"name"`
		Data struct {
			Bonded bool   `json:"bonded"`
			Mac    string `json:"mac"`
		} `json:"data"`
		Bond struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"bond,omitempty"`
		NativeVirtualNetwork      interface{}   `json:"native_virtual_network"`
		VirtualNetworks           []interface{} `json:"virtual_networks"`
		DisbondOperationSupported bool          `json:"disbond_operation_supported"`
		Href                      string        `json:"href"`
		NetworkType               string        `json:"network_type,omitempty"`
	} `json:"network_ports"`
	State                        string `json:"state"`
	AllowSameVlanOnMultiplePorts bool   `json:"allow_same_vlan_on_multiple_ports"`
	Sos                          string `json:"sos"`
	Href                         string `json:"href"`
}
