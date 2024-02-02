package netbox

import (
	"context"
	"fmt"
	"log"
	"strings"

	n "github.com/netbox-community/go-netbox/v3"
	"github.com/rogerscuall/crispy-enigma/pkg"
)

var status = "active"
var pageLimit = int32(100)

type Manufacturers struct {
	List []*n.ManufacturerRequest
}

func createManufacturer(nb *n.APIClient, manu n.ManufacturerRequest) error {
	_, _, err := nb.DcimAPI.DcimManufacturersCreate(context.Background()).ManufacturerRequest(manu).Execute()
	if err != nil {
		if !strings.Contains(err.Error(), "no value given for required property devicetype_count") {
			return fmt.Errorf("failed to create manufacturer %s: %w", manu.Name, err)
		}
	}
	fmt.Println("Manufacturer created: ", manu.Name)
	return nil
}

func findManufacturer(nb *n.APIClient, slug string) (fnd bool, id int32, err error) {
	reqM, _, err := nb.DcimAPI.DcimManufacturersList(context.Background()).Slug([]string{slug}).Limit(pageLimit).Execute()
	if err != nil {
		return fnd, id, fmt.Errorf("failed to find manufacturer %v: %w", &slug, err)
	}
	if len(reqM.Results) != 0 {
		fnd = true
		id = reqM.Results[0].Id
		fmt.Printf("Vendor: %s \tID: %v \n", reqM.Results[0].Name, id)
	}
	return fnd, id, nil
}

func createDeviceType(nb *n.APIClient, dt *n.WritableDeviceTypeRequest) error {
	apiRequest := nb.DcimAPI.DcimDeviceTypesCreate(context.Background()).WritableDeviceTypeRequest(*dt)
	_, _, err := apiRequest.Execute()
	if err != nil {
		if !strings.Contains(err.Error(), "no value given for required property ") {
			return fmt.Errorf("failed to create device type %s: %w", dt.Model, err)
		}
	}
	fmt.Println("Device Type created: ", dt.Model)
	return nil
}

func findDeviceType(nb *n.APIClient, slug string) (fnd bool, id int32, err error) {
	reqDT, _, err := nb.DcimAPI.DcimDeviceTypesList(context.Background()).Slug([]string{slug}).Limit(pageLimit).Execute()
	if err != nil {
		return fnd, id, fmt.Errorf("failed to find device type %v: %w", &slug, err)
	}
	if len(reqDT.Results) != 0 {
		fnd = true
		id = reqDT.Results[0].Id
		fmt.Printf("Device Type: %s \tID: %v \n", reqDT.Results[0].Model, id)
	}
	return fnd, id, nil
}

func createDeviceRole(nb *n.APIClient, dr n.WritableDeviceRoleRequest) error {
	apiRequest := nb.DcimAPI.DcimDeviceRolesCreate(context.Background()).WritableDeviceRoleRequest(dr)
	_, _, err := apiRequest.Execute()
	if err != nil {
		if !strings.Contains(err.Error(), "no value given for required property ") {
			return fmt.Errorf("failed to create device role %s: %w", dr.Name, err)
		}
	}
	fmt.Println("Device Role created: ", dr.Name)
	return nil
}

func findDeviceRole(nb *n.APIClient, slug string) (fnd bool, id int32, err error) {
	reqDr, _, err := nb.DcimAPI.DcimDeviceRolesList(context.Background()).Slug([]string{slug}).Limit(pageLimit).Execute()
	if err != nil {
		return fnd, id, fmt.Errorf("failed to find device role %v: %w", &slug, err)
	}
	if len(reqDr.Results) != 0 {
		fnd = true
		id = reqDr.Results[0].Id
		fmt.Printf("Device Role: %s \tID: %v \n", reqDr.Results[0].Name, id)
	}
	return fnd, id, nil
}

func createSite(nb *n.APIClient, s n.WritableSiteRequest) error {
	apiRequest := nb.DcimAPI.DcimSitesCreate(context.Background()).WritableSiteRequest(s)
	_, _, err := apiRequest.Execute()
	if err != nil {
		if !strings.Contains(err.Error(), "no value given for required property ") {
			return fmt.Errorf("failed to create site %s: %w", s.Name, err)
		}
	}
	fmt.Println("Site created: ", s.Name)
	return nil
}

