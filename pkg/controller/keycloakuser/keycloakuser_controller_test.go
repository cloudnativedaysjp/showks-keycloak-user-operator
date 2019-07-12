/*
Copyright 2019 TAKAISHI Ryo.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package keycloakuser

import (
	"github.com/Nerzal/gocloak"
	"github.com/cloudnativedaysjp/showks-keycloak-user-operator/pkg/keycloak"
	"github.com/cloudnativedaysjp/showks-keycloak-user-operator/pkg/mock"
	"github.com/golang/mock/gomock"
	corev1 "k8s.io/api/core/v1"
	"os"
	"testing"
	"time"

	showksv1beta1 "github.com/cloudnativedaysjp/showks-keycloak-user-operator/pkg/apis/showks/v1beta1"
	"github.com/onsi/gomega"
	"golang.org/x/net/context"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var c client.Client

var expectedRequest = reconcile.Request{NamespacedName: types.NamespacedName{Name: "foo", Namespace: "default"}}
var userKey = types.NamespacedName{Name: "foo", Namespace: "default"}

const timeout = time.Second * 5

var realm = "master"
var userName = "dummyUser"
var userID = "DUMMYDUMMY"
var userPassword = "DUMMYPASSWORD"

func newKeyCloakClientMock(controller *gomock.Controller) keycloak.KeyCloakClientInterface {
	c := mock_keycloak.NewMockKeyCloakClientInterface(controller)
	param := gocloak.GetUsersParams{Username: userName}
	emptyUsers := &[]gocloak.User{}
	users := &[]gocloak.User{
		{
			Username: userName,
			ID:       userID,
		},
	}
	first := c.EXPECT().GetUsers(realm, param).Return(emptyUsers, nil).Times(1)
	c.EXPECT().GetUsers(realm, param).Return(users, nil).After(first).Times(1)
	userParam := gocloak.User{
		Username: userName,
		Enabled:  true,
	}
	c.EXPECT().CreateUser(realm, userParam).Return(userID, nil).Times(1)
	c.EXPECT().SetPassword(realm, userID, userPassword).Return(nil).Times(1)
	user := &gocloak.User{Username: userName, ID: userID}
	c.EXPECT().GetUserByID(realm, userID).Return(user, nil).Times(1)

	return c
}

func TestReconcile(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	instance := &showksv1beta1.KeyCloakUser{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: "default",
		},
		Spec: showksv1beta1.KeyCloakUserSpec{
			UserName: userName,
			Realm:    realm,
		},
	}

	passwordSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "default",
		},
		Data: map[string][]byte{
			"password": []byte(userPassword),
		},
	}
	// Setup the Manager and Controller.  Wrap the Controller Reconcile function so it writes each request to a
	// channel when it is finished.
	mgr, err := manager.New(cfg, manager.Options{})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	c = mgr.GetClient()

	os.Setenv("KEYCLOAK_REALM", realm)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	kcClient := newKeyCloakClientMock(mockCtrl)

	recFn, requests := SetupTestReconcile(newReconciler(mgr, kcClient))
	g.Expect(add(mgr, recFn)).NotTo(gomega.HaveOccurred())

	stopMgr, mgrStopped := StartTestManager(mgr, g)

	defer func() {
		close(stopMgr)
		mgrStopped.Wait()
	}()

	err = c.Create(context.TODO(), passwordSecret)
	// The instance object may not be a valid object because it might be missing some required fields.
	// Please modify the instance object by adding required fields and then remove the following if statement.
	if apierrors.IsInvalid(err) {
		t.Logf("failed to create object, got an invalid object error: %v", err)
		return
	}

	// Create the KeyCloakUser object and expect the Reconcile and Deployment to be created
	err = c.Create(context.TODO(), instance)
	// The instance object may not be a valid object because it might be missing some required fields.
	// Please modify the instance object by adding required fields and then remove the following if statement.
	if apierrors.IsInvalid(err) {
		t.Logf("failed to create object, got an invalid object error: %v", err)
		return
	}
	g.Expect(err).NotTo(gomega.HaveOccurred())
	defer c.Delete(context.TODO(), instance)
	g.Eventually(requests, timeout).Should(gomega.Receive(gomega.Equal(expectedRequest)))

	kcUser := &showksv1beta1.KeyCloakUser{}
	g.Eventually(func() error { return c.Get(context.TODO(), userKey, kcUser) }, timeout).
		Should(gomega.Succeed())
	g.Expect(kcUser.Status.ID).To(gomega.Equal(userID))

	// Delete the Deployment and expect Reconcile to be called for Deployment deletion
	//g.Expect(c.Delete(context.TODO(), deploy)).NotTo(gomega.HaveOccurred())
	//g.Eventually(requests, timeout).Should(gomega.Receive(gomega.Equal(expectedRequest)))
	//g.Eventually(func() error { return c.Get(context.TODO(), depKey, deploy) }, timeout).
	//	Should(gomega.Succeed())

	// Manually delete Deployment since GC isn't enabled in the test control plane
	//g.Eventually(func() error { return c.Delete(context.TODO(), deploy) }, timeout).
	//	Should(gomega.MatchError("deployments.apps \"foo-deployment\" not found"))

}
