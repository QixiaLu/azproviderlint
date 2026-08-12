package azs005

// typed resources and data sources — correlated via ResourceType()

type CoveredTypedResource struct{}

func (r CoveredTypedResource) ResourceType() string {
	return "azurerm_covered_typed"
}

type CoveredTypedDataSource struct{}

func (r CoveredTypedDataSource) ResourceType() string {
	return "azurerm_covered_typed"
}

// CrossCoveredDataSource covers an untyped resource of the same name, exercising cross-flavour
// matching.
type CrossCoveredDataSource struct{}

func (r CrossCoveredDataSource) ResourceType() string {
	return "azurerm_cross_covered"
}

type MissingTypedResource struct{}

func (r MissingTypedResource) ResourceType() string {
	return "azurerm_missing_typed"
}

// RunCommandResource is an invoke-style resource; its name matches the action-style suffix
// list, so it is never reported despite having no data source.
type RunCommandResource struct{}

func (r RunCommandResource) ResourceType() string {
	return "azurerm_thing_run_command"
}

type MissingAppendedResource struct{}

func (r MissingAppendedResource) ResourceType() string {
	return "azurerm_missing_appended"
}

// framework resources and data sources — same ResourceType() correlation

type CoveredFrameworkResource struct{}

func (r CoveredFrameworkResource) ResourceType() string {
	return "azurerm_covered_framework"
}

type CoveredFrameworkDataSource struct{}

func (r CoveredFrameworkDataSource) ResourceType() string {
	return "azurerm_covered_framework"
}

type MissingFrameworkResource struct{}

func (r MissingFrameworkResource) ResourceType() string {
	return "azurerm_missing_framework"
}