func findSite(nb *n.APIClient, slug string) (fnd bool, id int32, err error) {
	reqSite, _, err := nb.DcimAPI.DcimSitesList(context.Background()).Slug([]string{slug}).Limit(pageLimit).Execute()
	if err != nil {
		return fnd, id, fmt.Errorf("failed to find site %v: %w", &slug, err)
	}
	if len(reqSite.Results) != 0 {
		fnd = true
		id = reqSite.Results[0].Id
		fmt.Printf("Site: %s \tID: %v \n", reqSite.Results[0].Name, id)
	}
	return fnd, id, nil
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func findDevice(nb *n.APIClient, name string) (fnd bool, err error) {
	rsp, _, err := nb.DcimAPI.DcimDevicesList(context.Background()).Name([]string{name}).Limit(pageLimit).Execute()
	if err != nil {
		return fnd, fmt.Errorf("failed to find device %s: %w", name, err)
	}
	if len(rsp.Results) != 0 {
		fnd = true
		fmt.Printf("Device: %q \tID: %v \n", name, rsp.Results[0].Id)
	}
	return fnd, nil
}

// func getDeviceIDs(nb *n.APIClient, in models.Device) (fnd bool, out *models.WritableDeviceWithConfigContext, err error) {
// 	rsp, err := nb.Dcim.DcimDevicesList(&dcim.DcimDevicesListParams{
// 		Context: context.TODO(),
// 		NameIe:  in.Name,
// 	}, nil)
// 	if err != nil {
// 		return fnd, out, fmt.Errorf("failed to find device %s: %w", *in.Name, err)
// 	}
// 	var id int64
// 	if len(rsp.Payload.Results) != 0 {
// 		fnd = true
// 		id = rsp.Payload.Results[0].ID
// 		fmt.Printf("Device: %q \tID: %v \n",
// 			strings.TrimSpace(*rsp.Payload.Results[0].Name), id)

// 		out = &models.WritableDeviceWithConfigContext{
// 			Name:       in.Name,
// 			ID:         id,
// 			DeviceRole: &rsp.Payload.Results[0].DeviceRole.ID,
// 			DeviceType: &rsp.Payload.Results[0].DeviceType.ID,
// 			Site:       &rsp.Payload.Results[0].Site.ID,
// 			Tags:       rsp.Payload.Results[0].Tags,
// 		}
// 		return fnd, out, nil
// 	}
// 	find, dr, err := findDeviceRole(nb, *in.DeviceRole.Slug)
// 	if err != nil || !find {
// 		return fnd, out, fmt.Errorf("failed to find device role id for %s: %w", *in.Name, err)
// 	}
// 	find, dt, err := findDeviceType(nb, *in.DeviceType.Slug)
// 	if err != nil || !find {
// 		return fnd, out, fmt.Errorf("failed to find device type id for %s: %w", *in.Name, err)
// 	}
// 	find, st, err := findSite(nb, *in.Site.Slug)
// 	if err != nil || !find {
// 		return fnd, out, fmt.Errorf("failed to find site id for %s: %w", *in.Name, err)
// 	}
// 	out = &models.WritableDeviceWithConfigContext{
// 		Name:       in.Name,
// 		Role:       &dr,
// 		DeviceRole: &dr,
// 		DeviceType: &dt,
// 		Site:       &st,
// 		Tags:       []*models.NestedTag{},
// 	}
// 	return fnd, out, nil
// }

func createResources(app *pkg.Application) error {
	//TODO: Create all these resources at once, for example create all the Manufactures once and check locally if they exist.
	// working with the manufacturer
	for _, device := range app.Devices {
		slug := strings.ToLower(device.Manufacturer)
		found, manID, err := findManufacturer(app.NetBoxclient, slug)
		if err != nil {
			return fmt.Errorf("error finding manufacturer %s: %w", device.Manufacturer, err)
		}

		if !found {
			man := n.NewManufacturerRequestWithDefaults()
			man.Name = device.Manufacturer
			man.Slug = slug
			man.Tags = []n.NestedTagRequest{}
			log.Println("Creating manufacturer: ", device.Manufacturer)
			err = createManufacturer(app.NetBoxclient, *man)
			if err != nil {
				return fmt.Errorf("error creating manufacturer %s: %w", device.Manufacturer, err)
			}
		}
		slugDeviceType := strings.ToLower(device.Model)
		found, deviceTypeID, err := findDeviceType(app.NetBoxclient, slugDeviceType)
		if err != nil {
			return fmt.Errorf("error finding device type %s: %w", device.Manufacturer, err)
		}
		if !found {
			log.Println("Creating device type: ", device.Manufacturer)
			deviceType := n.NewWritableDeviceTypeRequestWithDefaults()
			deviceType.Model = device.Model
			deviceType.Slug = slugDeviceType
			deviceType.Manufacturer = manID
			deviceType.Tags = []n.NestedTagRequest{}
			err = createDeviceType(app.NetBoxclient, deviceType)
			if err != nil {
				return fmt.Errorf("error creating device type %s: %w", device.Manufacturer, err)
			}
		}
		slugRole := strings.ToLower(device.DeviceRole)
		found, roleID, err := findDeviceRole(app.NetBoxclient, slugRole)
		if err != nil {
			return fmt.Errorf("error finding device role %s: %w", device.Manufacturer, err)
		}
		if !found {
			log.Println("Creating device role: ", device.Manufacturer)
			deviceRole := n.NewWritableDeviceRoleRequestWithDefaults()
			deviceRole.Name = device.DeviceRole
			deviceRole.Slug = slugRole
			// deviceRole.Color = "ffffff"
			deviceRole.Tags = []n.NestedTagRequest{}
			err = createDeviceRole(app.NetBoxclient, *deviceRole)
			if err != nil {
				return fmt.Errorf("error creating device role %s: %w", device.DeviceRole, err)
			}
		}
		slugSite := strings.ToLower(device.Site)
		found, siteID, err := findSite(app.NetBoxclient, slugSite)
		if err != nil {
			return fmt.Errorf("error finding site %s: %w", device.Manufacturer, err)
		}
		if !found {
			log.Println("Creating site: ", device.Site)
			site := n.NewWritableSiteRequestWithDefaults()
			site.Name = device.Site
			site.Slug = slugSite
			site.Tags = []n.NestedTagRequest{}
			err = createSite(app.NetBoxclient, *site)
			if err != nil {
				return fmt.Errorf("error creating site %s: %w", device.Manufacturer, err)
			}
		}
		found, err = findDevice(app.NetBoxclient, device.Hostname)
		if err != nil {
			return fmt.Errorf("error finding device %s: %w", device.Manufacturer, err)
		}
		if !found {
			createDevice := n.NewWritableDeviceWithConfigContextRequestWithDefaults()
			createDevice.SetName(device.Hostname)
			createDevice.SetName(device.Hostname)
			createDevice.SetSite(siteID)
			createDevice.SetRole(roleID)
			createDevice.SetDeviceType(deviceTypeID)
			apiRequest := app.NetBoxclient.DcimAPI.DcimDevicesCreate(context.Background()).WritableDeviceWithConfigContextRequest(*createDevice)
			_, _, err = apiRequest.Execute()
			if err != nil {
				log.Println(err)
				log.Println("Moving to the next device...")
			}
			fmt.Println("Device created: ", device.Hostname)
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

	//log.Printf("Device: %v", device)

	// found, devWithIDs, err := getDeviceIDs(app.NetBoxclient, device)
	// check(err)

	// ctx := context.Background()
	// if found {
	// 	res, err := app.NetBoxclient.Dcim.DcimDevicesRead(&dcim.DcimDevicesReadParams{
	// 		ID:      devWithIDs.ID,
	// 		Context: ctx,
	// 	}, nil)
	// 	check(err)
	// 	fmt.Println("Device already present: ", *res.Payload.Name)
	// 	return
	// }

	// created, err := app.NetBoxclient.Dcim.DcimDevicesCreate(&dcim.DcimDevicesCreateParams{
	// 	Context: ctx,
	// 	Data:    devWithIDs,
	// }, nil)
	// check(err)

	// res, err := app.NetBoxclient.Dcim.DcimDevicesRead(&dcim.DcimDevicesReadParams{
	// 	ID:      created.Payload.ID,
	// 	Context: ctx,
	// }, nil)
	// check(err)

	// fmt.Println("Device created: ", *res.Payload.Name)
	//}

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
