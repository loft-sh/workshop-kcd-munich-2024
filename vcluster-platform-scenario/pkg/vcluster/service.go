package vcluster

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/loft-sh/workshop-kcd-munich-2024/vcluster-platform-scenario/pkg/auth"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	DefaultVersion = "0.19.6"
	DefaultName    = "vcluster"
)

type Service struct {
	Clientset *kubernetes.Clientset
}

type VClusterInstallOptions struct {
	Values        map[string]any
	Namespace     string
	Name          string
	Version       string
	UseLocalChart bool
	Wait          bool
}

type VClusterHelmRelease struct {
	Name    string
	Version string
}

func (s *Service) Install(ctx context.Context, options VClusterInstallOptions) (_ VClusterHelmRelease, err error) {
	_, ok := auth.FromContext(ctx)
	if !ok {
		return VClusterHelmRelease{}, auth.ErrNotFound
	}

	name := cmp.Or(options.Name, DefaultName)

	// TODO(ThomasK33): Check if installation already exists --> if true, check ownership

	helmArgs := []string{
		"upgrade", "--install", name,
	}

	if options.UseLocalChart {
		// (ThomasK33): Install from local chart
		helmArgs = append(helmArgs, "./vcluster-oss/chart")
	} else {
		helmArgs = append(helmArgs,
			"vcluster",
			"--repo", "https://charts.loft.sh/",
			"--repository-config", "''",
		)
	}

	// TODO(ThomasK33): Set user as owner of vCluster

	// (ThomasK33): Add passed options to exec command as file
	f, err := os.CreateTemp("", fmt.Sprintf("values-%v-%v-", options.Namespace, name))
	if err != nil {
		return VClusterHelmRelease{}, err
	}
	defer os.Remove(f.Name())

	if options.Values != nil {
		if err := yaml.NewEncoder(f).Encode(options.Values); err != nil {
			return VClusterHelmRelease{}, err
		}
	}

	helmArgs = append(helmArgs, "--values", f.Name(), "--reuse-values")

	helmArgs = append(helmArgs, "--namespace", options.Namespace)

	version := cmp.Or(options.Version, DefaultVersion)
	helmArgs = append(helmArgs, "--version", version)

	if options.Wait {
		helmArgs = append(helmArgs, "--wait")
	}

	cmd := exec.CommandContext(
		ctx,
		"helm",
		helmArgs...,
	)
	output, err := cmd.CombinedOutput()

	defer func() {
		// Perform cleanup
		if err != nil {
			if uninstallErr := s.Uninstall(ctx, options.Namespace, name, options.Wait); uninstallErr != nil {
				err = errors.Join(err, uninstallErr)
			}
		}
	}()

	if err != nil {
		slog.Error("an error occurred while running the helm install command", "cmd", "helm "+strings.Join(helmArgs, " "), "err", err, "output", string(output))
		return VClusterHelmRelease{}, fmt.Errorf("helm install: %w", err)
	}

	release, _, err := s.Get(ctx, options.Namespace, name)
	if err != nil {
		return VClusterHelmRelease{}, err
	}

	return VClusterHelmRelease{
		Name:    name,
		Version: release.Version,
	}, nil
}

func (s *Service) Uninstall(ctx context.Context, namespace, name string, wait bool) error {
	_, ok := auth.FromContext(ctx)
	if !ok {
		return auth.ErrNotFound
	}

	name = cmp.Or(name, DefaultName)

	waitStr := ""
	if wait {
		waitStr = "--wait"
	}

	// TODO(ThomasK33): Check if user is owner of the vCluster

	cmd := exec.CommandContext(
		ctx,
		"helm",
		"uninstall", name, "-n", namespace, waitStr,
	)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

type HelmGetMetadata struct {
	DeployedAt time.Time `json:"deployedAt"`
	Name       string    `json:"name"`
	Chart      string    `json:"chart"`
	Version    string    `json:"version"`
	AppVersion string    `json:"appVersion"`
	Namespace  string    `json:"namespace"`
	Status     string    `json:"status"`
	Revision   int64     `json:"revision"`
}

func (s *Service) Get(ctx context.Context, namespace, name string) (HelmGetMetadata, map[string]any, error) {
	_, ok := auth.FromContext(ctx)
	if !ok {
		return HelmGetMetadata{}, nil, auth.ErrNotFound
	}

	name = cmp.Or(name, DefaultName)

	// TODO(ThomasK33): Check if user is owner of the vCluster

	output, err := exec.CommandContext(ctx, "helm", "get", "metadata", "-n", namespace, name, "-o", "json").CombinedOutput()
	if err != nil {
		return HelmGetMetadata{}, nil, fmt.Errorf("%w: %v", err, string(output))
	}

	metadata := HelmGetMetadata{}
	if err := json.Unmarshal(output, &metadata); err != nil {
		return HelmGetMetadata{}, nil, fmt.Errorf("unmarshal failed: %w", err)
	}

	valuesOutput, err := exec.CommandContext(ctx, "helm", "get", "values", "-n", namespace, name, "-o", "json").CombinedOutput()
	if err != nil {
		return HelmGetMetadata{}, nil, err
	}

	values := map[string]any{}
	if len(valuesOutput) != 0 {
		if err := json.Unmarshal(valuesOutput, &values); err != nil {
			return HelmGetMetadata{}, nil, err
		}
	}

	return metadata, values, nil
}

func (s *Service) Kubeconfig(ctx context.Context, namespace, name string) ([]byte, error) {
	_, ok := auth.FromContext(ctx)
	if !ok {
		return nil, auth.ErrNotFound
	}

	name = cmp.Or(name, DefaultName)

	// TODO(ThomasK33): Check if user is owner of the vCluster

	secret, err := s.Clientset.CoreV1().Secrets(namespace).Get(ctx, fmt.Sprintf("vc-%v", name), metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return secret.Data["config"], nil
}
