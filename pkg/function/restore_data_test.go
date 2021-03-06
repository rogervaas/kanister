package function

import (
	. "gopkg.in/check.v1"

	"github.com/kanisterio/kanister/pkg/param"
)

type RestoreDataTestSuite struct{}

var _ = Suite(&RestoreDataTestSuite{})

func newValidDeploymentTP() param.TemplateParams {
	return param.TemplateParams{
		Deployment: &param.DeploymentParams{
			Name:      "test-deployment",
			Namespace: "test-namespace",
			Pods: []string{
				"pod1",
				"pod2",
			},
			Containers: [][]string{{"test-container"}},
			PersistentVolumeClaims: map[string]map[string]string{
				"pod1": map[string]string{
					"pvc1": "path1",
				},
				"pod2": map[string]string{
					"pvc2": "path2",
				},
			},
		},
	}
}

func newInvalidDeploymentTP() param.TemplateParams {
	return param.TemplateParams{
		Deployment: &param.DeploymentParams{
			Name:      "test-deployment",
			Namespace: "test-namespace",
			Pods: []string{
				"pod1",
				"pod2",
			},
			Containers: [][]string{{"test-container"}},
			PersistentVolumeClaims: map[string]map[string]string{
				"pod1": map[string]string{
					"pvc1": "path1",
				},
			},
		},
	}
}

func newValidStatefulSetTP() param.TemplateParams {
	return param.TemplateParams{
		StatefulSet: &param.StatefulSetParams{
			Name:      "test-ss",
			Namespace: "test-namespace",
			Pods: []string{
				"pod1",
				"pod2",
			},
			Containers: [][]string{{"test-container"}},
			PersistentVolumeClaims: []map[string]string{
				map[string]string{
					"pvc1": "path1",
				},
				map[string]string{
					"pvc2": "path2",
					"pvc3": "path3",
				},
			},
		},
	}
}

func (s *RestoreDataTestSuite) TestFetchPodVolumesTest(c *C) {
	testCases := []struct {
		name       string
		tp         param.TemplateParams
		pod        string
		vols       map[string]string
		errChecker Checker
	}{
		{"Valid Deployment Pod", newValidDeploymentTP(), "pod1", map[string]string{"pvc1": "path1"}, IsNil},
		{"Valid StatefulSet Pod", newValidStatefulSetTP(), "pod2", map[string]string{"pvc2": "path2", "pvc3": "path3"}, IsNil},
		{"Invalid Deployment Pod", newValidDeploymentTP(), "pod3", nil, NotNil},
		{"Invalid StatefulSet Pod", newValidStatefulSetTP(), "pod4", nil, NotNil},
		{"Deployment Pod with no volumes", newInvalidDeploymentTP(), "pod2", nil, NotNil},
		{"Invalid Template Params", param.TemplateParams{}, "pod1", nil, NotNil},
	}
	for _, tc := range testCases {
		vols, err := fetchPodVolumes(tc.pod, tc.tp)
		c.Check(err, tc.errChecker, Commentf("Test: %s Failed!", tc.name))
		c.Check(vols, DeepEquals, tc.vols, Commentf("Test: %s Failed!", tc.name))
	}
}

func (s *RestoreDataTestSuite) TestValidateOptArgs(c *C) {
	testCases := []struct {
		name       string
		pod        string
		vols       map[string]string
		errChecker Checker
	}{
		{"Pod with nil Vols", "some-pod", nil, IsNil},
		{"Pod with empty Vols", "some-pod", map[string]string{}, IsNil},
		{"Vols with empty Pod", "", map[string]string{"pvc": "mount"}, IsNil},
		{"Empty Pod and Vols", "", map[string]string{}, NotNil},
		{"Empty Pod and nil Vols", "", nil, NotNil},
		{"Pod with Vols", "some-pod", map[string]string{"pvc": "mount"}, NotNil},
	}
	for _, tc := range testCases {
		err := validateOptArgs(tc.pod, tc.vols)
		c.Check(err, tc.errChecker, Commentf("Test: %s Failed!", tc.name))
	}
}
