package testing

import (
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
	"testing"
)

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockCreateResponse(t)

	options := &shares.CreateOpts{Size: 1, Name: "my_test_share", ShareProto: "NFS"}
	n, err := shares.Create(client.ServiceClient(), options).Extract()

	th.AssertNoErr(t, err)
	th.AssertEquals(t, n.Name, "my_test_share")
	th.AssertEquals(t, n.Size, 1)
	th.AssertEquals(t, n.ShareProto, "NFS")
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockDeleteResponse(t)

	result := shares.Delete(client.ServiceClient(), shareID)
	th.AssertNoErr(t, result.Err)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockGetResponse(t)

	s, err := shares.Get(client.ServiceClient(), shareID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, s.ID, shareID)
}

func TestListAllShort(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockListResponse(t)
	pages, err := shares.List(client.ServiceClient(), &shares.ListOpts{}).AllPages()
	th.AssertNoErr(t, err)
	act, err := shares.ExtractShares(pages)
	th.AssertNoErr(t, err)
	shortList := []shares.Share{
		{
			ID:   "d94a8548-2079-4be0-b21c-0a887acd31ca",
			Name: "My_share",
			Links: []map[string]string{
				{
					"href": "http://172.18.198.54:8786/v1/16e1ab15c35a457e9c2b2aa189f544e1/shares/d94a8548-2079-4be0-b21c-0a887acd31ca",
					"rel":  "self",
				},
				{
					"href": "http://172.18.198.54:8786/16e1ab15c35a457e9c2b2aa189f544e1/shares/d94a8548-2079-4be0-b21c-0a887acd31ca",
					"rel":  "bookmark",
				},
			},
		},
		{
			ID:   "406ea93b-32e9-4907-a117-148b3945749f",
			Name: "Share1",
			Links: []map[string]string{
				{
					"href": "http://172.18.198.54:8786/16e1ab15c35a457e9c2b2aa189f544e1/shares/d94a8548-2079-4be0-b21c-0a887acd31ca",
					"rel":  "bookmark",
				},
			},
		},
	}
	th.CheckDeepEquals(t, shortList, act)
}
