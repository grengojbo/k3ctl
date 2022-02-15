package client

import (
	"context"
	"time"

	k3s "github.com/grengojbo/k3ctl/pkg/k3s"
	"github.com/grengojbo/k3ctl/pkg/types"
	v1 "k8s.io/api/core/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type K8sClient struct {
	Config *rest.Config
	Clientset *kubernetes.Clientset
}
// NewClient return new cobernetes clent for clusterName config
func NewClient(clusterName string) (client K8sClient, err error) {
	
	// use the current context in kubeconfig
	config, err := k3s.BuildKubeConfigFromFlags(clusterName)
	if err != nil {
			return client, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err !=nil {
		return client, err
	}
	client.Clientset=clientset

	return client, err
}

// GetClusterStatus cluster status
func (c *K8sClient) GetClusterStatus() string {
	// _, err := c.Clientset.RESTClient().Get().Timeout(15 * time.Second).RequestURI("/readyz").DoRaw(context.TODO())
	_, err := c.Clientset.RESTClient().Get().Timeout(types.NodeWaitForLogMessageRestartWarnTime).RequestURI("/readyz").DoRaw(context.TODO())
	if err != nil {
		return types.ClusterStatusStopped
	}
	return types.ClusterStatusRunning
}

// IsReady To see if a Node is Ready
func IsReady(node *v1.Node) bool {
	var cond v1.NodeCondition
	for _, n := range node.Status.Conditions {
			if n.Type == v1.NodeReady {
					cond = n
					break
			}
	}

	return cond.Status == v1.ConditionTrue
}

// IsMaster To check if a Node is a master node, search if the label “node-role.kubernetes.io/master” is present.
func IsMaster(node *v1.Node) bool {
	if _, ok := node.Labels["node-role.kubernetes.io/master"]; ok {
			return true
	}
	return false
}

func GetStatus(node *v1.Node) (status string) {
	conditions := node.Status.Conditions
	for _, c := range conditions {
		if c.Type == v1.NodeReady {
			if c.Status == v1.ConditionTrue {
				return types.StatusRunning
			} else {
				return types.ClusterStatusStopped
			}
		}
	}
	return types.StatusFailed
}

// ListNodes Read more on Kubernets Nodes. In client-go, NodeInterface includes all the APIs to deal with Nodes.
// func ListNodes() ([]v1.Node, error) {
// 	nodes, err := c.Clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})

// 	if err != nil {
// 			return nil, err
// 	}

// 	return nodes.Items, nil
// }


// Age Calculate the age of a Node from the CreationTimestamp.
func (c *K8sClient) Age(node *v1.Node) uint {
	diff := uint(time.Now().Sub(node.CreationTimestamp.Time).Hours())
	return uint(diff / 24)
}

// GetClusterVersion get kube cluster version.
func GetClusterVersion(c *kubernetes.Clientset) string {
	v, err := c.DiscoveryClient.ServerVersion()
	if err != nil {
		return types.ClusterStatusUnknown
	}
	return v.GitVersion
}

// // DescribeClusterNodes describe cluster nodes.
// func DescribeClusterNodes(client *kubernetes.Clientset, instanceNodes []types.ClusterNode) ([]types.ClusterNode, error) {
// 	// list cluster nodes.
// 	timeout := int64(5 * time.Second)
// 	nodeList, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{TimeoutSeconds: &timeout})
// 	if err != nil || nodeList == nil {
// 		return nil, err
// 	}
// 	for _, node := range nodeList.Items {
// 		var internalIP, hostName string
// 		addressList := node.Status.Addresses
// 		for _, address := range addressList {
// 			switch address.Type {
// 			case v1.NodeInternalIP:
// 				internalIP = address.Address
// 			case v1.NodeHostName:
// 				hostName = address.Address
// 			default:
// 				continue
// 			}
// 		}
// 		for index, n := range instanceNodes {
// 			isCurrentInstance := false
// 			for _, address := range n.InternalIP {
// 				if address == internalIP {
// 					isCurrentInstance = true
// 					break
// 				}
// 			}
// 			if !isCurrentInstance {
// 				if n.InstanceID == node.Name {
// 					isCurrentInstance = true
// 				}
// 			}
// 			if isCurrentInstance {
// 				n.HostName = hostName
// 				n.Version = node.Status.NodeInfo.KubeletVersion
// 				n.ContainerRuntimeVersion = node.Status.NodeInfo.ContainerRuntimeVersion
// 				// get roles.
// 				labels := node.Labels
// 				roles := make([]string, 0)
// 				for role := range labels {
// 					if strings.HasPrefix(role, "node-role.kubernetes.io") {
// 						roleArray := strings.Split(role, "/")
// 						if len(roleArray) > 1 {
// 							roles = append(roles, roleArray[1])
// 						}
// 					}
// 				}
// 				if len(roles) == 0 {
// 					roles = append(roles, "<none>")
// 				}
// 				sort.Strings(roles)
// 				n.Roles = strings.Join(roles, ",")
// 				// get status.
// 				conditions := node.Status.Conditions
// 				for _, c := range conditions {
// 					if c.Type == v1.NodeReady {
// 						if c.Status == v1.ConditionTrue {
// 							n.Status = "Ready"
// 						} else {
// 							n.Status = "NotReady"
// 						}
// 						break
// 					}
// 				}
// 				instanceNodes[index] = n
// 				break
// 			}
// 		}
// 	}
// 	return instanceNodes, nil
// }