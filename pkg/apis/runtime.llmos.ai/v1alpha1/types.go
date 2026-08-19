package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RuntimeComponentType identifies the responsibility of a runtime component.
type RuntimeComponentType string

const (
	RuntimeComponentTypeBase        RuntimeComponentType = "Base"
	RuntimeComponentTypeDistributed RuntimeComponentType = "DistributedRuntime"
	RuntimeComponentTypeEngine      RuntimeComponentType = "InferenceEngine"
	RuntimeComponentTypeAdapter     RuntimeComponentType = "ModelAdapter"
)

// RuntimeComponent describes an independently versioned runtime building block.
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=rtc,scope=Cluster
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.phase"
type RuntimeComponent struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RuntimeComponentSpec   `json:"spec,omitempty"`
	Status RuntimeComponentStatus `json:"status,omitempty"`
}

type RuntimeComponentSpec struct {
	// +kubebuilder:validation:Enum=Base;DistributedRuntime;InferenceEngine;ModelAdapter
	Type RuntimeComponentType `json:"type"`
	Name string `json:"name"`
	Version string `json:"version"`
	Artifact RuntimeArtifact `json:"artifact"`
	Compatibility RuntimeCompatibility `json:"compatibility,omitempty"`
	Capabilities []string `json:"capabilities,omitempty"`
}

type RuntimeArtifact struct {
	// OCI image containing this component. Initially this can be used by an initContainer;
	// later it can be mounted directly using Kubernetes image volumes.
	Image string `json:"image"`
	// Path exported by the artifact, for example /opt/runtime/vllm.
	MountPath string `json:"mountPath,omitempty"`
	Digest string `json:"digest,omitempty"`
}

type RuntimeCompatibility struct {
	OS string `json:"os,omitempty"`
	Arch string `json:"arch,omitempty"`
	Python string `json:"python,omitempty"`
	CUDA string `json:"cuda,omitempty"`
	Torch string `json:"torch,omitempty"`
	NCCL string `json:"nccl,omitempty"`
	MinDriverVersion string `json:"minDriverVersion,omitempty"`
	GPUArchitectures []string `json:"gpuArchitectures,omitempty"`
}

type RuntimeComponentStatus struct {
	Phase string `json:"phase,omitempty"`
	Message string `json:"message,omitempty"`
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// ModelDeployment declares the desired model, engine and distributed execution strategy.
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=mdp
// +kubebuilder:printcolumn:name="Model",type="string",JSONPath=".spec.model.uri"
// +kubebuilder:printcolumn:name="Engine",type="string",JSONPath=".spec.engine.name"
// +kubebuilder:printcolumn:name="Executor",type="string",JSONPath=".spec.executor.name"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
type ModelDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ModelDeploymentSpec   `json:"spec,omitempty"`
	Status ModelDeploymentStatus `json:"status,omitempty"`
}

type ModelDeploymentSpec struct {
	Model ModelReference `json:"model"`
	Engine ComponentSelector `json:"engine"`
	Executor ComponentSelector `json:"executor,omitempty"`
	Runtime RuntimeRequirements `json:"runtime,omitempty"`
	Resources ModelResources `json:"resources"`
	Parallelism ParallelismSpec `json:"parallelism,omitempty"`
}

type ModelReference struct {
	URI string `json:"uri"`
	ServedName string `json:"servedName,omitempty"`
}

type ComponentSelector struct {
	Name string `json:"name"`
	Version string `json:"version,omitempty"`
	ComponentRef string `json:"componentRef,omitempty"`
}

type RuntimeRequirements struct {
	Python string `json:"python,omitempty"`
	CUDA string `json:"cuda,omitempty"`
	Torch string `json:"torch,omitempty"`
}

type ModelResources struct {
	GPUCount int32 `json:"gpuCount"`
	GPUProduct string `json:"gpuProduct,omitempty"`
	CPU corev1.ResourceList `json:"cpu,omitempty"`
}

type ParallelismSpec struct {
	TensorParallelSize int32 `json:"tensorParallelSize,omitempty"`
	PipelineParallelSize int32 `json:"pipelineParallelSize,omitempty"`
}

type ModelDeploymentStatus struct {
	Phase string `json:"phase,omitempty"`
	RuntimePlanRef string `json:"runtimePlanRef,omitempty"`
	Message string `json:"message,omitempty"`
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// RuntimePlan is the resolved, immutable-oriented execution plan produced for a ModelDeployment.
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=rtp
// +kubebuilder:printcolumn:name="Deployment",type="string",JSONPath=".spec.modelDeploymentRef"
// +kubebuilder:printcolumn:name="Resolved",type="string",JSONPath=".status.phase"
type RuntimePlan struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RuntimePlanSpec   `json:"spec,omitempty"`
	Status RuntimePlanStatus `json:"status,omitempty"`
}

type RuntimePlanSpec struct {
	ModelDeploymentRef string `json:"modelDeploymentRef"`
	ABI RuntimeABI `json:"abi"`
	Components []ResolvedRuntimeComponent `json:"components"`
	Resources ModelResources `json:"resources"`
	Parallelism ParallelismSpec `json:"parallelism,omitempty"`
}

type RuntimeABI struct {
	OS string `json:"os,omitempty"`
	Arch string `json:"arch,omitempty"`
	Python string `json:"python,omitempty"`
	CUDA string `json:"cuda,omitempty"`
	Torch string `json:"torch,omitempty"`
	NCCL string `json:"nccl,omitempty"`
}

type ResolvedRuntimeComponent struct {
	Type RuntimeComponentType `json:"type"`
	Name string `json:"name"`
	Version string `json:"version"`
	ComponentRef string `json:"componentRef"`
	Image string `json:"image"`
	Digest string `json:"digest,omitempty"`
	MountPath string `json:"mountPath,omitempty"`
}

type RuntimePlanStatus struct {
	Phase string `json:"phase,omitempty"`
	Reason string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}
