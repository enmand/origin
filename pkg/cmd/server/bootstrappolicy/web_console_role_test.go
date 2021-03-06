package bootstrappolicy

import (
	"strings"
	"testing"

	authorizationapi "github.com/openshift/origin/pkg/authorization/apis/authorization"
	"k8s.io/apimachinery/pkg/util/sets"
)

// NOTE: If this test fails, talk to the web console team to decide if your
// new role(s) should be visible to an end user in the web console.

var rolesToHide = sets.NewString(
	"cluster-admin",
	"cluster-debugger",
	"cluster-reader",
	"cluster-status",
	"registry-admin",
	"registry-editor",
	"registry-viewer",
	"self-access-reviewer",
	"self-provisioner",
	"storage-admin",
	"sudoer",
	"system:auth-delegator",
	"system:basic-user",
	"system:build-strategy-custom",
	"system:build-strategy-docker",
	"system:build-strategy-jenkinspipeline",
	"system:build-strategy-source",
	"system:discovery",
	"system:heapster",
	"system:image-auditor",
	"system:image-pruner",
	"system:image-signer",
	"system:kube-aggregator",
	"system:kube-controller-manager",
	"system:kube-dns",
	"system:kube-scheduler",
	"system:master",
	"system:node",
	"system:node-admin",
	"system:node-bootstrapper",
	"system:node-problem-detector",
	"system:node-proxier",
	"system:node-reader",
	"system:oauth-token-deleter",
	"system:openshift:template-service-broker",
	"system:openshift:templateservicebroker-client",
	"system:persistent-volume-provisioner",
	"system:registry",
	"system:router",
	"system:sdn-manager",
	"system:sdn-reader",
	"system:webhook",
)

func TestSystemOnlyRoles(t *testing.T) {
	show := sets.NewString()
	hide := sets.NewString()

	for _, role := range GetBootstrapClusterRoles() {
		if isControllerRole(&role) {
			continue // assume all controller roles can be ignored
		}
		if isSystemOnlyRole(role) {
			hide.Insert(role.Name)
		} else {
			show.Insert(role.Name)
		}
	}

	if !show.Equal(rolesToShow) || !hide.Equal(rolesToHide) {
		shouldNotShow := show.Difference(rolesToShow).List()
		shouldNotHide := hide.Difference(rolesToHide).List()
		t.Error("The list of expected end user roles has been changed.  Please discuss with the web console team to update role annotations.")
		if len(shouldNotShow) > 0 {
			t.Errorf("These roles are visible but not in rolesToShow: %v", shouldNotShow)
		}
		if len(shouldNotHide) > 0 {
			t.Errorf("These roles are hidden but not in rolesToHide: %v", shouldNotHide)
		}
	}
}

// this logic must stay in sync w/the web console for this test to be valid/valuable
// it is the same logic that is run on the membership page
func isSystemOnlyRole(role authorizationapi.ClusterRole) bool {
	return role.Annotations[roleSystemOnly] == roleIsSystemOnly
}

// helper so that roles following this pattern do not need to be manaully added
// to the hide list
func isControllerRole(role *authorizationapi.ClusterRole) bool {
	return strings.HasPrefix(role.Name, "system:controller:") ||
		strings.HasSuffix(role.Name, "-controller") ||
		strings.HasPrefix(role.Name, "system:openshift:controller:")
}
