package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/sharenetworks"
	"github.com/gophercloud/gophercloud/pagination"
)

func TestShareNetworkCreateDestroy(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create shared file system client: %v", err)
	}

	shareNetwork, err := CreateShareNetwork(t, client)
	if err != nil {
		t.Fatalf("Unable to create share network: %v", err)
	}

	newShareNetwork, err := sharenetworks.Get(client, shareNetwork.ID).Extract()
	if err != nil {
		t.Errorf("Unable to retrieve shareNetwork: %v", err)
	}

	if newShareNetwork.Name != shareNetwork.Name {
		t.Fatalf("Share network name was expeted to be: %s", shareNetwork.Name)
	}

	PrintShareNetwork(t, shareNetwork)

	defer DeleteShareNetwork(t, client, shareNetwork)
}

// Create a share network and update the name and description. Get the share
// network and verify that the name and description have been updated
func TestShareNetworkUpdate(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create shared file system client: %v", err)
	}

	shareNetwork, err := CreateShareNetwork(t, client)
	if err != nil {
		t.Fatalf("Unable to create share network: %v", err)
	}

	options := sharenetworks.UpdateOpts{
		Name:        "NewName",
		Description: "New share network description",
		NovaNetID:   "New_nova_network_id",
	}

	_, err = sharenetworks.Update(client, shareNetwork.ID, options).Extract()
	if err != nil {
		t.Errorf("Unable to update shareNetwork: %v", err)
	}

	newShareNetwork, err := sharenetworks.Get(client, shareNetwork.ID).Extract()
	if err != nil {
		t.Errorf("Unable to retrieve shareNetwork: %v", err)
	}

	if newShareNetwork.Name != options.Name {
		t.Fatalf("Share network name was expeted to be: %s", options.Name)
	}

	if newShareNetwork.Description != options.Description {
		t.Fatalf("Share network description was expeted to be: %s", options.Description)
	}

	if newShareNetwork.NovaNetID != options.NovaNetID {
		t.Fatalf("Share network Nova ID was expeted to be: %s", options.NovaNetID)
	}

	PrintShareNetwork(t, shareNetwork)

	defer DeleteShareNetwork(t, client, shareNetwork)
}

func TestShareNetworkListDetail(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create a shared file system client: %v", err)
	}

	allPages, err := sharenetworks.ListDetail(client, sharenetworks.ListOpts{}).AllPages()
	if err != nil {
		t.Fatalf("Unable to retrieve share networks: %v", err)
	}

	allShareNetworks, err := sharenetworks.ExtractShareNetworks(allPages)
	if err != nil {
		t.Fatalf("Unable to extract share networks: %v", err)
	}

	for _, shareNetwork := range allShareNetworks {
		PrintShareNetwork(t, &shareNetwork)
	}
}

// The test creates 2 shared networks and verifies that only the one(s) with
// a particular name are being listed
func TestShareNetworkListFiltering(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create a shared file system client: %v", err)
	}

	shareNetwork, err := CreateShareNetwork(t, client)
	if err != nil {
		t.Fatalf("Unable to create share network: %v", err)
	}
	defer DeleteShareNetwork(t, client, shareNetwork)

	shareNetwork, err = CreateShareNetwork(t, client)
	if err != nil {
		t.Fatalf("Unable to create share network: %v", err)
	}
	defer DeleteShareNetwork(t, client, shareNetwork)

	options := sharenetworks.ListOpts{
		Name: shareNetwork.Name,
	}

	allPages, err := sharenetworks.ListDetail(client, options).AllPages()
	if err != nil {
		t.Fatalf("Unable to retrieve share networks: %v", err)
	}

	allShareNetworks, err := sharenetworks.ExtractShareNetworks(allPages)
	if err != nil {
		t.Fatalf("Unable to extract share networks: %v", err)
	}

	for _, listedShareNetwork := range allShareNetworks {
		if listedShareNetwork.Name != shareNetwork.Name {
			t.Fatalf("The name of the share network was expected to be %s", shareNetwork.Name)
		}
		PrintShareNetwork(t, &listedShareNetwork)
	}
}

func TestShareNetworkListPagination(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create a shared file system client: %v", err)
	}

	shareNetwork, err := CreateShareNetwork(t, client)
	if err != nil {
		t.Fatalf("Unable to create share network: %v", err)
	}
	defer DeleteShareNetwork(t, client, shareNetwork)

	shareNetwork, err = CreateShareNetwork(t, client)
	if err != nil {
		t.Fatalf("Unable to create share network: %v", err)
	}
	defer DeleteShareNetwork(t, client, shareNetwork)

	count := 0

	err = sharenetworks.ListDetail(client, sharenetworks.ListOpts{Offset: 0, Limit: 1}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		_, err := sharenetworks.ExtractShareNetworks(page)
		if err != nil {
			t.Fatalf("Failed to extract share networks: %v", err)
			return false, err
		}

		return true, nil
	})
	if err != nil {
		t.Fatalf("Unable to retrieve share networks: %v", err)
	}

	if count < 2 {
		t.Fatal("Expected to get at least 2 pages")
	}

}
