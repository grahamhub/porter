package gorm_test

import (
	"os"
	"testing"
	"time"

	"github.com/porter-dev/porter/internal/adapter"
	"github.com/porter-dev/porter/internal/config"
	"github.com/porter-dev/porter/internal/models"
	ints "github.com/porter-dev/porter/internal/models/integrations"
	"github.com/porter-dev/porter/internal/repository"
	"github.com/porter-dev/porter/internal/repository/gorm"
)

type tester struct {
	repo         *repository.Repository
	key          *[32]byte
	dbFileName   string
	initUsers    []*models.User
	initProjects []*models.Project
	initGRs      []*models.GitRepo
	initRegs     []*models.Registry
	initClusters []*models.Cluster
	initHRs      []*models.HelmRepo
	initInfras   []*models.Infra
	initReleases []*models.Release
	initInvites  []*models.Invite
	initCCs      []*models.ClusterCandidate
	initKIs      []*ints.KubeIntegration
	initBasics   []*ints.BasicIntegration
	initOIDCs    []*ints.OIDCIntegration
	initOAuths   []*ints.OAuthIntegration
	initGCPs     []*ints.GCPIntegration
	initAWSs     []*ints.AWSIntegration
}

func setupTestEnv(tester *tester, t *testing.T) {
	t.Helper()

	db, err := adapter.New(&config.DBConf{
		EncryptionKey: "__random_strong_encryption_key__",
		SQLLite:       true,
		SQLLitePath:   tester.dbFileName,
	})

	if err != nil {
		t.Fatalf("%v\n", err)
	}

	err = db.AutoMigrate(
		&models.Project{},
		&models.Role{},
		&models.User{},
		&models.Session{},
		&models.GitRepo{},
		&models.Registry{},
		&models.Release{},
		&models.HelmRepo{},
		&models.Cluster{},
		&models.ClusterCandidate{},
		&models.ClusterResolver{},
		&models.Infra{},
		&models.GitActionConfig{},
		&models.Invite{},
		&ints.KubeIntegration{},
		&ints.BasicIntegration{},
		&ints.OIDCIntegration{},
		&ints.OAuthIntegration{},
		&ints.GCPIntegration{},
		&ints.AWSIntegration{},
		&ints.ClusterTokenCache{},
		&ints.RegTokenCache{},
		&ints.HelmRepoTokenCache{},
		&ints.GithubAppInstallation{},
	)

	if err != nil {
		t.Fatalf("%v\n", err)
	}

	var key [32]byte

	for i, b := range []byte("__random_strong_encryption_key__") {
		key[i] = b
	}

	tester.key = &key

	tester.repo = gorm.NewRepository(db, &key)
}

func cleanup(tester *tester, t *testing.T) {
	t.Helper()

	// remove the created file file
	os.Remove(tester.dbFileName)
}

func initUser(tester *tester, t *testing.T) {
	t.Helper()

	user := &models.User{
		Email:    "example@example.com",
		Password: "hello1234",
	}

	user, err := tester.repo.User.CreateUser(user)

	if err != nil {
		t.Fatalf("%v\n", err)
	}

	tester.initUsers = append(tester.initUsers, user)
}

func initMultiUser(tester *tester, t *testing.T) {
	t.Helper()

	user := &models.User{
		Email:    "example@example.com",
		Password: "hello1234",
	}

	user, err := tester.repo.User.CreateUser(user)

	if err != nil {
		t.Fatalf("%v\n", err)
	}

	tester.initUsers = append(tester.initUsers, user)

	user = &models.User{
		Email:    "example2@example.com",
		Password: "hello1234",
	}

	user, err = tester.repo.User.CreateUser(user)

	if err != nil {
		t.Fatalf("%v\n", err)
	}

	tester.initUsers = append(tester.initUsers, user)
}

func initProject(tester *tester, t *testing.T) {
	t.Helper()

	proj := &models.Project{
		Name: "project-test",
	}

	proj, err := tester.repo.Project.CreateProject(proj)

	if err != nil {
		t.Fatalf("%v\n", err)
	}

	tester.initProjects = append(tester.initProjects, proj)
}

func initProjectRole(tester *tester, t *testing.T) {
	t.Helper()

	role := &models.Role{
		Kind:      models.RoleAdmin,
		UserID:    tester.initUsers[0].Model.ID,
		ProjectID: tester.initProjects[0].Model.ID,
	}

	role, err := tester.repo.Project.CreateProjectRole(tester.initProjects[0], role)

	if err != nil {
		t.Fatalf("%v\n", err)
	}
}

