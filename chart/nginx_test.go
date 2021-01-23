package chart

import (
	"github.com/Waterdrips/helmunit/pkg/helmunit"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	"testing"
)

func Test_foo(t *testing.T) {
	want := v1.Ingress{
		TypeMeta:   metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "networking.k8s.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:                       "helmChart-nginx-example",
			Labels:                   map[string]string{"app.kubernetes.io/instance": "helmChart", "app.kubernetes.io/name": "nginx-example"} ,
		},
		Spec:       v1.IngressSpec{
			Rules: []v1.IngressRule{
				{Host: "chart-example.local" , IngressRuleValue: v1.IngressRuleValue{HTTP: &v1.HTTPIngressRuleValue{Paths: nil}}},
			},
		},
	}
	var got v1.Ingress

	if err := helmunit.Template("helmChart",
		"default",
		"./nginx-example",
		"templates/ingress.yaml",
		[]string{"./nginx-example/values.yaml"},
		nil,
		&got,
	); err != nil {
		t.Fatalf("got an error templating chart: %v", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("difference in chart and expected. \nGot:\n%+v\nWant:\n%+v\n", got, want)
	}
}
