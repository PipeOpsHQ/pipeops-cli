package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

// AgentClusterConfig holds the configuration retrieved from the cluster
type AgentClusterConfig struct {
	Token       string
	ClusterName string
	APIURL      string
}

// getAgentConfigFromCluster attempts to retrieve the agent configuration from the running cluster.
// It checks for the 'pipeops-agent-config' secret in the 'pipeops-system' namespace.
func getAgentConfigFromCluster() (*AgentClusterConfig, error) {
	// Check if kubectl is available
	if _, err := exec.LookPath("kubectl"); err != nil {
		return nil, fmt.Errorf("kubectl not found")
	}

	// 1. Get the secret in JSON format
	// kubectl get secret pipeops-agent-config -n pipeops-system -o json
	cmd := exec.Command("kubectl", "get", "secret", "pipeops-agent-config", "-n", "pipeops-system", "-o", "json")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get pipeops-agent-config secret: %w", err)
	}

	// 2. Parse the secret
	var secret struct {
		Data map[string][]byte `json:"data"`
	}
	if err := json.Unmarshal(output, &secret); err != nil {
		return nil, fmt.Errorf("failed to parse secret json: %w", err)
	}

	config := &AgentClusterConfig{}

	// 3. Extract and decode values
	if tokenData, ok := secret.Data["PIPEOPS_TOKEN"]; ok {
		config.Token = string(tokenData)
	} else if tokenData, ok := secret.Data["AGENT_TOKEN"]; ok { // Fallback for older versions
		config.Token = string(tokenData)
	}

	if clusterNameData, ok := secret.Data["PIPEOPS_CLUSTER_NAME"]; ok {
		config.ClusterName = string(clusterNameData)
	} else if clusterNameData, ok := secret.Data["CLUSTER_NAME"]; ok { // Fallback
		config.ClusterName = string(clusterNameData)
	}

	if apiURLData, ok := secret.Data["PIPEOPS_API_URL"]; ok {
		config.APIURL = string(apiURLData)
	}

	if config.Token == "" {
		return nil, fmt.Errorf("token not found in secret")
	}

	log.Println("[INFO] Detected existing PipeOps agent configuration in cluster")
	return config, nil
}