func initKubeIntegration(tester *tester, t *testing.T) {
	t.Helper()

	if len(tester.initProjects) == 0 {
		initProject(tester, t)
	}

	if len(tester.initUsers) == 0 {
		initUser(tester, t)
	}

	ki := &ints.KubeIntegration{
		Mechanism:  ints.KubeLocal,
		ProjectID:  tester.initProjects[0].ID,
		UserID:     tester.initUsers[0].ID,
		Kubeconfig: []byte("current-context: testing\n"),
	}

	ki, err := tester.repo.KubeIntegration.CreateKubeIntegration(ki)

	if err != nil {
		t.Fatalf("%v\n", err)
	}

	tester.initKIs = append(tester.initKIs, ki)
}

func initBasicIntegration(tester *tester, t *testing.T) {
	t.Helper()

	if len(tester.initProjects) == 0 {
		initProject(tester, t)
	}

	if len(tester.initUsers) == 0 {
		initUser(tester, t)
	}

	basic := &ints.BasicIntegration{
		ProjectID: tester.initProjects[0].ID,
		UserID:    tester.initUsers[0].ID,
		Username:  []byte("username"),
		Password:  []byte("password"),
	}

	basic, err := tester.repo.BasicIntegration.CreateBasicIntegration(basic)

	if err != nil {
		t.Fatalf("%v\n", err)
	}

	tester.initBasics = append(tester.initBasics, basic)
}

func initOIDCIntegration(tester *tester, t *testing.T) {
	t.Helper()

	if len(tester.initProjects) == 0 {
		initProject(tester, t)
	}

	if len(tester.initUsers) == 0 {
		initUser(tester, t)
	}

	oidc := &ints.OIDCIntegration{
		Client:       ints.OIDCKube,
		ProjectID:    tester.initProjects[0].ID,
		UserID:       tester.initUsers[0].ID,
		IssuerURL:    []byte("https://oidc.example.com"),
		ClientID:     []byte("exampleclientid"),
		ClientSecret: []byte("exampleclientsecret"),
		IDToken:      []byte("idtoken"),
		RefreshToken: []byte("refreshtoken"),
	}

	oidc, err := tester.repo.OIDCIntegration.CreateOIDCIntegration(oidc)

	if err != nil {
		t.Fatalf("%v\n", err)
	}

	tester.initOIDCs = append(tester.initOIDCs, oidc)
}

func initOAuthIntegration(tester *tester, t *testing.T) {
	t.Helper()

	if len(tester.initProjects) == 0 {
		initProject(tester, t)
	}

	if len(tester.initUsers) == 0 {
		initUser(tester, t)
	}

	oauth := &ints.OAuthIntegration{
		SharedOAuthModel: ints.SharedOAuthModel{
			ClientID:     []byte("exampleclientid"),
			AccessToken:  []byte("idtoken"),
			RefreshToken: []byte("refreshtoken"),
		},
		Client:    ints.OAuthGithub,
		ProjectID: tester.initProjects[0].ID,
		UserID:    tester.initUsers[0].ID,
	}

	oauth, err := tester.repo.OAuthIntegration.CreateOAuthIntegration(oauth)

	if err != nil {
		t.Fatalf("%v\n", err)
	}

	tester.initOAuths = append(tester.initOAuths, oauth)
}

func initGCPIntegration(tester *tester, t *testing.T) {
	t.Helper()

	if len(tester.initProjects) == 0 {
		initProject(tester, t)
	}

	if len(tester.initUsers) == 0 {
		initUser(tester, t)
	}

	gcp := &ints.GCPIntegration{
		ProjectID:    tester.initProjects[0].ID,
		UserID:       tester.initUsers[0].ID,
		GCPProjectID: "test-proj-123456",
		GCPUserEmail: "test@test.it",
		GCPKeyData:   []byte("{\"test\":\"key\"}"),
	}

	gcp, err := tester.repo.GCPIntegration.CreateGCPIntegration(gcp)

	if err != nil {
		t.Fatalf("%v\n", err)
	}

	tester.initGCPs = append(tester.initGCPs, gcp)
}

func initAWSIntegration(tester *tester, t *testing.T) {
	t.Helper()

	if len(tester.initProjects) == 0 {
		initProject(tester, t)
	}

	if len(tester.initUsers) == 0 {
		initUser(tester, t)
	}

	aws := &ints.AWSIntegration{
		ProjectID:          tester.initProjects[0].ID,
		UserID:             tester.initUsers[0].ID,
		AWSClusterID:       []byte("example-cluster-0"),
		AWSAccessKeyID:     []byte("accesskey"),
		AWSSecretAccessKey: []byte("secret"),
		AWSSessionToken:    []byte("optional"),
	}

	aws, err := tester.repo.AWSIntegration.CreateAWSIntegration(aws)

	if err != nil {
		t.Fatalf("%v\n", err)
	}

	tester.initAWSs = append(tester.initAWSs, aws)
}

