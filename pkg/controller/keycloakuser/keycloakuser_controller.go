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
	"context"
	"fmt"
	"github.com/Nerzal/gocloak"
	showksv1beta1 "github.com/cloudnativedaysjp/showks-keycloak-user-operator/pkg/apis/showks/v1beta1"
	"github.com/cloudnativedaysjp/showks-keycloak-user-operator/pkg/keycloak"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller")
var finalizerName = "finalizer.keycloakuser.showks.cloudnativedays.jp"

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new KeyCloakUser Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	kcBasePath := os.Getenv("KEYCLOAK_BASE_PATH")
	kcUsername := os.Getenv("KEYCLOAK_USERNAME")
	kcPassword := os.Getenv("KEYCLOAK_PASSWORD")
	kcRealm := os.Getenv("KEYCLOAK_REALM")
	c, err := keycloak.NewClient(kcBasePath, kcUsername, kcPassword, kcRealm)
	if err != nil {
		return err
	}
	return add(mgr, newReconciler(mgr, c))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, kcClient keycloak.KeyCloakClientInterface) reconcile.Reconciler {
	return &ReconcileKeyCloakUser{
		Client:   mgr.GetClient(),
		scheme:   mgr.GetScheme(),
		kcClient: kcClient,
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("keycloakuser-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to KeyCloakUser
	err = c.Watch(&source.Kind{Type: &showksv1beta1.KeyCloakUser{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create
	// Uncomment watch a Deployment created by KeyCloakUser - change this for objects you create
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &showksv1beta1.KeyCloakUser{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileKeyCloakUser{}

// ReconcileKeyCloakUser reconciles a KeyCloakUser object
type ReconcileKeyCloakUser struct {
	client.Client
	scheme   *runtime.Scheme
	kcClient keycloak.KeyCloakClientInterface
}

// Reconcile reads that state of the cluster for a KeyCloakUser object and makes changes based on the state read
// and what is in the KeyCloakUser.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=showks.cloudnativedays.jp,resources=keycloakusers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=showks.cloudnativedays.jp,resources=keycloakusers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=secret,verbs=get
func (r *ReconcileKeyCloakUser) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	fmt.Println("Reconcile")
	// Fetch the KeyCloakUser instance
	instance := &showksv1beta1.KeyCloakUser{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get user", "NamespacedName", request.NamespacedName.String())
		return reconcile.Result{}, err
	}
	if instance.ObjectMeta.DeletionTimestamp.IsZero() {
		if err := r.setFinalizer(instance); err != nil {
			log.Error(err, "Failed to set finalizer", "user", instance.Spec.UserName)
			return reconcile.Result{}, err
		}
	} else {
		return r.runFinalizer(instance)
	}

	var user *gocloak.User
	param := gocloak.GetUsersParams{Username: instance.Spec.UserName}
	fmt.Printf("param: %+v\n", param)
	users, err := r.kcClient.GetUsers(instance.Spec.Realm, param)
	if err != nil {
		return reconcile.Result{}, err
	}
	fmt.Printf("users: %+v\n", users)
	if len(*users) == 0 {
		userParam := gocloak.User{
			Username: instance.Spec.UserName,
			Enabled:  true,
		}
		id, err := r.kcClient.CreateUser(instance.Spec.Realm, userParam)
		if err != nil {
			log.Error(err, "Failed to create user", "user", instance.Spec.UserName)
			return reconcile.Result{}, err
		}
		user, err = r.kcClient.GetUserByID(instance.Spec.Realm, id)
		if err != nil {
			log.Error(err, "Failed to get user by id", "id", id)
			return reconcile.Result{}, err
		}

		passwordSecret := corev1.Secret{}
		err = r.Get(context.TODO(), types.NamespacedName{Name: instance.Spec.PasswordSecretName, Namespace: instance.Namespace}, &passwordSecret)
		if err != nil {
			log.Error(err, "Failed to get secret", "user", types.NamespacedName{Name: instance.Spec.PasswordSecretName, Namespace: instance.Namespace}.String())
			return reconcile.Result{}, err
		}
		password := string(passwordSecret.Data["password"])

		err = r.kcClient.SetPassword(instance.Spec.Realm, id, password)
		if err != nil {
			log.Error(err, "Failed to set password", "user", instance.Spec.UserName)
			return reconcile.Result{}, err
		}
	} else {
		a := *users
		user = &a[0]
	}

	instance.Status.ID = user.ID

	if err := r.Status().Update(context.Background(), instance); err != nil {
		log.Error(err, "Failed to update instance", "user", instance.Spec.UserName)
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileKeyCloakUser) setFinalizer(instance *showksv1beta1.KeyCloakUser) error {
	fmt.Println("setFinalizer")
	if !containsString(instance.ObjectMeta.Finalizers, finalizerName) {
		instance.ObjectMeta.Finalizers = append(instance.ObjectMeta.Finalizers, finalizerName)
		if err := r.Update(context.Background(), instance); err != nil {
			return err
		}
	}

	return nil
}

func (r *ReconcileKeyCloakUser) runFinalizer(instannce *showksv1beta1.KeyCloakUser) (reconcile.Result, error) {
	fmt.Println("runFinalizer")
	if containsString(instannce.ObjectMeta.Finalizers, finalizerName) {
		if err := r.deleteExternalDependency(instannce); err != nil {
			fmt.Printf("Failed to delete: %s\n", err)
			return reconcile.Result{}, err
		}

		instannce.ObjectMeta.Finalizers = removeString(instannce.ObjectMeta.Finalizers, finalizerName)
		if err := r.Update(context.Background(), instannce); err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileKeyCloakUser) deleteExternalDependency(instance *showksv1beta1.KeyCloakUser) error {
	fmt.Println("deleteExternalDependency")
	kcRealm := os.Getenv("KEYCLOAK_REALM")
	return r.kcClient.DeleteUser(kcRealm, instance.Status.ID)
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}
