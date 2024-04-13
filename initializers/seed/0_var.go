package seed

import "parameter-store-be/models"

var (
	defaultRoles = []models.Role{}
	// defaultPermissions  = []models.Permission{}
	defaultStages       = []models.Stage{}
	defaultEnvironments = []models.Environment{}

	//// sample
	sampleOrganizations = models.Organization{}
	sampleAdmin         = models.User{}
	sampleUsers1        = models.User{}
	sampleUsers2        = models.User{}
	sampleProjects      = []models.Project{}

	// test
	testProject  = models.Project{}
	testVersions = []models.Version{}
	testUser     = models.User{}
	testParams   = []models.Parameter{}
)
