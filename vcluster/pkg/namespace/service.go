package namespace

import (
	"cmp"
	"context"
	"errors"
	"fmt"

	"github.com/loft-sh/workshop-kcd-munich-2024/vcluster/pkg/auth"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var ErrUserNotFound = fmt.Errorf("namespace: %w", auth.ErrNotFound)

var LabelOwner = "workshop.loft.sh/kcd-munich-owner"

type Service struct {
	Clientset *kubernetes.Clientset
	Config    *rest.Config
}

func (s *Service) Create(ctx context.Context, name string) (_ *corev1.Namespace, err error) {
	user, ok := auth.FromContext(ctx)
	if !ok {
		return nil, ErrUserNotFound
	}

	namespace, err := s.Clientset.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				LabelOwner: user.Name,
			},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	// clean up in case of an error
	defer func() {
		if err != nil {
			delErr := s.Clientset.CoreV1().Namespaces().Delete(ctx, namespace.Name, metav1.DeleteOptions{})
			if delErr != nil {
				err = errors.Join(err, delErr)
			}
		}
	}()

	// Create service account and role bindings
	sa, err := s.Clientset.CoreV1().ServiceAccounts(namespace.Name).Create(ctx, &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("sa-%v", user.Name),
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	role, err := s.Clientset.RbacV1().Roles(namespace.Name).Create(ctx, &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name: "power-user",
		},
		Rules: []rbacv1.PolicyRule{
			{
				Verbs:     []string{"*"},
				APIGroups: []string{"*"},
				Resources: []string{"*"},
			},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	_, err = s.Clientset.RbacV1().RoleBindings(namespace.Name).Create(ctx, &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{Name: "power-user-binding"},
		Subjects: []rbacv1.Subject{
			{
				APIGroup:  sa.GroupVersionKind().Group,
				Kind:      "ServiceAccount",
				Name:      sa.Name,
				Namespace: sa.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: role.GroupVersionKind().Group,
			Kind:     "Role",
			Name:     role.Name,
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	// Create service account token
	_, err = s.Clientset.CoreV1().Secrets(namespace.Name).Create(ctx, &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "power-user-token",
			Annotations: map[string]string{
				corev1.ServiceAccountNameKey: sa.Name,
			},
		},
		Type: corev1.SecretTypeServiceAccountToken,
	}, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return namespace, nil
}

func (s *Service) Delete(ctx context.Context, name string) error {
	user, ok := auth.FromContext(ctx)
	if !ok {
		return ErrUserNotFound
	}

	namespace, err := s.Clientset.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if !isNamespaceOwner(user, namespace) {
		return auth.ErrUnauthorized
	}

	err = s.Clientset.CoreV1().Namespaces().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) List(ctx context.Context) ([]corev1.Namespace, error) {
	user, ok := auth.FromContext(ctx)
	if !ok {
		return nil, ErrUserNotFound
	}

	lo := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%v=%v", LabelOwner, user.Name),
	}

	namespaces, err := s.Clientset.CoreV1().Namespaces().List(ctx, lo)
	if err != nil {
		return nil, err
	}

	return namespaces.Items, nil
}

func (s *Service) Get(ctx context.Context, name string) (*corev1.Namespace, error) {
	user, ok := auth.FromContext(ctx)
	if !ok {
		return nil, ErrUserNotFound
	}

	namespace, err := s.Clientset.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if !isNamespaceOwner(user, namespace) {
		return nil, auth.ErrUnauthorized
	}

	return namespace, nil
}

type KubeconfigOptions struct {
	Server   string
	CaData   []byte
	Insecure bool
}

func (s *Service) GetKubeconfig(ctx context.Context, name string, options KubeconfigOptions) ([]byte, error) {
	user, ok := auth.FromContext(ctx)
	if !ok {
		return nil, ErrUserNotFound
	}

	namespace, err := s.Clientset.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if !isNamespaceOwner(user, namespace) {
		return nil, auth.ErrUnauthorized
	}

	secret, err := s.Clientset.CoreV1().Secrets(namespace.Name).Get(ctx, "power-user-token", metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	caCert := secret.Data["ca.crt"]

	clusters := make(map[string]*clientcmdapi.Cluster)
	clusters["default-cluster"] = &clientcmdapi.Cluster{
		Server:                   cmp.Or(options.Server, s.Config.Host),
		CertificateAuthorityData: caCert,
	}

	if options.CaData != nil {
		clusters["default-cluster"].CertificateAuthorityData = options.CaData
	}

	if options.Insecure {
		clusters["default-cluster"].InsecureSkipTLSVerify = true
		clusters["default-cluster"].CertificateAuthorityData = nil
	}

	contexts := make(map[string]*clientcmdapi.Context)
	contexts["default-context"] = &clientcmdapi.Context{
		Cluster:   "default-cluster",
		Namespace: namespace.Name,
		AuthInfo:  namespace.Name,
	}

	authinfos := make(map[string]*clientcmdapi.AuthInfo)
	authinfos[namespace.Name] = &clientcmdapi.AuthInfo{
		Token: string(secret.Data["token"]),
	}

	clientConfig := clientcmdapi.Config{
		Kind:           "Config",
		APIVersion:     "v1",
		Clusters:       clusters,
		Contexts:       contexts,
		CurrentContext: "default-context",
		AuthInfos:      authinfos,
	}

	kubeconfig, err := clientcmd.Write(clientConfig)
	if err != nil {
		return nil, err
	}

	return kubeconfig, nil
}

// isNamespaceOwner returns true if the passed namespace contains
// a LabelOwner annotation that's equal to the passed user's name.
func isNamespaceOwner(user *auth.User, namespace *corev1.Namespace) bool {
	return namespace.Labels != nil && namespace.Labels[LabelOwner] == user.Name
}