func initClusterCandidate(tester *tester, t *testing.T) {
	t.Helper()

	if len(tester.initProjects) == 0 {
		initProject(tester, t)
	}

	cc := &models.ClusterCandidate{
		AuthMechanism:     models.AWS,
		ProjectID:         tester.initProjects[0].ID,
		CreatedClusterID:  0,
		Resolvers:         []models.ClusterResolver{},
		Name:              "cluster-test",
		Server:            "https://localhost",
		ContextName:       "context-test",
		AWSClusterIDGuess: []byte("example-cluster-0"),
		Kubeconfig:        []byte("current-context: testing\n"),
	}

	cc, err := tester.repo.Cluster.CreateClusterCandidate(cc)

	if err != nil {
		t.Fatalf("%v\n", err)
	}

	tester.initCCs = append(tester.initCCs, cc)
}

func initCluster(tester *tester, t *testing.T) {
	t.Helper()

	if len(tester.initProjects) == 0 {
		initProject(tester, t)
	}

	if len(tester.initKIs) == 0 {
		initKubeIntegration(tester, t)
	}

	cluster := &models.Cluster{
		ProjectID:                tester.initProjects[0].ID,
		Name:                     "cluster-test",
		Server:                   "https://localhost",
		KubeIntegrationID:        tester.initKIs[0].ID,
		CertificateAuthorityData: []byte("-----BEGIN"),
	}

	cluster, err := tester.repo.Cluster.CreateCluster(cluster)

	if err != nil {
		t.Fatalf("%v\n", err)
	}

	tester.initClusters = append(tester.initClusters, cluster)
}

func initGitRepo(tester *tester, t *testing.T) {
	t.Helper()

	if len(tester.initProjects) == 0 {
		initProject(tester, t)
	}

	if len(tester.initOAuths) == 0 {
		initOAuthIntegration(tester, t)
	}

	gr := &models.GitRepo{
		ProjectID:          tester.initProjects[0].ID,
		RepoEntity:         "porter-dev",
		OAuthIntegrationID: tester.initOAuths[0].ID,
	}

	gr, err := tester.repo.GitRepo.CreateGitRepo(gr)

	if err != nil {
		t.Fatalf("%v\n", err)
	}

	tester.initGRs = append(tester.initGRs, gr)
}

func initRegistry(tester *tester, t *testing.T) {
	t.Helper()

	if len(tester.initProjects) == 0 {
		initProject(tester, t)
	}

	reg := &models.Registry{
		ProjectID: tester.initProjects[0].ID,
		Name:      "registry-test",
	}

	reg, err := tester.repo.Registry.CreateRegistry(reg)

	if err != nil {
		t.Fatalf("%v\n", err)
	}

	tester.initRegs = append(tester.initRegs, reg)
}

func initHelmRepo(tester *tester, t *testing.T) {
	t.Helper()

	if len(tester.initProjects) == 0 {
		initProject(tester, t)
	}

	hr := &models.HelmRepo{
		Name:      "helm-repo-test",
		RepoURL:   "https://example-repo.com",
		ProjectID: tester.initProjects[0].Model.ID,
	}

	hr, err := tester.repo.HelmRepo.CreateHelmRepo(hr)

	if err != nil {
		t.Fatalf("%v\n", err)
	}

	tester.initHRs = append(tester.initHRs, hr)
}

func initInfra(tester *tester, t *testing.T) {
	t.Helper()

	if len(tester.initProjects) == 0 {
		initProject(tester, t)
	}

	infra := &models.Infra{
		Kind:      models.InfraECR,
		ProjectID: tester.initProjects[0].Model.ID,
		Status:    models.StatusCreated,
	}

	infra, err := tester.repo.Infra.CreateInfra(infra)

	if err != nil {
		t.Fatalf("%v\n", err)
	}

	tester.initInfras = append(tester.initInfras, infra)
}

func initInvite(tester *tester, t *testing.T) {
	t.Helper()

	if len(tester.initProjects) == 0 {
		initProject(tester, t)
	}

	expiry := time.Now().Add(24 * time.Hour)

	invite := &models.Invite{
		Token:     "abcd",
		Expiry:    &expiry,
		Email:     "testing@test.it",
		ProjectID: 1,
	}

	invite, err := tester.repo.Invite.CreateInvite(invite)

	if err != nil {
		t.Fatalf("%v\n", err)
	}

	tester.initInvites = append(tester.initInvites, invite)
}

func initRelease(tester *tester, t *testing.T) {
	t.Helper()

	if len(tester.initProjects) == 0 {
		initProject(tester, t)
	}

	release := &models.Release{
		Name:         "denver-meister-dakota",
		Namespace:    "default",
		ProjectID:    1,
		ClusterID:    1,
		WebhookToken: "abcdefgh",
	}

	release, err := tester.repo.Release.CreateRelease(release)

	if err != nil {
		t.Fatalf("%v\n", err)
	}

	tester.initReleases = append(tester.initReleases, release)
}
