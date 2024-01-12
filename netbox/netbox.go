package netbox

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/go-resty/resty/v2"
	"github.com/netbox-community/go-netbox/v3/netbox/client"
	"github.com/netbox-community/go-netbox/v3/netbox/client/dcim"
	"github.com/netbox-community/go-netbox/v3/netbox/models"
	"github.com/rogerscuall/crispy-enigma/pkg"
)

var status = "active"
var pageLimit = int64(100)

type Manufacturers struct {
	List []models.Manufacturer
}

func createManufacturer(nb *client.NetBoxAPI, vnd models.Manufacturer) error {
	_, err := nb.Dcim.DcimManufacturersCreate(&dcim.DcimManufacturersCreateParams{
		Context: context.TODO(),
		Data:    &vnd,
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to create manufacturer %s: %w", vnd.Display, err)
	}
	//fmt.Println("Last Updated: ", crd.Payload.LastUpdated)
	return nil
}

func findManufacturer(nb *client.NetBoxAPI, slug string) (fnd bool, id int64, err error) {
	reqM := dcim.NewDcimManufacturersListParams().WithLimit(&pageLimit).WithSlug(&slug)
	rsp, err := nb.Dcim.DcimManufacturersList(reqM, nil)
	if err != nil {
		return fnd, id, fmt.Errorf("failed to find manufacturer %v: %w", &slug, err)
	}
	if len(rsp.Payload.Results) != 0 {
		fnd = true
		id = rsp.Payload.Results[0].ID
		fmt.Printf("Vendor: %s \tID: %v \n",
			*rsp.Payload.Results[0].Name, id)
	}
	return fnd, id, nil
}

type DeviceTypes struct {
	List []models.DeviceType
}

func createDeviceType(nb *client.NetBoxAPI, dt models.DeviceType) error {
	man := models.Manufacturer{
		Display: dt.Manufacturer.Display,
		Name:    dt.Manufacturer.Name,
		Slug:    dt.Manufacturer.Slug,
	}

	found, id, err := findManufacturer(nb, *man.Slug)
	if err != nil || !found {
		return fmt.Errorf("error finding manufacturer %s: %w", man.Display, err)
	}

	ndt := models.WritableDeviceType{
		Manufacturer: &id,
		Display:      dt.Display,
		Model:        dt.Model,
		Slug:         dt.Slug,
		Tags:         []*models.NestedTag{},
	}
	f := strfmt.NewFormats()
	err = ndt.Validate(f)
	if err != nil {
		return fmt.Errorf("failed to validate values for type %s: %w", *dt.Model, err)
	}
	//TODO: Is failing to create the device.
	_, err = nb.Dcim.DcimDeviceTypesCreate(&dcim.DcimDeviceTypesCreateParams{
		Context: context.TODO(),
		Data:    &ndt,
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to create device type %s: %w", *dt.Model, err)
	}
	//fmt.Println("Last Updated: ", crd.Payload.LastUpdated)
	return nil
}

func findDeviceType(nb *client.NetBoxAPI, slug string) (fnd bool, id int64, err error) {
	rsp, err := nb.Dcim.DcimDeviceTypesList(&dcim.DcimDeviceTypesListParams{
		Context: context.TODO(),
		SlugIe:  &slug,
	}, nil)
	if err != nil {
		return fnd, id, fmt.Errorf("failed to find device type %v: %w", &slug, err)
	}
	if len(rsp.Payload.Results) != 0 {
		fnd = true
		id = rsp.Payload.Results[0].ID
		fmt.Printf("Device Type: %q \tID: %v \n",
			strings.TrimSpace(*rsp.Payload.Results[0].Model), id)
	}
	return fnd, id, nil
}

type DeviceRoles struct {
	List []models.DeviceRole
}

func createDeviceRole(nb *client.NetBoxAPI, dr models.DeviceRole) error {
	f := strfmt.NewFormats()
	err := dr.Validate(f)
	if err != nil {
		return fmt.Errorf("failed to validate values for type %s: %w", *dr.Name, err)
	}
	_, err = nb.Dcim.DcimDeviceRolesCreate(&dcim.DcimDeviceRolesCreateParams{
		Context: context.TODO(),
		Data:    &dr,
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to create device role %s: %w", dr.Display, err)
	}
	//fmt.Println("Last Updated: ", crd.Payload.LastUpdated)
	return nil
}

func findDeviceRole(nb *client.NetBoxAPI, slug string) (fnd bool, id int64, err error) {
	rsp, err := nb.Dcim.DcimDeviceRolesList(&dcim.DcimDeviceRolesListParams{
		Context: context.TODO(),
		SlugIe:  &slug,
	}, nil)
	if err != nil {
		return fnd, id, fmt.Errorf("failed to find device role %v: %w", &slug, err)
	}
	if len(rsp.Payload.Results) != 0 {
		fnd = true
		id = rsp.Payload.Results[0].ID
		fmt.Printf("Site: %q \tID: %v \n",
			strings.TrimSpace(rsp.Payload.Results[0].Display), id)
	}
	return fnd, id, nil
}

type Sites struct {
	List []models.Site
}

func createSite(nb *client.NetBoxAPI, s models.Site) error {
	ns := models.WritableSite{
		Name:    s.Name,
		Display: s.Display,
		Slug:    s.Slug,
	}
	f := strfmt.NewFormats()
	err := ns.Validate(f)
	if err != nil {
		return fmt.Errorf("failed to validate values for site %s: %w", ns.Display, err)
	}

	_, err = nb.Dcim.DcimSitesCreate(&dcim.DcimSitesCreateParams{
		Context: context.TODO(),
		Data:    &ns,
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to create site %s: %w", ns.Display, err)
	}
	//fmt.Println("Last Updated: ", crd.Payload.LastUpdated)
	return nil
}

func findSite(nb *client.NetBoxAPI, slug string) (fnd bool, id int64, err error) {
	rsp, err := nb.Dcim.DcimSitesList(&dcim.DcimSitesListParams{
		Context: context.TODO(),
		SlugIe:  &slug,
	}, nil)
	if err != nil {
		return fnd, id, fmt.Errorf("failed to find site %v: %w", &slug, err)
	}
	if len(rsp.Payload.Results) != 0 {
		fnd = true
		id = rsp.Payload.Results[0].ID
		fmt.Printf("Device Role: %q \tID: %v \n",
			strings.TrimSpace(rsp.Payload.Results[0].Display), id)
	}
	return fnd, id, nil
}

func createToken(usr, pwd string, url *url.URL) (string, error) {
	client := resty.New()
	client.SetBaseURL("https://" + url.Host)

	body := fmt.Sprintf(`{"username":"%s", "password":"%s"}`, usr, pwd)

	result := make(map[string]interface{})
	_, err := client.R().
		SetResult(&result).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post("/api/users/tokens/provision/")

	if err != nil {
		return "", fmt.Errorf("error requesting a token: %w", err)
	}

	if val, ok := result["key"]; ok {
		return val.(string), nil
	}

	return "", fmt.Errorf("empty token")
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func getDeviceIDs(nb *client.NetBoxAPI, in models.Device) (fnd bool, out *models.WritableDeviceWithConfigContext, err error) {
	rsp, err := nb.Dcim.DcimDevicesList(&dcim.DcimDevicesListParams{
		Context: context.TODO(),
		NameIe:  in.Name,
	}, nil)
	if err != nil {
		return fnd, out, fmt.Errorf("failed to find device %s: %w", *in.Name, err)
	}
	var id int64
	if len(rsp.Payload.Results) != 0 {
		fnd = true
		id = rsp.Payload.Results[0].ID
		fmt.Printf("Device: %q \tID: %v \n",
			strings.TrimSpace(*rsp.Payload.Results[0].Name), id)

		out = &models.WritableDeviceWithConfigContext{
			Name:       in.Name,
			ID:         id,
			DeviceRole: &rsp.Payload.Results[0].DeviceRole.ID,
			DeviceType: &rsp.Payload.Results[0].DeviceType.ID,
			Site:       &rsp.Payload.Results[0].Site.ID,
			Tags:       rsp.Payload.Results[0].Tags,
		}
		return fnd, out, nil
	}
	find, dr, err := findDeviceRole(nb, *in.DeviceRole.Slug)
	if err != nil || !find {
		return fnd, out, fmt.Errorf("failed to find device role id for %s: %w", *in.Name, err)
	}
	find, dt, err := findDeviceType(nb, *in.DeviceType.Slug)
	if err != nil || !find {
		return fnd, out, fmt.Errorf("failed to find device type id for %s: %w", *in.Name, err)
	}
	find, st, err := findSite(nb, *in.Site.Slug)
	if err != nil || !find {
		return fnd, out, fmt.Errorf("failed to find site id for %s: %w", *in.Name, err)
	}
	out = &models.WritableDeviceWithConfigContext{
		Name:       in.Name,
		Role:       &dr,
		DeviceRole: &dr,
		DeviceType: &dt,
		Site:       &st,
		Tags:       []*models.NestedTag{},
	}
	return fnd, out, nil
}


func createResources(app *pkg.Application) error {
	//TODO: Create all these resources at once, for example create all the Manufactures once and check locally if they exist.
	// working with the manufacturer
	var manInput Manufacturers
	for _, device := range app.Devices {
		slug := strings.ToLower(device.Manufacturer)
		man := models.Manufacturer{
			Display: device.Manufacturer,
			Name:    &device.Manufacturer,
			Slug:    &slug,
		}
		manInput.List = append(manInput.List, man)
	}

	for _, vendor := range manInput.List {
		found, _, err := findManufacturer(app.NetBoxclient, *vendor.Slug)
		if err != nil {
			return fmt.Errorf("error finding manufacturer %s: %w", vendor.Display, err)
		}
		if !found {
			err = createManufacturer(app.NetBoxclient, vendor)
			if err != nil {
				return fmt.Errorf("error creating manufacturer %s: %w", vendor.Display, err)
			}
		}
	}

	// working with the device type

	var devTypes DeviceTypes
	for _, device := range app.Devices {
		slug := strings.ToLower(device.Model)
		maSlug := strings.ToLower(device.Manufacturer)
		dt := models.DeviceType{
			Manufacturer: &models.NestedManufacturer{
				Display: device.Manufacturer,
				Slug:    &maSlug,
				Name:    &device.Manufacturer,
			},
			Display: device.Model,
			Model:   &device.Model,
			Slug:    &slug,
		}
		devTypes.List = append(devTypes.List, dt)
	}

	for _, devType := range devTypes.List {
		found, _, err := findDeviceType(app.NetBoxclient, *devType.Slug)
		if err != nil {
			return fmt.Errorf("error finding device type %s: %w", devType.Display, err)
		}
		if !found {
			err = createDeviceType(app.NetBoxclient, devType)
			if err != nil {
				return fmt.Errorf("error creating device type %s: %w", devType.Display, err)
			}
		}
	}

	// working with the device role
	var devRoles DeviceRoles
	for _, device := range app.Devices {
		slug := strings.ToLower(device.DeviceRole)
		dr := models.DeviceRole{
			Display: device.DeviceRole,
			Slug:    &slug,
			Name:    &device.DeviceRole,
		}
		devRoles.List = append(devRoles.List, dr)

	}

	for _, devRole := range devRoles.List {
		found, _, err := findDeviceRole(app.NetBoxclient, *devRole.Slug)
		if err != nil {
			return fmt.Errorf("error finding device role %s: %w", devRole.Display, err)
		}
		if !found {
			err = createDeviceRole(app.NetBoxclient, devRole)
			if err != nil {
				return fmt.Errorf("error creating device role %s: %w", devRole.Display, err)
			}
		}
	}

	// working with the sites

	var devSites Sites

	for _, device := range app.Devices {
		slug := strings.ToLower(device.Site)
		site := models.Site{
			Display: device.Site,
			Slug:    &slug,
			Name:    &device.Site,
		}
		f := strfmt.NewFormats()
		err := site.Validate(f)
		if err != nil {
			return fmt.Errorf("failed to validate values for site %s: %w", site.Display, err)
		}
		devSites.List = append(devSites.List, site)
	}

	for _, devSite := range devSites.List {
		found, _, err := findSite(app.NetBoxclient, *devSite.Slug)
		if err != nil {
			return fmt.Errorf("error finding site %s: %w", devSite.Display, err)
		}
		if !found {
			err = createSite(app.NetBoxclient, devSite)
			if err != nil {
				return fmt.Errorf("error creating site %s: %w", devSite.Display, err)
			}
		}
	}

	return nil
}

func Work(app *pkg.Application) {

	err := createResources(app)
	if err != nil {
		panic(err)
	}

	// working with the devices

	var device models.Device

	for _, dev := range app.Devices {
		device = models.Device{
			Name:       &dev.Hostname,
			DeviceRole: &models.NestedDeviceRole{
				Slug: &dev.DeviceRole,

			},
			DeviceType: &models.NestedDeviceType{Slug: &dev.Model},
			Site:       &models.NestedSite{Slug: &dev.Site},
			Tags:       []*models.NestedTag{},
		}

		found, devWithIDs, err := getDeviceIDs(app.NetBoxclient, device)
		check(err)

		ctx := context.Background()
		if found {
			res, err := app.NetBoxclient.Dcim.DcimDevicesRead(&dcim.DcimDevicesReadParams{
				ID:      devWithIDs.ID,
				Context: ctx,
			}, nil)
			check(err)
			fmt.Println("Device already present: ", *res.Payload.Name)
			return
		}

		created, err := app.NetBoxclient.Dcim.DcimDevicesCreate(&dcim.DcimDevicesCreateParams{
			Context: ctx,
			Data:    devWithIDs,
		}, nil)
		check(err)

		res, err := app.NetBoxclient.Dcim.DcimDevicesRead(&dcim.DcimDevicesReadParams{
			ID:      created.Payload.ID,
			Context: ctx,
		}, nil)
		check(err)

		fmt.Println("Device created: ", *res.Payload.Name)
	}

	// found, devWithIDs, err := getDeviceIDs(nb, device)
	// check(err)

	// ctx := context.Background()
	// if found {
	// 	res, err := nb.Dcim.DcimDevicesRead(&dcim.DcimDevicesReadParams{
	// 		ID:      devWithIDs.ID,
	// 		Context: ctx,
	// 	}, nil)
	// 	check(err)
	// 	fmt.Println("Device already present: ", *res.Payload.Name)
	// 	return
	// }

	// created, err := nb.Dcim.DcimDevicesCreate(&dcim.DcimDevicesCreateParams{
	// 	Context: ctx,
	// 	Data:    devWithIDs,
	// }, nil)
	// check(err)

	// res, err := nb.Dcim.DcimDevicesRead(&dcim.DcimDevicesReadParams{
	// 	ID:      created.Payload.ID,
	// 	Context: ctx,
	// }, nil)
	// check(err)

	// fmt.Println("Device created: ", *res.Payload.Name)
}
