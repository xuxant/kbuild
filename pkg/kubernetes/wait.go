package kubernetes

import (
	"context"
	"fmt"
	"github.com/xuxant/kbuild/pkg/output/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"time"
)

func WaitForPodInitialized(pods corev1.PodInterface, podName string) error {
	log.Entry(context.TODO()).Infof("Waiting for %s to be initialized", podName)

	w, err := pods.Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("initializing pod watcher: %s", err)
	}

	return watchUntilTimeOut(10*time.Minute, w, func(event *watch.Event) (bool, error) {
		pod := event.Object.(*v1.Pod)
		if pod.Name != podName {
			return false, nil
		}

		for _, ic := range pod.Status.InitContainerStatuses {
			if ic.State.Running != nil {
				return true, nil
			}
		}
		return false, nil
	})
}

func watchUntilTimeOut(timeout time.Duration,
	w watch.Interface,
	condition func(event *watch.Event) (bool, error),
) error {
	ctx := context.TODO()
	ctx, cancleTimeOut := context.WithTimeout(ctx, timeout)

	defer cancleTimeOut()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case event := <-w.ResultChan():
			done, err := condition(&event)
			if err != nil {
				return err
			}
			if done {
				return nil
			}
		}
	}

	return nil
}
