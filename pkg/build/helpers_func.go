package build

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/xuxant/kbuild/pkg/constants"
	k8s "github.com/xuxant/kbuild/pkg/kubernetes"
	"github.com/xuxant/kbuild/pkg/output/log"
	"io"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"math/big"
	"time"
)

func addSecretVolume(pod *v1.Pod, secretName, dockerConfigPath, dockerSecretName string) {
	pod.Spec.Containers[0].VolumeMounts = append(pod.Spec.Containers[0].VolumeMounts, v1.VolumeMount{
		Name:      secretName,
		MountPath: dockerConfigPath,
	})

	pod.Spec.Volumes = append(pod.Spec.Volumes, v1.Volume{
		Name: secretName,
		VolumeSource: v1.VolumeSource{
			Secret: &v1.SecretVolumeSource{
				SecretName: dockerSecretName,
			},
		},
	})
}

func generateImageName(initial string) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	charsetLenght := big.NewInt(constants.RandomPodCharacter)

	randomString := make([]byte, constants.RandomPodCharacter)
	for i := range randomString {
		randomIndex, err := rand.Int(rand.Reader, charsetLenght)
		if err != nil {
			return "", err
		}
		randomString[i] = charset[randomIndex.Int64()]
	}

	return initial + "-" + string(randomString), nil

}

func (b *Builder) setupContext(pods corev1.PodInterface, name string) error {
	if err := k8s.WaitForPodInitialized(pods, name); err != nil {
		return fmt.Errorf("waiting for pod to initialize: %w", err)
	}

	attemp := 1
	maxAttemt := 5
	timeout, err := time.ParseDuration("5m")
	if err != nil {
		return fmt.Errorf("parsing timeout: %w", err)
	}

	ctx := context.TODO()
	ctx, cancleTimeOut := context.WithTimeout(ctx, timeout)

	defer cancleTimeOut()

	err = wait.PollUntilContextTimeout(ctx, time.Second, timeout*time.Duration(3), true, func(ctx context.Context) (bool, error) {
		if err := b.copyBuildContext(b.Target, name); err != nil {
			if errors.Is(ctx.Err(), context.Canceled) {
				return false, err
			}
			log.Entry(context.TODO()).Warnf("uploading build context failed. retrying (%d/%d): %v", attemp, maxAttemt, err)
			if attemp == maxAttemt {
				return false, err
			}
			attemp++
			return false, nil
		}
		return true, nil
	})

	return nil
}

func (b *Builder) copyBuildContext(artifact string, podName string) error {
	copyTimeOut, err := time.ParseDuration("5m")
	if err != nil {
		return err
	}

	ctx := context.TODO()
	ctx, cancleTimeOut := context.WithTimeout(ctx, copyTimeOut)

	defer cancleTimeOut()

	errs := make(chan error, 1)
	buildCtx, buildCtxWritter := io.Pipe()
	go func() {
		err := createTarContext(buildCtxWritter, b.BuildContext, b.Dockerfile)
		if err != nil {
			buildCtxWritter.CloseWithError(fmt.Errorf("creating build context: %w", err))
			errs <- err
			return
		}
		buildCtxWritter.Close()
	}()

}
