package add_qihoo_metadata

import (
	"fmt"
	"strings"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/logp"
	"github.com/elastic/beats/v7/libbeat/processors"
)

const (
	processorName                          = "add_qihoo_metadata"
	keyKubernetesPodName                   = "kubernetes.pod.name"
	keyKubernetesNamespace                 = "kubernetes.namespace"
	keyKubernetesContainerName             = "kubernetes.container.name"
	keyKubernetesAnnotationsControllerKind = "kubernetes.annotations.qihoo.cloud/controller-kind"
	keyKubernetesLabelAppName              = "kubernetes.labels.app"
)

func init() {
	processors.RegisterPlugin(processorName, newQihooMetadataProcessor)
}

type addQihooMetadata struct {
	log *logp.Logger
}

func newQihooMetadataProcessor(cfg *common.Config) (processors.Processor, error) {
	return &addQihooMetadata{logp.NewLogger(processorName)}, nil
}

func (d *addQihooMetadata) Run(event *beat.Event) (*beat.Event, error) {
	defaultControllerKind := "deployment"
	defaultTopicInfix := "docker"
	defaultDeploymentName := "none"
	defaultAppName := "none"

	kubernetesAnnotationsControllerKind, err := event.Fields.GetValue(keyKubernetesAnnotationsControllerKind)
	if err != nil {
		d.log.Debugf("Error while get %s fields. %s ,err %v", keyKubernetesAnnotationsControllerKind, event.Fields.String(),err)
	} else {
		defaultControllerKind = strings.ToLower(kubernetesAnnotationsControllerKind.(string))
	}

	if defaultControllerKind != "deployment" {
		defaultTopicInfix = defaultControllerKind
	}
	event.Fields["controller_kind"] = defaultControllerKind

	kubernetesLabelAppName, err := event.Fields.GetValue(keyKubernetesLabelAppName)
	if err != nil {
		d.log.Debugf("Error while get %s fields. %s ,err %v", keyKubernetesLabelAppName, event.Fields.String(),err)
	} else {
		defaultAppName = kubernetesLabelAppName.(string)
		if defaultControllerKind == "deployment" {
			defaultDeploymentName = defaultAppName
		}
	}

	// 为了兼容之前的deployment 和topic字段
	// 非 deployment 或者无法取到app label的deployment为none
	// topic 如果类型为deployment 则为 k8s_docker_{{deploymentName}}, 其他类型 则为 k8s_{{controllerKind}}_{{appName}}
	event.Fields["deployment"] = defaultDeploymentName
	event.Fields["topic"] = fmt.Sprintf("k8s_%s_%s", defaultTopicInfix, defaultAppName)
	event.Fields["app"] = defaultAppName

	kubernetesPodName, err := event.Fields.GetValue(keyKubernetesPodName)
	if err != nil {
		d.log.Debugf("Error while get %s fields. %s", keyKubernetesPodName, event.Fields.String())
	} else {
		event.Fields["pod"] = kubernetesPodName
	}

	kubernetesNamespace, err := event.Fields.GetValue(keyKubernetesNamespace)
	if err != nil {
		d.log.Debugf("Error while get %s fields. %s ", keyKubernetesNamespace, event.Fields.String())
	} else {
		event.Fields["namespace"] = kubernetesNamespace
	}

	kubernetesContainerName, err := event.Fields.GetValue(keyKubernetesContainerName)
	if err != nil {
		d.log.Debugf("Error while get %s fields. %s", keyKubernetesContainerName, event.Fields.String())
	} else {
		event.Fields["container_name"] = kubernetesContainerName
	}

	if value, err := event.Fields.GetValue("log.file.path"); err == nil {
		event.Fields["source"] = value
	}

	if value, err := event.Fields.GetValue("message2.log"); err == nil {
		event.Fields["log"] = value
	}


	//fmt.Println(event.Fields.String())

	return event, nil
}

func (d *addQihooMetadata) String() string {
	return processorName
}
