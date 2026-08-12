package azs005

import (
	"github.com/example/provider/pluginsdk"
	"github.com/example/provider/sdk"
)

type Registration struct{}

func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	resources := map[string]*pluginsdk.Resource{
		"azurerm_covered_untyped": untypedResource(), // covered by SupportedDataSources
		"azurerm_cross_covered":   untypedResource(), // covered by the typed DataSources entry
		"azurerm_missing_untyped": untypedResource(), // want `resource "azurerm_missing_untyped" has no corresponding data source`
	}

	if featureFlag() {
		resources["azurerm_missing_conditional"] = untypedResource() // want `resource "azurerm_missing_conditional" has no corresponding data source`
	}

	return resources
}

func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_covered_untyped": untypedResource(),
	}
}

func (r Registration) Resources() []sdk.Resource {
	out := []sdk.Resource{
		CoveredTypedResource{},
		MissingTypedResource{}, // want `resource "azurerm_missing_typed" has no corresponding data source`
		RunCommandResource{},   // action-style resource — never reported
	}

	if featureFlag() {
		out = append(out, MissingAppendedResource{}) // want `resource "azurerm_missing_appended" has no corresponding data source`
	}

	return out
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		CoveredTypedDataSource{},
		CrossCoveredDataSource{},
	}
}

func (r Registration) FrameworkResources() []sdk.FrameworkWrappedResource {
	return []sdk.FrameworkWrappedResource{
		CoveredFrameworkResource{},
		MissingFrameworkResource{}, // want `resource "azurerm_missing_framework" has no corresponding data source`
	}
}

func (r Registration) FrameworkDataSources() []sdk.FrameworkWrappedDataSource {
	return []sdk.FrameworkWrappedDataSource{
		CoveredFrameworkDataSource{},
	}
}

func untypedResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{}
}

func featureFlag() bool {
	return false
}
