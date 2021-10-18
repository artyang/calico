module github.com/kelseyhightower/confd

go 1.15

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/kelseyhightower/memkv v0.1.1
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.1
	github.com/projectcalico/api v0.0.0-20211015200158-15cbe046d41e
	github.com/projectcalico/libcalico-go v1.7.2-0.20211016010322-1d9f1da31ab7
	github.com/projectcalico/typha v0.7.3-0.20211016011846-0f06d9bab7c4
	github.com/sirupsen/logrus v1.4.2
	gopkg.in/go-playground/validator.v9 v9.28.0 // indirect
	k8s.io/api v0.22.0
	k8s.io/apimachinery v0.22.0
	k8s.io/client-go v0.22.0
)

replace github.com/sirupsen/logrus => github.com/projectcalico/logrus v1.0.4-calico
