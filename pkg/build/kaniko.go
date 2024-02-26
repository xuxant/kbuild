package build

import (
	"context"
	"fmt"
	"github.com/xuxant/kbuild/pkg/constants"
	"github.com/xuxant/kbuild/pkg/kaniko"
	k8s "github.com/xuxant/kbuild/pkg/kubernetes"
	"github.com/xuxant/kbuild/pkg/options"
	"github.com/xuxant/kbuild/pkg/output/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func InitializeBuild(tag, context, dockerfile string) error {
	var build Builder
	cfg := options.GetConfig()

	randStr, err := generateImageName("kaniko")
	if err != nil {
		return err
	}
	build.PodName = randStr
	build.BuildContext = context
	build.Namespace = cfg.Namespace
	build.Dockerfile = dockerfile
	build.Tag = tag

	err = build.InitializeKaniko(cfg)
	if err != nil {
		return err
	}

	return nil
}

func (b *Builder) InitializeKaniko(cfg options.RGlobalOptions) error {

	client, err := k8s.DefaultClient()
	if err != nil {
		return fmt.Errorf("getting kubernetes client: %w", err)
	}
	pods := client.CoreV1().Pods(cfg.Namespace)

	podSpec := b.podSpec()

	pod, err := pods.Create(context.TODO(), podSpec, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("creating kaniko pod: %w", err)
	}
	log.Entry(context.TODO()).Infof("created kaniko pod: %s", pod.Name)

	if err := b.setupContext(pods, pod.Name); err != nil {
		return fmt.Errorf("copying sources: %w", err)
	}
	return nil
}

func (b *Builder) podSpec() *v1.Pod {
	cfg := options.GetConfig()
	args := kaniko.Args(b.Tag, b.BuildContext)

	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				"build": "kaniko",
			},
			GenerateName: "kaniko-",
			Labels:       map[string]string{"build": "kaniko"},
			Namespace:    cfg.Namespace,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{{
				Name:            "kanilo",
				Image:           constants.DefaultKanikoImage,
				ImagePullPolicy: v1.PullIfNotPresent,
				Args:            args,
			}},
			RestartPolicy: v1.RestartPolicyNever,
		},
	}
	addSecretVolume(pod, "dockersecret", "/home/.docker/config.json", "docker")
	return pod
}
