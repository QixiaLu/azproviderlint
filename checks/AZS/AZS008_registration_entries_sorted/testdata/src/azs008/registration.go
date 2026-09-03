package azs008

import (
	"github.com/example/provider/framework"
	"github.com/example/provider/pluginsdk"
	"github.com/example/provider/sdk"
	"github.com/example/provider/typed"
)

type Registration struct{}

func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_availability_set":    nil,
		"azurerm_dedicated_host":      nil,
		"azurerm_disk_encryption_set": nil,
		"azurerm_managed_disk":        nil,
		"azurerm_virtual_machine":     nil,
	}
}

func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{}
}

func (r Registration) InvalidSupportedResources() map[string]*pluginsdk.Resource { // want `registration entries should be sorted alphabetically`
	return map[string]*pluginsdk.Resource{
		"azurerm_availability_set":    nil,
		"azurerm_dedicated_host":      nil,
		"azurerm_managed_disk":        nil,
		"azurerm_disk_encryption_set": nil,
		"azurerm_ssh_public_key":      nil,
	}
}

func (r Registration) SupportedResourcesViaVariable() map[string]*pluginsdk.Resource { // want `registration entries should be sorted alphabetically`
	lookup := map[string]string{
		"z": "last",
		"a": "first",
	}
	_ = lookup

	resources := map[string]*pluginsdk.Resource{
		"azurerm_availability_set":    nil,
		"azurerm_dedicated_host":      nil,
		"azurerm_managed_disk":        nil,
		"azurerm_disk_encryption_set": nil,
		"azurerm_ssh_public_key":      nil,
	}

	return resources
}

func (r Registration) SectionedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		// CDN
		"azurerm_cdn_profile": nil,

		// FrontDoor
		"azurerm_cdn_frontdoor_custom_domain": nil,
		"azurerm_cdn_frontdoor_endpoint":      nil,
		"azurerm_cdn_frontdoor_profile":       nil,
	}
}

func (r Registration) InvalidSectionedDataSources() map[string]*pluginsdk.Resource { // want `registration entries should be sorted alphabetically`
	return map[string]*pluginsdk.Resource{
		// CDN
		"azurerm_cdn_profile": nil,

		// FrontDoor
		"azurerm_cdn_frontdoor_profile":       nil,
		"azurerm_cdn_frontdoor_custom_domain": nil,
		"azurerm_cdn_frontdoor_endpoint":      nil,
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		ApiManagementResource{},
		WorkspaceResource{},
	}
}

func (r Registration) InvalidResources() []sdk.Resource { // want `registration entries should be sorted alphabetically`
	return []sdk.Resource{
		WorkspaceResource{},
		ApiManagementResource{},
	}
}

func (r Registration) ResourcesViaVariable() []sdk.Resource { // want `registration entries should be sorted alphabetically`
	resources := []sdk.Resource{
		WorkspaceResource{},
		ApiManagementResource{},
	}

	return resources
}

func (r Registration) InvalidPointerResources() []sdk.Resource { // want `registration entries should be sorted alphabetically`
	return []sdk.Resource{
		&WorkspaceResource{},
		&ApiManagementResource{},
	}
}

func (r Registration) QualifiedResources() []sdk.Resource {
	return []sdk.Resource{
		typed.ComputeResource{},
		typed.NetworkResource{},
	}
}

func (r Registration) InvalidQualifiedResources() []sdk.Resource { // want `registration entries should be sorted alphabetically`
	return []sdk.Resource{
		typed.NetworkResource{},
		typed.ComputeResource{},
	}
}

func (r Registration) FrameworkResources() []func() framework.Resource {
	return []func() framework.Resource{
		newApiManagementResource,
		newWorkspaceResource,
	}
}

func (r Registration) InvalidFrameworkResources() []func() framework.Resource { // want `registration entries should be sorted alphabetically`
	return []func() framework.Resource{
		newWorkspaceResource,
		newApiManagementResource,
	}
}

func (r Registration) WebsiteCategories() []string {
	return []string{
		"Compute",
		"Network",
	}
}

func (r Registration) InvalidWebsiteCategories() []string { // want `registration entries should be sorted alphabetically`
	return []string{
		"Network",
		"Compute",
	}
}

// CaseInsensitiveCategories is not reported: entries are ordered case-insensitively even though
// their casing differs.
func (r Registration) CaseInsensitiveCategories() []string {
	return []string{
		"apple",
		"Banana",
		"cherry",
	}
}

func (r Registration) InvalidCaseInsensitiveCategories() []string { // want `registration entries should be sorted alphabetically`
	return []string{
		"Cherry",
		"apple",
	}
}

func (r Registration) AttachedEntryComment() []string { // want `registration entries should be sorted alphabetically`
	return []string{
		"cherry",
		// zebra is special
		"zebra",
		"apple",
	}
}

// BlankLineCommentStartsSection: a comment with a blank line before it starts a new section even
// when a blank line also follows it, so these single-entry sections are each sorted and not flagged.
func (r Registration) BlankLineCommentStartsSection() []string {
	return []string{
		"zebra",

		// heading

		"apple",
	}
}

// UnresolvableEntriesNotReported has an entry whose sort key cannot be resolved, so the whole
// section is skipped rather than judged on the resolvable entries alone.
func (r Registration) UnresolvableEntriesNotReported() []sdk.Resource {
	return []sdk.Resource{
		WorkspaceResource{},
		makeResource("x"),
		ApiManagementResource{},
	}
}

func (r Registration) AppendedResourcesViaVariable() []sdk.Resource { // want `registration entries should be sorted alphabetically`
	resources := []sdk.Resource{
		WorkspaceResource{},
		ApiManagementResource{},
	}

	resources = append(resources, buildResources()...)

	return resources
}

// BraceSharingNotFixed reports the unsorted entries but offers no auto-fix because the opening and
// closing braces share entry lines and a whole-line rewrite would corrupt the source.
func (r Registration) BraceSharingNotFixed() []sdk.Resource { // want `registration entries should be sorted alphabetically`
	return []sdk.Resource{WorkspaceResource{},
		ApiManagementResource{}}
}

// PartiallyFixableSections leaves the same-line section unchanged and fixes the safe section.
func (r Registration) PartiallyFixableSections() []string { // want `registration entries should be sorted alphabetically`
	return []string{
		"zebra", "apple",

		"dog",
		"cat",
	}
}

// SpanningBlockCommentNotFixed cannot be safely reordered by whole source lines.
func (r Registration) SpanningBlockCommentNotFixed() []string { // want `registration entries should be sorted alphabetically`
	return []string{
		"zebra", /* note
		spanning */"apple",
	}
}

type ApiManagementResource struct{}
type WorkspaceResource struct{}

func (ApiManagementResource) ResourceType() string { return "azurerm_api_management" }
func (WorkspaceResource) ResourceType() string     { return "azurerm_workspace" }

func newApiManagementResource() framework.Resource { return nil }
func newWorkspaceResource() framework.Resource     { return nil }

func makeResource(string) sdk.Resource { return nil }

func buildResources() []sdk.Resource { return nil }
