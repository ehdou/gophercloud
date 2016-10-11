package shares

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToShareCreateMap() (map[string]interface{}, error)
}

// CreateOpts used for ...
type CreateOpts struct {
	ShareProto         string            `json:"share_proto"`
	Size               int               `json:"size"`
	Name               string            `json:"name,omitempty"`
	Description        string            `json:"description,omitempty"`
	DisplayName        string            `json:"display_name,omitempty"`
	DisplayDescription string            `json:"display_description,omitempty"`
	ShareType          string            `json:"share_type,omitempty"`
	VolumeType         string            `json:"volume_type,omitempty"`
	SnapshotID         string            `json:"snapshot_id,omitempty"`
	IsPublic           bool              `json:"is_public,omitempty"`
	MetaData           map[string]string `json:"metadata,omitempty"`
	ShareNetworkID     string            `json:"share_network_id,omitempty"`
	ConsistencyGroupID string            `json:"consistency_group_id,omitempty"`
	AvailabilityZone   string            `json:"availability_zone,omitempty"`
}

// ToShareCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToShareCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "share")
}

// Create will create a new Share based on the values in CreateOpts. To extract
// the Share object from the response, call the Extract method on the
// CreateResult.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToShareCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

// Delete will delete an existing Share with the given UUID.
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, id), nil)
	return
}

// Get will get a single share with given UUID
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}

// ListOptsBuilder allows extensions to add additional parameters to the List request
type ListOptsBuilder interface {
	ToShareListQuery() (string, error)
}

// ListOpts is used to filter the List() call
type ListOpts struct {
	// Admin-only option to list shares from all tenants
	AllTenants bool `q:"tenant_id"`
	// Name of the share
	Name string `q:"name"`
	// Status of the share
	Status string `q:"status"`
	// ShareServerID is the UUID of the share server
	ShareServerID string `q:"share_server_id"`
	// MetaData is key-value-pairs of custom metadata
	MetaData map[string]string `q:"metadata"`
	// ShareTypeID is the UUID of the share type
	ShareTypeID string `q:"share_type_id"`
	// SortKey defines the sorting key
	SortKey string `q:"sort_key"`
	// SortDir defines the sorting direction
	SortDir string `q:"sort_dir"`
	// SnapShotID is the UUID of the snapshot used to create this share
	SnapShotID string `q:"snapshot_id"`
	// ShareNetworkID is the UUID of the share network of this share
	ShareNetworkID string `q:"share_network_id"`
	// ProjectID is the UUID of the project in which the share belongs to
	ProjectID string `q:"project_id"`
	// IsPublic determines the level of visibility of the share
	IsPublic bool `q:"is_public"`
	// ConsistencyGroupID is the UUID of the consistency group to which the share belongs to
	ConsistencyGroupID string `q:"consistency_group_id"`
}

// ToShareListQuery converts a foo to bar
func (l ListOpts) ToShareListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(l)
	return q.String(), err
}

// List returns shares filtered by the conditions provided in ListOpts. If ListOpts
// is not nil, returns a list of shares with details.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return SharePage{pagination.SinglePageBase(r)}
	})
}
