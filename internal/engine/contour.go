package engine

import (
	"context"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	networking "k8s.io/client-go/kubernetes/typed/networking/v1"
	"k8s.io/client-go/rest"
	"strings"
)

const contourController = "projectcontour.io/snappcloud-ingress"

type Contour struct {
	k8sConfig rest.Config
}

func NewContour(k8s rest.Config) *Contour {
	return &Contour{
		k8sConfig: k8s,
	}
}

func (c *Contour) GetConfig() (map[string]interface{}, error) {
	networkingClientset, err := networking.NewForConfig(&c.k8sConfig)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	ingressClassList, err := networkingClientset.IngressClasses().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return extractContourIngressClass(ingressClassList.Items), nil
}

func extractContourIngressClass(iClassList []v1.IngressClass) map[string]interface{} {
	iClasses := make([]string, 0)

	for _, iClass := range iClassList {
		if strings.Contains(iClass.Spec.Controller, contourController) {
			iClasses = append(iClasses, iClass.Name)
		}
	}

	return map[string]interface{}{
		contourController: iClasses,
	}
}
